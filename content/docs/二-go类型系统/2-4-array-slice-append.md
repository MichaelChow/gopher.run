---
title: "2.4 array、slice、append"
date: 2024-12-22T04:12:00Z
draft: false
weight: 2004
---

# 2.4 array、slice、append

# array

- 数组是具有固定长度且拥有零个或多个相同数据类型元素的序列；由于固定长度，Go中很少直接使用array，更多使用动态长度的slice；slice底层由array实现；
- 数组在Go和C中的主要区别：
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
    - 数组的大小是其类型的一部分。类型 `[10]int` 和 `[20]int` 是不同的类型。
- **数组的长度是数组类型的一个组成部分**，[3]int、[4]int是两种不同的数组类型，长度需要在编译阶段确定，所以必须是常量表达式；
    - 在数组字面值中，如果在数组的**长度位置出现的是“...”省略号**，则表示数组的长度是根据初始化值的个数来计算
        ```go
        q := [...]int{1, 2, 3}
        r := [...]int{99: -1} // 前99个元素初始化为零值，最后一个初始化为-1
        ```
    - 数组元素可通过索引下标来访问，每个元素都被默认初始化为元素类型对应的零值，或使用字面量初始化；
        ```go
        var r [3]int = [3]int{1, 2}
        ```
- 如果一个数组的元素类型是可以相互比较的，数组类型相同（包括数组长度一致），则数组是可以相互比较==的，只有当两个数组的所有元素都是相等的时候数组才是相等的；
- 作为一个真实的例子，crypto/sha256包的Sum256函数对一个任意的byte slice类型的数据生成一个对应的消息摘要。消息摘要有256bit大小，因此对应[32]byte数组类型。如果两个消息摘要是相同的，那么可以认为两个消息本身也是相同（译注：理论上有HASH码碰撞的情况，但是实际应用可以基本忽略）；如果消息摘要不同，那么消息本身必然也是不同的。下面的例子用SHA256算法分别生成“x”和“X”两个信息的摘要：
- crypto/sha256： 两个消息虽然只有一个字符的差异，但是生成的消息摘要则几乎有一半的bit位是不相同的
    ```go
    import "crypto/sha256"
    func main() {
        c1 := sha256.Sum256([]byte("x"))  // 2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881, [32]uint8
        c2 := sha256.Sum256([]byte("X"))  // 4b68ab3847feda7d6c62c1fbcbeebfa35eab7351ed5e78f4ddadea5df64b8015
    }
    ```
- 函数参数传递的机制：当调用一个函数时，每个传入的参数都会创建一个副本，然后赋值给对应的函数变量**（**所以函数参数变量接收的是一个副本，并不是原始的参数）。当传递大的数组时将变得很低效，并且在函数内部对数组的任何修改都仅影响副本，而不是原始数组。这种情况下。Go把数组和其他的类型都看成值传递（**不同于其他语言中数组都是隐式的使用引用传递**）。
- **Go中可以显示的传递一个数组的指针给函数，这样在函数内部对数组的任何修改都会反应到原始数组上，这样很高效**；但由于数组的长度不可变特性，无法为数组添加或者删除元素，除了在SHA256、矩阵变换这类固定长度的情况下，很少直接使用数组；
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
# **slice**

- slice（切片,/slaɪs/）`[]T`，**表示一个拥有相同类型元素的可变长度的序列**，**看上去像没有长度的数组类型，**通过**对数组进行封装**，更通用、强大，可以用来访问**底层引用的数组**的元素（引用类型）；
- slice的三个属性：
    - pointer: 指向第一个可以从slice中访问的元素（不一定是数组的第一个元素）
    - len：数组的元素个数，不能超过cap
    - cap：slice的起始元素到底层数组的最后一个元素间元素的个数
