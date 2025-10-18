---
title: "3.1 CSP"
date: 2025-04-28T23:53:00Z
draft: false
weight: 3001
---

# 3.1 CSP



**并发程序指同时运行多个任务的程序**，如Web服务器可以一次处理数千个请求。CPU内核的个数每年变多，但是速度却没什么变化。

**从线性编程获取的直觉，反而让我们在并发编程上误入歧途。**

学习之初可以暂时假设goroutine类似于操作系统的线程，但实际上goroutine和线程之间在数量上有非常大的差别。



**Go语言中的并发程序可以用两种手段来实现：**

- **多线程共享内存**（传统的并发模型，如Java/Python/C++）：共享数据结构由锁保护，线程会争夺这些锁以访问数据。由于实现正确访问共享变量的复杂性而变得困难。
- **通信顺序进程**(**Communicating Sequential Process**,**CSP**)的**并发编程模型**：Go不鼓励显式使用锁来协调对共享数据的访问，而是鼓励使用 channels在不同的运行实例（goroutine）之间**传递共享值**。在任何给定时间，只有一个 goroutine 可以访问该值，**实际上它从未被不同的执行线程主动共享，从设计上确保了数据竞争不可能发生。**为了鼓励这种思维方式，将它简化为一个口号：
    > 💡 **Do not communicate by sharing memory; instead, share memory by communicating.**
    > - [https://go.dev/doc/effective_go#channels](https://go.dev/doc/effective_go#channels)
    > - [https://go.dev/blog/codelab-share](https://go.dev/blog/codelab-share)


**CSP模型的理解方法**：

考虑一个典型的单线程程序在一个 CPU 上运行，它不需要同步原语（如锁、信号量等）。

现在运行另一个这样的实例，它同样不需要同步。

现在让这两个实例进行通信；如果通信是同步器，那么仍然不需要其他同步。（通信操作本身提供了必要的同步）。

如，Unix 管道完美地符合这个模型。尽管 Go 的并发方法源于 Hoare 的通信顺序进程（CSP），但它也可以被视为 **Unix 管道的类型安全泛化**。

```shell
# 经典的Unix管道
# 管道自动处理了程序间的同步: cat 写入数据时，grep 自动等待; grep 处理完数据后，wc 才能接收
cat file.txt | grep "hello" | wc -l
```

<!-- 列布局开始 -->

```go
// **多线程共享内存模型:**
type Resource struct {
    url        string
    polling    bool
    lastPolled int64
}

type Resources struct {
    data []*Resource
    lock *sync.Mutex
}

func Poller(res *Resources) {
    for {
        // get the least recently-polled Resource
        // and mark it as being polled
        res.lock.Lock()
        var r *Resource
        for _, v := range res.data {
            if v.polling {
                continue
            }
            if r == nil || v.lastPolled < r.lastPolled {
                r = v
            }
        }
        if r != nil {
            r.polling = true
        }
        res.lock.Unlock()
        if r == nil {
            continue
        }

        // poll the URL

        // update the Resource's polling and lastPolled
        res.lock.Lock()
        r.polling = false
        r.lastPolled = time.Nanoseconds()
        res.lock.Unlock()
    }
}
```


---

```go
// **CSP模型:**
type Resource string

func Poller(in, out chan *Resource) {
    for r := range in {
        // poll the URL

        // send the processed Resource to out
        out <- r
    }
}
```

<!-- 列布局结束 -->

## 一、g**oroutines**

**为什么叫"Goroutine"：**

| 类型 | 描述 | 通信 | 
| --- | --- | --- | 
| **Process(进程)** | **独立的地址空间，资源隔离。**（进程A无法访问进程B的内存） | 通过系统API来通信，安全隔离但开销较大。 | 
| **Thread(线程)** | **操作系统级别的执行单元**，有固定栈大小。 | 通过**共享内存**来通信，用mutex确保并发写安全。 | 
| **Co****routine****(协程)** | 用户级线程，通常**需要手动调度**。 | 通过**共享内存**来通信，用mutex确保并发写安全。 | 
| **Go****routine** | （**/ˈɡoʊruːtiːn/，Rob Pike官方发音，**"勾-如-汀”），区别于上述术语。 | 通过通信（channel）来共享内存 | 

**goroutines和线程**

goroutine与操作系统(OS)线程的差异，本质上是属于量变

**Rob Pike:**

> Goroutine背后的含义是：它是一个coroutine，但是它在阻塞之后会转移到其它coroutine，同一线程上的其它coroutines也会转移，因此它们不会阻塞。
因此，从根本上讲Goroutines是coroutines的一个分支，可在足够多的操作线程上获得多路特性，不会有Goroutines会被其他coroutine阻塞。如果它们只是协作的话，只需一个线程即可。但是如果有很多IO操作的话，就会有许多操作系统动作，也就会有许多许多线程。但是Goroutines还是非常廉价的，它们可以有数十万之众，总体运行良好并只占用合理数量的内存，它们创建起来很廉价并有垃圾回收功能，一切都非常简单。



**动态栈/可增长的栈：**

- **每一个OS线程都有一个固定大小的内存块（通常为2MB）来做栈**，这个栈会用来存储当前正在被调用或挂起/临时暂停（指在调用其它函数时）的函数中的局部变量。
    - 2MB太大：对于一个小的goroutine, 2MB的栈是一个巨大的浪费，比如有的goroutine仅仅等待一个WaitGroup再关闭一个通道。在Go程序中，一次创建十万左右的goroutine也不罕见，对于这种情况，栈就太大了。
    - 2MB太小：另外，对于最复杂和深度递归的函数，固定大小的栈始终不够大。
    - 改变这个固定大小可以提高空间效率并允许创建更多的线程，或者也可以容许更深的递归函数，但无法同时做到上面的两点。
- **作为对比，一个goroutine在生命周期开始时只有一个很小的栈（典型情况下仅为2KB，比OS线程的栈缩小1024倍）**。与OS线程类似，goroutine的栈也用于存放那些正在执行或临时暂停的函数中的局部变量。但与OS线程不同的是，goroutine的栈不是固定大小的，它**可以按需增大和缩小**。goroutine的栈大小限制**可以达到1GB**，比线程典型的固定大小栈高几个数量级。当然，只有极少的goroutine会使用这么大的栈。


**goroutine调度：**

- OS线程会被OS内核来调度。每几毫秒，一个硬件计时器会中断处理器，这会调用一个叫作scheduler的内核函数。这个函数会暂停/挂起当前执行的线程，并将它的寄存器信息保存到内存中，检查线程列表并决定接下来运行哪一个线程，再从内存中恢复该线程的寄存器信息，然后恢复执行该线程的现场并开始执行线程。
- 因为OS线程是被内核来调度，所以控制权限从一个线程到另外一个线程**需要一个完整的上下文切换(context switch)**：即保存一个线程的状态到内存，再恢复另外一个线程的状态，最后更新调度器的数据结构。**这三步操作很慢，因为其局部性很差需要几次内存访问，并且会增加运行的cpu周期。**
- Go的runtime运行时包含一个自己的调度器，这个调度器使用一个称为**m:n调度的技术**（因为其会在n个操作系统线程上多工（调度）m个goroutine）。Go调度器的工作与内核调度器类似，但是这个调度器只关注单独的Go程序中的goroutine（译注：按程序独立）
- 与操作系统的线程调度器不同的是，Go调度器并不是用一个硬件定时器，而是被Go语言“建筑”本身进行调度的。例如当一个goroutine调用了time.Sleep，或者被channel调用或者mutex操作阻塞时，调度器会使其进入休眠并开始执行另一个goroutine，直到时机到了再去唤醒第一个goroutine。**因为这种调度方式不需要进入内核的上下文，所以重新调度一个goroutine比调度一个线程代价要低得多。**


**GOMAXPROCS：**

- Go的调度器使用了一个叫做GOMAXPROCS的变量来决定会有多少个操作系统的线程同时执行Go的代码。其默认的值是运行机器上的CPU的核心数，所以在一个有8个核心的机器上时，调度器一次会在8个OS线程上去调度GO代码。（GOMAXPROCS是前面说的m:n调度中的n）。
    - 在休眠中的或者在通信中被阻塞的goroutine是不需要一个对应的线程来做调度的。
    - 在I/O中或系统调用中或调用非Go语言函数时，是需要一个对应的操作系统线程的，但是GOMAXPROCS并不需要将这几种情况计算在内。
- 你可以用GOMAXPROCS的环境变量来显式地控制这个参数，或者也可以在运行时用runtime.GOMAXPROCS函数来修改它。我们在下面的小程序中会看到GOMAXPROCS的效果，这个程序会无限打印0和1。
    ```go
    for {
        go fmt.Print(0)
        fmt.Print(1)
    }
    $ GOMAXPROCS=1 go run hacker-cliché.go
    111111111111111111110000000000000000000011111...
    $ GOMAXPROCS=2 go run hacker-cliché.go
    010101010101010101011001100101011010010100110...
    ```
    - 在第一次执行时，最多同时只能有一个goroutine被执行。初始情况下只有main goroutine被执行，所以会打印很多1。过了一段时间后，GO调度器会将其置为休眠，并唤醒另一个goroutine，这时候就开始打印很多0了，在打印的时候，goroutine是被调度到操作系统线程上的。
    - 在第二次执行时，我们使用了两个操作系统线程，所以两个goroutine可以一起被执行，以同样的频率交替打印0和1。我们必须强调的是goroutine的调度是受很多因子影响的，而runtime也是在不断地发展演进的，所以这里的你实际得到的结果可能会因为版本的不同而与我们运行的结果有所不同


**Goroutine没有ID号：**

- 在大多数支持多线程的操作系统和程序语言中，当前的线程都有一个独特的身份（id），并且这个身份信息可以以一个普通值的形式被很容易地获取到，典型的可以是一个integer或者指针值。这种情况下我们做一个抽象化的thread-local storage（线程本地存储，多线程编程中不希望其它线程访问的内容）就很容易，只需要以线程的id作为key的一个map就可以解决问题，每一个线程以其id就能从中获取到值，且和其它线程互不冲突。
- goroutine没有可以被程序员获取到的身份（id）的概念。这一点是设计上故意而为之，由于thread-local storage总是会被滥用。
    - 比如说，一个web server是用一种支持tls的语言实现的，而非常普遍的是很多函数会去寻找HTTP请求的信息，这代表它们就是去其存储层（这个存储层有可能是tls）查找的。这就像是那些过分依赖全局变量的程序一样，会导致一种非健康的“距离外行为”，在这种行为下，一个函数的行为可能并不仅由自己的参数所决定，而是由其所运行在的线程所决定。因此，如果线程本身的身份会改变——比如一些worker线程之类的——那么函数的行为就会变得神秘莫测。
- Go鼓励更为简单的模式，这种模式下参数（译注：外部显式参数和内部显式参数。tls 中的内容算是"外部"隐式参数）对函数的影响都是显式的。这样不仅使程序变得更易读，而且会让我们自由地向一些给定的函数分配子任务时不用担心其身份信息影响行为。
- 你现在应该已经明白了写一个Go程序所需要的**所有语言特性信息**。


**goroutine的简单模型**：**它是在同一地址空间中（共享堆内存）与其他 goroutines 并发执行的****一个函数****。**

go语句使 函数或方法 在一个新创建的goroutine中调用，go语句本身的执行立即完成。

```go
func main() {
    // 启动一个goroutine
    go func() {
        fmt.Println("Hello from goroutine")
    }()
    
    // 主goroutine继续执行
    fmt.Println("Hello from main")
}
```

在函数或方法调用前加上 `go` 关键字前缀，可以在新的 goroutine 中运行该调用。当调用完成时，goroutine 会无声地退出。（这种效果类似于 Unix shell 的 `&` 符号，用于**在后台运行命令**。）

```go
go list.Sort()  // run list.Sort concurrently; don't wait for it.
```

在 goroutine 调用中，函数字面量会很有用。**在go中函数字面量是闭包：实现确保了函数所引用的变量在其活跃期间一直存在。**

```go
func Announce(message string, delay time.Duration) {
    go func() {
        time.Sleep(delay)
        fmt.Println(message)
    }()  // Note the parentheses - must call the function.
}
```



**轻量级设计**：除了栈空间的分配外几乎不增加额外开销。并且栈的初始大小很小，因此它们内存开销很小，并且可以根据需要通过分配（和释放）堆存储来增长（避免为每个goroutine预分配大栈）。

```go
// 传统线程：固定栈大小（通常1-8MB）
// Goroutine：动态栈，初始很小（2KB），按需增长

func recursiveFunction(n int) {
    if n <= 0 {
        return
    }
    // 每次递归调用，栈会自动增长
    recursiveFunction(n - 1)
}

// 可以轻松启动成千上万个goroutine
func main() {
    for i := 0; i < 100000; i++ {
        go func(id int) {  // go语句本身立即执行完成，不等待func的执行结束
            // 每个goroutine成本很低
            time.Sleep(time.Second)
        }(i)
    }
    // 传统线程无法做到这一点
}
```



**多路复用机制**：

Goroutines 被多路复用到多个操作系统线程上，所以如果一个线程应该阻塞，如在等待 I/O 时，其他线程仍然可以继续运行。它们的设计隐藏了许多线程创建和管理的复杂性。



你可能有几万个 goroutine，但**底层只会用少量（比如 CPU 核数）OS 线程并行跑**。

Go 的调度器采用 **GMP 模型**：是 **M:N 模型**：M 个 goroutine 映射到 N 个 OS 线程**。****goroutine 是运行时调度的轻量任务，OS 线程只是运行容器。**

- **G** = goroutine
- **M** = machine（**OS thread**）
- **P** = processor（调度上下文，控制可运行的 goroutine 队列）。`GOMAXPROCS` 控制 P（默认为CPU核数，可同时运行 goroutine 的逻辑 CPU 数量


8核 CPU → `GOMAXPROCS=8` （默认情况下）→ 至多 8 个 P 同时在工作。（所以开 **10万 goroutine** 而不会像 10万 OS 线程那样把系统撑爆）

同一时刻最多只有 8 个 goroutine 真正在并行运行，其他 goroutine 都在就绪队列里等待调度。

**并行（parallelism）** = 多核 CPU 上多个 goroutine 同时跑（真正的物理同时）。指在多个 CPU 上并行执行计算以提高效率。

**并发（concurrency）** = goroutine 数量远大于核数，**调度器快速切换，让你“感觉”它们同时在跑**。指将程序结构化为独立执行的组件。



**为什么默认是CPU核数？**

**CPU密集型任务**：线程数 = CPU核心数 最优化

**I/O密集型任务**：Go运行时自动创建更多线程处理阻塞

**避免上下文切换开销**：减少不必要的线程切换

**内存效率**：每个线程都有栈空间，控制线程数就是控制内存使用



```go
// Go运行时将多个goroutine映射到少量OS线程上
// 通常：goroutine数量 >> OS线程数量

func main() {
    // 启动1000个goroutine
    for i := 0; i < 1000; i++ {
        go func() {
            // 这些goroutine可能运行在同一个OS线程上
            time.Sleep(time.Millisecond * 100)
        }()
    }
}
```



**阻塞处理:**

当goroutine在**I/O/channel/mutex/sleep**上阻塞**，Go运行时调度器会将其挂起，从M（OS线程）上移走**（避免整个线程空转，是 Go 能高效支持 **几十万 goroutine 并发** 的关键），让其他就绪的goroutine来占用该OS线程继续执行。

如果是**syscall阻塞**（如文件读写）：runtime 可能会再起一个新的 OS 线程来替代被 syscall 卡住的线程，避免整体卡死。

```go
func worker(id int, ch chan string) {
    // 当这个goroutine阻塞在I/O时
    data := <-ch  // 阻塞点
    
    // 其他goroutine可以继续在同一个OS线程上运行
    fmt.Printf("Worker %d received: %s\n", id, data)
}
```

阻塞的goroutine恢复后，会被重新调度到可用的OS线程上。

**注意：**

- **纯计算阻塞**（例如死循环 `for {}`，不调用 runtime 的可中断点），调度器最初是无法切走的。但 Go 1.14 之后加入了 **抢占式调度**，运行时会在函数调用和安全点打断。
- **系统调用（syscall）** 可能导致线程级阻塞，这种情况下 Go runtime 可能会额外起一个 OS 线程，保证调度不被完全卡死。


**隐藏复杂性：**

简化的并发模型：

```go
// 传统方式：需要手动管理线程
// Go方式：专注于通信

func main() {
    ch := make(chan int)
    
    go func() {
        ch <- 42  // 发送数据
    }()
    
    value := <-ch  // 接收数据
    // 通信自动处理同步，无需锁
}
```

**自动调度:**

```go
// 开发者不需要关心：
// - 线程创建和销毁
// - 线程池管理
// - 负载均衡
// - 上下文切换

func main() {
    // 只需要关心业务逻辑
    go processData()
    go handleRequests()
    go monitorSystem()
    
    // Go运行时自动处理所有调度细节
}
```

**Python**（手动协作式）: 必须手动 `await`，否则不会切走。依赖事件循环asyncio.run。如果一个协程里跑阻塞代码，整个 loop 卡住。

```go
import asyncio

async def foo():
    print("start foo")
    await asyncio.sleep(1)  # 手动挂起
    print("end foo")

async def bar():
    print("bar")

asyncio.run(asyncio.gather(foo(), bar()))
```

**Rust**（手动协作式）：Rust 协程（future）只是状态机，不会自动切换，必须配合 runtime。`tokio` 提供线程池调度，但仍需 `await`。

```go
async fn foo() {
    println!("start foo");
    tokio::time::sleep(std::time::Duration::from_secs(1)).await;
    println!("end foo");
}
```

**JavaScript**: 遇到 `await Promise` 才会把控制权交还给 event loop。

```go
async function foo() {
  console.log("start foo");
  await new Promise(r => setTimeout(r, 1000)); // 手动挂起
  console.log("end foo");
}

async function bar() {
  console.log("bar");
}

foo();
bar();
```



| 语言 | 协程调度方式 | 是否抢占 | 手动挂起点 | 类似 Go goroutine 吗？ | 
| --- | --- | --- | --- | --- | 
| **Go** | 自动 M:N 调度器 | ✅ 支持 | ❌ | ✅ 是 goroutine 代表 | 
| **Python** | 协作式 (asyncio) | ❌ | ✅ (await) | ❌ | 
| **JavaScript** | 协作式 (event loop) | ❌ | ✅ (await) | ❌ | 
| **Lua** | 手动 yield/resume | ❌ | ✅ (yield) | ❌ | 
| **C#** | Task Scheduler + await | ❌ (无抢占) | ✅ (await) | 部分相似 | 
| **Rust** | Future + runtime | ❌ | ✅ (await) | ❌ | 



### **example:**clock

格式化模板限定为Mon Jan 2 03:04:05PM 2006 UTC-0700。有8个部分（周几、月份、一个月的第几天……）。可以以任意的形式来组合前面这个模板；出现在模板中的部分会作为参考来对时间格式进行输出。这是go语言和其它语言相比比较奇葩的一个地方。你需要记住格式化字符串是：**1月2日下午3点4分5秒零六年UTC-0700****（记忆：**1234567**）**，而不像其它语言那样Y-m-d H:i:s一样。

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

**阻塞执行（顺序编程）：**服务器顺序执行，第二个nc客户端接收不到时间；

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_25224637-29b5-8093-8a88-efbf87318d59.jpg)

**并发执行（并发编程）：**多个客户端可以同时接收到时间；

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_25224637-29b5-80f4-9abb-d941b13a55d2.jpg)





