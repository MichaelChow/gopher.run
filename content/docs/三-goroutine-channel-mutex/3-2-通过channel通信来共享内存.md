---
title: "3.2 通过channel通信来共享内存"
date: 2024-12-28T15:08:00Z
draft: false
weight: 3002
---

# 3.2 通过channel通信来共享内存



- **goroutine是Go程序并发的执行体，****channel（**/ˈtʃænl/，通道**）****是它们之间的连接**。 
- 通道是可以让**一个goroutine发送特定值到另一个goroutine的通信机制**。**每一个通道是一个具体类型的导管，叫作通道的元素类型**。如一个有int类型元素的通道写为chan int；
- **像map一样，通道是一个使用内置的make函数创建的数据结构的引用（**pointer、slice、map、function、**channel****为引用类型）。**
    - 当复制或者作为参数传递到一个函数时，复制的是引用，这样调用者和被调用者都引用同一份数据结构；
    - 和其他引用类型一样，通道的零值是nil；
    - **可比较性**：同种类型的通道可以使用==符号进行比较：**当二者都是同一通道数据的引用时，比较值为true；**通道也可以和nil进行比较；
    ```go
    ch := make(chan int) // ch的类型: 'chan int' ， unbuffered channel（无缓冲通道）
    ch = make(chan int, 0) // unbuffered channel
    ch = make(chan int, 3) // 通道容量3的缓冲通道
    ```
- 通道有三个主要操作，都使用`<-`运算符，send、receive统称为通信：
    - **发送(send)语句**：`ch <- x` ，从一个goroutine传输一个值到另一个在执行接收表达式的goroutine；
    - **接收(receive)语句**：`x = <- ch` ；`<-ch` 一个不使用接收结果的接收操作也是合法的；
    - **关闭(close)**：`close(ch)` 设置一个标志位来指示值当前已经发送完毕，这个通道后面没有值了；
        - 关闭后的再次**send**发送操作将导致panic宕机；
        - 在一个已经关闭的通道上进行**receive**接收操作，将获取所有已经发送的值，直到通道为空；这时任何接收操作会立即完成，同时获取到一个通道元素类型对应的零值；
## 1. 无缓冲通道/同步通道

- 无缓冲通道上的发送操作将会阻塞，直到另一个goroutine在对应的通道上执行完接收操作，这时值传送完成，两个goroutine都可以恢复继续执行；
- 如果接收操作先执行，接收方goroutine将阻塞，直到另一个goroutine在同一个通道上发送一个值；
- 使用无缓冲通道进行的通信导致发送和接收goroutine同步化（同步通道）。当一个值在无缓冲通道上传递时，接收值后发送方goroutine才被再次唤醒；
- ***happens before：***Go语言并发内存模型的一个关键术语
    - 在讨论并发的时候，当我们说x早于y发生时，不仅仅是说x发生的时间早于y，而是说确定它是这样，并且是可预期的，如更新变量；我们可以放心的依赖这个机制；
    - 当x既不比y早也不比y晚时，我们说x和y并发。这不意味着，x和y一定同时发生，只说明我们无法确定它们的顺序。（费曼学习法：量子不确定性）；
- 通过通道发送消息有两个重要的方面需要考虑：每一条消息有一个值，但有时候也强调通信本身、通信发生的时间，此时通常把消息叫作**事件(event)**（费曼学习法：OpenAPI WebHook来订单的事件）。当事件没有携带额外的信息时，它单纯的目的是进行同步。我们通过使用一个`struct{}`元素类型的通道来强调它，尽管通常使用bool或int类型的通道来做相同的事情，因为`done<-1`比`done<-struct{}{}`要短；
- 示例：*netcat3.go*
    ```go
    func main() {
        conn, err := net.Dial("tcp", "localhost:8000")
        if err != nil {
            log.Fatal(err)
        }
        done := make(chan struct{})
        go func() {   // go语句调用了一个函数字面量（匿名函数），这是Go语言中启动goroutine常用的形式
            io.Copy(os.Stdout, conn) // NOTE: ignoring errors
            log.Println("done")
            done <- struct{}{} // signal the main goroutine
        }()
        mustCopy(conn, os.Stdin)
        conn.Close()
        <-done // wait for background goroutine to finish
    }
    ```
