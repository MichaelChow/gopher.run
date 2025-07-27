---
title: "2.6 map、make、new"
date: 2025-04-03T08:00:00Z
draft: false
weight: 2006
---

# 2.6 map、make、new



- make: 内建函数，仅限于分配并初始化一个 **slice、map、chan类型** 的对象，make返回Type，new返回*Type
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
    func make(t Type, size ...IntegerType) Type
    // The new built-in function allocates memory. The first argument is a type,
    // not a value, and the value returned is a pointer to a newly
    // allocated zero value of that type.
    func new(Type) *Type
    ```


- Map(映射)：分配一个空映射，并有足够的空间来容纳指定数量的元素。size参数可以省略，此时会分配一个较小的初始大小。map是随机顺序；
    - **map**存储了键/值（key/value）的集合，对集合元素，提供常数时间的存、取或测试操作。键可以是任意类型，只要其值能用`==`运算符比较，最常见的例子是字符串；值则可以是任意类型。这个例子中的键是字符串，值是整数。
    - **内置函数 **`make`** 创建空 **`map`（译注：从功能和实现上说**，**`Go`** 的 **`map`** 类似于 **`Java`** 语言中的 **`HashMap`**，Python 语言中的 **`dict`，`Lua` 语言中的`table`，通常使用`hash`实现。遗憾的是，对于该词的翻译并不统一，**数学界术语为映射（**注释：如MyBatis中的Mapper**），**而计算机界众说纷纭莫衷一是。为了防止对读者造成误解，保留不译。
    - `map`中不存在某个键时不用担心，首次读到新行时，等号右边的表达式`counts[line]`的值将被**计算为其类型的零值**，对于`int`即`0`。
    - `map`的迭代顺序并不确定：从实践来看，**该顺序随机，每次运行都会变化**(实测是这样)。这种设计是有意为之的，因为**能防止开发的程序依赖特定遍历顺序**，而这是无法保证的。（译注：具体可以参见这里[https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order](https://stackoverflow.com/questions/11853396/google-go-lang-assignment-order)）map的顺序取决于使用的hash函数，hash函数为了修复DOS拒绝服务攻击做了随机化处理。[https://github.com/golang/go/issues/2630](https://github.com/golang/go/issues/2630)。


- 哈希表（map）是一种巧妙并且实用的数据结构，它是一个**无序的key/value对的集合**，其所有的key都是不同的，然后通过给定的key可以在**常数时间复杂度内检索、更新或删除对应的value**。
- 在Go语言中，一个map map[K]V，就是**一个哈希表的引用（引用类型）**，其中K和V分别对应key和value。map中所有的key都有相同的类型，所有的value也有着相同的类型，但key和value之间可以是不同的数据类型。
    - K对应的key必须是支持==比较运算符的数据类型（所以map可以通过测试key是否相等来判断是否已经存在）。虽然浮点数类型也是支持相等运算符比较的，但是将浮点数用做key类型则是一个坏的想法，正如第三章提到的，最坏的情况是可能出现的NaN和任何浮点数都不相等。
    - 对于V对应的value数据类型则没有任何的限制。
    - map的零值是nil，即没有引用任何hash表
    ```go
    ages := make(map[string]int) // **make 创建一个空map map[]（非nil）**，len(ages) = 0
    ages := map[string]int{}    // **map字面量创建一个空 map[]（非nil）**，len(ages) = 0
    var ages map[string]int    // 初始化为零值是nil，len(ages) = 0
    ages := map[string]int{   // 用map字面值的语法创建一个非空map
    		"alice":   31,
    		"charlie": 34,
    }
    ```
- map中的元素并不是一个变量
- 增删改查：
    ```go
    ages["alice"]  // 通过key对应的下标访问元素；如果key不存在，那么将返回该map元素类型的零值0（V为int）
    delete(ages, "alice")  // 数删除元素
    ```
    - 查找、删除、len和range都可以安全在nil值的map上，它们的行为和一个空map类似；但向一个nil值的map存入元素奖导致一个panic异常
    ```go
    var age map[string]int  // nil
    ages["carol"] = 21 // panic: assignment to entry in nil map
    ```
- **map元素不是一个变量，不可以获取它的地址**。**禁止对map元素取地址的原因**：**map可能随着元素数量的增长而重新分配更大的内存空间（已有的元素重新散列到新的存储位置），导致之前的地址失效**。
    ```go
    _ := &ages["bob"]  // compile error: cannot take address of map element
    ```


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

### dedup.go

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



### charcount.go

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

### graph.go

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



- map是方便而强大的内建数据结构，它可以关联不同类型的值。其键可以是任何相等性操作符支持的类型， 如整数、浮点数、复数、字符串、指针、接口（只要其动态类型支持相等性判断）、结构以及数组。 切片不能用作映射键，因为它们的相等性还未定义。与切片一样，映射也是引用类型。 若将映射传入函数中，并更改了该映射的内容，则此修改对调用者同样可见。
映射可使用一般的复合字面语法进行构建，其键-值对使用逗号分隔，因此可在初始化时很容易地构建它们。

```go
var timeZone = map[string]int{
	"UTC":  0*60*60,
	"EST": -5*60*60,
	"CST": -6*60*60,
	"MST": -7*60*60,
	"PST": -8*60*60,
}

