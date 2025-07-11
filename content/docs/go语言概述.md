---
title: "Go语言概述"
description: "关于Go语言概述的详细说明"
date: 2025-07-05T13:25:00Z
draft: false

---

# Go语言概述

> 💡 推荐学习方法：
> - 细读《GO程序设计语言》第1章-入门、第2章-程序结构；该书由Go语言核心团队成员 [Alan A. A. Donovan](https://github.com/adonovan) 和 [K&R](https://en.wikipedia.org/wiki/K%26R) 中的 [Brian W. Kernighan](http://www.cs.princeton.edu/~bwk/) 合作编写。
> - 充分利用Go的标准函数库以**Go语言自己的特性和风格**来编写程序，避免按照自己曾经熟悉的Java风格、Python风格套路，写Go语言程序；避免将现有的C++或Java程序直译为Go程序；

**为什么一门语言是这样子的**：编程语言的演化跟**生物物种的演化**类似，一个成功的编程语言的后代一般都会继承它们祖先的优点；当然有时**多种语言杂合**也可能会产生**令人惊讶的特性**；还有一些**激进的新特性**可能并没有先例。

![](/images/22724637-29b5-8056-a8d3-d3c8b9937033/image_14024637-29b5-807d-b402-ef52beada612.jpg)



- **C语言祖先**：**Go = C + GC + Goroutine**。从C语言继承了相似的**表达式语法、控制流结构、基础数据类型、调用参数传值、指针**等很多思想，还有C语言一直所看中的**编译后机器码的运行效率**以及和现有操作系统的无缝**适配**。
- **Pascal语言祖先**：[Niklaus Wirth](https://en.wikipedia.org/wiki/Niklaus_Wirth)所设计的Pascal语言，Modula-2语言激发了**包**的概念，Oberon语言摒弃了模块接口文件和模块实现文件之间的区别。第二代的Oberon-2语言直接影响了**包的导入和声明的语法**，还有Oberon语言的**面向对象特性所提供的方法的声明语法**等。
- **CSP（顺序通信进程）语言祖先**：灵感来自于贝尔实验室的[Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)于1978年发表的鲜为外界所知的关于并发研究的基础文献**顺序通信进程**（communicating sequential processes)，缩写为CSP。在CSP中，程序是一组中间**没有共享状态的平行运行的处理过程**，它们之间**使用管道进行通信和控制同步**。不过[Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)的CSP只是一个用于描述并发性基本概念的描述语言，并不是一个可以编写可执行程序的通用编程语言。“不要通过共享内存来通信，而应通过通信来共享内存。” **Rob Pike**和其他人开始不断尝试将[CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes)引入实际的编程语言中。他们第一次尝试引入[CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes)特性的编程语言叫[Squeak](http://doc.cat-v.org/bell_labs/squeak/)（老鼠间交流的语言），是一个提供鼠标和键盘事件处理的编程语言，它的管道是静态创建的。然后是改进版的[Newsqueak](http://doc.cat-v.org/bell_labs/squeak/)语言，提供了类似C语言语句和表达式的语法和类似Pascal语言的推导语法。Newsqueak是一个带垃圾回收的纯函数式语言，它再次针对键盘、鼠标和窗口事件管理。但是在Newsqueak语言中管道是动态创建的，属于第一类值，可以保存到变量中。在Plan9操作系统中，这些优秀的想法被吸收到了一个叫Alef的编程语言中。Alef试图将Newsqueak语言改造为系统编程语言，但是因为**缺少垃圾回收机制而导致并发编程很痛苦**。在Alef之后还有一个叫Limbo的编程语言，Go语言从其中借鉴了很多特性。具体请参考Pike的讲稿：[http://talks.golang.org/2012/concurrency.slide#9](http://talks.golang.org/2012/concurrency.slide#9)
- iota语法是从APL语言借鉴，词法作用域与嵌套函数来自于Scheme语言。
- Go语言的的创新：**切片**为动态数组提供了有效的**随机存取的性能**，这可能会让人联想到链表的底层的共享机制。Go语言新发明的defer语句。


**Go语言：**

- 在上个世纪70年代，贝尔实验室的 [Ken Thompson](http://genius.cat-v.org/ken-thompson/) 和 [Dennis M. Ritchie](http://genius.cat-v.org/dennis-ritchie/)合作发明了 [UNIX](http://doc.cat-v.org/unix/) 操作系统，同时 [Dennis M. Ritchie](http://genius.cat-v.org/dennis-ritchie/) 为了解决[UNIX](http://doc.cat-v.org/unix/)系统的移植性问题而发明了C语言。
- 2009年，在[UNIX](http://doc.cat-v.org/unix/)和C语言发明40年之后，当时已经在 Google 工作的 [Ken Thompson](http://genius.cat-v.org/ken-thompson/)（汤普森）和 [Rob Pike](http://genius.cat-v.org/rob-pike/)（Go语言项目总负责人，与Ken Thompson共同发明了UTF-8字符集规范）、[Robert Griesemer](http://research.google.com/pubs/author96.html)（设计了V8 JavaScript引擎和Java HotSpot虚拟机）一起合作，为了解决在21世纪多核和网络化环境下越来越复杂的编程问题而发明了 Go语言。
    ![](/images/22724637-29b5-8056-a8d3-d3c8b9937033/image_1c224637-29b5-80c5-b592-f55dcabede97.jpg)
- 从 Go 语言库早期代码库日志可以看出它的演化历程（ Git 用 `git log --before={2008-03-03} --reverse` 命令查看）。从早期提交日志中也可以看出，Go 语言是从[Ken Thompson](http://genius.cat-v.org/ken-thompson/)发明的B 言、[Dennis M. Ritchie](http://genius.cat-v.org/dennis-ritchie/)发明的C语言逐步演化过来的，是C语言家族的成员，因此很多人将Go语言称为21世纪的C语言。
    ![](/images/22724637-29b5-8056-a8d3-d3c8b9937033/image_1c224637-29b5-800f-88fe-e4cdde882305.jpg)
- 纵观这几年来的发展趋势，Go语言已经成为云计算、云存储时代最重要的基础编程语言，包括[docker](https://github.com/docker)、[K8s](https://github.com/kubernetes/kubernetes)、[v2ray](https://github.com/v2fly/v2ray-core)、[go-ethereum](https://github.com/ethereum/go-ethereum)、[hugo](https://github.com/gohugoio/hugo)、istio、etcd 、prometheus。正在使用的Go语言的公司如下，**Google**在2010年后开始将基础设施迁移到Go。
![](/images/22724637-29b5-8056-a8d3-d3c8b9937033/image_1c224637-29b5-80dd-a2d0-f7e42558e70b.jpg)

- 在国内，**七牛云**是第一家大面积采用Go语言开发的公司，时间还在Go1.0正式发布之前。许式伟也是中国第一个知名的Go语言布道师。在2015年之前，许式伟和七牛云团队是国内Go语言社区推广的主力。
- **字节跳动**在2012年创业团队早期使用Python技术栈做web后端服务，到2014年业务体量迅速增长遇到Python性能问题，逐步有团队开始尝试用Go。发现学习成本低，开发和部署非常简单，顺带解决了Python的依赖库版本问题。随着字节内部基于Go自研的RPC框架（[**Kitex**](https://github.com/cloudwego/kitex)）和HTTP框架（[**Hertz**](https://github.com/cloudwego/hertz)）的推广，逐步将Python微服务全面重写为Go版本，到2020年前后微服务数量达到5万+。在Go的sort优化上，字节使用了pdqsort 算法 + Go 1.18 泛型，实现了一个比标准库 API 在几乎所有情况下快 2x ~ 60x 的算法库。论文地址：[https://arxiv.org/pdf/2106.05123.pdf](https://arxiv.org/pdf/2106.05123.pdf)
- **哔哩哔哩**创业团队早期使用PHP语言开发，后哔哩哔哩的中台技术逐步切换到Node、后台技术为了更高的并发和稳定性逐步切换到Java。这导致了哔哩哔哩的技术架构混乱，早期几乎天天故障，难以维护；统一技术栈的背景下最终选择使用更能满足哔哩哔哩需求的Go语言重写，其研发总监毛剑是一位Go语言的忠实布道者。[https://github.com/Tencent/secguide](https://github.com/Tencent/secguide)
- **腾讯**随着云计算和大数据相关业务的迅速发展，逐渐尝试使用Go语言。如Go语言代码安全指南：[https://github.com/Tencent/secguide/blob/main/Go%E5%AE%89%E5%85%A8%E6%8C%87%E5%8D%97.md](https://github.com/Tencent/secguide/blob/main/Go%E5%AE%89%E5%85%A8%E6%8C%87%E5%8D%97.md)
- Go官网：[https://go.dev/](https://go.dev/)、 [https://github.com/golang/go](https://github.com/golang/go)
- 目前，Go每半年发布一个二级版本。Go发行说明：[Go Release History](https://go.dev/doc/devel/release)、[中文说明](https://github.com/guyan0319/golang_development_notes/blob/master/zh/1.6.md)。其中2019年前后，Go语言三位设计者之一的Robert Griesemer提交的通过try内置函数来简化错误处理就被否决了。
    | **Release** | **Release date** | **Status** | **Maintenance end** | **Feature** | 
    | --- | --- | --- | --- | --- | 
    | go1.24 | 2025-02-11 |   |   |   | 
    | go1.23 | 2024-08-13 |   |   |   | 
    | go1.22 | 2024-02-06 |   | Q1 2025 |   | 
    | go1.21 | 2023-08-08 |   | Q3 2024 |   | 
    | go1.20 | 2023-02-01 | End-of-Life | Q1 2024 |   | 
    | go1.19 | 2022-08-02 | End-of-Life | 2023-08-08 |   | 
    | go1.18 | 2022-03-15 | End-of-Life | 2023-02-01 | 实现**泛型** | 
    | go1.17 | 2021-08-16 | End-of-Life | 2022-08-02 |   | 
    | go1.16 | 2021-02-16 | End-of-Life | 2022-03-15 |   | 
    | go1.15 | 2020-08-11 | End-of-Life | 2021-08-16 |   | 
    | go1.14 | 2020-02-25 | End-of-Life | 2021-02-16 |   | 
    | go1.13 | 2019-09-03 | End-of-Life | 2020-08-11 | 淘汰Go1.7增加的导入二进制包（不直接导入Go源码，以保护商业代码）的特性（伪需求，实际很少使用）。 | 
    | go1.12 | 2019-02-25 | End-of-Life | 2020-02-25 | 实现官方的包依赖管理工具（Go Module），淘汰GOPATH和vendor的版本管理。 | 
    | go1.11 | 2018-08-24 | End-of-Life | 2019-09-03 |   | 
    | go1.10 | 2018-02-16 | End-of-Life | 2019-02-25 |   | 
    | go1.9 | 2017-08-24 | End-of-Life | 2018-08-24 | 引入了类型别名的特性：type T1 = T2 | 
    | go1.8 | 2017-02-16 | End-of-Life | 2018-02-16 |   | 
    | go1.7 | 2016-08-15 | End-of-Life | 2017-08-24 | 标准库引入了context包，简化对于处理单个请求的多个Goroutine之间与请求域的数据、超时和退出等操作；<br/>如gRPC以及很多Web框架都通过context来控制Goroutine的生命周期。 | 
    | go1.6 | 2016-02-17 | End-of-Life | 2017-02-16 | GC的Stop the world的停顿时间由Go1.4之前的几秒-几十秒降低到10毫秒以内。 | 
    | go1.5 | 2015-08-19 | End-of-Life | 2016-08-15 | 实现自举，**编译器**和**运行时**完全用Go编写（不再使用C语言，带来到好处如GC可以精确知道每个内存指针的状态）。 | 
    | go1.4 | 2014-12-10 | End-of-Life | 2016-02-17 |   | 
    | go1.3 | 2014-06-18 | End-of-Life | 2015-08-19 |   | 
    | go1.2 | 2013-12-01 | End-of-Life | 2014-12-10 | 给切片语法增加了容量的控制，避免不同的切片不小心越界访问有着相同底层数组的其它切片的内存。 | 
    | go1.1 | 2013-05-13 | End-of-Life | 2014-06-18 |   | 
    | go1 | 2012-03-28 | End-of-Life | 2013-12-01 |   | 
### GO语言特性

所有的编程语言都反映了**语言设计者对编程哲学的反思**，通常包括之前的语言所暴露的一些**不足地方的改进****。**Go项目是在Google公司维护超级复杂的几个软件系统时遇到的一些共性问题的反思。

- **容易忘记简洁的内涵：**通过增加功能、选项和配置是修复问题的最快的途径，但这**慢慢地增加了其他部分的复杂性**（如抖音修复越权方案通过引入SDK鉴权，由于复杂度高而推动困难），“软件的复杂性是乘法级相关的”（[Rob Pike](http://genius.cat-v.org/rob-pike/)）。简洁的设计需要在工作开始的时候**舍弃不必要的想法**，并且在软件的生命周期内严格区别**好的改变和坏的改变**。通过足够的努力，一个好的改变可以**在不破坏原有完整概念的前提下保持自适应**，正如[Fred Brooks](http://www.cs.unc.edu/~brooks/)所说的“概念完整性”；而一个坏的改变则不能达到这个效果，它们仅仅是**通过肤浅的和简单的妥协来破坏原有设计的一致性**。只有通过简洁的设计，才能让一个系统保持**稳定、安全和持续的进化**（备注：基于简洁设计原则，所以很多开发者提出的issue中特性没有被采纳）。Go语言有一个被称之为 **“没有废物” 的宗旨**，就是将一切没有必要的东西都去掉。
- Go语言本身只有很少的特性
    - 拥有自动垃圾回收、一个包系统、函数作为一等公民（函数在语言中具有与其他数据类型相同的地位）、词法作用域、系统调用接口、只读的UTF8字符串等。
    - **强类型语言，不允许隐式的类型转换**（避免C/C++开发中的bug）。Go语言有足够的**类型系统**以避免动态语言中那些粗心的类型错误，但是，Go语言的类型系统相比传统的强类型语言又要简洁很多。虽然，有时候这会导致一个**“无类型”**的抽象类型概念，但是Go语言程序员并不需要像C++或Haskell程序员那样纠结于具体类型的安全属性。在实践中，Go语言简洁的类型系统给程序员带来了更多的安全性和更好的运行时性能。
    - 没有构造函数和析构函数
    - 没有运算符重载
    - 没有默认参数
    - 没有继承
    - 没有异常
    - 没有宏
    - 没有函数修饰
    - 没有线程局部存储
    - 承诺保证向后兼容：用之前的Go语言编写程序，可以用新版本的Go语言编译器和标准库直接构建而不需要修改代码。
    - Go语言鼓励局部的重要性。它的内置数据类型和大多数的标准库数据结构都经过精心设计而**避免显式的初始化或隐式的构造函数**，因为很少的内存分配和内存初始化代码被隐藏在库代码中了。
    - Go语言的**聚合类型（结构体和数组）**可以直接操作它们的元素，只需要更少的存储空间、更少的内存写操作，而且**指针操作比其他间接操作的语言也更有效率****。**
    - Go语言提供了基于CSP的并发特性支持。Go语言的动态栈使得**轻量级线程**goroutine**的初始栈可以很小**，因此，创建一个goroutine的代价很小，进而创建**百万级的**goroutine完全是可行的。goroutine” 的发音为：英 /ɡəˈruːtiːn/；美 /ɡəˈruːtiːn/ Go [性能说明](https://learnku.com/docs/the-way-to-go/go-performance-description/3580)。 Go通过以下的 Logo 来展示它的速度，并以**囊地鼠 (Gopher， /ˈɡofɚ/) 作为它的吉祥物**。
        ![](/images/22724637-29b5-8056-a8d3-d3c8b9937033/image_16908c4f-aa32-48cf-bc6d-b35fd00d457f.jpg)
    - Go语言的标准库提供了清晰的构建模块和公共接口，包含**I/O操作、文本处理、图像、密码学、网络和分布式应用程序**等，并支持许多标准化的文件格式和编解码协议。**库和工具使用了大量的约定来减少额外的配置和解释****，从而最终简化程序的逻辑**，而且，**每个Go程序结构都是如此的相似**，因此，Go程序也很容易学习。使用Go语言自带工具构建Go语言项目只需要使用文件名和标识符名称，一个偶尔的特殊注释来确定所有的库、可执行文件、测试、基准测试、例子、以及特定于平台的变量、项目的文档等；**Go语言源代码本身就包含了构建规范**。
    - 基本结构、数据类型、函数与主流**命令式编程语言类似****。**方法、接口、并发、包、测试和反射等语言特性是Go语言特有的。
    - **Go语言的面向对象机制**与一般语言不同。它没有类和类层次结构（没有继承）；仅仅**通过组合（而不是继承）简单的对象来构建复杂的对象**。方法不仅可以**定义在结构体上**，而且，可以定义在任何**用户自定义的类型**上；**具体类型和抽象类型（接口）之间的关系是隐式的**，所以很多**类型的设计者可能并不知道该类型到底实现了哪些接口**。
    - 利用Go自带的工具，可使用单个命令完成编译、测试、基准测试、代码格式化、文档以及其他诸多任务。
    - 在单元测试上，Go语言的工具和标准库中集成了轻量级的测试功能，避免了强大但复杂的测试框架。
    - 反射是一个强大的编程工具，一种程序在运行期间审视自己的能力，不过要谨慎地使用；利用反射机制能实现一些重要的Go语言库函数，展示了反射的强大用法。
    - 在必要时，可以使用unsafe包绕过Go语言安全的类型系统。
### GO开发环境

- 推荐直接使用Goland，教育版本的License。
    - 设置国内代理：go env -w GOPROXY=https://goproxy.cn
- 不建议使用vscode：
    1. 安装go官网的安装包，新版本不再需要配置gopath。go env看相关配置参数。
    1. 安装VS Code；
    1. 安装VS Code的官方的“GO”插件；安装成功后VS Code会弹框提示自动安装`gopls`（提供函数文档弹框）等依赖，点“Install ALl”）；
        1. 翻墙方式1: 在打开OpenWrt的全局代理或配置VS Code Proxy的基础上
        1. 翻墙方式2: go env -w GOPROXY=[https://goproxy.cn](https://goproxy.cn/)
        1. 方式3：连手机的热点
    1. 使用 `go install` 命令安装 Go 包或二进制文件，Go 会将编译后的二进制文件放到 $HOME/go/bin，需要手动添加环境变量：`export PATH=$PATH:$HOME/go/bin`
    1. 安装go package：七牛云的镜像[https://goproxy.cn/](https://goproxy.cn/)
        1. export GOPROXY=https://goproxy.cn
        1. go get [gopl.io/ch5/links](http://gopl.io/ch5/links)
### Go程序的基本组成结构

Go语言的中的大程序都从小的基本组件构建而来：

- **变量**存储值；
- **基本类型**通过数据和结构体进行**聚合**；
- 简单的**表达式**通过加和减等操作合并成大的；
- 表达式通过if、for、switch、defer等**控制语句**来决定**执行顺序**；
- 语句被组织成**函数**用于隔离和复用；
- 函数被组织成**源文件和包**；