### **example：并发的Echo服务**reverb

go后跟的函数的参数表达式求职会在go语句自身（这里是main goroutine）中执行。

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



## **二、channel**

**goroutine是Go程序并发的执行体，****channel（**/ˈtʃænl/，通道**）****是它们之间的连接**，可以让**一个goroutine发送特定值到另一个goroutine的通信机制。**

当**复制或者作为参数传递到一个函数时，复制的是引用**，这样调用者和被调用者都引用同一份数据结构，零值是nil。



和map一样，chan使用内置函数make分配，生成的值作为底层数据结构的引用**（**pointer、slice、map、function、**channel****为引用类型），chan元素有具体的类型（chan int, 类似 []int），零值为nil。**如果提供可选的整型参数，它将设置通道的缓冲区大小。默认值是零，表示无缓冲区或同步通道。

每一个通道是一个具体类型的chan，叫作通道的元素类型，如一个有int类型元素的通道写为chan int。

```go
ch := make(chan int) // ch的类型: 'chan int' ， unbuffered channel（无缓冲通道）。等价于：ch = make(chan int, 0) // unbuffered channel
ch = make(chan int, 3) // 通道容量3的缓冲通道
```



**可比较性**：同种类型的通道可以使用==符号进行比较：**当二者都是同一通道数据的引用时，比较值为true；**通道也可以和nil进行比较。