```

赋值和获取映射值的语法类似于数组，不同的是映射的索引不必为整数。

```go
offset := timeZone["EST"]

```

若试图通过映射中不存在的键来取值，就会返回与该映射中项的类型对应的零值。 例如，若某个映射包含整数，当查找一个不存在的键时会返回 `0`。 集合可实现成一个值类型为 `bool` 的映射。将该映射中的项置为 `true` 可将该值放入集合中，此后通过简单的索引操作即可判断是否存在。

```go
attended := map[string]bool{
	"Ann": true,
	"Joe": true,
	...
}

if attended[person] { // 若某人不在此映射中，则为 false
	fmt.Println(person, "正在开会")
}

```

有时你需要区分某项是不存在还是其值为零值。如对于一个值本应为零的 `"UTC"` 条目，也可能是由于不存在该项而得到零值。你可以使用多重赋值的形式来分辨这种情况。

```go
var seconds int
var ok bool
seconds, ok = timeZone[tz]

```

显然，我们可称之为“逗号 ok”惯用法。在下面的例子中，若 `tz` 存在， `seconds` 就会被赋予适当的值，且 `ok` 会被置为 `true`； 若不存在，`seconds` 则会被置为零，而 `ok` 会被置为 `false`。

```go
func offset(tz string) int {
	if seconds, ok := timeZone[tz]; ok {
		return seconds
	}
	log.Println("unknown time zone:", tz)
	return 0
}

```

若仅需判断映射中是否存在某项而不关心实际的值，可使用[空白标识符](https://go-zh.org/doc/effective_go.html#%E7%A9%BA%E7%99%BD) （`_`）来代替该值的一般变量。

```go
_, present := timeZone[tz]

```

要删除映射中的某项，可使用内建函数 `delete`，它以映射及要被删除的键为实参。 即便对应的键不在该映射中，此操作也是安全的。

```go
delete(timeZone, "PDT")  // 现在用标准时间
```



## new

另一个创建变量的方法是调用内建的new函数。表达式new(T)将创建一个T类型的**匿名变量**，初始化为T类型的零值，然后返回变量地址，返回的指针类型为`*T`。用new创建变量和普通变量声明语句方式创建变量没有什么区别，new函数类似是一种**语法糖****（一种简洁的写法，编译器会自动转换）**，而不是一个新的基础概念。

```go
p := new(int)   // p, *int 类型, 指向匿名的 int 变量
fmt.Println(p, *p) // 0x14000098020 0
*p = 2          // 设置 int 匿名变量的值为 2
fmt.Println(p, *p) // 0x14000098020 2

// 下面的两个newInt函数有着相同的行为：
func newInt() *int {
    return new(int)
}

func newInt() *int {
    var dummy int
    return &dummy
}
```

当然也可能有特殊情况：**如果两个类型都是空的（****说类型的大小是0****）**，例如`struct{}`和`[0]int`，**有可能有相同的地址（依赖具体的语言实现）**（译注：请谨慎使用大小为0的类型，因为如果类型的大小为0的话，可能导致Go语言的自动垃圾回收器有不同的行为，具体请查看`runtime.SetFinalizer`函数相关文档）。

```go
// new study
package main

import "fmt"

func main() {
	// 特殊情况：如果类型的大小都为0，地址可能会相同
	p, q := new([0]int), new([0]int)
	fmt.Println(&p, *p, &q, *q) // 0x14000058030 [] 0x14000058038 []
	fmt.Println(p == q, &p, &q) // true 0x14000058030 0x14000058038，为什么是true?
}
```



```shell
// new study
package main

import "fmt"

