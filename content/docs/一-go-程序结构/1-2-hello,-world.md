---
title: "1.2 Hello, World"
date: 2024-11-30T04:24:00Z
draft: false
weight: 1002
---

# 1.2 Hello, World

## 一、hello world

### 源码

以下是Go语言版本的hello world，“hello world”案例首次出现于1978年出版的C语言圣经《The C Programming Language》

```go
// Helloworld is our first Go program.
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
```



### 源码解读

- **main package & main function**：每个源文件都以一条`package`声明语句开始，通过包（package）组织****(包类似于其它语言里的**库（libraries）/模块（modules）**）。package main 定义了一个**独立可执行的程序，而不是一个库**。在`main`里的`main`函数也很特殊，它是**整个程序执行时的入口**（译注：C 系语言差不多都这样）;
- **Unicode原生支持**：Go的三位设计者中**2位为UTF-8的设计者**;
- **import**：紧跟着必须告诉编译器的一系列导入（import）的包。Go语言的代码：一个包由位于单个目录下的一个或多个`.go`源代码文件组成，目录定义包的作用。**Go严格要求：缺少了必要的包或者导入了不需要的包，程序都无法编译通过**，避免了程序开发过程中引入未使用的包。`goimports` 自动导入；
- **gofmt**：Go语言在代码格式上**采取了很强硬的态度，****避免了无尽的无意义的琐碎争执**（译注：也导致了Go语言的[TIOBE](https://www.tiobe.com/tiobe-index/)排名较低，因为缺少撕逼的话题）
    - 编译器会主动**把【特定符号清单】后的换行符转换为分号**，包括：行末的**标识符、整数、浮点数、虚数、字符或字符串文字、关键字**`**break**`**、**`**continue**`**、**`**fallthrough**`**或 **`**return**`** 中的一个、运算符和分隔符 **`**++**`**、**`**--**`**、**`**)**`**、**`**]**`** 或 **`**}**`** 。**
    - 所以除非一行上有多条语句，否则不需要在语句或者声明的末尾添加分号（编辑器保存时会自动删除行末的分号，自动执行`gofmt`自动格式化代码）。更重要的是，这样可以做多种自动源码转换，如果放任Go语言代码格式，这些转换就不大可能了。
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_15724637-29b5-805e-afe9-c7722c11b100.jpg)
- **Go标准库**：提供了184个包，以支持常见功能，如输入、输出、排序以及文本处理；
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_15724637-29b5-80b5-a492-c5ebcda00da7.jpg)


### go build & go run

- `**go build helloworld.go**`：编译；由于每个go程序都通过goroutine运行，go的二进制文件体积很大；**体积对比**：C版本 33k，go版本 2.1M（**65倍**）；
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_a143a20e-a8b3-4ce7-a8e9-fb59bfd1bd92.jpg)
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_1c324637-29b5-80e3-b843-c4d8567bac22.jpg)
    - C版本 hello，world: gcc -o hello hello.c
        ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_1c324637-29b5-8080-b53b-c246fa8be498.jpg)
- `**go run helloworld.go**`：编译后运行；


## 二、echo

### 源码

大多数的程序都是：处理输入、产生输出；输入包括文件、网络连接、其它程序的输出、敲键盘的用户、命令行参数或其它类似输入源