- 一个底层数组可以对应多个slice，这些slice可以引用数组的任何位置，彼此之间的元素还可以重叠；
    ```go
    months := [...]string{1: "January", 2: "February", 3: "March", 4: "April", 5: "May", 6: "June", 7: "July", 8: "Augest", 9: "September", 10: "October", 11: "November", 12: "December"}  // **[13]string 模拟slice引用的底层array**
    Q2 := months[4:7]    //  **[]string slice,poiner=, len=3, cap=9**Q2[:10] // // panic: runtime error: slice bounds out of range [:10] with capacity 9
    summer := months[6:9]  // **[]string slice, poiner= , len=3, cap=7**
    ```
    ![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/22974f28-ad48-43ce-b343-60527b2a2b8a/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466YLFT3IBK%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005845Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIQCWDNEhlR8AZkN5sj7vGhPL4l6TpKQr7waJxT5%2FDntZeQIgKLs5pAe%2FnXnEKFNNAnD%2Fprt1I%2BPHFiytKf9zvFyYPA4qiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDAoW8sZdR4e7hAxfkCrcA6au%2FA8zQUSe88lEnmBOvrGVgVKjlgTPzPLV2iNJDEEPYRV6q44UVpdNCknQzN6ma8f7Gwfg%2FWYTXV5436TCKUQBc8MpKD%2B9plFd72djLFDspHLqlPloQpdI5Vz90THlB7TqRY4J2O5j1a4G7U0y61iokXgHLQZxHdcyPXImlM53aBYZ21M3y%2Bpku7CKvx0rhV0%2FlvHTlojWESp4tTjSC6ZcVYFFMH9NkJaibocW0PzZO6W6j9L0SsA%2BXbgE2AOI%2F%2ByT%2Bw1nh%2BbVQo%2BV0%2B4AMshBvqWU0FSFR8UGlrm5bS3e7mAdsg9CDtm6Jol9i34wisL5b8JIVwV9eAftOH%2BCTiO%2F15WbIxirEuNlImAhZJ2jHOtjwXxC0oK8nh7isnPYRlwcu3lou4Ub8SSmL3sq55NVB0sjJE%2B2mE2Fj3cLsHN7qVGiVFaED7AU20kgvCWHqSdHtZHJw5pyp5%2B4RXpNoNnOqyOUQfIYtslqx5m5cAjIxIReKafhawEm%2Fs1dYnYQd8ImiNg78VUjGfZ4fBLeADlW8WtbCULlya5uxyiXI55pcVmJKtyHSZPwTfQa1wjvPkV5C5pGk5bEmC9LZ%2FiD%2Bb7GhU%2BtplOMsbaId4CHqlh87iD7RDsIvItBoNSXMIe668MGOqUByYGaZ2pmgBPHYG%2FQIcljsn7YCHNqX9GI%2FzlvHfGwm6KNtYFwzrnol%2FjqGivisvbEhJ3S9rKlFhuR4I%2BUIg5C5tgUq%2FYuIEKZQlOWzbEVUjRKHhtQOmz3Wxi76kuRO4DsJVCgfpSMhbsZzqDC6NLAT7Te5Sw9JbguGwk5FNxPOV2CpBhBjbOS6rhl14K26vFUXiGzorgZO%2FtslQC8MaYrycoAsecO&X-Amz-Signature=fd9cfe5dedb7ca00831d3302b0d9b60f045b86bd713c43ddac63f3db45b2a0b6&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)
- reverse
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
- **slice不直接支持比较运算符：**标准库提供了高度优化的**bytes.Equal函数**来判断两个[]byte是否相等。但对于其他类型的slice，我们必须自己展开每个元素进行比较，运行的时间并不比支持==操作的数组或字符串更多。
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
    - 原因1：一个slice的元素是间接引用的，当slice声明为[]interface{}时，slice的元素甚至可以是自身。虽然有很多办法处理这种情形，但是没有一个是简单有效的。
    - 原因2：因为slice的元素是**间接引用的**，**slice本身的值（不是元素的值）****在不同的时刻可能包含不同的元素**，因为底层数组的元素可能会被修改。而如Go语言中map的key只做简单的浅拷贝，它要求key在整个生命周期内保持不变性（译注：如slice扩容，就会导致其本身的值/地址变化）。而用深度相等判断的话，显然在map的key这种场合不合适。对于像指针或chan之类的引用类型，==相等测试可以判断两个是否是引用相同的对象。一个针对slice的浅相等测试的==操作符可能是有一定用处的，也能临时解决map类型的key问题，**但是slice和数组不同的相等测试行为会让人困惑**。因此，安全的做法是**直接禁止slice之间的比较操作**。
    - slice唯一合法的比较操作是和nil比较，如：
        ```go
        if summer == nil { /* ... */ }
        ```