### **管道/Pipeline**

- **管道(pipeline)**：通道可以用来连接goroutine，上一个通道的输出是下一个通道的输入；
    - 第一个goroutine是counter，产生一个0, 1, 2，…的整数序列；然后通过一个管道发送给第二个goroutine（叫square），计算数值的平方；然后将结果通过另一个通道发送给第三个goroutine（叫printer），接收值并输出它们。
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/c20c500e-7819-49e7-83f2-7e75f6dea228/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4663MIDRUW3%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005602Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQDtJMNb4o8mzx%2BtIYHkgiesScztn0CNJe6ISGerZ%2BRVPgIhAPhVCCuMq%2BvThJYPpgYRHTOUkJD4lORpvZTHZWm3mABOKogECJn%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgzzZ%2FNfadZa7nwf0fIq3ANer2Oh%2Fer3APiCuu9Ibc%2Fe%2BOcIyjf9uT0l%2FrRK%2F8lvkmVO%2BkW1Bv5IRmYxsQzjXNfPvq6IzCi2BYPFWt6MDm3ix%2BNhs8XoYl1a1Wa4YBvLD6HLj%2FRx0EUAsxKOwrFvWf2h9HEWWoV2ccUOdcIr%2Fa6%2Bw71x6%2Fyt0Hd5eVByx%2BtjRWZq2JDkhQDwDRRaSbJbTpR3ccDMhKVQim3skzSqI7fVKbRuQpKFYuEs1EqygszhAhVDdGG9Fmr36AHIt8OPeNYWzUwd1EMgcfssNudcefFuUaQVUW6MMoV9WcqnGZVOaZ1QKctfSj8cm%2B%2BSa6vwXAMmGsigf6C7EigaWhHsaJGFzavOMnnhTgTuBFEhN8hqqUWK66gfuaTzKEOabdgIxlmnV381IGnW0uEwu4okLh%2B2TJRQJdigwNVlGB7lB06MSn7lGz3xedDWeFj3S4GNqe00YwhtLFVDtMnY0DPFiOJcAqWOXLgcMlxWyoQckK3mqTzKnLLAoJ%2BUMXUXTQ8zghNhQ1atj9bLKrFKJy3Jmwe49abVTHw1tDXj2pG%2BmHPyrH2q1eOTRhJGSE1gbD3K4yRb32KdSfO5tRBfhBSG85BggoAiDcE0FkjS3vPAhaOetYpM3L6pDrAYQ3sfNTDXuuvDBjqkAUzb8oIQsKzbRTAO47zN6xu%2BoN6ASfDqUdLNIeGYv2xHf7pZj5TKZnvLS18PNyK%2FFECuEn4AeacPuztG4b3xXb8MVAj%2FKutIAcfQ90fZ90aqCnoRj%2BGcFLrMQFQfEpp44OAcqhYRY6bcdD4Vt93aUyKbZtFQx57nDTL9HFFjAPBnNyQIAhAmUB08U3IZ97RHfNVNgM299XjpfxIZCLZkP06b8zDH&X-Amz-Signature=25515b265ff5c9231151bf9a3f4e4ff6caabb2b83484bd5bf2285f1a5772cc18&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    - 像这样的管道出现在长期运行的服务器程序中，其中通道用于在包含无限循环的goroutine之间整个生命周期中的通信。
    - 没有一个直接的方式来判断是否通道已经关闭，但是这里有接收操作的一个变种，它产生两个结果：接收到的通道元素，以及**一个布尔值（通常称为ok），它为true的时候代表接收成功**，**false表示当前的接收操作在一个关闭的并且读完的通道上。**