```go
// Echo prints its command-line arguments.
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

### 源码解读

- **包注释/文档注释**：**按照惯例在每个包的包声明前，****以程序名开头，从整体角度对程序做个简要描述**。
- `os`包以与平台无关的方式提供了一些与操作系统交互的函数和变量。`os.Args`变量是一个**字符串（string）的切片（slice）**（/slaɪs/，**动态数组**）。和大多数编程语言类似，区间索引时，Go 语言里也采用**左闭右开，**因为**这样可以简化逻辑。**`s[m:n]`，省略m和n时，会默认传入 `0` 或 `len(s)`。
- 导入多个包，习惯上用括号把它们**括起来写成列表形式**；`gofmt`工具格式化时**按照字母顺序对包名排序。**
- `var` 声明：如果变量没有显式初始化，则被隐式地赋予其类型的**零值（zero value）**，**数值类型是 **`**0**`**，字符串类型是空字符串 **`**""**`
- 符号`:=`是**短变量声明（short variable declaration）**的一部分，这是定义一个或多个变量并根据它们的初始值为这些变量赋予适当类型的语句。只能用于在函数体中，而不能用在包级别；
- 对数值类型，Go语言提供了常规的数值和逻辑运算符。而对 `string` 类型，`+` 运算符连接字符串（译注：和 C++ 或者 JavaScript 是一样的）。等价于：`s=s+sep+os.Args[i]`。
- 自增语句`i++`：给`i`加`1`；这和 `i+=1` 以及 `i=i+1` 都是等价的。**表达式是赋值=的右边部分，而语句是独立完整一条。****Go有意将**`**i++**`**设计成语句而不是表达式，简化了语言规范，提高了代码可读性，符合 Go 的"简单明确"设计哲学**。Go中`**j=i++**`** 非法，**而且 `++` 和 `--` 都只能放在变量名后面，因此`**--i**`**也非法。**
- `s += sep + arg`：
    - 创建了一个**新的字符串**，包含 s + sep + os.Args[i]
        ```shell
        // 假设循环执行过程
        s = ""                    // 空字符串
        s = "" + " " + "arg1"     // 新字符串 " arg1"，旧字符串 "" 可回收
        s = " arg1" + " " + "arg2" // 新字符串 " arg1 arg2"，旧字符串 " arg1" 可回收
        s = " arg1 arg2" + " " + "arg3" // 新字符串 " arg1 arg2 arg3"，旧字符串 " arg1 arg2" 可回收
        ```
    - 将新字符串赋值给变量 s
    - **旧的字符串内容**变成了"垃圾"，等待gc垃圾回收。**如果连接涉及的数据量很大，这种方式代价高昂（注释：会产生大量的垃圾，进而产生大量的gc）**；Go 使用**并发标记清除（Concurrent Mark and Sweep）**垃圾回收器：
        - **自动检测**：GC 会自动检测不再被引用的内存
        - **并发执行**：GC 与程序并发运行，减少停顿时间
        - **适时回收**：在合适的时机回收垃圾内存
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
    - `for` range：在字符串或切片等数据类型的区间（range）上遍历
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
        - 每次循环迭代，`range` 产生一对值；**索引(数组下标)、在该索引处的元素值**。
        - 这个例子不需要索引，但`**range**`** 的语法要求，要处理元素，必须处理索引**。一种思路是把索引赋值给一个临时变量（如 `temp`）然后忽略它的值，但 **Go 语言不允许使用无用的局部变量（local variables），因为这会导致编译错误**。（注释：这种强制要求节省了不必要的局部变量内存空间）
        - **空标识符**（blank identifier，即`_`）：**在任何****语法需要变量名但程序逻辑不需要的时候****（如：在循环里），用于丢弃不需要的循环索引，并保留元素值**。大多数的 Go 程序员都会像上面这样使用`range`和`_`写`echo` 程序，因为隐式地而非显式地索引`os.Args`，容易写对。
- **使用显示的初始化来说明初始化变量的重要性；使用隐式的初始化来表明初始化变量不重要；**
    ```go
    s := ""  // v1： 短变量声明，不能用于包级别变量，只能用在函数内部
    var s string  // v2：**依赖于字符串的默认初始化零值机制****，被初始化为 ""**。~~
    var s = ""~~ // v3： 当声明多个变量时用到~~
    var s string = ""~~ // v4：当变量类型与初值类型不同时使用
    ```


### 性能优化

源码：一种简单且高效的实现是：使用 `strings` 包的 `Join` 函数

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



- 上述的函数注释会出现在IDE中的提示框
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_23924637-29b5-8017-97bd-d374e9a10ade.jpg)
- 上述三个版本因为循环变量s的gc耗时，循环1亿次的性能对比相差高达一倍：
    ![](/images/14e24637-29b5-8021-9e7e-e1ba4d8be658/image_23824637-29b5-8084-964f-d8a2912aa3e1.jpg)


源码解读：看看标准库的实现

strings.Join源码：

```go
// Join concatenates the elements of its first argument to create a single string. The separator
// string sep is placed between elements in the resulting string.
func Join(elems []string, sep string) string {
	// 边界检查
	switch len(elems) {
	case 0:
		return ""  // 空切片返回空字符串
	case 1:
		return elems[0] // 只有一个元素直接返回
	}

	var n int  // 提前计算输出字符串的总长度 = 分隔符长度 × (元素个数 - 1) + 所有元素长度
	if len(sep) > 0 { 
		if len(sep) >= maxInt/(len(elems)-1) {   // 边界检查：检查整数溢出，防止内存分配过大
			panic("strings: Join output length overflow")
		}
		n += len(sep) * (len(elems) - 1)
	}
	for _, elem := range elems {
		if len(elem) > maxInt-n {  // 边界检查：检查整数溢出，防止内存分配过大
			panic("strings: Join output length overflow")
		}
		n += len(elem)
	}

	var b Builder
	b.Grow(n)    // 一次性预分配足够的内存空间，避免多次内存分配
	b.WriteString(elems[0])  // 直接写入预分配的内存，写入第一个元素
	for _, s := range elems[1:] {
		b.WriteString(sep)  // 继续写入同一块内存，直接追加，避免了复制写入分隔符
		b.WriteString(s)   // 写入元素
	}
	return b.String()
}
```



strings.Builder:

```go
// A Builder is used to efficiently build a string using [Builder.Write] methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Builder.
type Builder struct {
	addr *Builder // of receiver, to detect copies by value. 用于防止值复制，确保 Builder 只能通过指针使用

	// External users should never get direct access to this buffer, since
	// the slice at some point will be converted to a string using unsafe, also
	// data between len(buf) and cap(buf) might be uninitialized.
	buf []byte.    // 底层字节切片，实际存储字符串内容
}

