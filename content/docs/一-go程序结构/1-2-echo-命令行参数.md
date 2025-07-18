---
title: "1.2 echo 命令行参数"
date: 2025-03-27T06:01:00Z
draft: false
weight: 1002
---

# 1.2 echo 命令行参数

大多数的程序都是：

- 处理输入（文件、网络连接、其它程序的输出、敲键盘的用户、命令行参数或其它类似输入源）；
- 产生输出（“计算”的定义）；


`os`包以与平台无关的方式提供了一些与操作系统交互的函数和变量。

`os.Args`变量是一个**字符串（string）的切片（slice）**（/slaɪs/，**动态数组**）。和大多数编程语言类似，区间索引时，Go 语言里也采用**左闭右开，**因为**这样可以简化逻辑。**`s[m:n]`，省略m和n时，会默认传入 `0` 或 `len(s)`。

```go
// echo prints its command-line arguments.
package main

import (
    "fmt"
    "os"
)

func main() {
    var s, sep string
    for i := 1; i < len(os.Args); i++ {
        s += sep + os.Args[i]
        sep = " "
    }
    fmt.Println(s)
}
```

- 注释语句：**惯例，我们在每个包的包声明前添加注释**；注释以程序名开头，从整体角度对程序做个简要描述。[Untitled](https://www.notion.so/1502463729b5803ba1e5d5d2967ff16c) 
- 导入多个包，习惯上用括号把它们**括起来写成列表形式**；`gofmt`工具格式化时**按照字母顺序对包名排序**
- `var` 声明：如果变量没有显式初始化，则被隐式地赋予其类型的**零值（zero value）**，**数值类型是 **`**0**`**，字符串类型是空字符串 **`**""**`
- 符号`:=`是**短变量声明（short variable declaration）**的一部分，这是定义一个或多个变量并根据它们的初始值为这些变量赋予适当类型的语句。只能用于在函数体中，而不能用在包级别；
- 对数值类型，Go语言提供了常规的数值和逻辑运算符。而对 `string` 类型，`+` 运算符连接字符串（译注：和 C++ 或者 JavaScript 是一样的）。等价于：`s=s+sep+os.Args[i]`。
- 自增语句`i++`：给`i`加`1`；这和 `i+=1` 以及 `i=i+1` 都是等价的。**它们****是语句****，而不像 C 系的其它语言那样是****表达式**。**所以 **`**j=i++**`** 非法****；**而且 `++` 和 `--` 都只能放在变量名后面，因此`**--i**`**也非法****。（****表达式是赋值=的右边部分，而语句是独立完整一条****）**
- Go语言**只有**`**for**`**循环这一种循环语句**。`for` 循环有多种形式；`**for**`** 循环三个部分不需括号包围**。由于++为【特定符号清单】，结尾会自动加分号而导致编译错误，所以**左大括号必须和 **`***post***`** 语句在同一行**。
    ```go
    for initialization; condition; post {
        // zero or more statements
    }
    ```
    - for循环的这三个部分每个都可以省略，如果省略`initialization`和`post`，分号也可以省略：
        ```go
        // a traditional "while" loop
        for condition {
            // ...
        }
        ```
    - 如果连 `condition` 也省略了，像下面这样：这就变成一个无限循环，尽管如此，还可以用其他方式终止循环，如一条`break`或`return` 语句。
        ```go
        // a traditional infinite loop
        for {
            // ...
        }
        ```


`for` range：在字符串或切片等数据类型的区间（range）上遍历

```go
// echo prints its command-line arguments.
package main

import (
    "fmt"
    "os"
)

func main() {
    s, sep := "", ""
    for _, arg := range os.Args[1:] {
        s += sep + arg
        sep = " "
    }
    fmt.Println(s)
}
```

- 每次循环迭代，`range` 产生一对值；**索引、在该索引处的元素值**。
- 这个例子不需要索引，但`**range**`** 的语法要求，要处理元素，必须处理索引**。一种思路是把索引赋值给一个临时变量（如 `temp`）然后忽略它的值，但 **Go 语言不允许使用无用的局部变量（local variables），因为这会导致编译错误**。（注释：这种强制要求节省了不必要的局部变量内存空间）
- **空标识符**（blank identifier，即`_`）：**在任何****语法需要变量名但程序逻辑不需要的时候****（如：在循环里），****用于丢弃不需要的循环索引****，并保留元素值**。大多数的 Go 程序员都会像上面这样使用`range`和`_`写`echo` 程序，因为隐式地而非显式地索引`os.Args`，容易写对。
- **使用显示的初始化来说明初始化变量的重要性；使用隐式的初始化来表明初始化变量不重要；**
    ```go
    s := ""  // v1： 短变量声明，不能用于包级别变量，只能用在函数内部
    var s string  // v2：**依赖于字符串的默认初始化零值机制****，被初始化为 ""**。~~
    var s = ""~~ // v3： 当声明多个变量时用到~~
    var s string = ""~~ // v4：当变量类型与初值类型不同时使用
    ```
- 每次循环迭代字符串`s` 原来的内容已经不再使用**，****将在适当时机对它进行垃圾回收****。**如果连接涉及的数据量很大，这种方式代价高昂（注释：会产生大量的垃圾，进而产生大量的gc）；


一种简单且高效的解决方案是使用 `strings` 包的 `Join` 函数：

```go
// Echo prints its command-line arguments.

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}
```



上述三个版本因为循环变量s的gc耗时，循环1亿次的性能对比相差高达一倍：

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/874cdc3d-c372-42bb-b178-6b22047715e0/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466Q7E76THJ%2F20250718%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250718T165135Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEHgaCXVzLXdlc3QtMiJGMEQCIEnRgatAz%2FvbRp81kZiDDP64372Er9LaugPTWrdaaprZAiACMsnLkeO3w0UmcpGx2WRMgDNg4VBpmixgVLn4t%2BnnOSqIBAiQ%2F%2F%2F%2F%2F%2F%2F%2F%2F%2F8BEAAaDDYzNzQyMzE4MzgwNSIMh8hi%2BdCc0zeNLfLcKtwDKj4yVNFfPhV4FXGK0%2Bf7NpNEmuTUBR8EvMV1U7S978smQVrTxeapsRoqhS4%2F6q52MK6CTvZHiLStodxhcRguwK4ihQnGdXGbv24qw7qkmF7m7ZJI7kNTzvg3zFB9oFfw4JqKTgtKxj0ZPBDyFOK7QS3SX3IHRI%2FvbySYvhy8mP8rYGVmOp%2FntD7CHZVY53u9t%2FKlLYQ615mZruXFAzVjjvP6XviiHaKIg0uXAKoix%2FCEVYuKDxcQLR6nF3Zfnije7I5K3errP63oxFAVzsGpUsvLeIyDwMKs8scU%2BzRsc688cOvTOyBEx3NcIKjwC35yeiHZbi3gqDOLanQ5zpEoCtSplgjOIpiA6cyhcgIvUm1l0KfV%2BXkqZB5M%2BuShtZ3Nut%2BvJF3ov5EtvHeaLigaPs%2F5v8lOyTWMlMMu8hZcGFqg90oLRSXhZv3OWLNPVKWNFuJ9WwaJ7WZ%2FkmLWSZ5OYKzPSaulBAzNSTFo3EkVk7DP%2BS5KlcrGKCWKLr%2FF4nlN90IFlMfPWlNGUXjtUFfoMDF6ze%2BCy0nZl3uEPcJWmOTjAE3AX84ohb7kt1ip%2FsKtMLYzzBilhNhIywbxVUN2CnHZkrx7OY1xCli6zLymqfO40L4IGoWDA%2BHdMR8wkM7pwwY6pgF2daIfuRNNIWoqU4jvZaGj1K%2BEnDRp%2F3m2W%2BGBy0BcTkDoxxHEFi%2FZL6k9XFE7JiOjWTOENBI%2B%2FGdfIjvJU5IRC24RKah%2FjPvlP6j7%2B1m4ZjFnbq6KsVQd4I4G1LFZAOWDn6T%2F0MO1Bv1RSMlz%2FJZWeh1Ah20FESFke4pwv44fl%2F9DK%2B0%2FSSZAdkqLaNQgmJfM4BxrvIUhfeWBaqpZReO4lrDAslcl&X-Amz-Signature=efd2be1854ccdedced8a09e203e0603ccbd5ee0fd794e402f8a5c55b357da9c7&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)



