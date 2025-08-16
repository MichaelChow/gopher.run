---
title: "2.3 slice与append"
date: 2024-12-22T04:12:00Z
draft: false
weight: 2003
---

# 2.3 slice与append

## 一、array

数组：具有固定长度且拥有零个或多个**相同数据类型元素的序列****。**

**由于固定长度，Go中很少直接使用array**，更多使用动态长度的slice，而**slice底层由array实现。**



数组在Go和C中的主要区别：

- **数组是值**。将一个数组赋予另一个数组会复制其所有元素。数组为值的属性很有用，但代价高昂；若你想要C那样的行为和效率，你可以传递一个指向该数组的指针。但这并不是Go的习惯用法，切片才是。
    ```go
    func Sum(a *[3]float64) (sum float64) {
    	for _, v := range *a {
    		sum += v
    	}
    	return
    }
    array := [...]float64{7.0, 8.5, 9.1}
    x := Sum(&array)  // 注意显式的取址操作
    ```
- 特别地，若将某个数组传入某个函数，它将接收到该数组的一份**副本**而非指针。
- 数组的长度**是其类型的一个组成部分**。类型 `[10]int` 和 `[20]int` 是不同的类型，长度需要在编译阶段确定，所以必须是常量表达式。
    - 在数组字面值中，如果在数组的**长度位置出现的是“...”省略号**，则**表示数组的长度是根据初始化值的个数来计算**
        ```go
        q := [...]int{1, 2, 3}
        r := [...]int{99: -1} // 前99个元素初始化为零值，最后一个初始化为-1
        ```
    - 数组元素可通过索引下标来访问，每个元素都被默认初始化为元素类型对应的零值，或使用字面量初始化；
        ```go
        var r [3]int = [3]int{1, 2}
        ```


**数组的可比较性 comparable**：如果一个数组的元素类型是可以相互比较的，数组类型相同（包括数组长度一致），则数组是可以相互比较==的，只有当两个数组的所有元素都是相等的时候数组才是相等的；

- crypto/sha256包的Sum256函数对一个任意的byte slice类型的数据生成一个对应的消息摘要。消息摘要有256bit大小，因此对应[32]byte数组类型。如果两个消息摘要是相同的，那么可以认为两个消息本身也是相同（译注：理论上有HASH码碰撞的情况，但是实际应用可以基本忽略）；如果消息摘要不同，那么消息本身必然也是不同的。下面的例子用SHA256算法分别生成“x”和“X”两个信息的摘要：
- crypto/sha256：两个消息虽然只有一个字符的差异，但是生成的消息摘要则几乎有一半的bit位是不相同的
    ```go
    import "crypto/sha256"
    func main() {
        c1 := sha256.Sum256([]byte("x"))  // 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881, [32]uint8
        c2 := sha256.Sum256([]byte("X"))  // 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
    }
    ```


**数组作为函数参数传递的机制**：**值传递，而非引用传递**。当调用一个函数时，每个传入的参数都会创建一个**副本**，然后赋值给对应的函数变量**（**所以函数参数变量接收的是一个副本，并不是原始的参数）。

当**传递大的数组时将变得很低效，并且在函数内部对数组的任何修改都仅影响副本，而不是原始数组。**这种情况下。Go把数组和其他的类型都看成**值传递**（**不同于其他语言中数组都是隐式的使用引用传递**）。

**Go中可以显示的传递一个数组的指针给函数，****这样在函数内部对数组的任何修改都会反应到原始数组上，这样很高效**；但由于数组的长度不可变特性，无法为数组添加或者删除元素，除了在SHA256、矩阵变换这类固定长度的情况下，很少直接使用数组；

```go
// 用于给[32]byte类型的数组清零
func zero(ptr *[32]byte) {
    for i := range ptr {
        ptr[i] = 0
    }
}
```

```go
// zero函数更简洁版本：
func zero(ptr *[32]byte) {
    *ptr = [32]byte{}  // 数组字面值[32]byte{}就可以生成一个32字节的数组。而且每个数组的元素都是零值初始化，也就是0。
}
```



## 二、pointer