- 示例：pipeline1.go
    ```go
    func main() {
        naturals := make(chan int)
        squares := make(chan int)
        go func() {            // goroutine1：计数器Counter
            for x := 0; ; x++ {
                naturals <- x
            }
        }()
        
        go func() {           // goroutine2：求平方Squarer
            for {
                x := <-naturals
                squares <- x * x
            }
        }()
        
        for {                  // main goroutine：打印Printer (in main goroutine)
            fmt.Println(<-squares)
        }
    }
    ```
- **示例：pipeline2.go**
    - 因为上面的for死循环语法比较笨拙，而模式又比较通用，所以**Go提供了range循环语法以在通道上迭代**。这个语法更方便接收在通道上所有发送的值，接收完最后一个值后关闭循环。
        ```go
        func main() {
            naturals := make(chan int)
            squares := make(chan int)
            // Counter
            go func() {
                for x := 0; x < 100; x++ {
                    naturals <- x
                }
                close(naturals)
            }()
            
              
        		go func() {            // Squarer
        		    for {
        		        x, ok := <-naturals     // channel被关闭并且没有值可接收时跳出循环
        		        if !ok {
        		            break // channel was closed and drained
        		        }
        		        squares <- x * x
        		    }
        		    close(squares)
        		}()
            go func() {           
                for x := range naturals {    // 
                    squares <- x * x
                }
                close(squares)
            }()
            // Printer (in main goroutine)
            for x := range squares {
                fmt.Println(x)
            }
        }
        ```


- 通道的close：
    - 关闭每一个通道不是必需的，**close操作只用于 断言/通知接收方goroutine 不再向channel发送新的数据****，所以仅仅在发送方的goroutine上才能调用close函数；**因此close一个只接收的channel将是一个编译错误；
    - 通道也是可以**通过GC垃圾回收器在它没有被引用时回收它**（而不是根据它是否关闭）；（**不要将这个close操作和对于文件的close操作混淆，**当结束的时候对每一个文件调Close方法是非常重要的）
    - 试图关闭一个已经关闭的通道会导致panic宕机，就像关闭一个nil值的空通道也会导致panic宕机。关闭通道还可以作为一个广播机制；
### **单向通道类型**

- 当一个channel作为一个函数参数时，它一般总是被专门用于只发送或者只接收。为了防止被滥用，Go语言的类型系统提供了单方向的channel类型，分别用于只发送或只接收的channel。
    - `chan<- int`类型：表示一个只发送int的channel，只能发送不能接收。
    - `<-chan int`类型：表示一个只接收int的channel，只能接收不能发送。（箭头`<-`和关键字chan的相对位置表明了channel的方向。）这种限制将在编译期检测。
- 任何双向通道向单向通道变量的赋值操作都将导致该**隐式转换**，**但反过来是不行的**，一旦有一个像`chan<- int`这样的单向通道，是没有办法通过它获取到引用同一个数据结构的`chan int`数据类型的
    ```go
    func counter(out chan<- int) {
        for x := 0; x < 100; x++ {
            out <- x
        }
        close(out)
    }
    func squarer(out chan<- int, in <-chan int) {
        for v := range in {
            out <- v * v
        }
        close(out)
    }
    func printer(in <-chan int) {
        for v := range in {
            fmt.Println(v)
        }
    }
    func main() {
        naturals := make(chan int)
        squares := make(chan int)
        go counter(naturals)   // naturals的类型将隐式地从chan int转换成chan<- int。
        go squarer(squares, naturals)
        printer(squares)  // 类型将隐式地转换为<-chan int类型只接收型的channel
    }
    ```
## **2. 缓冲通道**