**通道的三个主要操作**：

send、receive都使用`<-`运算符（简化到一个运算符），**统称为通信**：

- **发送(send)语句**：`ch <- x` ，从一个goroutine传输一个值到另一个在执行接收表达式的goroutine。
- **接收(receive)语句**：`x = <- ch` ；`<-ch` 一个不使用接收结果的接收操作也是合法的。
- **关闭(close)**：`close(ch)` ，**关闭通道的入口**。设置一个标志位来指示值当前已经发送完毕，这个通道后面没有值了。
    - 关闭后的再次**send**发送操作、关闭一个已经关闭的channel、关闭一个nil值的空channel都将导致panic；
    - 在一个已经关闭的通道上进行**receive**接收操作，将**获取所有已经发送的值**，直到通道为空；这时任何接收操作会**立即完成**，同时获取到一个通道元素类型对应的**零值**；
    - **仅能close发送方的**channel**，**close一个只接收的channel将是一个编译错误；
    - **不要将channel的close和文件的close操作混淆：**当结束的时候对每一个文件调Close方法是非常重要的，但channel可以不用close，channel的close**只用于 断言/通知 接收方goroutine不再向channel发送新的数据；**
    - **GC垃圾回收器 在channel没有被引用时回收它**（而不是根据它是否close）；
    - **close channel**还可以作为一个广播机制；


### unbuffered channel

**unbuffered channel（无缓冲通道/同步通道）**：将通信（值的交换）与同步（确保两个goroutines计算处于已知状态，类似：goroutines间的信息同步，如传统的共享变量）结合起来。

1. **一个goroutine的发送操作（ch <- 0）会阻塞在原地，直到另一个goroutine在对应的通道上执行完接收操作（<-ch）。****这时值传送完成，两个goroutine都可以恢复继续执行。**
1. 如果接收操作先执行，接收方goroutine将一直阻塞，直到另一个goroutine在同一个通道上发送一个值。
使用无缓冲通道进行的通信**实现****发送goroutine和接收goroutine的同步化**（同步通道）。当一个值在无缓冲通道上传递时，接收值后发送方goroutine才被再次唤醒。



**happens before（**Go语言并发内存模型的一个关键术语）：

**”x早于y发生“**：x发生的时间100%的可预期的确定早于y，可放心的依赖这个机制。

**”x既不比y早也不比y晚“：**x和y并发，**无法确定它们的执行顺序，**不一定是x和y同一时刻执行，依赖这类情况可能产生并发问题。



**消息事件：**

每一条消息有一个值，但有时候更强调通信发生的时刻，成为**消息事件(event)**（类似：OpenAPI WebHook来新订单的事件）。这时消息事件不需要携带额外的信息，**降级为仅仅用作两个goroutine之间的同步**。通常直接用`done <- true` `done <- 1`、`done <- struct{}{}` 。

**example**：netcat

```go
conn, err := net.Dial("tcp", "localhost:8000")
if err != nil {
	log.Fatal(err)
}
done := make(chan struct{})
go func() {   // Go语言中启动goroutine常用的形式: go语句调用了一个函数字面量（匿名函数）
	io.Copy(os.Stdout, conn) // NOTE: ignoring errors
	log.Println("done")
	done <- struct{}{} // signal the main goroutine
}()
mustCopy(conn, os.Stdin)
conn.Close()
<-done // wait for background goroutine to finish
```



