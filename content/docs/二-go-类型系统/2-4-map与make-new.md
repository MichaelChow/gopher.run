---
title: "2.4 map与make、new"
date: 2025-04-03T08:00:00Z
draft: false
weight: 2004
---

# 2.4 map与make、new

## make & new

内建函数，仅限于分配并初始化一个 **slice、map、chan类型** 的对象，make返回Type，new返回*Type。

**内置函数 **`make`** 创建空 **`map`（译注：从功能和实现上说**，**`Go`** 的 **`map`** 类似于 **`Java`** 语言中的 **`HashMap`**，Python 语言中的 **`dict`，`Lua` 语言中的`table`，通常使用`hash`实现。）

```go
// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//
//   - Slice: The size specifies the length. The capacity of the slice is
//     equal to its length. A second integer argument may be provided to
//     specify a different capacity; it must be no smaller than the
//     length. For example, make([]int, 0, 10) allocates an underlying array
//     of size 10 and returns a slice of length 0 and capacity 10 that is
//     backed by this underlying array.
//   - Map: An empty map is allocated with enough space to hold the
//     specified number of elements. The size may be omitted, in which case
//     a small starting size is allocated.
//   - Channel: The channel's buffer is initialized with the specified
//     buffer capacity. If zero, or the size is omitted, the channel is
//     unbuffered.
func make(t Type, size ...IntegerType) **Type****// 返回值，值已经初始化。
//**make([]int, 10, 100) // 分配一个具有100个 int 的数组空间，接着创建一个长度为10， 容量为100并指向该数组中前10个元素的切片结构。****
// The new built-in function allocates memory. The first argument is a type,
// not a value, and the value returned is a pointer to a newly
// allocated zero value of that type.
func new(Type) ***Type****//********返回地址，值初始化为零值。实际使用较少。**new函数类似是一种**语法糖****（一种简洁的写法，编译器会自动转换）。**对于结构体来说，直接用字面量语法创建新变量的方法会更灵活。
```



`new` 和 `make` 之间的区别：

```go
var p *[]int = new([]int)       // 分配切片结构；*p == nil；**基本没用**
var v  []int = make([]int, 100) // 切片 v 现在引用了一个具有 100 个 int 元素的新数组

// 没必要的复杂：
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// 习惯用法：
v := make([]int, 100)
```



当然也可能有特殊情况：**如果两个类型都是空的（说类型的大小是0）**，例如`struct{}`和`[0]int`，**有可能有相同的地址（依赖具体的语言实现）**（译注：请谨慎使用大小为0的类型，因为如果类型的大小为0的话，可能导致Go语言的自动垃圾回收器有不同的行为，具体请查看`runtime.SetFinalizer`函数相关文档）。

```go
// new study
package main

import "fmt"

func main() {
	// 特殊情况：如果类型的大小都为0，地址可能会相同
	p, q := new([0]int), new([0]int) // 0x14000058030 [] 0x14000058038 []
	fmt.Println(p == q) // false，fmt.Println导致了两种不同的比较结果，可能为特定的编译器和优化设置下发生的罕见情况； true 0x14000058030 0x14000058038，为什么是true?
}
```



### **Newxxx() 构造函数/工厂函数**

```go
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	return &File{fd: fd, name: name}
}
```

**复合字面量：**

```go
a := [...]string   {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string      {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

```



由于new只是一个预定义的函数，**它并不是一个关键字**（Java中为关键字），因此我们可以将new名字重新定义为别的类型。如下面的例子：

```go
// /usr/local/go/src/errors/errors.go
// New returns an error that formats as the given text.
// Each call to New returns a distinct error value even if the text is identical.
func New(text string) error {
	return &errorString{text}
}
```

```go
// 由于new被定义为int类型的变量名，因此在delta函数内部是无法使用内置的new函数的。
func delta(old, new int) int { return new - old }
```



## map

哈希表（map）是一种巧妙并且实用的数据结构，它是一个**无序的key/value对的集合**，其所有的key都是不同的，然后通过给定的key可以在**常数时间复杂度内检索、更新或删除对应的value**。

在Go语言中，一个map[K]V**底层引用了一个哈希表（引用类型）**，其中K和V分别对应key和value。map中所有的key都有相同的类型，所有的value也有着相同的类型，但key和value之间可以是不同的数据类型。作为引用类型，与slice一样，将map传入函数，并更改了内容，此修改对调用者同样可见。

