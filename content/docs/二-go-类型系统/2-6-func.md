---
title: "2.6 func"
date: 2024-12-28T15:07:00Z
draft: false
weight: 2006
---

# 2.6 func



语句被组织成**函数**用于封装和复用。

## 一、func**类型 (Function Types)**

> 💡 **func类型****：即函数的签名，由 形参列表的参数类型 + 返回值列表的参数类型 组成。**和函数名、参数名无关。为引用类型，零值为nil。

```go
// 函数类型的基本语法
type FunctionType func(参数列表) 返回值列表
```



### 函数声明

函数声明（**func declaration**）包括：函数名、形式参数列表、返回值列表（可选）、函数体。

- 形式参数列表：指定了一组**局部变量**的参数名和参数类型，其值由调用者传递的实际参数赋值而来。
- 返回值列表：指定了返回值的类型，可像形参一样命名；**命名的返回值会声明为一个局部变量，初始化为其类型的零值；**当函数存在返回列表时，必须显式地以return语句结束，除非函数明确不会走完整个执行流程（如在函数中抛出宕机异常或者函数体内存在一个没有break退出条件的无限for循环）


函数形参以及命名返回值同属于**函数最外层作用域的局部变量**。

```go
// example
func add(x int, y int) int   {return x + y}  // 函数类型：func(int, int) int
```

**Go没有默认参数值的概念**，也不能指定实参名，所以除了用于文档说明之外，**形参和返回值的命名不会对调用方有任何影响。**

**实参是按****值****传递的：**函数接收到的是每个**实参的副本，**修改函数的形参变量并不会影响到调用者提供的实参。

但当实参包含**引用类型**（pointer、slice、map、func、channel)，那么当函数使用形参变量时就有可能**会间接地修改实参变量。**



有些函数的声明没有函数体，说明这个函数使用 除了Go以外的语言 实现（如**assembly language(汇编语言)**）。

```go
func Sin(x float64) float //implemented in **assembly language(汇编语言)**
```

### **递归调用**

函数可以递归调用（**可以直接或间接的调用自己**）。

递归是一种实用的技术，**可以处理许多带有递归特性的数据结构**。