// copyCheck - 防止开发者错误地复制 Builder（应该用指针）
func (b *Builder) copyCheck() {
    if b.addr == nil {
        // 初始化 addr 字段
        b.addr = (*Builder)(abi.NoEscape(unsafe.Pointer(b)))
    } else if b.addr != b {
        // 检测到值复制，panic
        panic("strings: illegal use of non-zero Builder copied by value")
    }
}

//  转换为字符串
func (b *Builder) String() string {
    return unsafe.String(unsafe.SliceData(b.buf), len(b.buf))  // Builder 性能优势的关键：使用 unsafe.String 直接将字节切片转换为字符串；**零拷贝：不复制数据，直接共享底层内存**；
}

// **预分配内存**
func (b *Builder) Grow(n int) {
    b.copyCheck()
    if n < 0 {
        panic("strings.Builder.Grow: negative count")
    }
    if cap(b.buf)-len(b.buf) < n {
        b.grow(n)
    }
}

func (b *Builder) grow(n int) {
    // **使用 2倍扩容 策略，减少扩容次数**；新容量 = 2*当前容量 + 需要的空间；确保至少有 n 字节的额外空间
    buf := bytealg.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
    copy(buf, b.buf)  // 复制现有数据
    b.buf = buf       // 更新引用
}

// Write 系列方法
func (b *Builder) Write(p []byte) (int, error) {
    b.copyCheck()
    b.buf = append(b.buf, p...)
    return len(p), nil
}
func (b *Builder) WriteString(s string) (int, error) {
    b.copyCheck()
    b.buf = append(b.buf, s...)
    return len(s), nil
}
func (b *Builder) WriteRune(r rune) (int, error) {
    b.copyCheck()
    n := len(b.buf)
    b.buf = utf8.AppendRune(b.buf, r)
    return len(b.buf) - n, nil
}
```



**内存使用对比：**

```go
// 传统方式：多次分配
elems := []string{"a", "b", "c", "d"}
s := ""
s += "a"    // 分配 1 字节
s += "b"    // 分配 2 字节，释放 1 字节
s += "c"    // 分配 3 字节，释放 2 字节
s += "d"    // 分配 4 字节，释放 3 字节
// 总共分配 10 字节，产生 6 字节垃圾