K对应的**key必须是支持==比较运算符的数据类型**（所以map可以通过测试key是否相等来判断是否已经存在）。虽然浮点数类型也是支持相等运算符比较的，但是将浮点数用做key类型则是一个坏的想法，正如第三章提到的，最坏的情况是可能出现的NaN和任何浮点数都不相等。

对于V对应的value数据类型则没有任何的限制。

map的零值是nil，**即没有引用任何hash表。**

```go
ages := make(map[string]int) // **make 创建一个空map[]（非nil）**，len(ages) = 0
ages := map[string]int{}    // **map字面量创建一个空map[]（非nil）**，len(ages) = 0
var ages map[string]int    // **初始化为零值nil**，len(ages) = 0

ages := map[string]int{   // 用map字面值的语法创建一个非空map
		"alice":   31,
		"charlie": 34,
}
```





**增删改查**：

**下标查找、delete、len和range都可以安全在nil值的map上**，它们的行为和一个空map类似；但向一个nil值的map存入元素将导致一个panic异常。

`**map**`**中不存在某个键时不用担心，首次读到新行时，等号右边的表达式**`**counts[line]**`**的值将被计算为其类型的零值，对于**`**int**`**即**`**0**`**。**

```go
var age map[string]int  // nil
ages["carol"] = 21 // panic: assignment to entry in nil map
ages["alice"]  // 通过key对应的下标访问元素；如果key不存在，那么将返回该map元素类型的零值0（V为int）
delete(ages, "alice")  // 删， 即便对应的键不在该映射中，此操作也是安全的。
```

Map(映射)：分配一个空映射，并有足够的空间来容纳指定数量的元素。size参数可以省略，此时会分配一个较小的初始大小。

map是随机顺序，`map`的迭代顺序**每次运行都会变化**（实测）。这种设计是有意为之的，因为**能防止开发的程序依赖特定遍历顺序**，而这是无法保证的。（译注：具体可以参见这里[https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order](https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order)）map的顺序取决于底层使用的hash函数，hash函数为了修复DOS拒绝服务攻击做了随机化处理。。



**map元素不是一个变量，不可以获取它的地址**。

**禁止对map元素取地址的原因**：**map可能随着元素数量的增长而重新分配更大的内存空间（已有的元素重新散列到新的存储位置），导致之前的地址失效**。

```go
_ := &ages["bob"]  // compile error: cannot take address of map element
```





**for range 迭代器遍历map：**

```go
// 遍历map中全部的key/value对：可使用range风格的for循环遍历key/value对（同slice遍历语法）
	for name, age := range ages {
		fmt.Printf("%s\t%d\n", name, age)
	}

	// map的迭代顺序是不确定的，并且不同的哈希函数实现可能导致不同的遍历顺序。
	// 在实践中，遍历的顺序是随机的，每一次遍历的顺序都不相同。这是故意的，每次都使用随机的遍历顺序可以强制要求程序不会依赖具体的哈希函数实现。

	// 如果一定要按顺序遍历key/value对，我们必须显式地对map的key进行排序，可以使用sort包的Strings函数对字符串slice进行排序：
	// 创建一个空slice，由于name的最终大小完全确定，容量直接设为map的长度
	names := make([]string, 0, len(ages))
	// 取出map中的所有key
	for name := range ages {
		names = append(names, name)
	}
	// 对key进行排序
	sort.Strings(names)
	// 遍历排序后的slice中的key，根据key取出map中的value
	for _, name := range names {
		fmt.Printf("%s\t%d\n", name, ages[name])
	}
	equal(map[string]int{"A": 0}, map[string]int{"B": 42})

}

// 和slice一样，map之间也不能进行相等比较；唯一的例外是和nil进行比较。
// 要判断两个map是否包含相同的key和value，我们必须通过一个循环实现：

func equal(x, y map[string]int) bool {
	if len(x) != len(y) {
		return false
	}
	for k, xv := range x {
		// 如果key不存在，那么将返回该map元素类型的零值。（这个规则很实用）
		// 但如果元素是一个数字，就需要区分一个已经存在元素的0和不存在元素返回的默认零值的0。
		// 所以我们需要使用一个额外的布尔值变量ok（一般命名为ok）来表示元素是否存在。
		if yv, ok := y[k]; !ok || yv != xv {
			return false
		}
	}
	return true
}
```



### example：dedup