- 测试一个slice是否是空的，使用len(s) == 0来判断，而不应该用s == nil来判断。
- 一个nil值的slice的行为和其它任意0长度的slice一样,所有的Go语言函数应该以相同的方式对待nil值的slice和0长度的slice(除了文档已经明确说明的地方): 如reverse(nil)也是安全的。除了文档已经明确说明的地方
    ```go
    // 一个零值的slice等于nil。一个nil值的slice并没有底层数组。一个nil值的slice的长度和容量都是0
    // 但也有非nil值的slice的长度和容量也是0的，如[]int{}或make([]int, 3)[3:]。与任意类型的nil值一样，我们可以用[]int(nil)类型转换表达式来生成一个对应类型slice的nil值。
    var s []int    // len(s) == 0, s == nil
    s = nil        // len(s) == 0, s == nil
    s = []int(nil) // len(s) == 0, s == nil
    s = []int{}    // len(s) == 0, s != nil
    ```
- 内置的make函数创建一个指定元素类型、长度和容量的slice。
    - 容量部分可以省略，在这种情况下，容量将等于长度。
    - 在底层，make创建了一个匿名的数组变量，然后返回一个slice；
    - 只有通过返回的slice才能引用底层匿名的数组变量。
    ```go
    // 在第一种语句中，slice是整个数组的view。
    make([]T, len)
    // 在第二个语句中，slice只引用了底层数组的前len个元素，但是容量将包含整个的数组。额外的元素是留给未来的增长用的。
    make([]T, len, cap) // same as make([]T, cap)[:len]
    ```
### **切片**

因此，`Read` 函数可接受一个切片实参 而非一个指针和一个计数；切片的长度决定了可读取数据的上限。以下为 `os` 包中 `File` 类型的 `Read` 方法签名:

```go
func (file *File) Read(buf []byte) (n int, err error)

```

该方法返回读取的字节数和一个错误值（若有的话）。若要从更大的缓冲区 `b` 中读取前32个字节，只需对其进行**切片**即可。

```go
	n, err := f.Read(buf[0:32])

```

这种切片的方法常用且高效。若不谈效率，以下片段同样能读取该缓冲区的前32个字节。

```go
	var n int
	var err error
	for i := 0; i < 32; i++ {
		nbytes, e := f.Read(buf[i:i+1])  // 读取一个字节
		if nbytes == 0 || e != nil {
			err = e
			break
		}
		n += nbytes
	}

```

只要切片不超出底层数组的限制，它的长度就是可变的，只需将它赋予其自身的切片即可。 切片的**容量**可通过内建函数 `cap` 获得，它将给出该切片可取得的最大长度。 



以下是将数据追加到切片的函数。若数据超出其容量，则会重新分配该切片。返回值即为所得的切片。 该函数中所使用的 `len` 和 `cap` 在应用于 `nil` 切片时是合法的，它会返回0.

```go
func Append(slice, data[]byte) []byte {
	l := len(slice)
	if l + len(data) > cap(slice) {  // 重新分配
		// 为了后面的增长，需分配两份。
		newSlice := make([]byte, (l+len(data))*2)
		// copy 函数是预声明的，且可用于任何切片类型。
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:l+len(data)]
	for i, c := range data {
		slice[l+i] = c
	}
	return slice
}

```

最终我们必须返回切片，因为尽管 `Append` 可修改 `slice` 的元素，但切片自身（其运行时数据结构包含指针、长度和容量）是通过值传递的。

向切片追加东西的想法非常有用，因此有专门的内建函数 `append`。 要理解该函数的设计，我们还需要一些额外的信息，我们将稍后再介绍它。

### **二维切片**

Go的数组和切片都是一维的。要创建等价的二维数组或切片，就必须定义一个数组的数组， 或切片的切片，就像这样：

```go
type Transform [3][3]float64  // 一个 3x3 的数组，其实是包含多个数组的一个数组。
type LinesOfText [][]byte     // 包含多个字节切片的一个切片。

```

由于切片长度是可变的，因此其内部可能拥有多个不同长度的切片。在我们的 `LinesOfText` 例子中，这是种常见的情况：每行都有其自己的长度。

