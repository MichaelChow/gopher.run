---
title: "3.6 Goroutines和线程"
date: 2025-05-16T14:33:00Z
draft: false
weight: 3006
---

# 3.6 Goroutines和线程

goroutine与操作系统(OS)线程的差异，本质上是属于量变

### **动态栈/可增长的栈**

- **每一个OS线程都有一个固定大小的内存块（通常为2MB）来做栈**，这个栈会用来存储当前正在被调用或挂起/临时暂停（指在调用其它函数时）的函数中的局部变量。
    - 2MB太大：对于一个小的goroutine, 2MB的栈是一个巨大的浪费，比如有的goroutine仅仅等待一个WaitGroup再关闭一个通道。在Go程序中，一次创建十万左右的goroutine也不罕见，对于这种情况，栈就太大了。
    - 2MB太小：另外，对于最复杂和深度递归的函数，固定大小的栈始终不够大。
    - 改变这个固定大小可以提高空间效率并允许创建更多的线程，或者也可以容许更深的递归函数，但无法同时做到上面的两点。
- **作为对比，一个goroutine在生命周期开始时只有一个很小的栈（典型情况下仅为2KB，比OS线程的栈缩小1024倍）**。与OS线程类似，goroutine的栈也用于存放那些正在执行或临时暂停的函数中的局部变量。但与OS线程不同的是，goroutine的栈不是固定大小的，它**可以按需增大和缩小**。goroutine的栈大小限制**可以达到1GB**，比线程典型的固定大小栈高几个数量级。当然，只有极少的goroutine会使用这么大的栈。
### goroutine调度

- OS线程会被OS内核来调度。每几毫秒，一个硬件计时器会中断处理器，这会调用一个叫作scheduler的内核函数。这个函数会暂停/挂起当前执行的线程，并将它的寄存器信息保存到内存中，检查线程列表并决定接下来运行哪一个线程，再从内存中恢复该线程的寄存器信息，然后恢复执行该线程的现场并开始执行线程。
- 因为OS线程是被内核来调度，所以控制权限从一个线程到另外一个线程**需要一个完整的上下文切换(context switch)**：即保存一个线程的状态到内存，再恢复另外一个线程的状态，最后更新调度器的数据结构。**这三步操作很慢，因为其局部性很差需要几次内存访问，并且会增加运行的cpu周期。**
- Go的runtime运行时包含一个自己的调度器，这个调度器使用一个称为**m:n调度的技术**（因为其会在n个操作系统线程上多工（调度）m个goroutine）。Go调度器的工作与内核调度器类似，但是这个调度器只关注单独的Go程序中的goroutine（译注：按程序独立）
- 与操作系统的线程调度器不同的是，Go调度器并不是用一个硬件定时器，而是被Go语言“建筑”本身进行调度的。例如当一个goroutine调用了time.Sleep，或者被channel调用或者mutex操作阻塞时，调度器会使其进入休眠并开始执行另一个goroutine，直到时机到了再去唤醒第一个goroutine。**因为这种调度方式不需要进入内核的上下文，所以重新调度一个goroutine比调度一个线程代价要低得多。**
### GOMAXPROCS

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


### **Goroutine没有ID号**

- 在大多数支持多线程的操作系统和程序语言中，当前的线程都有一个独特的身份（id），并且这个身份信息可以以一个普通值的形式被很容易地获取到，典型的可以是一个integer或者指针值。这种情况下我们做一个抽象化的thread-local storage（线程本地存储，多线程编程中不希望其它线程访问的内容）就很容易，只需要以线程的id作为key的一个map就可以解决问题，每一个线程以其id就能从中获取到值，且和其它线程互不冲突。
- goroutine没有可以被程序员获取到的身份（id）的概念。这一点是设计上故意而为之，由于thread-local storage总是会被滥用。
    - 比如说，一个web server是用一种支持tls的语言实现的，而非常普遍的是很多函数会去寻找HTTP请求的信息，这代表它们就是去其存储层（这个存储层有可能是tls）查找的。这就像是那些过分依赖全局变量的程序一样，会导致一种非健康的“距离外行为”，在这种行为下，一个函数的行为可能并不仅由自己的参数所决定，而是由其所运行在的线程所决定。因此，如果线程本身的身份会改变——比如一些worker线程之类的——那么函数的行为就会变得神秘莫测。
- Go鼓励更为简单的模式，这种模式下参数（译注：外部显式参数和内部显式参数。tls 中的内容算是"外部"隐式参数）对函数的影响都是显式的。这样不仅使程序变得更易读，而且会让我们自由地向一些给定的函数分配子任务时不用担心其身份信息影响行为。
- 你现在应该已经明白了写一个Go程序所需要的**所有语言特性信息**。