**变量存储值，变量被称为****可寻址的值****，对应一个保存了变量对应类型值的****内存地址**。

**一个指针对应变量的内存地址。**通过指针可以绕过变量的名字**直接读或更新对应变量的值。**通过变量名或表达式访问（如x[i]或x.f），必定能接受`&`取地址操作。

对一个变量取地址(p := &v)、复制指针，都是**为原变量创建了新的别名**。`*p`就是变量v的别名。

**指针特别有价值的地方在于我们可以不用变量名的情况下，直接访问一个变量。**但因为此，**要找到一个变量的所有访问者并不容易**，我们必须知道变量全部的别名（译注：**这是Go语言的垃圾回收器所做的工作**）。

不仅仅是指针会创建别名，很多**其他引用类型也会创建别名**，例如slice、map和chan，甚至结构体、数组和接口**都会创建所引用变量的别名**。



对于聚合类型每个成员，如结构体的每个字段、或者是数组的每个元素也都是对应一个变量，因此可以被取地址。

```go
x := 1          // &x: 读取x变量的内存地址；p指针指向变量、p指针保存了x变量的内存地址，**其数据类型为 *int** （**指向int类型的指针**）
p := &x         // p, of type *int, points to x， *p: 读取p指针指向的变量的值
fmt.Println(*p) // "1"   // 因为*p对应一个变量，所以可以赋值语句更新指针所指向的变量的值。
*p = 2          // equivalent to x = 2
fmt.Println(x)  // "2"
```



Go中**函数中返回局部变量的地址也是安全的**。因为指针p依然引用这个变量，Go中的GC使用的简单的标记清除算法的可达性树法不会识别为垃圾变量，但这会导致内存泄露。

```go
// point study
package main

import "fmt"

func main() {
	// 在局部变量地址&v被返回之后依然有效，因为指针p依然引用这个变量。
	var p = f()
	fmt.Println(p)
	fmt.Println(f() == f()) // "false"

}

func f() *int {
	v := 1
	return &v
}
```

**如果将指针作为参数调用函数，可以在函数中通过该指针来更新变量的值**。（译注：这是对C语言中`++v`操作的模拟，这里只是为了说明指针的用法，incr函数模拟的做法并不推荐）：

```go
// incr study
package main

import "fmt"

func main() {
	v := 1
	incr(&v)              // v = 2,return 2
	fmt.Println(incr(&v)) // v = 3 return 3
}

func incr(p *int) int {
	*p++ // 非常重要：只是增加p指向的变量的值，并不改变p指针！！！
	return *p
}
```


nil是一个**预声明的标识符，**表示**指针、通道、函数、接口、映射或切片类型**的零值

```go
// src/builtin/builtin.go
// nil is a predeclared identifier representing the zero value for a
// pointer, channel, func, interface, map, or slice type.
var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
// Type只是为了文档目的，它是一个占位符，代表任何Go类型。在给定的函数调用中代表相同的类型。
type Type int
```

### example：echo

标准库中的flag包的关键技术为大量使用到了指针，它使用命令行参数来设置对应变量的值。

为了说明这一点，在早些的echo版本中，就包含了两个可选的命令行参数：`-n`用于忽略行尾的换行符，`-s sep`用于指定分隔字符（默认是空格）。下面这是第四个版本，对应包路径为gopl.io/ch2/echo4。



调用flag.Bool函数会创建一个新的对应布尔型标志参数的变量。它有三个属性：第一个是命令行标志参数的名字“n”，然后是该标志参数的默认值（这里是false），最后是该标志参数对应的描述信息。

如果用户在命令行输入了一个无效的标志参数，或者输入`-h`或`-help`参数，那么将打印所有标志参数的名字、默认值和描述信息。