```go
text := LinesOfText{
	[]byte("Now is the time"),
	[]byte("for all good gophers"),
	[]byte("to bring some fun to the party."),
}

```

有时必须分配一个二维数组，例如在处理像素的扫描行时，这种情况就会发生。 我们有两种方式来达到这个目的。一种就是独立地分配每一个切片；而另一种就是只分配一个数组， 将各个切片都指向它。采用哪种方式取决于你的应用。若切片会增长或收缩， 就应该通过独立分配来避免覆盖下一行；若不会，用单次分配来构造对象会更加高效。 以下是这两种方法的大概代码，仅供参考。首先是一次一行的：

```go
// 分配顶层切片。
picture := make([][]uint8, YSize) // 每 y 个单元一行。
// 遍历行，为每一行都分配切片
for i := range picture {
	picture[i] = make([]uint8, XSize)
}

```

现在是一次分配，对行进行切片：

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



## **append**

append函数对于理解slice底层是如何工作的非常重要。

```go
// 由九个rune字符构成的slice
// 这个特殊的问题我们可以通过Go语言内置的[]rune("Hello, 世界")转换操作完成
var runes []rune
for _, r := range "Hello, 世界" {
    runes = append(runes, r)
}
fmt.Printf("%q\n", runes) // "['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']"
```

### append.go

```go
// Append illustrates the behavior of the built-in append function.
// See page 88.

package main

import "fmt"

func appendslice(x []int, y ...int) []int {
	var z []int
	zlen := len(x) + len(y)
	if zlen <= cap(x) {
		// There is room to expand the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.
		// Grow by doubling, for amortized linear complexity.
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	copy(z[len(x):], y)
	return z
}

// append函数对于理解slice底层是如何工作的非常重要，appendInt专门用于处理[]int类型的slice
// 函数将y添加到slice x中
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	// 必须先检测slice x底层数组是否有足够的容量来保存新添加的元素y
	// 如果有，直接在原有的底层数组x之上扩展slice的len(x)。因此输出的z直接共享输入的x相同的底层数组。
	if zlen <= cap(x) {
		// There is room to grow.  Extend the slice.
		z = x[:zlen]
	} else {
		// There is insufficient space.  Allocate a new array.
		// Grow by doubling, for amortized linear complexity.
		// 如果没有足够的增长空间：
		// 为了提高内存使用效率，新分配的数组一般略大于保存x和y所需要的最低大小。
		// 通过在每次扩展数组时直接将长度翻倍，从而避免了多次内存分配，也确保了添加单个元素的平均时间是一个常数时间。
		zcap := zlen
		if zcap < 2*len(x) {
			zcap = 2 * len(x)
		}
		// 分配一个新的数组z，将长度设置为zlen，容量设置为zcap。因此输出的z与输入的x引用的将是不同的底层数组。
		z = make([]int, zlen, zcap)
		// 虽然通过循环复制元素更直接，不过内置的copy函数可以方便地将一个slice复制另一个相同类型的slice。
		// copy函数的第一个参数是要复制的目标slice，第二个参数是源slice，目标和源的位置顺序和dst = src赋值语句是一致的。两个slice可以共享同一个底层数组，甚至有重叠也没有问题。copy函数将返回成功复制的元素的个数（我们这里没有用到），等于两个slice中较小的长度，所以我们不用担心覆盖会超出目标slice的范围。
		copy(z, x) // a built-in function; see text
	}
	//  最后将新添加的y元素复制到新扩展的空间z[len(x)]，并返回slice。
	z[len(x)] = y
	return z
}

func main() {
	var x, y []int
	for i := 0; i < 10; i++ {
		y = appendInt(x, i)
		fmt.Printf("i=%d  cap=%d\t y=%v\n", i, cap(y), y)
		x = y
	}
}

/*
//!+output
i=0  cap=1       y=[0]
i=1  cap=2       y=[0 1]
i=2  cap=4       y=[0 1 2]
i=3  cap=4       y=[0 1 2 3]
i=4  cap=8       y=[0 1 2 3 4]
i=5  cap=8       y=[0 1 2 3 4 5]
i=6  cap=8       y=[0 1 2 3 4 5 6]
i=7  cap=8       y=[0 1 2 3 4 5 6 7]
i=8  cap=16      y=[0 1 2 3 4 5 6 7 8]
i=9  cap=16      y=[0 1 2 3 4 5 6 7 8 9]
//!-output
*/
```