// Builder 方式：一次分配
var b Builder
b.Grow(4)   // 一次性分配 4 字节
b.WriteString("a")  // 写入
b.WriteString("b")  // 写入
b.WriteString("c")  // 写入
b.WriteString("d")  // 写入
// 总共分配 4 字节，无垃圾产生
```



Go中将字符串设计为不可变的，优先保证安全性，再通过strings.Builder构建器优化性能。这意味着一切对字符串的操作都转换为：重新赋值，创建新字符串。

```go
s := "hello"
s[0] = 'H'  // 编译错误：cannot assign to s[0]
s = "Hello" // 重新赋值，创建新字符串
```

| 语言 | 字符串特性 | 主要优势 | 性能优化方案 | 
| --- | --- | --- | --- | 
| C++ | 可变 | 性能高、内存效率 | 直接操作 | 
| Rust | 可变 | 内存安全、性能高 | 直接操作 | 
| Go | 不可变 | 线程安全、内存共享 | strings.Builder | 
| Java | 不可变 | 线程安全、缓存友好 | StringBuilder | 
| Python | 不可变 | 简单、安全 | join() 方法 | 
| JavaScript | 不可变 | 简单、安全 | 数组 join() | 
| C# | 不可变 | 线程安全 | StringBuilder | 

> 更多内容，在后续的string部分深入。



## 三、dup

### 源码

文件处理类程序都有相似的结构：**一个处理输入的循环，在每个元素上执行计算处理，在处理的同时或最后产生输出**。如文件的拷贝、打印、搜索、排序、统计等。

下面模拟Unix的`uniq`命令，其寻找相邻的重复行，并打印标准输入中多次出现的行，以重复次数开头。

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



### 源码解读-fmt.print

fmt 包的设计体现了 Go 语言的几个重要特性：

1. **接口驱动**：通过 Stringer、Formatter 等接口支持自定义格式化
1. **性能优化**：使用对象池、快速路径等减少分配
1. **类型安全**：通过反射和类型断言处理各种类型
1. **错误处理**：完善的 panic 恢复和错误报告
1. **易用性**：简洁的 API 设计，支持多种格式化选项
**1. 核心接口定义**

**State 接口**

```go
type State interface {
    Write(b []byte) (n int, err error)  // 写入格式化输出
    Width() (wid int, ok bool)          // 获取宽度设置
    Precision() (prec int, ok bool)     // 获取精度设置
    Flag(c int) bool                    // 检查标志位
}
```

**Formatter 接口**

```go
type Formatter interface {
    Format(f State, verb rune)  // 自定义格式化方法
}
```

**Stringer 接口**

```go
type Stringer interface {
    String() string  // 默认字符串表示
}
```

**GoStringer 接口**

```go
type GoStringer interface {
    GoString() string  // Go 语法表示
}
```

**2. 核心数据结构**

**pp 结构体（打印状态）**

```go
type pp struct {
    buf buffer           // 输出缓冲区
    
    arg any              // 当前参数
    value reflect.Value  // 反射值
    
    fmt fmt              // 格式化器
    
    reordered bool       // 是否重排序
    goodArgNum bool      // 参数编号是否有效
    panicking bool       // 是否在 panic 中
    erroring bool        // 是否在错误处理中
    wrapErrs bool        // 是否包装错误
    wrappedErrs []int    // 包装的错误索引
}
```

**buffer 类型**

```go
type buffer []byte