```go
// Echo4 prints its command-line arguments.
// See page 33.
package main

import (
	"flag"
	"fmt"
	"strings"
)

// 定义命令行参数名、默认值、和描述信息。
var n = flag.Bool("n", false, "omit trailing newline") // *bool类型的指针，默认为false不换行（终端会打印出一个%作为无换行符的标识）
var sep = flag.String("s", " ", "separator")           // *string类型的指针，默认为空格

func main() {
	// 解析命令行参数，更新每个标志参数对应变量的值（之前是默认值）。
	// 解析命令行参数时遇到错误，默认将打印相关的提示信息，然后调用os.Exit(2)终止程序。
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep)) // 打印命令行参数，*sep 是一个字符串指针，它的值是通过命令行参数 -s 指定的分隔符。

	if !*n { // 如果命令行参数 -n 没有指定，就打印一个换行符
		fmt.Println()
	}
}
```



## **三、slice**

slice（切片,/slaɪs/）`[]T`：**表示一个拥有相同类型元素的可变长度的序列****，**用来访问**底层引用和封装的数组**的元素（引用类型）；

slice的三个属性：

- pointer: 指向**第一个可以从slice中访问的元素**（数组中的任意一个元素）
- len：slice的元素个数，不能超过cap
- cap：slice的起始元素到底层数组的最后一个元素间元素的个数
多个slice可以引用同一个底层数组及其任何位置：

```go
months := [...]string{1: "January", 2: "February", 3: "March", 4: "April", 5: "May", 6: "June", 7: "July", 8: "Augest", 9: "September", 10: "October", 11: "November", 12: "December"}  // **[13]string 模拟slice引用的底层array**

Q2 := months[4:7]    //  **[]string slice,poiner=, len=3, cap=9**Q2[:10] // // panic: runtime error: slice bounds out of range [:10] with capacity 9

summer := months[6:9]  // **[]string slice, poiner= , len=3, cap=7**
```

![](/images/16424637-29b5-8082-a2c5-c95383fbf49c/image_16424637-29b5-8037-840d-dd93459bc516.jpg)



**example**：reverse

```go
//  reverse a slice of ints in place.
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}

a := [...]int{0, 1, 2, 3, 4, 5}
reverse(a[:])
fmt.Println(a) // "[5 4 3 2 1 0]"


// 一种将slice元素循环向左旋转n个元素的方法是三次调用reverse反转函数
// 第一次是反转开头的n个元素，然后是反转剩下的元素，最后是反转整个slice的元素。
//（如果是向右循环旋转，则将第三个函数调用移到第一个调用位置就可以了。）
s := []int{0, 1, 2, 3, 4, 5}
// Rotate s left by two positions.
reverse(s[:2])
reverse(s[2:])
reverse(s)
fmt.Println(s) // "[2 3 4 5 0 1]"
```



**slice不直接支持比较运算符：**

- 原因1：一个slice的元素是间接引用的，当slice声明为[]interface{}时，slice的元素甚至可以是自身。虽然有很多办法处理这种情形，但是没有一个是简单有效的。
- 原因2：因为slice的元素是**间接引用的**，**slice本身的值（不是元素的值）****在不同的时刻可能包含不同的元素**，因为底层数组的元素可能会被修改。而如Go语言中map的key只做简单的浅拷贝，它要求key在整个生命周期内保持不变性（译注：如slice扩容，就会导致其本身的值/地址变化）。而用深度相等判断的话，显然在map的key这种场合不合适。对于像指针或chan之类的引用类型，==相等测试可以判断两个是否是引用相同的对象。一个针对slice的浅相等测试的==操作符可能是有一定用处的，也能临时解决map类型的key问题，**但是slice和数组不同的相等测试行为会让人困惑**。因此，安全的做法是**直接禁止slice之间的比较操作**。
标准库提供了高度优化的**bytes.Equal函数**来判断两个[]byte是否相等。

对于其他类型的slice，我们**必须自己展开每个元素进行比较**，运行的时间并不比支持==操作的数组或字符串更多。

slice唯一合法的比较操作是和nil比较，如：

```go
if summer == nil { /* ... */ }
```

```go
func equal(x, y []string) bool {
    if len(x) != len(y) {    // Go语言排错式的风格
        return false
    }
    for i := range x {
        if x[i] != y[i] {   // Go语言排错式的风格
            return false
        }
    }
    return true            // Go语言排错式的风格，正常执行的语句不被if缩进
}
```



测试一个slice是否是空的，**使用len(s) == 0来判断**，**而不应该用s == nil来判断**。

