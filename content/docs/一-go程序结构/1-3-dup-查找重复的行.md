---
title: "1.3 dup 查找重复的行"
date: 2025-03-27T12:15:00Z
draft: false
weight: 1003
---

# 1.3 dup 查找重复的行

对文件做拷贝、打印、搜索、排序、统计或类似事情的程序都有一个差不多的程序结构：

- 一个处理输入的循环，在每个元素上执行计算处理，在处理的同时或最后产生输出。
- 灵感来自于Unix的`uniq`命令，其寻找相邻的重复行。
- 打印标准输入中多次出现的行，以重复次数开头。该程序将引入`if`语句，`map`数据类型以及`bufio` 包。
- builtin.go 中对内建函数做了声明和文档注释，如make注释翻译：
    - make内置函数分配并初始化一个slice、map或chan类型的对象（仅限这三种类型）。
    - 和new类似，第一个参数是类型而非值。
    - 与new不同的是，make的返回类型与其参数类型T相同，而不是指向它的指针*T。具体结果取决于类型：
    - Slice(切片)：size参数指定长度。切片的容量等于其长度。可以提供第二个整数参数来指定不同的容量，该容量必须不小于长度。例如，make([]int, 0, 10)会分配一个大小为10的底层数组，并返回一个长度为0、容量10的切片，该切片由这个底层数组支持。
    - Map(映射)：分配一个空映射，并有足够的空间来容纳指定数量的元素。size参数可以省略，此时会分配一个较小的初始大小。map是随机顺序；
        - **map** 存储了键/值（key/value）的集合，对集合元素，提供常数时间的存、取或测试操作。键可以是任意类型，只要其值能用`==`运算符比较，最常见的例子是字符串；值则可以是任意类型。这个例子中的键是字符串，值是整数。
        - **内置函数 **`**make**`** 创建空 **`**map**`（译注：从功能和实现上说**，**`**Go**`** 的 **`**map**`** ****类似于 **`**Java**`** 语言中的 **`**HashMap**`**，****Python 语言中的 **`**dict**`，`Lua` 语言中的`table`，通常使用`hash`实现。遗憾的是，对于该词的翻译并不统一，**数学界术语为映射（**注释：如MyBatis中的Mapper**），**而计算机界众说纷纭莫衷一是。为了防止对读者造成误解，保留不译。
    - Channel(通道)：通道的缓冲区用指定的缓冲区容量初始化。如果为零或省略size参数，则通道是无缓冲的。
```go
// Dup prints the text of each line that appears more than
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%s\t%d\n", line, n)
		}
	}}
```