```go
// Dedup prints only one instance of each line; duplicates are removed.
// See page 96.
package main

import (
	"bufio"
	"fmt"
	"os"
)

// Go语言中并没有提供set类型，但是map类型的键是不重复的，因此我们可以用map实现set的功能：
// Go程序员将这种忽略value的map当做一个字符串集合

// 读取多行输入，但只打印第一次出现的行
// 通过map来表示所有的输入行所对于的set集合，以确保已经出现过的行不会被重复打印
func main() {
	seen := make(map[string]bool) // a set of strings
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}

	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}
}
```



```go
// 程序使用map来记录提交相同的字符串列表的次数。
package main

import "fmt"

// 有时候我们需要一个map或set的key是slice类型，但是map的key必须是可比较的类型，但是slice并不满足这个条件。
// 我们可以绕过这个限制：map的key依然用可比较类型（如string），但每次存取前先将slice转为string类型。
// 1. 定义一个辅助函数k（key），将string slice转为string，作为map的key（string类型）。确保只有x和y相等时k(x) == k(y)才成立。
// 这里使用了fmt.Sprintf函数将字符串列表转换为一个字符串。以用于map的key。通过%q参数忠实地记录每个字符串元素的信息。
func k(list []string) string { return fmt.Sprintf("%q", list) }

// 2. 创建一个key为string类型的map m。在每次对map操作时m[key]，处理成m[k(list)](用k辅助函数将slice转化为string类型)。
var m = make(map[string]int)

func Add(list []string)       { m[k(list)]++ }
func Count(list []string) int { return m[k(list)] }

// 使用同样的技术可以处理任何不可比较的key类型，而不仅仅是slice类型。
// 这种技术对于想使用自定义key比较函数的时候也很有用
// 如在比较字符串的时候忽略大小写。同时，辅助函数k(x)也不一定是字符串类型，它可以返回任何可比较的类型，例如整数、数组或结构体等。
func main() {
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "C"})
	Add([]string{"A", "B", "D"})
	fmt.Println(Count([]string{"A", "B", "C"})) // 2

}

```



### example：charcount

```go
// Charcount computes counts of Unicode characters.
// See page 97.
// 程序用于统计输入中每个Unicode码点出现的次数
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	counts := make(map[rune]int) // counts of Unicode characters
	// UTF-8编码的长度总是从1到utf8.UTFMax（最大是4个字节），使用数组将更有效
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters

	in := bufio.NewReader(os.Stdin)
	for {
		// ReadRune方法执行UTF-8解码并返回三个值：解码的rune字符的值，字符UTF-8编码后的长度，和一个错误值。
		// 我们可预期的错误值只有对应文件结尾的io.EOF。
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		// 如果输入的是无效的UTF-8编码的字符，返回的将是unicode.ReplacementChar表示无效字符，并且编码长度是1。
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	// 读取本书的英文版原稿，统计的不同UTF-8编码长度的字符的数目：
	// 	len count
	// 1   765391
	// 2   60
	// 3   70
	// 4   0
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
}
```

### example：graph

```go
// Graph shows how to use a map of maps to represent a directed graph.
// See page 99.
package main

import "fmt"

// Map的value类型也可以是一个聚合类型（如是一个map或slice）。
// 从概念上讲，graph将一个字符串类型的key，映射到一组相关的字符串集合，它们指向新的graph的key。
var graph = make(map[string]map[string]bool)

func addEdge(from, to string) {
	// addEdge函数惰性初始化map是一个惯用方式，也就是说在每个值首次作为key时才初始化。
	// addEdge函数显示了如何让map的零值也能正常工作；即使from到to的边不存在，graph[from][to]依然可以返回一个有意义的结果。
	edges := graph[from]
	if edges == nil {
		edges = make(map[string]bool)
		graph[from] = edges
	}
	edges[to] = true
}

func hasEdge(from, to string) bool {
	return graph[from][to]
}

func main() {
	addEdge("a", "b")
	addEdge("c", "d")
	addEdge("a", "d")
	addEdge("d", "a")
	fmt.Println(hasEdge("a", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("a", "d"))
	fmt.Println(hasEdge("d", "a"))
	fmt.Println(hasEdge("x", "b"))
	fmt.Println(hasEdge("c", "d"))
	fmt.Println(hasEdge("x", "d"))
	fmt.Println(hasEdge("d", "x"))
}
```