**nil值的slice & 长度和容量为0的slice**：**一个nil值的slice的行为和其它任意0长度的slice一样，所有的Go语言函数应该以相同的方式对待nil值的slice和0长度的slice**(除了文档已经明确说明的地方): 如reverse(nil)也是安全的。除了文档已经明确说明的地方

```go
// slice（引用类型）的零值为nil：**一个nil值的slice并没有引用任何底层数组**。一个nil值的slice的长度和容量都是0。
// 非nil值但长度和容量也是0的slice：如**[]int{}或make([]int, 3)[3:]**。
// 与任意类型的nil值一样，我们可以用[]int(nil)类型转换表达式来生成一个对应类型slice的nil值。
var s []int    // len(s) == 0, s == nil
s = nil        // len(s) == 0, s == nil
s = []int(nil) // len(s) == 0, s == nil
s = []int{}    // len(s) == 0, s != nil
```



内置的make函数：创建一个指定元素类型、长度和容量的slice。

容量部分可以省略，在这种情况下，容量将等于长度。

在底层，make创建了一个匿名的数组变量，然后返回一个slice。只有通过返回的slice才能引用底层匿名的数组变量。

```go
// 在第一种语句中，slice是整个数组的view。
make([]T, len)
// 在第二个语句中，slice只引用了底层数组的前len个元素，但是容量将包含整个的数组。额外的元素是留给未来的增长用的。
make([]T, len, cap) // same as make([]T, cap)[:len]
```



**二维切片:**切片的切片

```go
type Transform [3][3]float64  // 一个 3x3 的数组，其实是包含多个数组的一个数组。
type LinesOfText [][]byte     // 包含多个字节切片的一个切片。
```

有时必须分配一个二维数组，如在处理像素的扫描行时。

独立地分配每一个切片，一次一行：

```go
// 分配顶层切片。
picture := make([][]uint8, YSize) // 每 y 个单元一行。
// 遍历行，为每一行都分配切片
for i := range picture {
	picture[i] = make([]uint8, XSize)
}

```

一次分配，对行进行切片：

```go
// 分配顶层切片，和前面一样。
picture := make([][]uint8, YSize) // 每 y 个单元一行。
// 分配一个大的切片来保存所有像素
pixels := make([]uint8, XSize*YSize) // 拥有类型 []uint8，尽管图片是 [][]uint8.
// 遍历行，从剩余像素切片的前面切出每行来。
for i := range picture {
	picture[i], pixels = pixels[:XSize], pixels[XSize:]
}
```

若切片会增长或收缩，就应该通过独立分配来避免覆盖下一行；若不会，用单次分配来构造对象会更加高效。 



## 四、append

### 源码解读

**核心功能**：append是内置函数，**用于向切片末尾追加元素。**

**容量处理**：如果切片有足够容量，直接重新切片；如果容量不足，会分配新的底层数组。

**返回值：**返回更新后的切片，必须存储返回值，通常覆盖原变量。

```go
// src/builtin/builtin.go
// The append built-in function appends elements to the end of a slice. If
// it has sufficient capacity, **the destination is resliced to accommodate the
// new elements.** If it does not, a new underlying array will be allocated.
// Append returns the updated slice. It is therefore necessary to store the
// result of append, often in the variable holding the slice itself:
//
//	slice = append(slice, elem1, elem2)
//	slice = append(slice, anotherSlice...)
//
// As a special case, it is legal to append a string to a byte slice, like this:
//
//	slice = append([]byte("hello "), "world"...)
func append(slice []Type, elems ...Type) []Type
```



**工作原理：**

1. 容量充足时：&slice[0]指针地址相同，说明没有重新分配数组
```go
// 创建有足够容量的切片
slice := make([]int, 0, 5)  // len=0, cap=5, **&slice[0]: ptr=0xc000018180**

// 追加元素
slice = append(slice, 1, 2, 3)

fmt.Printf("After:  len=%d, cap=%d, ptr=%p\n", len(slice), cap(slice), &slice[0])  // len=3, cap=5, **&slice[0]: ptr=0xc000018180**
```