在i=3次的迭代，x包含了[0 1 2]三个元素，但是容量是4。因此可以简单将新的元素添加到末尾，不需要新的内存分配。新的y的长度和容量都是4，并且y和x引用着相同的底层数组，如图4.2所示。

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/934cf07d-9dc0-4fa1-bd2f-91db21925046/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466XPYEOVWR%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005840Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIAtKGmMBwUITjCt8ep2A90CH3e%2B6l8vWELtrjmlEH50eAiEAkVz2hmyeSFRbUxsW8vnh3tD6InrXjSfVPwufrzBAzWQqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDLBn07Lwkx24iWdJbCrcA4i0ERwfX3OkhHYRzF7JoczNIxtAg4QvdGRy5Cw8I05CKLbcXNeHmy6l89nVBeWLWuqe0aDcSIRzey6XPLlefDCo5IpHL0ni5KrmH3s%2FVoMCPqtUxgC6dv5hx1%2FukrQwj9P2NPEnEb012qQ9QSCdsNAyCpe%2BOC%2FtFDy7xm2AUYSxUOdyNkHJyIQkU6Yanpb%2F5uX4vqS%2FaQXlUp1gYrdJFPjePj5POyBnifgvICeNBlKCyKHLggKxyB4nPjoktGXZHsSQkbr8DkjZeQ0lOnsrd2qFKm81fYVt%2BsHpMch3ZarZWRojfP2UXTRziunLscr8L77TAzlmoJf8JMvZOo1jxqMx4t69ymBknoMaPv3ucnlrxyrvo1XS3nEiflW07AM9nPM%2FM38V4vm2qvnMYMjBCMA8%2B8m3YODc9x7H3VwkyUYG%2BXMpqAMWgXbSLOSmDvjkn2FXJRIf0QHilZv05W4XlDp%2F3XcENNQ5Wz80q7X%2FIhcWc55X9bboJs4qzzYvUcsxQ%2FjVXv1qCMo%2FWBknvcnTlHkouqBKz98yFnmpLglY6GFb6BborWKM1%2FQRhIflGTRrqz8PrYeBS%2BqZiW0Q77P1w6Sw4ksoX165pz0vOeqP8D66AnZSTy6zr8dDt4ITMK%2B668MGOqUBuOKrba0Z3wmIotOTWVfqSz4P5IOLSQuc%2FUbYAefEg1MK%2FJwYmqFFfKUV6TO02cVx5BGaWBmai1s%2F55D2lshT2n3PTb%2BZ2E9NQUxb91rnyfz5IQHTZWIANbXaipBlLA03ZF9QmVnmFj8or4dr3xQevL%2FKCxzOSrSnTNDMcDVDOft3iC%2Boi1u0RwTGcKEaBWnteDye4dCNQdWeMCxogfyvSSzua3FL&X-Amz-Signature=4aa0c8da81268a17718c4784dc72068acdcf3f449639a39e8af0f2b186735249&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

