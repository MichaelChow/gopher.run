---
title: "3.1 goroutine"
date: 2025-04-28T23:53:00Z
draft: false
weight: 3001
---

# 3.1 goroutine

- 发音：/ˈɡoʊruːtiːn/ （国际音标）"**够-如-挺**”
    - "go" 发 /ɡoʊ/（类似中文"够"的发音）
    - "routine" 发 /ruːtiːn/（类似中文"如听"的快速连读）
- 并发编程表现为 程序由若干个自主的活动单元组成，如Web服务器可以一次处理数千个请求；传统的批处理任务——读取数据、计算、将结果输出——也使用并发来隐藏I/O操作的延迟；CPU内核的个数每年变多，但是速度没什么变化；
- 并发编程在本质上也比顺序编程要困难一些，从顺序编程获取的直觉让我们在并发编程上加倍地迷茫；学习之初可以暂时假设goroutine类似于操作系统的线程，goroutine和线程之间在数量上有非常大的差别；
- Go有两种并发编程的风格：
    - goroutine、通道(channel)：支持通信顺序进程(Communicating Sequential Process,CSP)的并发模式，在不同的执行体(goroutine)之间传递值，但是变量本身局限于单一的执行体；
    - 共享内存多线程的传统模型：和在其他主流语言中使用线程类似；
- 当一个程序启动时，只有一个goroutine来调用main函数，称它为 main **goroutine；**
- 新的goroutine通过在函数或方法调用前，加go关键字前缀进行创建。**go语句使函数或方法在一个新创建的goroutine中调用，go语句本身的执行立即完成**。
    ```go
    f()
    go f() // go语句本身的执行立即完成，并不等待f()的return
    ```
- 示例：spinner.go
    -  **main函数返回时，所有的goroutine都暴力地直接终结，然后程序退出；**
    - **除了从main返回或者退出程序之外，没有程序化的方法让一个goroutine来停止另一个，但有办法和goroutine通信来要求它自己停止;**
    ```go
    func main() {    // main goroutine将计算菲波那契数列的第45个元素值。
        go spinner(100 * time.Millisecond)
        const n = 45
        fibN := fib(n) // slow
        fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN) 
    func spinner(delay time.Duration) {
        for {
            for _, r := range `-\|/` {
                fmt.Printf("\r%c", r)
                time.Sleep(delay)
            }
        }
    }
    func fib(x int) int {
        if x < 2 {
            return x
        }
        return fib(x-1) + fib(x-2)
    }
    ```