**串联的channels（pipeline）**

通过Channels可连接将多个goroutine（类似：langchain。Eino的compose.graph通过channel连接。）一个Channel的输出作为下一个Channel的输入。这种串联的Channels就是所谓的管道（pipeline）。

下面的程序用两个channels将三个goroutine串联起来：

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_19524637-29b5-80eb-93d8-fa3a355c77d6.jpg)

- goroutine1-counter：产生一个0, 1, 2，…的整数序列；
- goroutine2-square：计算数值的平方；
- goroutine3-printer：接收值并打印；


像这样的管道出现在长期运行的服务器程序中，其中通道用于在包含无限循环的goroutine之间整个生命周期中的通信。

没有一个直接的方式来判断是否通道已经关闭，但是这里有接收操作的一个变种，它产生两个结果：接收到的通道元素，以及**一个布尔值（通常称为ok），它为true的时候代表接收成功**，**false表示当前的接收操作在一个关闭的并且读完的通道上。**



**example**：**pipeline**

```go
naturals := make(chan int)
squares := make(chan int)

// goroutine1-counter
go func() {
        for x := 0; x < 100; x++ {
            naturals <- x
        }
        close(naturals)
}()

// goroutine2-square    
go func() {            
    for {
        x, ok := <-naturals     // channel被关闭并且没有值可接收时跳出循环
        if !ok {
            break // channel was closed and drained
        }
        squares <- x * x
    }
    close(squares)
}()

// goroutine3-printer
go func() {           
    for x := range naturals {
        squares <- x * x
    }
    close(squares)
}()

// Printer (in main goroutine)
for x := range squares {  // **for range：支持在channel上的迭代器。**更方便接收在通道上所有发送的值，接收完最后一个值后关闭循环。
    fmt.Println(x)
}
```



**单向通道类型:**

当一个channel作为一个函数参数时，它**一般总是被专门用于只发送或者只接收**。

为了防止被滥用，Go语言的类型系统提供了单方向的channel类型，分别用于只发送或只接收的channel。箭头`<-`和关键字chan的相对位置表明了channel的方向，这种限制将在编译期检测。

`chan<- int`类型：表示一个只发送int的channel，只能发送不能接收。

`<-chan int`类型：表示一个只接收int的channel，只能接收不能发送。



**可赋值性**：双向通道可以赋值给单向通道变量（**隐式转换）**，**但单向通道（**`chan<- int`**）不可以赋值给双向通道（**`chan int`**）。**

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



### buffered channel

`ch = make(chan string, 3)`缓冲通道有一个元素队列，队列的最大长度在创建的时候通过make的容量参数来设置。



**通道的缓冲区解耦了发送goroutine和接收goroutine：**

- **发送操作**：缓冲通道上的发送操作在队列的尾部插入一个元素；如果通道满了，发送操作会阻塞所在的goroutine直到另一个goroutine对它进行接收操作来留出可用的空间 （降级为无缓冲通道/同步通道？）。
- **接收操作：**从队列的头部移出一个元素；如果通道是空的，执行接收操作的goroutine阻塞，直到另一个goroutine在通道上发送数据；
<!-- 列布局开始 -->

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_19524637-29b5-8038-9876-c97102bcb179.jpg)


---

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_19524637-29b5-8067-a0c4-eb8c7dbda2a8.jpg)


---

![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_19524637-29b5-808a-94f9-fc7854cae2fa.jpg)

<!-- 列布局结束 -->



在某些特殊情况下，程序需要知道通道缓冲区存的容量和元素个数：

- 通道缓冲区的容量： `fmt.Println(cap(ch))`
- 通道缓冲区的有效元素个数： `fmt.Println(len(ch))` ； 在并发程序中元素个数会随着接收操作而立即失效，但是它对某些故障诊断和性能优化会有帮助；


**goroutines泄漏****（类似内存泄漏）：**

如果使用一个无缓冲通道，**两个比较慢的goroutine由于发送响应结果到channel的时候，没有goroutine来接收而将被永远卡住的bug。**

因为语法简单，Go新手粗暴地将缓冲通道作为队列在单个goroutine中使用，但是这是个严重错误。**如果仅仅需要一个简单的队列，使用slice创建一个就可以。**

**channel和goroutine的调度深度关联，如果没有另一个goroutine从通道进行接收，发送者（也许是整个程序）有被****永久阻塞****的风险。**

**和回收变量不同，泄漏的goroutines不会自动回收**，因此必须确保每个goroutine在不再需要的时候可以自动结束。

example：

```go
func mirroredQuery() string {
    responses := make(chan string, 3)
    go func() { responses <- request("asia.gopl.io") }()       // 并发地向三个镜像站点发出请求，三个镜像站点分散在不同的地理位置，它们分别将收到的响应发送到带缓存channel
    go func() { responses <- request("europe.gopl.io") }()
    go func() { responses <- request("americas.gopl.io") }()
    return <-responses // return the quickest response      // 最后接收者只接收第一个收到的（最快的）响应，mirroredQuery函数可能在另外两个响应慢的镜像站点响应之前就返回了结果。（多个goroutines并发地向同一个channel发送数据，或从同一个channel接收数据都是常见的用法。）
}

func request(hostname string) (response string) { /* ... */ }
```



unbuffered channel和buffered channel的选择、buffered channel 容量cap大小的选择，都会对程序的正确性产生影响。

unbuffered channel提供强同步保障，**因为每一次发送都需要和一次对应的接收同步。**

对于buffered channel，这些操作则是解耦的；如果我们知道要发送的值数量的上限，通常会创建一个容量是使用上限的缓冲通道，在接收第一个值前就完成所有的发送。**在内存无法提供缓冲容量的情况下，可能导致程序死锁。**

**组装流水线**是对于通道和goroutine合适的类比：

- `make(chan int)`**：**想象蛋糕店里的三个厨师，在生产线上，在把每一个蛋糕传递给下一个厨师之前，一个烤，一个加糖衣，一个雕刻。在空间比较小的厨房，每一个厨师完成一个蛋糕流程，必须等待下一个厨师准备好接受它；
- `make(chan int, 1)`: 如果在厨师之间有可以放一个蛋糕的位置，一个厨师可以将制作好的蛋糕放到这里，然后立即开始制作下一个，这类似于使用一个容量1的缓冲通道。只要厨师们以相同的速度工作，大多数工作就可以快速处理，**消除他们各自之间的速率差异**。
- `make(chan int, 3)`: 如果在厨师之间有更多的空间(更长的缓冲区)，就可以**消除更大的暂态速率波动而不影响组装流水线**，比如当一个厨师稍作休息时，后面再抓紧跟上进度。
- 另一方面，如果生产线的上游持续比下游快，**缓冲区满的时间占大多数**。如果后续的流程更快，**缓冲区通常是空的**。这时缓冲区的存在是没有价值的。
- **创建另外一个goroutine而使用同一个通道来通信**：如果第二段更加复杂，一个厨师可能跟不上第一个厨师的供应，或者跟不上第三个厨师的需求。为了解决这个问题，我们可以雇用另一个厨师来帮助第二段流程，独立地执行同样的任务。


### channel

channel是一等公民，可以像其他任何值一样被分配和传递。这一特性常用于实现安全的并行多路分解。

下述段代码是一个带有**限流、并行、非阻塞特性**的 RPC 系统的框架，而且其中没有出现任何互斥锁。

```go
queue chan *Request // chan结构体，结构体里有chan

type Request struct {
    args        []int
    f           func([]int) int
    resultChan  chan int
}

func sum(a []int) (s int) {
    for _, v := range a {
        s += v
    }
    return
}

request := &Request{[]int{3, 4, 5}, sum, make(chan int)}
// Send request
clientRequests <- request
// Wait for response.
fmt.Printf("answer: %d\n", <-request.resultChan)
```

```go
// 在服务器端，唯一变化的是处理函数。
func handle(queue chan *Request) {
    for req := range queue {
        req.resultChan <- req.f(req.args)
    }
}
```

## 三、example

### **并发的循环**迭代

生成一批全尺寸图片的缩略图：很明显，处理文件的顺序没有关系，因为每一个缩放操作和其他的操作独立。通过并行可以利用多核CPU的计算能力拉伸图像，隐藏文件I/O产生的延迟。