在i=4次的迭代，已经没有新的空余空间了，因此appendInt函数分配一个容量为8的底层数组，将x的4个元素[0 1 2 3]复制到新空间的开头，然后添加新的元素i=4。新的y的长度是5，容量是8；后面有3个空闲的位置，往后3次迭代都不需要分配新的空间。当前迭代中，y和x是对应不同底层数组的view。如图4.3所示。

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/766620db-3307-447d-9e24-cf700f8aebf5/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466XPYEOVWR%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005840Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJHMEUCIAtKGmMBwUITjCt8ep2A90CH3e%2B6l8vWELtrjmlEH50eAiEAkVz2hmyeSFRbUxsW8vnh3tD6InrXjSfVPwufrzBAzWQqiAQImf%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgw2Mzc0MjMxODM4MDUiDLBn07Lwkx24iWdJbCrcA4i0ERwfX3OkhHYRzF7JoczNIxtAg4QvdGRy5Cw8I05CKLbcXNeHmy6l89nVBeWLWuqe0aDcSIRzey6XPLlefDCo5IpHL0ni5KrmH3s%2FVoMCPqtUxgC6dv5hx1%2FukrQwj9P2NPEnEb012qQ9QSCdsNAyCpe%2BOC%2FtFDy7xm2AUYSxUOdyNkHJyIQkU6Yanpb%2F5uX4vqS%2FaQXlUp1gYrdJFPjePj5POyBnifgvICeNBlKCyKHLggKxyB4nPjoktGXZHsSQkbr8DkjZeQ0lOnsrd2qFKm81fYVt%2BsHpMch3ZarZWRojfP2UXTRziunLscr8L77TAzlmoJf8JMvZOo1jxqMx4t69ymBknoMaPv3ucnlrxyrvo1XS3nEiflW07AM9nPM%2FM38V4vm2qvnMYMjBCMA8%2B8m3YODc9x7H3VwkyUYG%2BXMpqAMWgXbSLOSmDvjkn2FXJRIf0QHilZv05W4XlDp%2F3XcENNQ5Wz80q7X%2FIhcWc55X9bboJs4qzzYvUcsxQ%2FjVXv1qCMo%2FWBknvcnTlHkouqBKz98yFnmpLglY6GFb6BborWKM1%2FQRhIflGTRrqz8PrYeBS%2BqZiW0Q77P1w6Sw4ksoX165pz0vOeqP8D66AnZSTy6zr8dDt4ITMK%2B668MGOqUBuOKrba0Z3wmIotOTWVfqSz4P5IOLSQuc%2FUbYAefEg1MK%2FJwYmqFFfKUV6TO02cVx5BGaWBmai1s%2F55D2lshT2n3PTb%2BZ2E9NQUxb91rnyfz5IQHTZWIANbXaipBlLA03ZF9QmVnmFj8or4dr3xQevL%2FKCxzOSrSnTNDMcDVDOft3iC%2Boi1u0RwTGcKEaBWnteDye4dCNQdWeMCxogfyvSSzua3FL&X-Amz-Signature=3ce52f41dd3da680885e3d45aa74823708df47b9de4e7d407329435d8314e2ab&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)





内置的append函数可能使用比appendInt更复杂的内存扩展策略。因此，通常我们并不知道append调用是否导致了内存的重新分配，因此我们也不能确认新的slice和原始的slice是否引用的是相同的底层数组空间。同样，我们不能确认在原先的slice上的操作是否会影响到新的slice。因此，通常是将append返回的结果直接赋值给输入的slice变量：

```go
runes = append(runes, r)
```



更新slice变量不仅对调用append函数是必要的，实际上对应任何可能导致长度、容量或底层数组变化的操作都是必要的。要正确地使用slice，需要记住尽管底层数组的元素是间接访问的，但是slice对应结构体本身的指针、长度和容量部分是直接访问的。要更新这些信息需要像上面例子那样一个显式的赋值操作。从这个角度看，slice并不是一个纯粹的引用类型，它**实际上是一个类似下面结构体的聚合类型**：

```go
type IntSlice struct {
    ptr      *int
    len, cap int
}
```

我们的appendInt函数每次只能向slice追加一个元素，但是内建函数append函数则可以追加多个元素，甚至追加一个slice。

```go
var x []int
x = append(x, 1)
x = append(x, 2, 3)
x = append(x, 4, 5, 6)
x = append(x, x...) // append the slice x
fmt.Println(x)      // "[1 2 3 4 5 6 1 2 3 4 5 6]"

```

通过下面的小修改，我们可以达到append函数类似的功能。其中在appendInt函数参数中的最后的“...”省略号表示接收变长的参数为slice。（我们将在5.7节详细解释这个特性）

```go
func appendInt(x []int, y ...int) []int {
    var z []int
    zlen := len(x) + len(y)
    // ...expand z to at least zlen...
    copy(z[len(x):], y)
    return z
}

```

为了避免重复，和前面相同的代码并没有显示。



### **追加**

现在我们要对内建函数 `append` 的设计进行补充说明。`append` 函数的签名不同于前面我们自定义的 `Append` 函数。大致来说，它就像这样：

```go
func append(slice []T, 元素 ...T) []T
```

其中的

*T*

为任意给定类型的占位符。实际上，你无法在Go中编写一个类型

```go
T
```

由调用者决定的函数。这也就是为何

```go
append
```

为内建函数的原因：它需要编译器的支持。