- **示例:**clock.go
    - 格式化模板限定为Mon Jan 2 03:04:05PM 2006 UTC-0700。有8个部分（周几、月份、一个月的第几天……）。可以以任意的形式来组合前面这个模板；出现在模板中的部分会作为参考来对时间格式进行输出。
    - 这是go语言和其它语言相比比较奇葩的一个地方。你需要记住格式化字符串是：**1月2日下午3点4分5秒零六年UTC-0700****（记忆：**1234567**）**，而不像其它语言那样Y-m-d H:i:s一样
        ```go
        // Clock1 is a TCP server that periodically writes the time.
        package main
        func main() {
            listener, err := net.Listen("tcp", "localhost:8000")  // 监听8000端口
            if err != nil {
                log.Fatal(err)
            }
            for {
                conn, err := listener.Accept()  // Accept方法会直接阻塞，直到一个新的连接被创建，然后会返回一个net.Conn对象来表示这个连接
                if err != nil {
                    log.Print(err) // e.g., connection aborted
                    continue
                }
                go handleConn(conn) // 仅仅在 函数调用的地方**增加go关键字**，让每一次handleConn的调用都进入自己的一个独立的goroutine内执行
            }
        }
        func handleConn(c net.Conn) {
            defer c.Close()   // 关闭服务器侧的连接，然后返回到主函数，继续等待下一个连接请求
            for {   // 死循环会一直执行，直到写入失败，如可能的原因是客户端主动断开连接
        		    // 
                _, err := io.WriteString(c, time.Now().Format("\r15:04:05"))  // 由于net.Conn实现了io.Writer接口，我们可以直接向其写入内容。 \r (回车,Carriage Return，CR): 将光标回到当前行的行首(而不会换到下一行),之后的输出会把之前的输出覆盖
                if err != nil {
                    return // e.g., client disconnected
                }
                time.Sleep(1 * time.Second)
            }
        }
        ```
    - **阻塞执行/顺序编程：服务器顺序执行，第二个nc客户端接收不到时间；**
        ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/15467f0b-76de-44dd-ae9a-238eb80604b3/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4662KFCKZAR%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005559Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIECQ4Uz3xTB3Eq4pHae4AYjXw0QiwgRjnA2vOfp95PquAiEA%2BM0UkwhXAPf93tekMSMNZqdoK2SP7l1OYyMmhzJR5kYqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDNbH1LRiQBh9PzrDrSrcA5d2xdFIj8IM8lbfWS3gbZCMhJXLkcBGrJcVyOx3itimdslw7bsOA%2BlXufI%2BrJSdcSxJnEnHUFASNPpXARxbpr9a1OEqJgHErccNCl1U6hIBlij%2BpOlz119rx2h55tnQ56evN7VDVwXDH92qREpSTMB8c2V5KliYKNiNvUuSMk787CxY1qYauVZFFEvlAoR1hjdbKg3EMnjzD37BMp7snBZugp9oosG3FwlyB0KM8rdA7tBKeXZM7KjJ%2BrtzJ7KHBMpsQFKMZala4ICG19iGw9pWXUzz9Qwghc2yj4pWnD9huquMHbIQ9QejxF88TgNMj9h7vZE2k5gF%2FNz%2F8VvwxoYC9rLUdTk4JcNJMNQAQP3j7DwfJz%2BDMTRPXg1veBtRWH8wUckF2JmGgoE0ins6UsE8TcCgex0Nt447bjs0a7IaxXYRinffsIqeIEhmAvmjwYQbO%2FyUL6YJbP3QgaO23IiNeJ1rYIyEgR5h%2B0KHbj73h1tEqpKv4wrWHxtUUvlCddpKS%2FxXJku3ffYP9n%2FSp8uadce88x17RxuVRiP91VKvgJEy1QnSDSFrhVVrjCEZYdQ8K9cXNFQ2q27Tl8qKxxB41YmEZWzlvbjfXLQqOf1iLGzMFCl%2BAH5pq1I1MIq668MGOqUBRdNhHO7j78P0UToI6C7j2XLd%2BD5NQyCtH2y8WH%2FQ6m56yBz4n8Et5RWHKwATBA%2FLqO5rwrLLnD2fg3elpCiGmr7wdCSZ7AkcImhl7CJATshafcHWpQ0Q7yzRIoOzkK0Tfu7TTDHxULbavkngAgqwufr0lyYbfu0aKtzsmPAHWKHHoNyk08QJ6kmo1Qp95cC2wbB2geSLCI%2B6Nl%2BZGP8vzf68gsBx&X-Amz-Signature=30d31754f20a40a854ece62157935a81c755c6e34c069e070c6cf5ce767f57b1&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
    - **并发执行/并发编程：多个客户端可以同时接收到时间；**
        ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/bb3fab5e-7a8b-4cd6-a3a1-1c765bf02310/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB4667SH3XXPW%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005559Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIHbiyQeRuw1RWJCk%2FGxRO3Lu5v7feNcODB%2Bbqy7beQ6LAiEArf9Mr%2FFEv3ZWLk%2BnYQkv6IXQvzCgQTbARCCHkLGlOzcqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDKRZwL4aK3PMAFxTbSrcA3IUyr5B9gC%2FeiQNoRe9Y5CqY4hV9cwBfibxc0nIkq%2F8qsUJztB01Vw0NpvZo%2FbPiBNU7NFVcebqb1wAh1iqpgvvaezOAvEN7ux5wh09DsPNF6VR8yPLgUu0VVD2ULTRAg%2BbFmfYpNF%2Fw6FQGz4%2F9OinVyFtWoKF4t16eNfwhPQx3aqfss2yXuznAM1f05z0e5IWwZ7sMV5IHeNUklQQBC8ZhWK3eq7WveCQ5xkNhx3%2B5WKCOTNcS%2Fs6WGRCo%2BQ3czeUt0kathnP8zJwfFVpvbHPUrnMvgfbQLJE9r6HO%2FuE7ieL%2B0YT8NjyiE%2FpzT6ldNidX94hWn%2BdLWyZolaJ9kt7lYIhksE%2B2cJ3qF%2BWOIlce6RE9y%2BihWyQiNmmTdFoI6mlPYkLUQTNIY5FjH7Ab2LkHlEQJ6AJbz0oLK6FsSyODuN%2Bw6chMVEJveP3b4bSx83yEqX7YZYNmsAMNW5GOJQZr5cEeTnib8PWRLa0Yty3ATgbRC5CSW%2Br90Juw2LYzcH5tCYlzkTjo4YnF9ndBKJEedjh%2Fwr9swvYtdJF6vKSGoBqYlWZJJi9H%2Fhh8NT%2FeL8AlSL1mDxo6YHHhjY8R8ThpIfVJSl9%2FteM8%2BnvJemlJnclIPFiZ4rJy3hqMLe668MGOqUBotjCFI20WKCRN36hSOSxSUp7IsJXn1dD%2BrEVcGQ%2FKZPhsIShjVBM6uH8tCBdgyRvMAza7R2vuQ3pL5m%2BAn4Ug5eRemE5adkRdko5fwuKJfmcGuTcp845KScylQl0WKXUerNvbrMUqD5d1JARGPaI64QZNjz58qyqQTHn78X12bA3XY9rW4ICRow1nKDiysoHaERQTYqGD8tjFHiCc5UHINgZpqtE&X-Amz-Signature=13549389c58ee6f8104f62d976b5c61c387438e2b1785f0b8baff6158f0429ef&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
- **示例: 并发的Echo服务**reverb.go
    - **go后跟的函数的参数表达式会在go语句自身执行时(main goroutine中)被求值**
    ```go
    func echo(c net.Conn, shout string, delay time.Duration) {
        fmt.Fprintln(c, "\t", strings.ToUpper(shout))
        time.Sleep(delay)    // 由于这里设置了dalay，客户端多次发送不同的消息，所有echo的回显会顺序的执行，程序非常慢
        fmt.Fprintln(c, "\t", shout)
        time.Sleep(delay)
        fmt.Fprintln(c, "\t", strings.ToLower(shout))
    }
    func handleConn(c net.Conn) {
        input := bufio.NewScanner(c)
        for input.Scan() {
            go echo(c, input.Text(), 1*time.Second)    // 所有echo并发的执行，程序非常快
        }
        // NOTE: ignoring potential errors from input.Err()
        c.Close()
    }
    ```