1. **容量不足时：**&slice[0]指针地址不同，说明重新分配了数组
```go
// 创建容量不足的切片
slice := make([]int, 0, 2)  //  len=0, cap=2, **&slice[0]: ptr=0xc000018180**

// 追加超过容量的元素
slice = append(slice, 1, 2, 3, 4, 5)  // After:  len=5, cap=6, **&slice[0]: ptr=0xc00001a0c0**
```

1. **陷阱：必须接收返回值**
**由于通常并不知道某次append调用是否重新分配了内存**，不能确认新的slice和原始的slice是否引用的是相同的底层数组空间，不能确认在原先的slice上的操作是否会影响到新的slice。

**因此，通常是将append返回的结果直接赋值给输入的slice变量。**实际上对应任何可能导致长度、容量或底层数组变化的操作都是必要的。

```go
slice := []int{1, 2, 3}
    
// 错误：没有接收返回值
// append(slice, 4)  // 编译警告：结果被丢弃

// 正确：接收返回值
slice = append(slice, 4)
fmt.Println(slice)  // [1 2 3 4]
```



要正确地使用slice，需要记住尽管底层数组的元素是间接访问的，但是**slice对应结构体本身的指针、长度和容量部分是直接访问的**。要更新这些信息需要像上面例子那样一个**显式的赋值操作**。

从这个角度看，slice并不是一个纯粹的引用类型，它**实际上是一个类似下面结构体的聚合类型**：

```go
type slice struct {
	ptr    unsafe.Pointer  *// 指向底层数组*
	len    int            *// 长度*
	cap    int            *// 容量*
}
```



append函数对于理解slice底层是如何工作的非常重要，这里的appendInt专门用于处理[]int类型的slice。

```go
// gopl.io/ch4/append
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1 
	if zlen <= cap(x) { // slice容量充足时
		// There is room to grow.  Extend the slice.
		z = x[:zlen]
	} else { // slice容量不足时
		zcap := zlen
		if zcap < 2*len(x) { // 容量double扩充，避免了多次内存分配， 分摊线性复杂性，确保了添加单个元素的平均时间是一个常数时间。
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap) // 新分配内存
		copy(z, x) // copy内置函数，比通过循环复制元素更方便，返回成功复制的元素的个数
	}
	z[len(x)] = y //  最后将新添加的y元素复制到新扩展的空间z[len(x)]，并返回slice。
	return z
}
```



<!-- 列布局开始 -->

![](/images/16424637-29b5-8082-a2c5-c95383fbf49c/image_16624637-29b5-80fb-b24a-c8cf1feb4605.jpg)


---

![](/images/16424637-29b5-8082-a2c5-c95383fbf49c/image_25124637-29b5-800b-878a-c2f78247c821.jpg)

<!-- 列布局结束 -->

```go
var x, y []int
for i := 0; i < 10; i++ {
	y = appendInt(x, i)
	fmt.Printf("i=%d  cap=%d\t y=%v\n", i, cap(y), y)
	x = y
}

/*
//!+output
i=0  cap=1       y=[0]
i=1  cap=2       y=[0 1]
i=2  cap=4       y=[0 1 2]
i=3  cap=4       y=[0 1 2 3] // 容量充足：x包含了[0 1 2]三个元素，但是容量是4。y和x引用着相同的底层数组。
i=4  cap=8       y=[0 1 2 3 4]  // 容量不足：分配一个容量为8的底层数组，将x的4个元素[0 1 2 3]复制到新空间的开头，然后添加新的元素i=4。**y和x是对应不同底层数组的view。**
i=5  cap=8       y=[0 1 2 3 4 5]
i=6  cap=8       y=[0 1 2 3 4 5 6]
i=7  cap=8       y=[0 1 2 3 4 5 6 7]
i=8  cap=16      y=[0 1 2 3 4 5 6 7 8]
i=9  cap=16      y=[0 1 2 3 4 5 6 7 8 9]
//!-output
*/
```



### 性能优化