像这样由一些完全独立的子问题组成的问题称为高度并行的问题。



**匿名函数中的循环变量快照问题**: 

循环变量f是被所有的匿名函数值所共享，且会被连续的循环迭代所更新的。

当新的goroutine开始执行字面函数时，for循环 **可能已经更新了f并且开始了另一轮的迭代**或者**（更有可能的）已经结束了整个循环**，所以**当这些goroutine开始读取f的值时，它们所看到的值****已经是slice的最后一个元素****了**；

**bug**: 直接使用外层闭包中声明的循环变量f

```go
for _, f := range filenames {
    go func() {
        thumbnail.ImageFile(f)  // bug：启动了所有的goroutine，每一个文件名对应一个，**但没有等待它们一直到执行完毕而立即返回了**。
        // ...
    }()
}
```

**fix**: 显式参数传递循环变量f `go func(f string) {}(f)`：可以确保当go语句执行的时候，使用f的当前循环的值。

```go
func makeThumbnails3(filenames []string) {
    ch := make(chan struct{})
    for _, f := range filenames {
        go func(f string) {      		// 两条语句封装成一个匿名函数
            thumbnail.ImageFile(f) 
            ch <- struct{}{}      // 向一个共享的channel中发送事件（“事件”在channel中有介绍）
        }(f)                     // 注意我们将f的值作为一个显式的变量传给了函数，而不是在循环的闭包中声明；
    }
    // Wait for goroutines to complete.
    for range filenames {
        <-ch
    }
}
```



**goroutine泄漏：**

**bug**：可能导致整个程序卡住或者系统内存耗尽（oom: out of memory）

```go
// makeThumbnails4 makes thumbnails for the specified files in parallel.
// It returns an error if any step failed.
func makeThumbnails4(filenames []string) error {
    errors := make(chan error)
    for _, f := range filenames {
        go func(f string) {
            _, err := thumbnail.ImageFile(f)
            errors <- err    // 无缓存chan，goroutine一直在等待chan清空后写入
        }(f)
    }
    for range filenames {
        if err := <-errors; err != nil {
            return err // 当遇到第一个非nil的错误时，它将错误返回给调用者，这样没有goroutine继续从errors返回通道上进行接收。导致后续的所有发送者**goroutine被永久的阻塞**
        }
    }
    return nil
}
```

**fix**：

**方案1**：使用一个有足够容量的缓冲通道，这样没有worker goroutine在发送消息时候阻塞；

```go
func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
    type item struct {
        thumbfile string
        err       error
    }
    ch := make(chan item, len(filenames))   // 文件总数长度的缓冲通道
    for _, f := range filenames {
        go func(f string) {
            var it item
            it.thumbfile, it.err = thumbnail.ImageFile(f)
            ch <- it
        }(f)
    }
    for range filenames {
        it := <-ch
        if it.err != nil {
            return nil, it.err
        }
        thumbfiles = append(thumbfiles, it.thumbfile)
    }
    return thumbfiles, nil
}
```



**方案2**：在主goroutine返回第一个错误的同时，创建另一个goroutine来接收完通道；

这个版本里没有把文件名放在slice里，而是通过`filenames <-chan string`传过来，所以我们无法对循环的次数进行预测。

**在不知道迭代次数的情况下，下述代码结构是通用的、符合习惯的、地道的并行的循环迭代模板。**

**使用 sync.WaitGroup /sɪŋk/ 计数器类型**: **一个可以被多个goroutine安全地操作的计数器，计数当前时刻在跑的****goroutine数量**。在每一个goroutine启动前递增计数，在每一个goroutine结束时递减计数。

注意其Add和Done方法的不对称性：

- Add递增计数器，它必须在工作goroutine开始之前执行，而不是在中间。另一方面，不能保证Add会在关闭者goroutine调用Wait之前发生。
- Add有一个参数，但Done没有，它等价于Add(-1)。使用defer来确保在发送错误的情况下计数器可以递减。
sizes通道将每一个文件的大小带回主goroutine，它使用range循环进行接收然后计算总和。注意，在closer goroutine中，在关闭sizes通道之前，等待所有的工作者结束。这里两个操作（等待和关闭）必须和在sizes通道上面的迭代并行执行。

考虑替代方案：如果我们将等待操作放在循环之前的main goroutine中，因为通道会满，它将永不结束；如果放在循环后面，它将不可达，因为没有任何东西可用来关闭通道，循环可能永不结束；

```go
func makeThumbnails6(filenames <-chan string) int64 {
    sizes := make(chan int64)
    var wg sync.WaitGroup // number of working goroutines
    for f := range filenames {
        wg.Add(1)  // Add是为计数器加一，必须在**worker goroutine**开始之前调用，而不是在goroutine中 (否则的话我们没办法确定Add是在"closer" goroutine调用Wait之前被调用)
        go func(f string) {     // worker
            defer wg.Done()   // Done却没有任何参数；其实它和Add(-1)是等价的；使用defer来确保计数器即使是在出错的情况下依然能够正确地被减掉。
            thumb, err := thumbnail.ImageFile(f)
            if err != nil {
                log.Println(err)
                return
            }
            info, _ := os.Stat(thumb) // OK to ignore error
            sizes <- info.Size()
        }(f)
    }

    go func() {        // closer goroutine，让其在所有worker goroutine们结束之后再关闭sizes channel的
        wg.Wait()
        close(sizes)
    }()
    var total int64
    for size := range sizes {    // sizes channel携带了每一个文件的大小到main goroutine，在main goroutine中使用了range loop来计算总和。
        total += size
    }
    return total
}
```



**makeThumbnails6函数中的事件序列：**

- 垂直线表示goroutine。细片段表示休眠，粗片段表示活动。
- 斜箭头表示goroutine通过事件进行了同步。时间从上向下流动。
- 注意，主goroutine把大多数时间花在range循环休眠上，等待工作者发送值或等待closer来关闭通道。
![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_19624637-29b5-805f-9ff3-caa0cb0ffdbd.jpg)



### 使用channel实现信号量

**并发的Web爬虫：**将bfs(广度优先)算法来抓取整个网站的crawl改造成并发运行。

**crawl1:**

```go
func crawl(url string) []string {
    fmt.Println(url)
    list, err := links.Extract(url)
    if err != nil {
        log.Print(err)
    }
    return list
}

func main() 
    worklist := make(chan []string)  // 待抓取的URL列表，这里channel代替slice来做这个队列

    go func() { worklist <- os.Args[1:] }() {    // 另启一个goroutine，避免死锁（也可用缓冲通道）。
    seen := make(map[string]bool)
    for list := range worklist {
        for _, link := range list {
            if !seen[link] {
                seen[link] = true
                go func(link string) {
                    worklist <- crawl(link)
                }(link)    // 显示传参，避免循环变量捕获
            }
        }
    }
}
```

发送给任务列表的命令行参数必须在它自己的goroutine中运行来避免死锁，死锁是一种卡住的情况，其中主goroutine和一个爬取goroutine同时发送给对方但是双方都没有接收。另一个可选的方案是使用缓冲通道

**bug**：**无限制的并行通常不是一个好的主意，因为系统中总有物理极限（如：对于计算型应用CPU的核数，对于磁盘I/O操作磁头和磁盘的个数，下载流所使用的网络带宽，或者Web服务本身的容量）；**

```go
$ go build gopl.io/ch8/crawl1
$ ./crawl1 http://gopl.io/
http://gopl.io/
https://golang.org/help/
...
2015/07/15 18:22:12 Get ...: dial tcp: lookup blog.golang.org: no such host  // 对一个可靠的域名出现了解析失败
2015/07/15 18:22:12 Get ...: dial tcp 23.21.222.120:443: socket: too many open files  // 同时创建了太多的网络连接，超过了程序能打开文件数的限制
...
```



**fix**：**重写crawl2：**

根**据资源可用情况限制并发的个数，以匹配合适的并行度：**如限制对于links.Extract的同时调用不超过n个；

**计数信号量：使用容量为n的缓冲通道来建立一个并发原语**。概念上，对于缓冲通道中的n个空闲槽，每一个代表一个令牌，持有者可以执行。