- 缓冲通道有一个元素队列，队列的最大长度在创建的时候通过make的容量参数来设置：`ch = make(chan string, 3)` 创建一个可以容纳三个字符串的缓冲通道，并指向它的引用；
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/dd5e3997-a8a9-4347-a84a-321634b82ce6/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667L2UZBAN%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005608Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIQC6sLM7ge1VKXP9RlnUk7zDsnKsG3kWxHpSBfOzYn0xtQIgXNbrxa3sOwJtFD4Iwh3S9prTONe06AU8O2EArdlr%2B34qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDEiocnVmPjybUZ5p7ircAze7OU8H5JvsKtaKbvp4Yf2XrYuO9d1iWfrFEnwexILOp6Hqq1q0%2FzbFoGNjkRGfDSkwLNaRqJmrKcHq8Ml6GTEQLYFoof8V7%2B4jUDSyHfHG0%2F%2FlnveibndpGFdBUTe3blJgrC2z%2BOHl4UImrkqQRzURZnIDIA9TL4sRCXdGdNBz2LOaFJzzKi5XUTfXQp%2BY0hZmD%2BHLX5WS8oh%2BcrsxOIkv92cvbOG3fN5YvNQBRIF8vcPw6Eo4VSL6R%2FnMXuf4HP4qH%2Bf%2Bpt5VeEzqTvWXPHWBEcfhAawU8cKu68g6%2B3nEH3GOBaondcYbP3Hvmqtfy%2BqR12kJrr0iPSaq%2ByAMH7QfGByTPX4DdZyRNxGu6hw3Awn8gmSmRn9Md6WihY%2BydCpeOvAmCJWqTjzUyYh0wjzoMXdW2YSrLsgML1h2PeYjoJWto9amnPLNCA0eBtrJby3C0A9%2BHlWONa0s0yvakTPLuAvSX768T5dTq6a9so%2BgpAa2Mus2jyODDiNbMbQlsvDj5ka%2FpnC87lBkLdIhQzwzc8qJfhHx253wI2FHcjoufrHVebzVsYDiNgbFqQ0CwShuP3TioVhefWnWtLVWXP4B0%2BIPDX9H6WzQslTx%2FgShrHg%2FxbdNbSTfM%2FIeMMG668MGOqUB92L65d121htNhrGhKZzew71ORmpxsryGFhqP9EnwjIe0cfAiH5Uj9tepyTsML4QFh6cMroAG3lp2KTUUhLcVwXy0xuIhOJ9%2BoG9GLVrpIUu%2FhJawAO%2BpDN9y3SiO%2BrH%2FC4jjdnvvvjDMEILtymAj1xMibdcrhKm6QypVKg9tcGoEkHNzaR6ZUFNydrbiLXuzSAx7gtQVx7Tl14SvDSAkp6xhsIHs&X-Amz-Signature=0ed78da8c0387ce17b0f69b6dc3f964a40d33ab0fb78180d0a18044cf282d496&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/040e71d9-174d-4e9e-8eab-06a3cb8c0007/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667L2UZBAN%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005608Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIQC6sLM7ge1VKXP9RlnUk7zDsnKsG3kWxHpSBfOzYn0xtQIgXNbrxa3sOwJtFD4Iwh3S9prTONe06AU8O2EArdlr%2B34qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDEiocnVmPjybUZ5p7ircAze7OU8H5JvsKtaKbvp4Yf2XrYuO9d1iWfrFEnwexILOp6Hqq1q0%2FzbFoGNjkRGfDSkwLNaRqJmrKcHq8Ml6GTEQLYFoof8V7%2B4jUDSyHfHG0%2F%2FlnveibndpGFdBUTe3blJgrC2z%2BOHl4UImrkqQRzURZnIDIA9TL4sRCXdGdNBz2LOaFJzzKi5XUTfXQp%2BY0hZmD%2BHLX5WS8oh%2BcrsxOIkv92cvbOG3fN5YvNQBRIF8vcPw6Eo4VSL6R%2FnMXuf4HP4qH%2Bf%2Bpt5VeEzqTvWXPHWBEcfhAawU8cKu68g6%2B3nEH3GOBaondcYbP3Hvmqtfy%2BqR12kJrr0iPSaq%2ByAMH7QfGByTPX4DdZyRNxGu6hw3Awn8gmSmRn9Md6WihY%2BydCpeOvAmCJWqTjzUyYh0wjzoMXdW2YSrLsgML1h2PeYjoJWto9amnPLNCA0eBtrJby3C0A9%2BHlWONa0s0yvakTPLuAvSX768T5dTq6a9so%2BgpAa2Mus2jyODDiNbMbQlsvDj5ka%2FpnC87lBkLdIhQzwzc8qJfhHx253wI2FHcjoufrHVebzVsYDiNgbFqQ0CwShuP3TioVhefWnWtLVWXP4B0%2BIPDX9H6WzQslTx%2FgShrHg%2FxbdNbSTfM%2FIeMMG668MGOqUB92L65d121htNhrGhKZzew71ORmpxsryGFhqP9EnwjIe0cfAiH5Uj9tepyTsML4QFh6cMroAG3lp2KTUUhLcVwXy0xuIhOJ9%2BoG9GLVrpIUu%2FhJawAO%2BpDN9y3SiO%2BrH%2FC4jjdnvvvjDMEILtymAj1xMibdcrhKm6QypVKg9tcGoEkHNzaR6ZUFNydrbiLXuzSAx7gtQVx7Tl14SvDSAkp6xhsIHs&X-Amz-Signature=1fe9b6842b58dea884fd3745fcc568d14e7f636a313e2ae2c2a039e42404dea6&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/e9a36d67-7b31-47c3-8029-86ce3d8951b6/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667L2UZBAN%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005608Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIQC6sLM7ge1VKXP9RlnUk7zDsnKsG3kWxHpSBfOzYn0xtQIgXNbrxa3sOwJtFD4Iwh3S9prTONe06AU8O2EArdlr%2B34qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDEiocnVmPjybUZ5p7ircAze7OU8H5JvsKtaKbvp4Yf2XrYuO9d1iWfrFEnwexILOp6Hqq1q0%2FzbFoGNjkRGfDSkwLNaRqJmrKcHq8Ml6GTEQLYFoof8V7%2B4jUDSyHfHG0%2F%2FlnveibndpGFdBUTe3blJgrC2z%2BOHl4UImrkqQRzURZnIDIA9TL4sRCXdGdNBz2LOaFJzzKi5XUTfXQp%2BY0hZmD%2BHLX5WS8oh%2BcrsxOIkv92cvbOG3fN5YvNQBRIF8vcPw6Eo4VSL6R%2FnMXuf4HP4qH%2Bf%2Bpt5VeEzqTvWXPHWBEcfhAawU8cKu68g6%2B3nEH3GOBaondcYbP3Hvmqtfy%2BqR12kJrr0iPSaq%2ByAMH7QfGByTPX4DdZyRNxGu6hw3Awn8gmSmRn9Md6WihY%2BydCpeOvAmCJWqTjzUyYh0wjzoMXdW2YSrLsgML1h2PeYjoJWto9amnPLNCA0eBtrJby3C0A9%2BHlWONa0s0yvakTPLuAvSX768T5dTq6a9so%2BgpAa2Mus2jyODDiNbMbQlsvDj5ka%2FpnC87lBkLdIhQzwzc8qJfhHx253wI2FHcjoufrHVebzVsYDiNgbFqQ0CwShuP3TioVhefWnWtLVWXP4B0%2BIPDX9H6WzQslTx%2FgShrHg%2FxbdNbSTfM%2FIeMMG668MGOqUB92L65d121htNhrGhKZzew71ORmpxsryGFhqP9EnwjIe0cfAiH5Uj9tepyTsML4QFh6cMroAG3lp2KTUUhLcVwXy0xuIhOJ9%2BoG9GLVrpIUu%2FhJawAO%2BpDN9y3SiO%2BrH%2FC4jjdnvvvjDMEILtymAj1xMibdcrhKm6QypVKg9tcGoEkHNzaR6ZUFNydrbiLXuzSAx7gtQVx7Tl14SvDSAkp6xhsIHs&X-Amz-Signature=2c88e55f7f8db75a6af5db2f00e89e18bc63c52008d94159bdfd83538f9309e4&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    - 缓冲通道上的发送操作在队列的尾部插入一个元素；如果通道满了，发送操作会阻塞所在的goroutine直到另一个goroutine对它进行接收操作来留出可用的空间 （降级为无缓冲通道/同步通道？）；
    - 接收操作从队列的头部移除一个元素；如果通道是空的，执行接收操作的goroutine阻塞，直到另一个goroutine在通道上发送数据；
    - 通过这个方式，**通道的缓冲区解耦了发送goroutine和接收goroutine**；
