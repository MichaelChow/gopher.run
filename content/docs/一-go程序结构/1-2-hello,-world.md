---
title: "1.2 Hello, World"
date: 2024-11-30T04:24:00Z
draft: false
weight: 1002
---

# 1.2 Hello, World

# 一、hello world

以下是Go语言版本的hello world，“hello world”案例首次出现于1978年出版的C语言圣经《The C Programming Language》

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, 世界")
}
```



hello world源码解读：

- **原生支持 Unicode**：Go可以处理全世界任何语言的文本（Go的三位设计者中**2位为UTF-8的设计者**）;
- **main package & main function**：每个源文件都以一条`package`声明语句开始，通过包（package）组织****(包类似于其它语言里的**库（libraries）/模块（modules）**）。package main 定义了一个**独立可执行的程序，而不是一个库**。在`main`里的`main`函数也很特殊，它是**整个程序执行时的入口**（译注：C 系语言差不多都这样）;
- `goimports`：紧跟着必须告诉编译器的一系列导入（import）的包。Go语言的代码：一个包由位于单个目录下的一个或多个`.go`源代码文件组成，目录定义包的作用。Go严格要求缺少了必要的包或者**导入了不需要的包，程序都无法编译通过**，避免了程序开发过程中引入未使用的包。`goimports` 自动导入；
- `gofmt`：Go 语言在代码格式上**采取了很强硬的态度**；`gofmt` 自动格式化代码；编译器会主动**把【特定符号清单】后的换行符转换为分号**，包括：行末的**标识符、整数、浮点数、虚数、字符或字符串文字、关键字**`**break**`**、**`**continue**`**、**`**fallthrough**`**或 **`**return**`** 中的一个、运算符和分隔符 **`**++**`**、**`**--**`**、**`**)**`**、**`**]**`** 或 **`**}**`** 。**所以除非一行上有多条语句，否则不需要在语句或者声明的末尾添加分号（编辑器保存时会自动删除行末的分号，自动执行`gofmt`格式化代码）。 **避免了无尽的无意义的琐碎争执**（译注：也导致了 Go 语言的 [TIOBE](https://www.tiobe.com/tiobe-index/)排名较低，因为缺少撕逼的话题）。更重要的是，这样可以做多种自动源码转换，如果放任Go语言代码格式，这些转换就不大可能了。
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_15724637-29b5-805e-afe9-c7722c11b100.jpg)
- **Go标准库**：提供了100多个包，以支持常见功能，如输入、输出、排序以及文本处理；
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_15724637-29b5-80b5-a492-c5ebcda00da7.jpg)




# 二、go build & go run

- `**go build helloworld.go**`：编译；由于每个go程序都通过goroutine运行，go的二进制文件体积很大；**体积对比**：C版本 33k，go版本 2.1M（**65倍**）；
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_a143a20e-a8b3-4ce7-a8e9-fb59bfd1bd92.jpg)
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_1c324637-29b5-80e3-b843-c4d8567bac22.jpg)
    - C版本 hello，world: gcc -o hello hello.c
        ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_1c324637-29b5-8080-b53b-c246fa8be498.jpg)
- `**go run helloworld.go**`：编译后运行；