通过发送一个值到通道中来领取令牌，从通道中接收一个值来释放令牌，创建一个新的空闲槽。这保证了在没有接收操作的时候，最多同时有n个发送。（尽管使用已填充槽比令牌更直观，但使用空闲槽在创建通道缓冲区之后可以省掉填充的过程）。

**因为通道的元素类型在这里不重要，所以我们使用struct{}，它所占用的空间大小是0。**

使用令牌的获取和释放操作来包括对links.Extract函数的调用，这样保证最多同时20个调用可以进行。

**保持信号量操作离它所约束的I/O操作越近越好**（这是一个好的实践）：

重写crawl函数，将对links.Extract的调用操作**用获取、释放token的操作包裹起来**，来确保同一时间对其只有20个调用。信号量数量和其能操作的IO资源数量应保持接近。

```go
var tokens = make(chan struct{}, 20)   // **计数信号量 tokens，确保并发请求限制在20个以内**

func crawl(url string) []string {
    fmt.Println(url)
    tokens <- struct{}{} // acquire a token
    list, err := links.Extract(url)
    <-tokens // release the token
    if err != nil {
        log.Print(err)
    }
    return list
}

func main() {
    worklist := make(chan []string)
    var n int // 计数器n跟踪发送到任务列表中的任务个数。每次知道一个条目被发送到任务列表时，就递增变量n；

    // Start with the command-line arguments.
    n++                       // n的第一次递增是在发送初始化命令行参数之前
    go func() { worklist <- os.Args[1:] }()

    // Crawl the web concurrently.
    seen := make(map[string]bool)
    for ; n > 0; n-- {   // 主循环从n减到0，这时再没有任务需要完成
        list := <-worklist
        for _, link := range list {
            if !seen[link] {        		// 为了让程序终止，当任务列表为空且爬取goroutine都结束以后，需要从主循环退出
                seen[link] = true 
                n++                      // n的第二次及以后的递增
                go func(link string) {
                    worklist <- crawl(link)
                }(link)
            }
        }
    }
}
```



**fix2：另一个方案。**使用原来的crawl函数，它没有计数信号量，但是通过20个长期存活/常驻的爬虫goroutine来调用它，这样确保最多20个HTTP请求并发执行；

爬取goroutine使用同一个通道unseenLinks进行接收。主goroutine负责对从任务列表接收到的条目进行去重，然后发送每一个没有爬取过的条目到unseenLinks通道，然后被爬取goroutine接收。

seen map被限制在主goroutine里面，它仅仅需要被这个goroutine访问。与其他形式的信息隐藏一样，范围限制可以帮助我们推导程序的正确性。如，局部变量不能在声明它的函数之外通过名字引用；没有从函数中逃逸的变量不能从函数外面访问；一个对象的封装域只能被对象自己的方法访问。所有的场景中，信息隐藏帮助限制程序不同部分之间不经意的交互。

crawl发现的链接通过精心设计的goroutine发送到任务列表来避免死锁。

```go
// 所有的爬虫goroutine现在都是被同一个channel - unseenLinks喂饱的了。主goroutine负责拆分它从worklist里拿到的元素，然后把没有抓过的经由unseenLinks channel发送给一个爬虫的goroutine。

// seen这个map被限定在main goroutine中；也就是说这个map只能在main goroutine中进行访问。类似于其它的信息隐藏方式，这样的约束可以让我们从一定程度上保证程序的正确性。例如，内部变量不能够在函数外部被访问到；变量（§2.3.4）在没有发生变量逃逸（译注：局部变量被全局变量引用地址导致变量被分配在堆上）的情况下是无法在函数外部访问的；一个对象的封装字段无法被该对象的方法以外的方法访问到。在所有的情况下，信息隐藏都可以帮助我们约束我们的程序，使其不发生意料之外的情况。

// crawl函数爬到的链接在一个专有的goroutine中被发送到worklist中来避免死锁。
func main() {
    worklist := make(chan []string)  // lists of URLs, may have duplicates
    unseenLinks := make(chan string) // de-duplicated URLs

    // Add command-line arguments to worklist.
    go func() { worklist <- os.Args[1:] }()

    // Create 20 crawler goroutines to fetch each unseen link.
    for i := 0; i < 20; i++ {
        go func() {
            for link := range unseenLinks {
                foundLinks := crawl(link)
                go func() { worklist <- foundLinks }()
            }
        }()
    }

    // The main goroutine de-duplicates worklist items
    // and sends the unseen ones to the crawlers.
    seen := make(map[string]bool)
    for list := range worklist {
        for _, link := range list {
            if !seen[link] {
                seen[link] = true
                unseenLinks <- link
            }
        }
    }
}
```



### **使用**`select case`多路复用

**example**：火箭发射的倒计时 countdown1.go

```go
func main() {
    fmt.Println("Commencing countdown.")
    // ime.Tick函数返回一个 <-chan Time 类型，程序会周期性地像一个节拍器一样向这个channel发送事件
    tick := time.Tick(1 * time.Second)
    for countdown := 10; countdown > 0; countdown-- {
        fmt.Println(countdown)
        <-tick
    }
    launch()
}

abort := make(chan struct{})
go func() {   // // 启动一个goroutine，这个goroutine会尝试从标准输入中读入一个单独的byte并且，如果成功了，会向名为abort的channel发送一个值。
    os.Stdin.Read(make([]byte, 1)) // read a single byte
    abort <- struct{}{}
}()
```



**select语句**: 专用于channel的**switch**语句，允许一个goroutine同时等待/监听多个channel的读/写操作，但能随机执行一个匹配的case（防止饥饿）。有一个最后的（可选的）default分支。

- 每一个情况指定一次通信（在一些通道上进行发送或接收操作）和关联的一段代码块。
- 接收表达式操作可能出现在它本身上，像第一个情况，或者在一个短变量声明中，像第二个情况；第二种形式可以让你引用所接收的值。
- **对于没有对应情况的select, select{}将永远等待：**select一直等待，直到一次通信来告知有一些情况可以执行。然后，它进行这次通信，执行此情况所对应的语句；其他的通信将不会发生。
- **如果多个case同时满足，select随机选择一个，这样保证每一个通道有相同的机会被选中：**在前一个例子中增加缓冲区的容量，会使输出变得不可确定，**因为当缓冲既不空也不满的情况，相当于select语句在扔硬币做选择。**
- **通道的零值nil通道有时候很有用**：**因为在nil通道上发送和接收将永远阻塞，对于select语句中的情况，如果其通道是nil，它将永远不会被选择**。这次让我们用nil来开启或禁用特性所对应的情况，比如超时处理或者取消操作，响应其他的输入事件或者发送事件。
```go
select {
	case <-ch1:
	    // ...
	case x := <-ch2:
	    // ...use x...
	case ch3 <- y:
	    // ...
	default:
	    // ...
}
```

```go
// 下面的select语句会一直等待直到两个事件中的一个到达，无论是abort事件或者一个10秒经过的事件。如果10秒经过了还没有abort事件进入，那么火箭就会发射
func main() {
    // ...create abort channel...

    fmt.Println("Commencing countdown.  Press return to abort.")
    select {
    case <-time.After(10 * time.Second):
        // Do nothing.
    case <-abort:
        fmt.Println("Launch aborted!")
        return
    }
    launch()
}
```



下面这个例子更微妙：ch这个channel的buffer大小是1，所以会交替的为空或为满，所以只有一个case可以进行下去，无论i是奇数或者偶数，它都会打印0 2 4 6 8。

```go
ch := make(chan int, 1)
for i := 0; i < 10; i++ {
    select {
    case x := <-ch:    // ch为空时，跳过这个case；ch有值时，执行这个case；
        fmt.Println(x) // "0" "2" "4" "6" "8"
    case ch <- i:
    }
}
```



**example**：火箭发射的倒计时 countdown1.go：让发射程序打印倒计时。select语句使每一次迭代使用1s来等待中止；gopl.io/ch8/countdown3

time.Tick函数的行为很像创建一个goroutine在循环里面调用time.Sleep，然后在它每次醒来时发送事件。当上面的倒计时函数返回时，它停止从tick通道中接收事件，但是计时器goroutine还在运行，徒劳地向一个没有goroutine在接收的通道不断发送（发生goroutine泄漏）；