func (b *buffer) write(p []byte) {
    *b = append(*b, p...)
}

func (b *buffer) writeString(s string) {
    *b = append(*b, s...)
}
```

**3. 主要函数族**

***动词****（verb）*：f：format、ln：line

```go
%d          十进制整数
%x, %o, %b  十六进制，八进制，二进制整数。 % x会在字节之间插入空格
%f, %e, %g  浮点数： 3.141593（定点表示法） 3.141593e+00(科学计数法)  3.141592653589793(%g。默认使用%f，当指数部分小于 -4 或者大于或等于精度，则使用 %e)
%t          布尔：true或false
%c          字符（rune） (Unicode码点)
%s          字符串
%q          带双引号的字符串"abc"或带单引号的字符'c'，%#q 会尽可能使用反引号
%v          打印**任意值value，****变量的自然形式（natural format）**
%T          变量的类型
%%          字面上的百分号标志（无操作数）
```

```go
type T struct {
	a int
	b float64
	c string
}
t := &T{ 7, -2.35, "abc\tdef" }
fmt.Printf("%v\n", t)   // &{7 -2.35 abc   def}
fmt.Printf("%+v\n", t)  // &{a:7 b:-2.35 c:abc     def}
fmt.Printf("%#v\n", t)  // &main.T{a:7, b:-2.35, c:"abc\tdef"}
fmt.Printf("%#v\n", timeZone) // map[string] int{"CST":-21600, "PST":-28800, "EST":-18000, "UTC":0, "MST":-25200}
```

%v的性能开销：

1. **反射开销**
    ```go
    func (p *pp) printArg(arg any, verb rune) {
        // 当 %v 遇到复杂类型时，会走到这里
        default:
            if !p.handleMethods(verb) {
                // **使用反射 - 性能瓶颈**
                p.printValue(reflect.ValueOf(f), verb, 0)
            }
    }
    // 反射调用示例
    func (p *pp) printValue(value reflect.Value, verb rune, depth int) {
        switch value.Kind() {
        case reflect.Struct:
            // **每个字段都需要反射访问**
            for i := 0; i < value.NumField(); i++ {
                field := value.Field(i)  // 反射调用
                p.printValue(field, verb, depth+1)
            }
        case reflect.Map:
            // **map 需要排序和反射访问**
            sorted := fmtsort.Sort(value)  // 反射排序
            for _, kv := range sorted {
                p.printValue(kv.Key, verb, depth+1)   // 反射
                p.printValue(kv.Value, verb, depth+1) // 反射
            }
        }
    }
    ```
1. **接口方法调用开销: String() 方法调用**
```go
func (p *pp) handleMethods(verb rune) (handled bool) {
    switch verb {
    case 'v', 's':
        if stringer, ok := p.arg.(Stringer); ok {
            // 每次都会额外调用 String() 方法 - 虚函数调用开销
            p.fmtString(stringer.String(), verb)
            return true
        }
    }
    return false
}

// 慢：每次都会额外调用 String() 方法
fmt.Printf("%v", person)

// 快：直接访问字段
fmt.Printf("Name: %s, Age: %d", person.Name, person.Age)
```

1. **内存分配问题**
```go
// 临时对象创建
func (p *pp) printValue(value reflect.Value, verb rune, depth int) {
    case reflect.Slice, reflect.Array:
        // 每个元素都可能创建临时对象
        for i := 0; i < value.Len(); i++ {
            elem := value.Index(i)  // 可能创建临时 Value
            p.printValue(elem, verb, depth+1)
        }
}

// 字符串拼接开销：每次 %v 都可能涉及字符串拼接
func (p *pp) fmtString(v string, verb rune) {
    case 'v':
        if p.fmt.sharpV {
            p.fmt.fmtQ(v)  // 可能添加引号等字符
        } else {
            p.fmt.fmtS(v)  // 直接写入
        }
}
```





**Printf 系列（格式化打印）**

```go
func Printf(format string, a ...any) (n int, err error) {
    return Fprintf(os.Stdout, format, a...)
}