- 在某些特殊情况下，程序需要知道通道缓冲区存的容量和元素个数：
    - 通道缓冲区的容量： `fmt.Println(cap(ch))`
    - 通道缓冲区的有效元素个数： `fmt.Println(len(ch))` ； 在并发程序中元素个数会随着接收操作而立即失效，但是它对某些故障诊断和性能优化会有帮助；
- **因为语法简单，Go新手粗暴地将缓冲通道作为队列在单个goroutine中使用，但是这是个严重错误**。**通道和goroutine的调度深度关联，如果没有另一个goroutine从通道进行接收，发送者（也许是整个程序）有被****永久阻塞****的风险。****如果仅仅需要一个简单的队列，使用slice创建一个就可以；**
- 示例：**goroutines泄漏****（费曼学习法：类似内存泄漏）：**如果使用一个无缓冲通道，**两个比较慢的goroutine由于发送响应结果到通道的时候没有goroutine来接收而将被永远卡住的bug；**
    - **和回收变量不同，泄漏的goroutines不会自动回收**，因此必须确保每个goroutine在不再需要的时候可以自动结束；
    ```go
    func mirroredQuery() string {
        responses := make(chan string, 3)
        go func() { responses <- request("asia.gopl.io") }()       // 并发地向三个镜像站点发出请求，三个镜像站点分散在不同的地理位置，它们分别将收到的响应发送到带缓存channel
        go func() { responses <- request("europe.gopl.io") }()
        go func() { responses <- request("americas.gopl.io") }()
        return <-responses // return the quickest response      // 最后接收者只接收第一个收到的（最快的）响应，mirroredQuery函数可能在另外两个响应慢的镜像站点响应之前就返回了结果。（顺便说一下，多个goroutines并发地向同一个channel发送数据，或从同一个channel接收数据都是常见的用法。）
    }
    func request(hostname string) (response string) { /* ... */ }
    ```