- Tick函数很方便使用，但是它仅仅在应用的整个生命周期中都需要时才合适。否则，我们需要使用这个模式：
    ```go
    ticker := time.NewTicker(1 * time.Second)
    <-ticker.C    // receive from the ticker's channel
    ticker.Stop() // cause the ticker's goroutine to terminate
    ```
有时候我们试图在一个通道上发送或接收，但不想在通道没有准备好的情况下被阻塞（**非阻塞通信）**。这使用select语句也可以做到。select可以有一个默认情况，它用来指定在没有其他的通信发生时可以立即执行的动作。

- 上面的select语句从尝试从abort通道中接收一个值，如果没有值，它什么也不做。这是一个非阻塞的接收操作；重复这个动作称为对通道轮询：
    ```go
    select {
    case <-abort:
        fmt.Printf("Launch aborted!\n")
        return
    default:
        // do nothing
    }
    ```


### **示例: 并发的目录遍历**

- 示例：报告一个或多个目录的磁盘使用情况（类似于UNIX du命令）
    ```go
    // walkDir recursively walks the file tree rooted at dir
    // and sends the size of each found file on fileSizes.
    func walkDir(dir string, fileSizes chan<- int64) {
        for _, entry := range dirents(dir) {
            if entry.IsDir() {
                subdir := filepath.Join(dir, entry.Name())
                walkDir(subdir, fileSizes).  // 每个子目录递归调用自身
            } else {
                fileSizes <- entry.Size()  // 向fileSizes channel发送一条消息，内容为文件的字节大小的
            }
        }
    }
    // dirents returns the entries of directory dir.
    func dirents(dir string) []os.FileInfo {
        entries, err := ioutil.ReadDir(dir).  // 返回一个os.FileInfo类型的slice，os.FileInfo类型也是os.Stat这个函数的返回值
        if err != nil {
            fmt.Fprintf(os.Stderr, "du1: %v\n", err)
            return nil
        }
        return entries
    }
    ```
    - main函数
        - 后台的goroutine调用walkDir来遍历命令行给出的每一个路径并最终关闭fileSizes这个channel。
        - 主goroutine会对其从channel中接收到的文件大小进行累加，并输出其和
        ```go
        package main
        import (
            "flag"
            "fmt"
            "io/ioutil"
            "os"
            "path/filepath"
        )
        func main() {
            // Determine the initial directories.
            flag.Parse()
            roots := flag.Args()
            if len(roots) == 0 {
                roots = []string{"."}
            }
            // Traverse the file tree.
            fileSizes := make(chan int64)
            go func() {    // 后台goroutine
                for _, root := range roots {
                    walkDir(root, fileSizes)
                }
                close(fileSizes)
            }()
            // Print the results.
            var nfiles, nbytes int64
            for size := range fileSizes {
                nfiles++    // 直接++，Go的零值初始化的机制带来的更简洁的代码
                nbytes += size
            }
            printDiskUsage(nfiles, nbytes)
        }
        func printDiskUsage(nfiles, nbytes int64) {
            fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
        }
        ```
    - 思考：这个程序会在打印其结果之前运行很长时间，如果在运行的时候能够让我们知道处理进度的话想必更好。但如果简单地把printDiskUsage函数调用移动到循环里会导致其打印出成百上千的输出；
    - 下面这个du的变种会间歇打印内容，不过只有在调用时提供了-v的flag才会显示程序进度信息
        - 主goroutine现在使用了计时器来每500ms生成事件，然后用select语句来等待文件大小的消息来更新总大小数据，或者一个计时器的事件来打印当前的总大小数据；
        - 如果-v的flag在运行时没有传入的话，tick这个channel会保持为nil，这样在select里的case也就相当于被禁用了
        - 由于我们的程序不再使用range循环，第一个select的case必须显式地判断fileSizes的channel是不是已经被关闭了，这里可以用到channel接收的二值形式。如果channel已经被关闭了的话，程序会直接退出循环。
            - **这里的break语句用到了标签break，这样可以同时终结select和for两层循环**；如果没有用标签就break的话只会退出内层的select循环，而外层的for循环会使之进入下一轮select循环。
        ```go
        var verbose = flag.Bool("v", false, "show verbose progress messages")
        func main() {
            // ...start background goroutine...
            // Print the results periodically.
            var tick <-chan time.Time
            if *verbose {
                tick = time.Tick(500 * time.Millisecond)
            }
            var nfiles, nbytes int64
            
        loop:
            for {
                select {
                case size, ok := <-fileSizes:
                    if !ok {
                        break loop // fileSizes was closed**。标签break，这样可以同时终结select和for两层循环**
                    }
                    nfiles++
                    nbytes += size
                case <-tick:
                    printDiskUsage(nfiles, nbytes)
                }
            }
            printDiskUsage(nfiles, nbytes) // final totals
        }
        ```
    - 思考：并发调用walkDir，从而发挥磁盘系统的并行性能：对每一个walkDir的调用创建一个新的goroutine。
        - 它使用sync.WaitGroup来对仍旧活跃的walkDir调用进行计数，另一个goroutine会在计数器减为零的时候将fileSizes这个channel关闭；
        - 由于这个程序在高峰期会创建成百上千的goroutine，我们需要修改dirents函数，用计数信号量来阻止他同时打开太多的文件，就像前面的并发爬虫一样；
        ```go
        func main() {
            // ...determine roots...
            // Traverse each root of the file tree in parallel.
            fileSizes := make(chan int64)
            var n sync.WaitGroup
            for _, root := range roots {
                n.Add(1)
                go walkDir(root, &n, fileSizes)
            }
            go func() {
                n.Wait()
                close(fileSizes)
            }()
            // ...select loop...
        }   
        func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
            defer n.Done()
            for _, entry := range dirents(dir) {
                if entry.IsDir() {
                    n.Add(1)
                    subdir := filepath.Join(dir, entry.Name())
                    go walkDir(subdir, n, fileSizes)
                } else {
                    fileSizes <- entry.Size()
                }
            }
        }
        // sema is a counting semaphore for limiting concurrency in dirents.
        var sema = make(chan struct{}, 20)
        // dirents returns the entries of directory dir.
        func dirents(dir string) []os.FileInfo {
            sema <- struct{}{}        // acquire token
            defer func() { <-sema }() // release token
            // ...
        ```


### **示例：并发的退出**

- 有时候我们需要通知goroutine停止它当前的任务；如：一个Web服务器对客户请求处理到一半的时候客户端断开了；
- **一个goroutine无法直接终止另一个，因为这样会让所有goroutine之间的共享变量状态处于不确定状态**。在8.7节的火箭发射程序中，我们给abort通道发送一个值，倒计时goroutine把它理解为停止自己的请求。但是怎样才能取消两个或者指定个数的goroutine呢？
    - 一个可能是给abort通道发送和要取消的goroutine同样多的事件。如果一些goroutine已经自己终止了，这样计数就多了，然后发送过程会卡住。如果那些goroutine可以自我繁殖，数量又会太少，其中一些goroutine依然不知道要取消。
    - 通常，任何时刻都很难知道有多少goroutine正在工作。更多情况下，当一个goroutine从abort通道接收到值时，它利用这个值，这样其他的goroutine接收不到这个值。对于取消操作，我们需要一个可靠的机制通过一个通道广播一个事件，这样goroutine们能够看到这条事件消息，并且在事件完成之后，可以知道这件事已经发生过
    - **回忆一下，当一个通道关闭 且已取完所有发送的值 之后，接下来的接收操作立即返回，得到零值。**我们可以利用它创建一个广播机制：不在通道上发送值，而是用关闭一个通道来进行广播