1. **预分配容量**
```go
// 预分配容量，避免频繁重新分配
slice := make([]int, 0, 1000)
    
fmt.Printf("Initial: len=%d, cap=%d\n", len(slice), cap(slice))

// 追加1000个元素
for i := 0; i < 1000; i++ {
	slice = append(slice, i)
}

fmt.Printf("Final: len=%d, cap=%d\n", len(slice), cap(slice))  // 输出：Final: len=1000, cap=1000，没有重新分配数组
```

**容量增长策略：**通常是翻倍增长，但具体策略可能因Go版本而异。

```go
slice := make([]int, 0, 1)
    
fmt.Printf("Initial: len=%d, cap=%d\n", len(slice), cap(slice))

// 观察容量增长
for i := 0; i < 20; i++ {
	slice = append(slice, i)
	fmt.Printf("After %d: len=%d, cap=%d\n", i+1, len(slice), cap(slice)). // 1 -> 2 -> 4 -> 8 -> 16 -> 32
}
```

**2. 批量操作，避免循环中的append**

```go
slice := []int{1, 2, 3}
    
// 批量追加，减少函数调用次数
newElements := []int{4, 5, 6, 7, 8}
slice = append(slice, newElements...)

// 而不是逐个追加
// for _, elem := range newElements {
//     slice = append(slice, elem)
// }


// 不好的做法：在循环中频繁append
var slice []int
for i := 0; i < 1000; i++ {
    slice = append(slice, i)  // 可能多次重新分配
}
    
// 好的做法：预分配容量
slice = make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i)  // 不会重新分配
}
```





### example

旋转slice、反转slice或在slice原有内存空间修改元素。

nonempty：给定一个字符串列表，下面的nonempty函数将在原有slice内存空间之上返回不包含空字符串的列表：

```go
// Nonempty is an example of an in-place slice algorithm.
// nonempty returns a slice holding only the non-empty strings.
// The underlying array is modified during the call.
// 在原有string slice内存空间之上返回不包含空字符串的列表
// 比较微妙的地方是，输入的slice和输出的slice共享一个底层数组。
// 这可以避免分配另一个数组，不过原来的数据将可能会被覆盖
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}

func main() {
	data := []string{"one", "", "three"}
	fmt.Printf("%q\n", nonempty(data)) // `["one" "three"]`
	// 这可以避免分配另一个数组，不过原来的数据将可能会被覆盖
	fmt.Printf("%q\n", data) // `["one" "three" "three"]`
	// 因此我们通常会这样使用nonempty函数，和append函数类似
	data = nonempty(data)
	fmt.Printf("%q\n", data) // `["one" "three"]`
}

// nonempty函数也可以使用append函数实现
func nonempty2(strings []string) []string {
	out := strings[:0] // zero-length slice of original
	for _, s := range strings {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
```



stack：一个slice可以用来模拟一个stack

```go
func main() {
	// 无论如何实现，以这种方式重用一个slice一般都要求最多为每个输入值产生一个输出值
	// 事实上很多这类算法都是用来过滤或合并序列中相邻的元素。这种slice用法是比较复杂的技巧，虽然使用到了slice的一些技巧，但是对于某些场合是比较清晰和有效的。
	// 一个slice可以用来模拟一个stack。最初给定的空slice对应一个空的stack，然后可以使用append函数将新的值压入stack：
	var stack = make([]int, 0)
	v := 1
	stack = append(stack, v) // push v
	// stack的顶部位置对应slice的最后一个元素：
	top := stack[len(stack)-1] // top of stack
	// 通过收缩stack可以弹出栈顶的元素
	stack = stack[:len(stack)-1] // pop
	fmt.Println(top)

	s := []int{5, 6, 7, 8, 9}
	fmt.Println(remove(s, 2)) // [5 6 8 9]
	s = []int{5, 6, 7, 8, 9}
	fmt.Println(remove2(s, 2)) // [5 6 9 8]
}

func remove(slice []int, i int) []int {
	// // 要删除slice中间的某个元素并保存原有的元素顺序，可以通过内置的copy函数将后面的子slice向前依次移动一位完成：
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]

}

// 如果删除元素后不用保持原来顺序的话，我们可以简单的用最后一个元素覆盖被删除的元素：
func remove2(slice []int, i int) []int {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
```