`append` 会在切片末尾追加元素并返回结果。我们必须返回结果， 原因与我们手写的 `Append` 一样，即底层数组可能会被改变。以下简单的例子

```go
x := []int{1,2,3}
x = append(x, 4, 5, 6)
fmt.Println(x)

```

将打印 `[1 2 3 4 5 6]`。因此 `append` 有点像 `Printf` 那样，可接受任意数量的实参。

但如果我们要像 `Append` 那样将一个切片追加到另一个切片中呢？ 很简单：在调用的地方使用 `...`，就像我们在上面调用 `Output` 那样。以下代码片段的输出与上一个相同。

```go
x := []int{1,2,3}
y := []int{4,5,6}
x = append(x, y...)
fmt.Println(x)

```

如果没有 `...`，它就会由于类型错误而无法编译，因为 `y` 不是 `int` 类型的。



## **Slice内存技巧**

旋转slice、反转slice或在slice原有内存空间修改元素。



### nonempty.go

给定一个字符串列表，下面的nonempty函数将在原有slice内存空间之上返回不包含空字符串的列表：

```go
// Nonempty is an example of an in-place slice algorithm.
// See page 91.
package main

import "fmt"

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



### stack.go

一个slice可以用来模拟一个stack

```go
package main

import "fmt"

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





让我系统性地分析前面遇到的pathTitles切片共享bug：

## **Go语言切片共享Bug深度分析 🔍**

### **1. 问题表现**

- ❌ Excel输出中大量重复的"职人信息"
- ❌ 所有API的路径标题都变成了相同内容
- ❌ 实际上是最后处理的节点路径覆盖了所有之前的记录
**2. Bug根本原因分析**

**Go切片的底层机制goApply to slice.go**

```go
// 问题代码模式
currentPath := append(pathTitles, node.Title)  // 🚨 危险！
```

**切片结构：**

```go
type slice struct {
ptr    unsafe.Pointer  *// 指向底层数组*
len    int            *// 长度*
cap    int            *// 容量*
}
```

**append函数的行为**
1. **容量足够时** - 直接修改原底层数组

1. **容量不足时** - 分配新数组并复制
**3. 具体发生的问题**