func Fprintf(w io.Writer, format string, a ...any) (n int, err error) {
    p := newPrinter()
    p.doPrintf(format, a)
    n, err = w.Write(p.buf)
    p.free()
    return
}

func Sprintf(format string, a ...any) string {
    p := newPrinter()
    p.doPrintf(format, a)
    s := string(p.buf)
    p.free()
    return s
}
```

**Print 系列（默认格式）**

```go
func Print(a ...any) (n int, err error) {
    return Fprint(os.Stdout, a...)
}

func Fprint(w io.Writer, a ...any) (n int, err error) {
    p := newPrinter()
    p.doPrint(a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

**Println 系列（带换行）**

```go
func Println(a ...any) (n int, err error) {
    return Fprintln(os.Stdout, a...)
}

func Fprintln(w io.Writer, a ...any) (n int, err error) {
    p := newPrinter()
    p.doPrintln(a)
    n, err = w.Write(p.buf)
    p.free()
    return
}
```

**4. 核心处理逻辑**

**doPrintf - 格式化字符串解析**

```go
func (p *pp) doPrintf(format string, a []any) {
    end := len(format)
    argNum := 0
    
    for i := 0; i < end; {
        // 1. 处理普通文本
        lasti := i
        for i < end && format[i] != '%' {
            i++
        }
        if i > lasti {
            p.buf.writeString(format[lasti:i])
        }
        
        if i >= end {
            break
        }
        
        // 2. 处理格式说明符
        i++ // 跳过 %
        
        // 3. 解析标志位
        p.fmt.clearflags()
        for ; i < end; i++ {
            c := format[i]
            switch c {
            case '#': p.fmt.sharp = true
            case '0': p.fmt.zero = true
            case '+': p.fmt.plus = true
            case '-': p.fmt.minus = true
            case ' ': p.fmt.space = true
            default:
                // 快速路径：简单动词
                if 'a' <= c && c <= 'z' && argNum < len(a) {
                    p.printArg(a[argNum], rune(c))
                    argNum++
                    i++
                    continue
                }
                break
            }
        }
        
        // 4. 解析宽度和精度
        // 5. 解析动词
        // 6. 打印参数
    }
}
```

**printArg - 参数打印**

```go
func (p *pp) printArg(arg any, verb rune) {
    p.arg = arg
    p.value = reflect.Value{}
    
    if arg == nil {
        switch verb {
        case 'T', 'v':
            p.fmt.padString(nilAngleString)
        default:
            p.badVerb(verb)
        }
        return
    }
    
    // 特殊处理
    switch verb {
    case 'T':
        p.fmt.fmtS(reflect.TypeOf(arg).String())
        return
    case 'p':
        p.fmtPointer(reflect.ValueOf(arg), 'p')
        return
    }
    
    // 类型分发
    switch f := arg.(type) {
    case bool:
        p.fmtBool(f, verb)
    case int:
        p.fmtInteger(uint64(f), signed, verb)
    case string:
        p.fmtString(f, verb)
    case []byte:
        p.fmtBytes(f, verb, "[]byte")
    default:
        // 尝试接口方法
        if !p.handleMethods(verb) {
            // 使用反射
            p.printValue(reflect.ValueOf(f), verb, 0)
        }
    }
}
```

**5. 类型格式化方法**

**整数格式化**

```go
func (p *pp) fmtInteger(v uint64, isSigned bool, verb rune) {
    switch verb {
    case 'v':
        if p.fmt.sharpV && !isSigned {
            p.fmt0x64(v, true)
        } else {
            p.fmt.fmtInteger(v, 10, isSigned, verb, ldigits)
        }
    case 'd':
        p.fmt.fmtInteger(v, 10, isSigned, verb, ldigits)
    case 'b':
        p.fmt.fmtInteger(v, 2, isSigned, verb, ldigits)
    case 'o', 'O':
        p.fmt.fmtInteger(v, 8, isSigned, verb, ldigits)
    case 'x':
        p.fmt.fmtInteger(v, 16, isSigned, verb, ldigits)
    case 'X':
        p.fmt.fmtInteger(v, 16, isSigned, verb, udigits)
    }
}
```

**字符串格式化**

```go
func (p *pp) fmtString(v string, verb rune) {
    switch verb {
    case 'v':
        if p.fmt.sharpV {
            p.fmt.fmtQ(v)  // 带引号
        } else {
            p.fmt.fmtS(v)  // 普通字符串
        }
    case 's':
        p.fmt.fmtS(v)
    case 'x':
        p.fmt.fmtSx(v, ldigits)  // 十六进制
    case 'X':
        p.fmt.fmtSx(v, udigits)  // 大写十六进制
    case 'q':
        p.fmt.fmtQ(v)  // 带引号
    }
}
```

**6. 接口方法处理**

**handleMethods - 接口方法调用**

```go
func (p *pp) handleMethods(verb rune) (handled bool) {
    if p.erroring {
        return
    }
    
    // 1. 检查 Formatter 接口
    if formatter, ok := p.arg.(Formatter); ok {
        handled = true
        defer p.catchPanic(p.arg, verb, "Format")
        formatter.Format(p, verb)
        return
    }
    
    // 2. 检查 GoStringer 接口（%#v）
    if p.fmt.sharpV {
        if stringer, ok := p.arg.(GoStringer); ok {
            handled = true
            defer p.catchPanic(p.arg, verb, "GoString")
            p.fmt.fmtS(stringer.GoString())
            return
        }
    }
    
    // 3. 检查 Stringer 和 error 接口
    switch verb {
    case 'v', 's', 'x', 'X', 'q':
        switch v := p.arg.(type) {
        case error:
            handled = true
            defer p.catchPanic(p.arg, verb, "Error")
            p.fmtString(v.Error(), verb)
            return
        case Stringer:
            handled = true
            defer p.catchPanic(p.arg, verb, "String")
            p.fmtString(v.String(), verb)
            return
        }
    }
    return false
}
```

**7. 性能优化**

**对象池复用**

```go
var ppFree = sync.Pool{
    New: func() any { return new(pp) },
}

func newPrinter() *pp {
    p := ppFree.Get().(*pp)
    p.panicking = false
    p.erroring = false
    p.wrapErrs = false
    p.fmt.init(&p.buf)
    return p
}

func (p *pp) free() {
    if cap(p.buf) > 64*1024 {
        p.buf = nil
    } else {
        p.buf = p.buf[:0]
    }
    // ... 清理其他字段
    ppFree.Put(p)
}
```

**8. 错误处理**

**panic 恢复**

```go
func (p *pp) catchPanic(arg any, verb rune, method string) {
    if err := recover(); err != nil {
        if v := reflect.ValueOf(arg); v.Kind() == reflect.Pointer && v.IsNil() {
            p.buf.writeString(nilAngleString)
            return
        }
        
        if p.panicking {
            panic(err)
        }
        
        p.buf.writeString(percentBangString)
        p.buf.writeRune(verb)
        p.buf.writeString(panicString)
        p.buf.writeString(method)
        p.buf.writeString(" method: ")
        p.panicking = true
        p.printArg(err, 'v')
        p.panicking = false
        p.buf.writeByte(')')
    }
}
```



```go
type MyString string

func (m MyString) String() string {
	return fmt.Sprintf("MyString=%s", m) // 错误：调用 Sprintf 来构造 String 方法，会无限递归你的的 String 方法
}

// 解法：将该实参转换为基本的字符串类型，它没有这个方法
func (m MyString) String() string {
	return fmt.Sprintf("MyString=%s", string(m)) // 可以：注意转换
}
```



更多详情请参阅源码和`godoc`对`fmt`包的说明文档。