- dup的修改版本 *gopl.io/ch8/du4*：
    - 第一步，创建一个取消通道，在它上面不发送任何值，但是它的关闭表明程序需要停止它正在做的事情。也定义了一个工具函数cancelled，在它被调用的时候检测或轮询取消状态；
        ```go
        var done = make(chan struct{})
        func cancelled() bool {   // 工具函数
            select {
            case <-done:
                return true
            default:
                return false
            }
        }
        ```
    - 第二步：创建一个读取标准输入的goroutine，这是一个比较典型的连接到终端的程序。一旦开始读取任何输入（如用户按回车键）时，这个goroutine通过关闭done通道来广播取消事件
        ```go
        // Cancel traversal when input is detected.
        go func() {
            os.Stdin.Read(make([]byte, 1)) // read a single byte
            close(done)
        }()
        ```
    - 第三步：让goroutine来响应取消操作。在主goroutine中，添加第三个情况到select语句中，它尝试从done通道接收。如果选择这个情况，函数将返回，但是在返回之前它必须耗尽fileSizes通道，丢弃它所有的值，直到通道关闭。做这些是为了保证所有的walkDir调用可以执行完，不会卡在向fileSizes通道发送消息上。
        ```go
        for {
            select {
            case <-done:
                // Drain fileSizes to allow existing goroutines to finish.
                for range fileSizes {
                    // Do nothing.
                }
                return
            case size, ok := <-fileSizes:
                // ...
            }
        }
        ```
    - walkDir goroutine在开始的时候轮询取消状态，如果设置状态，什么都不做立即返回。它让在取消后创建的goroutine什么都不做：
        - 在walkDir循环中来进行取消状态轮询也许是划算的，它避免在取消后创建新的goroutine。取消需要权衡：更快的响应通常需要更多的程序逻辑变更入侵。确保在取消事件以后没有更多昂贵的操作发生，可能需要更新代码中很多的地方，但通常我们可以通过在少量重要的地方检查取消状态来达到目的。
        ```go
        func walkDir(dir string, n *sync.WaitGroup, fileSizes chan<- int64) {
            defer n.Done()
            if cancelled() {
                return
            }
            for _, entry := range dirents(dir) {
                // ...
            }
        }
        ```
    - 现在，当取消事件发生时，所有的后台goroutine迅速停止，然后main函数返回。当然，当main返回时，程序随之退出，不过这里没有谁在后面通知main函数来进行清理。在测试中有一个技巧：如果在取消事件到来的时候main函数没有返回，执行一个panic调用，然后运行时将转储程序中所有goroutine的栈。如果主goroutine是最后一个剩下的goroutine，它需要自己进行清理。但如果还有其他的goroutine存活，它们可能还没有合适地取消，或者它们已经取消，可是需要的时间比较长；多一点调查总是值得的。崩溃转储信息通常含有足够的信息来分辨这些情况。
- 程序的一点性能剖析揭示了它的瓶颈在于dirents中获取信号量令牌的操作。下面的select让取消操作的延迟从数百毫秒减为几十毫秒：
    ```go
    func dirents(dir string) []os.FileInfo {
        select {
        case sema <- struct{}{}: // acquire token
        case <-done:
            return nil // cancelled
        }
        defer func() { <-sema }() // release token
        // ...read directory...
    }
    ```


### **示例: 并发的聊天服务器**

- 示例: 并发的聊天服务器，可以在几个用户之间相互广播文本消息。这个程序里有4个goroutine。
    - 主goroutine：监听端口，接受连接客户端的网络连接。对每一个连接，它创建一个新的handleConngoroutine，就像本章开始时并发回声服务器中那样。
        ```go
        func main() {
            listener, err := net.Listen("tcp", "localhost:8000")
            if err != nil {
                log.Fatal(err)
            }
            go broadcaster()
            for {
                conn, err := listener.Accept()
                if err != nil {
                    log.Print(err)
                    continue
                }
                go handleConn(conn)
            }
        }
        ```
    - 广播器(broadcaster) goroutine：使用局部变量clients来记录当前连接的客户端集合。其记录的内容是每一个客户端的消息发出channel的“资格”信息
        - 广播者监听两个全局的通道entering和leaving，通过它们通知客户的到来和离开，如果它从其中一个接收到事件，它将更新clients集合。
        - 如果客户离开，那么它关闭对应客户对外发送消息的通道。广播者也监听来自messages通道的事件，所有的客户都将消息发送到这个通道。当广播者接收到其中一个事件时，它把消息广播给所有客户。
        ```go
        type client chan<- string // an outgoing message channel
        var (
            entering = make(chan client)
            leaving  = make(chan client)
            messages = make(chan string) // all incoming client messages
        )
        func broadcaster() {
            clients := make(map[client]bool) // all connected clients
            for {
                select {
                case msg := <-messages:
                    // Broadcast incoming message to all
                    // clients' outgoing message channels.
                    for cli := range clients {
                        cli <- msg
                    }
                case cli := <-entering:
                    clients[cli] = true
                case cli := <-leaving:
                    delete(clients, cli)
                    close(cli)
                }
            }
        }
        ```
    - handleConn函数创建一个对外发送消息的新通道，然后通过entering通道通知广播者新客户到来。接着，它读取客户发来的每一行文本，通过全局接收消息通道将每一行发送给广播者，发送时在每条消息前面加上发送者ID作为前缀。一旦从客户端读取完毕消息，handleConn通过leaving通道通知客户离开，然后关闭连接。
        ```go
        func handleConn(conn net.Conn) {
            ch := make(chan string) // outgoing client messages
            go clientWriter(conn, ch)
            who := conn.RemoteAddr().String()
            ch <- "You are " + who
            messages <- who + " has arrived"
            entering <- ch
            input := bufio.NewScanner(conn)
            for input.Scan() {
                messages <- who + ": " + input.Text()
            }
            // NOTE: ignoring potential errors from input.Err()
            leaving <- ch
            messages <- who + " has left"
            conn.Close()
        }
        func clientWriter(conn net.Conn, ch <-chan string) {
            for msg := range ch {
                fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
            }
        }
        ```
    - 另外，handleConn函数还为每一个客户创建了写入(clientWriter)goroutine，它接收消息，广播到客户的发送消息通道中，然后将它们写到客户的网络连接中。客户写入者的循环在广播者收到leaving通知并且关闭客户的发送消息通道后终止。
    - 每一个连接里面有一个连接处理(handleConn)goroutine和一个客户写入(clientWriter)goroutine。广播器(broadcaster)是关于如何使用select的一个规范说明，因为它需要对三种不同的消息进行响应。
- 当有n个客户session在连接的时候，程序并发运行着2n+2个相互通信的goroutine，它不需要隐式的加锁操作（参考9.2节）。clients map限制在广播器这一个goroutine中被访问，所以不会并发访问它。唯一被多个goroutine共享的变量是通道以及net.Conn的实例，它们又都是并发安全的。关于限制、并发安全，以及跨goroutine的变量共享的含义，将在下一章进行更多的讨论。
    ![](/images/1e324637-29b5-8036-b50f-dc1145654220/image_1e524637-29b5-8064-ba14-d87a51a2c2c8.jpg)


### 经典并发面试题：两个goroutine交替依次打印12

**要求**：两个goroutine交替依次打印10组12，12121212121212121212。只能使用无缓冲chan，不能使用time.sleep这样的硬等待。

**解法**：

1. **同步**：利用无缓冲chan的阻塞特性，一个goroutine的发送操作（ch <- 0）会阻塞在原地，直到另一个goroutine在对应的通道上执行完接收操作（<-ch）；接受操作同理；
1. **先打印1**：先打印1打印1的goroutine在循环体前单独发送ch <- 0；
1. **死锁**：循环的最后一个打印2后的ch <- 0，需要另一个goroutine在循环外接收<-ch（否则阻塞导致死锁）；
1. **main函数过早终止**：需要sync.WaitGroup{}，wg.Add(2)，defer wg.Done() 确保主函数等待两个 goroutine 都执行完毕后在结束（否则两个goroutine都没启动）；
参考：[https://github.com/lifei6671/interview-go/blob/master/question/q001.md](https://github.com/lifei6671/interview-go/blob/master/question/q001.md)

```shell
package main

import (
	"fmt"
	"sync"
)

func print12(wg *sync.WaitGroup) {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		fmt.Print(1)
		ch <- 0
		for i := 1; i < 10; i++ {
			<-ch
			fmt.Print(1)
			ch <- 0
		}
		<-ch

	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-ch
			fmt.Print(2)
			ch <- 0
		}
	}()

}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	print12(&wg)
	wg.Wait()
}

```