[golang.org/x/](http://golang.org/x/)... （如网络、国际化语言处理、移动平台、图片处理、加密功能以及开发者工具）都**由Go团队负责设计和维护，但**并不属于标准库，原因是**它们还在开发当中，或者很少被Go程序员使用。**



**可变长的函数调用栈：**

许多编程语言使用固定长度的函数调用栈（大小在64KB到2MB之间）。递归的深度会受限于固定长度的栈大小，所以当进行深度递归调用时必须谨防栈溢出。固定长度的栈甚至会造成一定的安全隐患。

相比固定长的栈，**Go的实现使用了可变长度的栈，栈的大小会随着使用而增长，可达到1GB左右的上限。这使得我们可以安全地使用递归而不用担心溢出问题。**



**example**: visit爬虫

函数递归调用，遍历HTML的节点树，从每一个anchor元素的href属性获得link,将这些links存入字符串数组中，并返回这个字符串数组。

```go
// visit appends to links each link found in n and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)  // 递归append()
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		 = visit(links, c) // 为了遍历结点n的所有后代结点，每次遇到n的孩子结点时，visit递归的调用自身（逻辑完全一样）。这些孩子结点存放在FirstChild链表中。
	}
	return links
}
```

```go
package html

type Node struct { 
	Type                    NodeType
	Data                    string
	Attr                    []Attribute
	FirstChild, NextSibling *Node  // 递归结构
}

type NodeType int32

const (
	ErrorNode NodeType = iota
	TextNode
	DocumentNode
	ElementNode
	CommentNode
	DoctypeNode
)

type Attribute struct {
	Key, Val string
}

func Parse(r io.Reader) (*Node, error)
```

```shell
doc, err := html.Parse(os.Stdin)
if err != nil {
	fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
	os.Exit(1)
}
for _, link := range visit(nil, doc) {
	fmt.Println(link)
}

/*
// $ ../../ch1/10.fetch/fetch https://www.taobao.com | ./findlinks1
https://bk.taobao.com/k/taobaowangdian_457/
https://www.tmall.com/
https://bk.taobao.com/k/zhibo_599/
https://bk.taobao.com/k/wanghong_598/
https://bk.taobao.com/k/zhubo_601/
...
//
```



### **多值返回**

Go（与众不同的特性之一）函数和方法可返回多个值：计算结果、错误值 或 是否调用正确的布尔值。

这可改善C中一些笨拙的习惯：将错误值返回（例如用 `-1` 表示 `EOF`）和修改通过地址传入的实参。**在C中，写入操作发生的错误会用一个负数标记，而错误码会隐藏在某个不确定的位置。**

而在Go中，`Write` 会返回写入的字节数**以及**一个错误： “是的，您写入了一些字节，但并未全部写入，因为设备已满”。 

正如文档所述，它返回写入的字节数，并在`n != len(b)` 时返回一个非 `nil` 的 `error` 错误值。 这是一种常见的编码风格。

```go
// /usr/local/go/src/os/file.go
// Write writes len(b) bytes from b to the File.
// It returns the number of bytes written and an error, if any.
// Write returns a non-nil error when n != len(b).
func (f *File) Write(b []byte) (n int, err error) {
	if err := f.checkValid("write"); err != nil {
		return 0, err
	}
	n, e := f.write(b)
	if n < 0 {
		n = 0
	}
	if n != len(b) {
		err = io.ErrShortWrite
	}

	epipecheck(f, e)

	if e != nil {
		err = f.wrapErr("write", e)
	}

	return n, err
}
```



我们可以采用一种简单的方法。来避免为模拟引用参数而传入指针。 以下简单的函数可从字节数组中的特定位置获取其值，并返回该数值和下一个位置。

```go
func nextInt(b []byte, i int) (int, int) {
	for ; i < len(b) && !isDigit(b[i]); i++ {
	}
	x := 0
	for ; i < len(b) && isDigit(b[i]); i++ {
		x = x*10 + int(b[i]) - '0'
	}
	return x, i
}
```

你可以像下面这样，通过它扫描输入的切片 `b` 来获取数字。

```go
	for i := 0; i < len(b); {
		x, i = nextInt(b, i)
		fmt.Println(x)
	}
```



Go的GC机制将**回收未使用的****内存**，但**不能回收未使用的操作系统资源（如打开的文件、网络连接）**，**必须显式地关闭它们**。

```go
resp.Body.Close()
```



良好的名称可以使得返回值更加有意义，尤其在一个函数返回多个结果且类型相同时。

**可命名的结果形参，起到文档的作用**，使代码更加简短清晰：如nexPos一看就知道返回的 `int` 就值如其意了。

```go
func nextInt(b []byte, pos int) (value, nextPos int) {
}
```

```go
func Size(rect image.Rectangle) (width, height int)
func Split(path string) (dir, file string)
func HourMinSec(t time.Time) (hour, minute, second int)
```



**按照惯例，函数的最后一个bool类型的返回值表示函数是否运行成功，error类型的返回值代表函数的错误信息，它们都无需再使用变量名解释。**



**bare return （裸返回）/ber/：**

如果返回值列表均为命名返回值，那么该函数的return语句可以省略操作数，代码更简洁。

默认按照返回值列表的次序，返回所有的返回值。**但是使得代码可读性很差**。

```go
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
		//  **return 0,0,err（Go会将返回值 words和images在函数体的开始处，根据它们的类型，将其初始化为0） // 等价代码**
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
	// return words, images, err // 等价代码
}
```

### **可变参数**

可变参数函数：**参数数量可变的函数**。

声明时需要在参数列表的**最后一个参数类型之前加上省略符号“...”**，表示该函数会接收任意数量的该类型参数。

常被用于格式化字符串: 函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串

```go
// Printf：首先接收一个必备的参数format string，之后接收任意个数的后续参数a ...anys。
func Printf(format string, a ...any) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}
```

```go
// errorf：构造了一个以行号开头的，经过格式化的错误信息
// **interface{} 表示函数的最后一个参数可以接收任意类型**
func errorf(linenum int, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}
```

```go
// any is an alias for interface{} and is equivalent to interface{} in all ways.
type any = interface{}
```



**在函数体中，vals被看作是类型为[] int的切片。（所以也是语法糖？）**

```go
// sum可以接收任意数量的int型参数
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}
```



**调用者****隐式的创建一个数组****，并将原始参数复制到数组中。****再把数组的一个切片作为参数传给被调用函数****。**

```go
fmt.Println(sum(1, 2, 3, 4))
```



**可变参数函数和以切片作为参数的函数****是不同的函数类型**

```go
func([]int)
func(...int)
```

如果原始参数已经是切片类型，只需在最后一个参数后加上省略符，即可将切片的元素进行传递sum函数。

```go
values := []int{1, 2, 3, 4}
fmt.Println(sum(values...)) // "10"
```



## 二、func值** (Function Values)**

### func值** **

Go中函数是**一等公民（first-class values），**可以和其他值一样使用。（而Java中没有独立的函数，只能作为方法在类中。）

Go使用闭包（closures）技术实现函数值，Go程序员也把函数值叫做闭包。

调用值为nil的函数值会引起panic错误。除了和nil比较外，不可比较，所以不能作为map的key。

```go
var f func(int) int // 声明一个变量f，其类型为func(int) int的函数类型，值被初始化为零值nul。
f(3) // f为nil，引发panic: runtime error: invalid memory address or nil pointer dereference

if f != nil {
		f(3)
}
```

```go
func square(n int) int     { return n * n }

f := square
fmt.Println(f(3)) // "9"
```



函数变量使得函数不仅将数据进行参数化，还将函数的行为当作参数进行传递。

strings.Map对字符串中的每个字符调用add1函数，并将每个add1函数的返回值组成一个新的字符串返回给调用者。

```go
func add1(r rune) rune { return r + 1 }

fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
```



```go
// 1. 函数可以赋值给变量
var fn func(int) int
fn = func(x int) int { return x * 2 }

// 2. 函数可以作为参数传递
result := applyFunction(fn, 5)
fmt.Println(result)  // 10

// 3. 函数可以作为返回值
multiplier := createMultiplier(3)
fmt.Println(multiplier(4))  // 12

// 4. 函数可以存储在数据结构中
functions := map[string]func(int) int{
	"double": func(x int) int { return x * 2 },
	"square": func(x int) int { return x * x },
	"addOne": func(x int) int { return x + 1 },
}

for name, fn := range functions {
	fmt.Printf("%s(5) = %d\n", name, fn(5))
}
```



### example: outline

5.2节的findLinks函数使用了辅助函数visit，遍历和操作了HTML页面的所有结点。

使用函数值，我们**可以将遍历结点的逻辑和操作结点的逻辑分离**，使得我们可以复用遍历的逻辑，从而对结点进行不同的操作。

- 该函数接收2个函数作为参数，分别在结点的孩子被访问前和访问后调用。这样的设计给调用者更大的灵活性。
- forEachNode针对每个结点x，都会调用pre(x)和post(x)。pre和post都是可选的。
- 遍历孩子结点之前，pre被调用；遍历孩子结点之后，post被调用
```go
func outline(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	forEachNode(doc, startElement, endElement)

	return nil
}


func forEachNode(n *html.Node, pre, post func(n *html.Node)) {   // 函数变量
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}

var depth int

// startElemen和endElement两个函数用于输出HTML元素的开始标签和结束标签<b>...</b>

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		// 上面的代码利用fmt.Printf的一个小技巧控制输出的缩进。
		// %*s中的*会在字符串之前填充一些空格。在例子中，每次输出会先填充depth*2数量的空格，再输出""，最后再输出HTML标签。
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}

```



### 匿名函数** (Anonymous Functions)**

**命名函数只能****声明在包级别的作用域****，而使用****函数字面量****可在任何表达式内指定函数变量。**

**函数字面量：**在func关键字后面没有函数的名称，是一个表达式，它的值称为匿名函数。

**函数字面量在我们需要调用的时候才定义。**

**通过函数字面量这种定义的函数在同一个****词法块内****，因此里层的函数可以使用外层函数中的变量。**

```go
strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")
```

**引用类型（可能引用某些外层函数的变量）**：函数值不仅是一段代码还可以**拥有状态**：**里层的匿名函数能够获取和更新外层squares函数的局部变量x**。这些**隐藏的变量引用就是我们把函数归类为引用类型，而且函数变量无法进行比较的原因。**



### **闭包（Closure）**

**闭包**是一个**引用了其外部作用域中的变量的函数值**。/'kloʒɚ/ n. 关闭；终止，结束 vt. 使终止

由于外部变量在闭包中被引用，无法被GC回收，外部变量将一直保持“存活”（类似全局变量），后续调用都会直接继承原来的值（不再是无状态的）？？

**example：**

```go
// 闭包的基本特征
func createCounter() func() int {
    count := 0  // 外部变量
    return func() int {    // 返回的函数引用了外部变量count
        count++  // 访问并修改外部变量。同一个词法块内。
        return count
    }
}

func main() {
    counter := createCounter()  // counter为引用了外部作用域变量的函数类型的函数值，即闭包。count隐藏在counter中？？
    fmt.Println(counter())  // 1
    fmt.Println(counter())  // 2   **每次调用都会保持count上一次调用的状态。**createCounter返回后，变量count仍然隐式的存在于counter中，变量的生命周期不由它的作用域决定。
    fmt.Println(counter())  // 3
}
```



**闭包的实际应用:**

1. **状态保持:**外部变量在闭包中"存活”，有了记忆。
1. 工厂函数
1. 配置和选项模式
1. 中间件handler和装饰器
1. 事件处理和回调
1. 函数式编程


### **闭包的捕获迭代变量内存地址的陷阱**

函数变量（引用类型）使用的循环变量的内存地址，该地址的值被循环不断的更新，直到最后一次循环的值。

等到延迟到最后才执行的**函数变量、goruntine的go语句、defer语句**时，执行的结果会不符合预期。

```go
// 这个问题不仅存在基于range的循环，在下面的例子中，对循环变量i的使用也存在同样的问题：
var rmdirs []func()
dirs := tempDirs()
for i := 0; i < len(dirs); i++ {
    os.MkdirAll(dirs[i], 0755) // OK
    rmdirs = append(rmdirs, func() {
        os.RemoveAll(dirs[i]) // NOTE: incorrect!
    })
}
```



dir在for循环引进的一个块作用域内进行声明。在循环里创建的所有函数变量共享相同的变量(一个可访问的存储位置，而不是固定的值）。

**dir变量的值在不断地迭代中更新，因此当调用清理函数时，dir变量已经被每一次的for循环更新多次，dir变量的实际取值是最后一次迭代时的值，**所以所有的os.RemoveAll调用最终都试图删除最后一个目录。

```go
var rmdirs []func() 
for _, dir := range tempDirs() {
    dir := dir // 每次循环单独声明一个变量dir，值只不过是dir的一个副本，这看起来有些奇怪却是一个关键性的声明
    os.MkdirAll(dir, 0755)
    rmdirs = append(rmdirs, func() {
        os.RemoveAll(dir)
    })
}

for _, rmdir := range rmdirs {
    rmdir() // clean up
}
```



## 三、函数式编程（Functional Programming）

函数式编程是一种编程范式，它将计算过程看作是数学函数的求值，**避免使用可变状态和可变数据**。

1. **纯函数（Pure Functions）：**没有副作用，如打印、修改全局变量。小函数易于组合
```go
// ✅ 纯函数：相同输入总是产生相同输出，无副作用
func add(a, b int) int {
    return a + b
}

func square(x int) int {
    return x * x
}

// ❌ 非纯函数：有副作用
func addWithSideEffect(a, b int) int {
    fmt.Println("Adding:", a, b)  // 副作用：打印
    globalCounter++               // 副作用：修改全局状态
    return a + b
}
```

1. **不可变性（Immutability）**
```go
// ✅ 不可变：不修改原始数据
func doubleSlice(slice []int) []int {
    result := make([]int, len(slice))
    for i, v := range slice {
        result[i] = v * 2
    }
    return result
}

// ❌ 可变：修改原始数据
func doubleSliceMutable(slice []int) {
    for i := range slice {
        slice[i] *= 2  // 修改原始数据
    }
}
```



1. **函数作为参数**
```go
// 高阶函数：接受函数作为参数
func mapSlice(slice []int, fn func(int) int) []int {
    result := make([]int, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

func filterSlice(slice []int, predicate func(int) bool) []int {
    var result []int
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

func reduceSlice(slice []int, fn func(int, int) int, initial int) int {
    result := initial
    for _, v := range slice {
        result = fn(result, v)
    }
    return result
}
```

1. **函数作为返回值**
```go
// 函数工厂：返回函数
func createMultiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

func createAdder(addend int) func(int) int {
    return func(x int) int {
        return x + addend
    }
}

func main() {
    double := createMultiplier(2)
    addFive := createAdder(5)
    
    fmt.Println(double(3))   // 6
    fmt.Println(addFive(3))  // 8
}
```



**函数式编程的优势:**

1. **可读性：清晰的数据流**
```go
// 函数式风格：清晰的数据流
func processUsers(users []User) []string {
    return mapSlice(
        filterSlice(users, func(u User) bool {
            return u.Age >= 18
        }),
        func(u User) string {
            return u.Name
        },
    )
}

// 命令式风格：需要跟踪状态
func processUsersImperative(users []User) []string {
    var result []string
    for _, user := range users {
        if user.Age >= 18 {
            result = append(result, user.Name)
        }
    }
    return result
}
```

1. **可测试性: 纯函数易于测试**
```go
// 纯函数易于测试
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, expected int
    }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    
    for _, test := range tests {
        result := add(test.a, test.b)
        if result != test.expected {
            t.Errorf("add(%d, %d) = %d, want %d", 
                test.a, test.b, result, test.expected)
        }
    }
}
```



1. **不可变数据天然的并发安全**
```go
// 不可变数据天然并发安全
func processConcurrently(data []int) []int {
    chunks := chunkSlice(data, 4)
    results := make(chan []int, len(chunks))
    
    for _, chunk := range chunks {
        go func(c []int) {
            // 处理数据，不修改原始数据
            processed := mapSlice(c, func(x int) int {
                return x * x
            })
            results <- processed
        }(chunk)
    }
    
    var finalResult []int
    for i := 0; i < len(chunks); i++ {
        result := <-results
        finalResult = append(finalResult, result...)
    }
    
    return finalResult
}
```



**Go 中的函数式编程特性:**

**闭包（Closures）**

```go
func createCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    counter := createCounter()
    fmt.Println(counter())  // 1
    fmt.Println(counter())  // 2
    fmt.Println(counter())  // 3
}
```

**匿名函数:**

```go
// 立即执行函数
result := func(x, y int) int {
    return x + y
}(3, 4)

// 函数作为值
add := func(x, y int) int {
    return x + y
}
fmt.Println(add(3, 4))
```

**方法链:**

```go
type StringProcessor struct {
    value string
}

func (sp StringProcessor) ToUpper() StringProcessor {
    return StringProcessor{strings.ToUpper(sp.value)}
}

func (sp StringProcessor) Trim() StringProcessor {
    return StringProcessor{strings.TrimSpace(sp.value)}
}

func (sp StringProcessor) String() string {
    return sp.value
}

func main() {
    result := StringProcessor{"  hello world  "}.
        ToUpper().
        Trim()
    
    fmt.Println(result)  // "HELLO WORLD"
}
```



**性能开销：****函数式编程可能产生更多内存分配**

```go
func processDataFunctional(data []int) []int {
    // 每次操作都可能创建新的切片
    return mapSlice(
        filterSlice(data, isEven),
        square,
    )
}

func processDataImperative(data []int) []int {
    // 原地操作，更高效
    result := make([]int, 0, len(data))
    for _, v := range data {
        if v%2 == 0 {
            result = append(result, v*v)
        }
    }
    return result
}
```



**实际项目中的应用：**

**Web 框架中的中间件**

```go
type Middleware func(http.HandlerFunc) http.HandlerFunc

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        fmt.Printf("Request processed in %v\n", time.Since(start))
    }
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !isAuthenticated(r) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

func applyMiddleware(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}

func main() {
    handler := func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, World!")
    }
    
    finalHandler := applyMiddleware(handler, loggingMiddleware, authMiddleware)
    http.HandleFunc("/", finalHandler)
    http.ListenAndServe(":8080", nil)
}
```