[Go提供了两种分配原语，即内建函数 new 和 make。 它们所做的事情不同，所应用的类型也不同。它们可能会引起混淆，但规则却很简单。 让我们先来看看 new。这是个用来分配内存的内建函数， 但与其它语言中的同名函数不同，它不会初始化内存，只会将内存置零。 也就是说，new(T) 会为类型为 T 的新项分配已置零的内存空间， 并返回它的地址，也就是一个类型为 *T 的值。用Go的术语来说，它返回一个指针， 该指针指向新分配的，类型为 T 的零值。](https://www.notion.so/1ca2463729b5806b8e6eff8cf792e30e#15a2463729b580d08f95d213f69cb8a3) 

[make 分配](https://www.notion.so/1502463729b58036a56bcad9774b31a5#1502463729b58059b7a8da5603c15ea4) **内置函数make（区别于new()），只用来创建map、切片、和信道**，并返回类型为 T（而非 *T 指针）的一个已初始化（而非置零）的值。



- `map`中不含某个键时不用担心，首次读到新行时，等号右边的表达式`counts[line]`的值将被**计算为其类型的零值**，对于`int`即`0`。
- `map`的迭代顺序并不确定：从实践来看，**该顺序随机，每次运行都会变化**(实测是这样)。这种设计是有意为之的，因为**能防止开发的程序依赖特定遍历顺序**，而这是无法保证的。（译注：具体可以参见这里[https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order](https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order)）map的顺序取决于使用的hash函数，hash函数为了修复DOS拒绝服务攻击做了随机化处理。[https://github.com/golang/go/issues/2630](https://github.com/golang/go/issues/2630)。
- `bufio`包使处理输入和输出方便又高效。`Scanner` 类型是该包最有用的特性之一，它读取输入并将其拆成行或单词；通常是处理行形式的输入最简单的方法。程序使用短变量声明创建 `bufio.Scanner` 类型的变量 `input`。该变量从程序的标准输入中读取内容。每次调用`input.Scan()`，即读入下一行，并移除行末的换行符；读取的内容可以调用 `input.Text()` 得到。
- `Scan`函数在读到一行时返回 `true`，不再有输入时返回 `false`。
- 类似于 C 或其它语言里的 `printf` 函数，`fmt.Printf` 函数对一些表达式产生格式化输出。该函数的首个参数是个格式字符串，指定后续参数被如何格式化。各个参数的格式取决于**“转换字符”**（conversion character），形式为百分号后跟一个字母。举个例子，`%d` 表示以十进制形式打印一个整型操作数，而`%s`则表示把字符串型操作数的值展开。
- `Printf` 有一大堆这种转换，Go程序员称之为***动词****（verb）*。下面的表格虽然远不是完整的规范，但展示了可用的很多特性：
    ```go
    %d          十进制整数
    %x, %o, %b  十六进制，八进制，二进制整数。
    %f, %e, %g  浮点数： 3.141593（定点表示法） 3.141593e+00(科学计数法)  3.141592653589793(%g。默认使用%f，当指数部分小于 -4 或者大于或等于精度，则使用 %e)
    %t          布尔：true或false
    %c          字符（rune） (Unicode码点)
    %s          字符串
    %q          带双引号的字符串"abc"或带单引号的字符'c'
    %v          变量的自然形式（natural format）
    %T          变量的类型
    %%          字面上的百分号标志（无操作数）
    ```
- 按照惯例：以字母 `f`（`format`）结尾的格式化函数**：**如fmt.Printf、log.Printf 和、fmt.Errorf，都采用 `fmt.Printf` 的格式化准则（默认情况下，`Printf` 不会换行）。以`ln`（`line`）结尾的格式化函数：如fmt.Println，以跟`%v`差不多的方式格式化参数，并在最后添加一个换行符。
    ```go
    input := bufio.NewScanner(os.Stdin)
    ```




```go
// Dup prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    counts := make(map[string]int)
    files := os.Args[1:]
    if len(files) == 0 {
        countLines(os.Stdin, counts)
    } else {
        for _, arg := range files {
            f, err := os.Open(arg)
            if err != nil {
                fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
                continue
            }
            countLines(f, counts)
            f.Close()
        }
    }
    for line, n := range counts {
        if n > 1 {
            fmt.Printf("%d\t%s\n", n, line)
        }
    }
}

func countLines(f *os.File, counts map[string]int) {
    input := bufio.NewScanner(f)
    for input.Scan() {
        counts[input.Text()]++
    }
    // NOTE: ignoring potential errors from input.Err()
}
```



[`os.Open`](http://os.open/): 函数返回的第一个值是被打开的文件（`*os.File`），其后被`Scanner`读取。返回的第二个值是内置`error`类型的值。

如果`err`等于**内置值**`**nil**`**（译注：相当于其它语言里的 **`**NULL**`**）**，那么文件被成功打开。读取文件，直到文件结束，然后调用 `Close` 关闭该文件，并释放占用的所有资源。

如果 `err` 的值不是 `nil`，**说明打开文件时出错了**。这种情况下，错误值描述了所遇到的问题。我们的错误处理非常简单，只是使用 `Fprintf` 与表示任意类型默认格式值的动词 `%v`，向标准错误流打印一条信息，然后`dup`继续处理下一个文件；`continue` 语句直接跳到 `for` 循环的下个迭代开始执行。



为了使示例代码保持合理的大小，本书开始的一些示例有意简化了错误处理，显而易见的是，应该检查 `os.Open` 返回的错误值，然而，使用`input.Scan`读取文件过程中，不大可能出现错误，因此我们忽略了错误处理。我们会在跳过错误检查的地方做说明。



注意dup.go中`countLines` 函数在其声明前被调用：**函数和包级别的变量（package-level entities）可以任意顺序声明，并不影响其被调用**。（译注：**最好还是遵循一定的规范**）。



`map`（counts）是一个由 `make` 函数创建的数据结构的**引用**。`map`作为参数传递给某函数时（countLines(f, counts)），该函数接收这个引用的一份拷贝（copy，或译为**副本**），被调用函数对 `map` 底层数据结构的任何修改，调用者函数都可以通过持有的 `map` 引用看到。在我们的例子中，`countLines` 函数向 `counts` 插入的值，也会被 `main` 函数看到。（译注：**类似于 C++ 里的****引用传递**，实际上指针是另一个指针了，但内部存的值指向同一块内存）。



`dup` 的前两个版本以"流”模式读取输入，并根据需要拆分成多个行。理论上，这些程序可以处理任意数量的输入数据。还有另一个方法，就是一次性把全部输入数据读到内存中，一次分割为多行，然后处理它们。



高级（higher-level）函数` bufio.Scanner`、`ioutil.ReadFile` 和 `ioutil.WriteFile` ，底层都使用的（lower-level）函数， `*os.File` 的 `Read` 和 `Write` 方法。

```go
// Dup3 prints the count and text of lines that
// See page 12.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		// ioutil.ReadFile is deprecated: As of Go 1.16, this function simply calls [os.ReadFile]
		// 一次性把全部输入数据读到内存中
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		// 把字符串按换行符切割成行的切片
		// 与字节切片（byte slice，类似java的byte[]）转成string ，然后拼接成字符串。
		// 相反： strings.Join(os.Args[1:], " ")
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%s: %d\n", line, n)
			}
		}
	}
}
```