func main() {
	// 特殊情况：如果类型的大小都为0，地址可能会相同
	p, q := new([0]int), new([0]int)
	fmt.Println(p == q) // false，fmt.Println导致了两种不同的比较结果，可能为特定的编译器和优化设置下发生的罕见情况；
}
```



new函数使用通常相对比较少，因为对于结构体来说，直接用字面量语法创建新变量的方法会更灵活**；**由于new只是一个预定义的函数，它并不是一个关键字（Java中为关键字），因此我们可以将new名字重新定义为别的类型。例如下面的例子：

```go
// 由于new被定义为int类型的变量名，因此在delta函数内部是无法使用内置的new函数的。
func delta(old, new int) int { return new - old }
```



### `**new**`** 分配**

Go提供了两种分配原语，即内建函数 `new` 和 `make`。 它们所做的事情不同，所应用的类型也不同。它们可能会引起混淆，但规则却很简单。 让我们先来看看 `new`。这是个用来分配内存的内建函数， 但与其它语言中的同名函数不同，它不会**初始化**内存，只会将内存**置零**。 也就是说，`new(T)` 会为类型为 `T` 的新项分配已置零的内存空间， 并返回它的地址，也就是一个类型为 `*T` 的值。用Go的术语来说，它返回一个指针， 该指针指向新分配的，类型为 `T` 的零值。

既然 `new` 返回的内存已置零，那么当你设计数据结构时， 每种类型的零值就不必进一步初始化了，这意味着该数据结构的使用者只需用 `new` 创建一个新的对象就能正常工作。例如，`bytes.Buffer` 的文档中提到“零值的 `Buffer` 就是已准备就绪的缓冲区。" 同样，`sync.Mutex` 并没有显式的构造函数或 `Init` 方法， 而是零值的 `sync.Mutex` 就已经被定义为已解锁的互斥锁了。

“零值属性”可以带来各种好处。考虑以下类型声明。

```go
type SyncedBuffer struct {
	lock    sync.Mutex
	buffer  bytes.Buffer
}

```

`SyncedBuffer` 类型的值也是在声明时就分配好内存就绪了。后续代码中， `p` 和 `v` 无需进一步处理即可正确工作。

```go
p := new(SyncedBuffer)  // type *SyncedBuffer
var v SyncedBuffer      // type  SyncedBuffer

```

### **构造函数与复合字面**

有时零值还不够好，这时就需要一个初始化构造函数，如来自 `os` 包中的这段代码所示。

```go
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := new(File)
	f.fd = fd
	f.name = name
	f.dirinfo = nil
	f.nepipe = 0
	return f
}

```

这里显得代码过于冗长。我们可通过**复合字面**来简化它， 该表达式在每次求值时都会创建新的实例。

```go
func NewFile(fd int, name string) *File {
	if fd < 0 {
		return nil
	}
	f := File{fd, name, nil, 0}
	return &f
}

```

请注意，返回一个局部变量的地址完全没有问题，这点与C不同。该局部变量对应的数据 在函数返回后依然有效。实际上，每当获取一个复合字面的地址时，都将为一个新的实例分配内存， 因此我们可以将上面的最后两行代码合并：

```go
	return &File{fd, name, nil, 0}

```

复合字面的字段必须按顺序全部列出。但如果以 **字段**`:`**值** 对的形式明确地标出元素，初始化字段时就可以按任何顺序出现，未给出的字段值将赋予零值。 因此，我们可以用如下形式：

```go
	return &File{fd: fd, name: name}

```

少数情况下，若复合字面不包括任何字段，它将创建该类型的零值。表达式 `new(File)` 和 `&File{}` 是等价的。

复合字面同样可用于创建数组、切片以及映射，字段标签是索引还是映射键则视情况而定。 在下例初始化过程中，无论 `Enone`、`Eio` 和 `Einval` 的值是什么，只要它们的标签不同就行。

```go
a := [...]string   {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
s := []string      {Enone: "no error", Eio: "Eio", Einval: "invalid argument"}
m := map[int]string{Enone: "no error", Eio: "Eio", Einval: "invalid argument"}

```

### `**make**`** 分配**

再回到内存分配上来。内建函数 `make(T, `*args*`)` 的目的不同于 `new(T)`。它只用于创建切片、映射和信道，并返回类型为 `T`（而非 `*T`）的一个**已初始化** （而非**置零**）的值。 出现这种用差异的原因在于，这三种类型本质上为引用数据类型，它们在使用前必须初始化。 例如，切片是一个具有三项内容的描述符，包含一个指向（数组内部）数据的指针、长度以及容量， 在这三项被初始化之前，该切片为 `nil`。对于切片、映射和信道，`make` 用于初始化其内部的数据结构并准备好将要使用的值。例如，

```go
make([]int, 10, 100)

```

会分配一个具有100个 `int` 的数组空间，接着创建一个长度为10， 容量为100并指向该数组中前10个元素的切片结构。（生成切片时，其容量可以省略，更多信息见切片一节。） 与此相反，`new([]int)` 会返回一个指向新分配的，已置零的切片结构， 即一个指向 `nil` 切片值的指针。

下面的例子阐明了 `new` 和 `make` 之间的区别：

```go
var p *[]int = new([]int)       // 分配切片结构；*p == nil；基本没用
var v  []int = make([]int, 100) // 切片 v 现在引用了一个具有 100 个 int 元素的新数组

// 没必要的复杂：
var p *[]int = new([]int)
*p = make([]int, 100, 100)

// 习惯用法：
v := make([]int, 100)

```

请记住，`make` 只适用于映射、切片和信道且不返回指针。若要获得明确的指针， 请使用 `new` 分配内存。