- 无缓冲通道和缓冲通道的选择、缓冲通道容量大小的选择，都会对程序的正确性产生影响。
    - 无缓冲通道提供强同步保障，因为每一次发送都需要和一次对应的接收同步；
    - 对于缓冲通道，这些操作则是解耦的;如果我们知道要发送的值数量的上限，通常会创建一个容量是使用上限的缓冲通道，在接收第一个值前就完成所有的发送。**在内存无法提供缓冲容量的情况下，可能导致程序死锁。**
    - 通道的缓冲也可能影响程序的性能：组装流水线是对于通道和goroutine合适的比喻
        - 想象蛋糕店里的三个厨师，在生产线上，在把每一个蛋糕传递给下一个厨师之前，一个烤，一个加糖衣，一个雕刻。在空间比较小的厨房，每一个厨师完成一个蛋糕流程，必须等待下一个厨师准备好接受它；这个场景类似于使用**无缓冲通道**来通信。
        - 如果在厨师之间有可以放一个蛋糕的位置，一个厨师可以将制作好的蛋糕放到这里，然后立即开始制作下一个，这类似于使用一个容量1的缓冲通道。只要厨师们以相同的速度工作，大多数工作就可以快速处理，**消除他们各自之间的速率差异**。
        - 如果在厨师之间有更多的空间——更长的缓冲区——就可以**消除更大的暂态速率波动而不影响组装流水线**，比如当一个厨师稍作休息时，后面再抓紧跟上进度。
        - 另一方面，如果生产线的上游持续比下游快，**缓冲区满的时间占大多数**。如果后续的流程更快，**缓冲区通常是空的**。这时缓冲区的存在是没有价值的。
        - 如果第二段更加复杂，一个厨师可能跟不上第一个厨师的供应，或者跟不上第三个厨师的需求。为了解决这个问题，我们可以雇用另一个厨师来帮助第二段流程，独立地执行同样的任务。这个类似于创建另外一个goroutine使用同一个通道来通信。