![](https://prod-files-secure.s3.us-west-2.amazonaws.com/3bd3cf7e-0f8f-40af-acf7-9f45a802bdba/f6339cf6-7bc6-417f-9426-bf7a21e72652/image.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAZI2LB466QZZO75ND%2F20250719%2Fus-west-2%2Fs3%2Faws4_request&X-Amz-Date=20250719T005840Z&X-Amz-Expires=3600&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEID%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEaCXVzLXdlc3QtMiJIMEYCIQD80Pf0HyjU%2BqBefY7651%2F1ir88h1%2BERsQkJLYCmcBQ7AIhAOxqh7ISSt7kv4ePbmqwgqGq8ngRrDf9ntCndjdA0GG8KogECJn%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FwEQABoMNjM3NDIzMTgzODA1IgwtDzyfPaHmS90QyNcq3AOuQAd7%2BaVd7VJbJpstEm8o2vpU%2Bi7kWQRIihe9PGRl9XdJYHGmdQSxIyt4DBOGXWm9Ld7q4%2FkAxIXDlwUBZVJ5rm%2FKLc%2FyX5vC598vpbCZTWQAQK%2F89PDumAOI5RLa7cpKIzDdqSD%2B1jKoSTThv3QNrIlM83IblX4q2z0c8STNvEMcKHYjhHg8NbaSrx334u%2BYlNqHsI0hBM5mHLgCXZjatBzihtnDoRZv2kaCbx%2FoFSrPYdp%2FiXbvpGeM4AcuU%2B8T6mv1ZMtWiD%2Fbg07D2%2F7DsUbam78h%2BGkm6ilnMit3BvZzCIK0DzaRh1zplr8Pe32TwHefIukkgHPEgRsw8HAjPXBUFjNVkmJLllHuSoLI6wtKV4IUI6L%2F%2BUhk0EH%2BH%2FvQXn7WY80C3t1Yo%2Fx6pfzgz4R8Kdyfy5jriXbHiApEmhaWNUppjVRQhQzXnDQswBJlzIZVzBTe2g7g1Ey11oxjTzUAMb82AE6ZdWOWgQ4yE4pwGPk9XauPS%2FHHuu3NQyZ1Q4MmfXjBRW7F9eSA1gl85gTwuTVeeqfP4ls18VMJkuSD6O9l0k99X7MdH%2FNuTCE0026445sCUpOCz4uAJqz%2B%2BoLxKzMQ6tjOs%2FHfLaEOqWjekiDsDKUe5bicEDDMuuvDBjqkAWiuvxYSBYfSrwzfpz8Cbm%2BunZSKln3ME1aPJX11FhCjI%2FvC7sKEYKa%2FXDFrZJCawjRgsNYUjaMzCSlWBDxp7mNU8B1jMUfA2%2BMZsnwxdPsfSq13%2FodTNnVYdTtkLLvvviP0%2Fl7gSkOiATjguINT2AyZ5%2Bdi03ubgRgsdOBe6UrhbL%2FwV4x1kEuAvHRRj5eUNz7LlGX0TAgSso4PRk6tAmZIf66i&X-Amz-Signature=83186e6fe0a817e55f5f9b373a2b2e445e6147e6537637913fb398c9dbb364a3&X-Amz-SignedHeaders=host&x-amz-checksum-mode=ENABLED&x-id=GetObject)

**4. 涉及的代码位置**

**问题代码分布：**

```go
// ❌ dyls/dyls.go:60
currentPath := append(pathTitles, node.Title)

// ❌ dyds/dyds.go:97  
currentPath := append(pathTitles, node.Name)

// ❌ dyds/dyds.go:127
doc.PathTitles = append(currentPath, doc.APIName)

// ❌ mt/utils.go (多处)
currentPath := append(pathTitles, node.Title)
```

**5. Bug影响范围**

**6. 解决方案演进**

**第一阶段：发现问题goApply to slice.go*****// 用户反馈：Excel中大量"职人信息"重复// 定位到：切片共享导致的路径污染***

**第二阶段：直接修复goApply to slice.go**

```go
// ✅ 直接复制方案
currentPath := make([]string, len(pathTitles), len(pathTitles)+1)
copy(currentPath, pathTitles)  
currentPath = append(currentPath, node.Title)
```

**第三阶段：工具化重构goApply to slice.go**

```go
// ✅ 创建utils/slice.go
func AppendToSliceCopy(src []string, elements ...string) []string {
    dst := make([]string, len(src), len(src)+len(elements))
    copy(dst, src)
    return append(dst, elements...)
}

// ✅ 统一使用
```

**7. 为什么这个Bug很常见**

**Go新手常犯错误：**

1. **误解append行为** - 以为总是创建新切片
1. **忽略容量影响** - 不知道容量足够时会修改原数组
1. **递归场景复杂** - 多层递归让问题更难发现
1. **并发加剧问题** - 多goroutine时race condition
**类似的危险模式：goApply to slice.go*****// 用户反馈：Excel中大量"职人信息"重复// 定位到：切片共享导致的路径污染***

**8. 最佳实践总结**

**安全的切片操作：goApply to slice.go*****// ✅ 直接复制方案*****currentPath := make([]string, len(pathTitles), len(pathTitles)+1)copy(currentPath, pathTitles)  currentPath = append(currentPath, node.Title)**

### **代码审查要点：**

- 🔍 所有append调用都需要审查
- 🔍 递归函数中的切片传递
- 🔍 goroutine间的数据共享
- 🔍 函数参数是否会被修改
**9. 预防措施**

**工具函数统一：goApply to slice.go*****// ✅ 创建utils/slice.go*****func AppendToSliceCopy(*****src***** []string, *****elements***** ...string) []string {    dst := make([]string, len(src), len(src)+len(elements))    copy(dst, src)    return append(dst, elements...)}*****// ✅ 统一使用*****currentPath := utils.AppendToSliceCopy(pathTitles, node.Title)**

**编程习惯：**
1. **默认复制** - 除非确定安全，否则总是复制切片
1. **工具函数** - 使用经过验证的工具函数
1. **单元测试** - 特别测试边界情况和并发场景
1. **代码审查** - 重点关注切片操作

这个bug很好地说明了**Go语言"简单但不简单"的特点** - 看似简单的append操作背后有复杂的内存管理机制！



