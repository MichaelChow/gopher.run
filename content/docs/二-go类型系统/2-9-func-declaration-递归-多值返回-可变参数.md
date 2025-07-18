---
title: "2.9 func declaration、递归、多值返回、可变参数"
date: 2024-12-28T15:07:00Z
draft: false
weight: 2009
---

# 2.9 func declaration、递归、多值返回、可变参数

## 函数声明

- 语句被组织成**函数**用于隔离和复用；
- 函数声明包括：
    - 一个名字
    - 一个形式参数列表：指定了一组**局部变量**的参数名和参数类型，其值由调用者传递的实际参数赋值而来
    - 一个可选的返回值列表：指定了返回值的类型，可像形参一样命名；**命名的返回值会声明为一个局部变量，初始化为其类型的零值；**
        - 当函数存在返回列表时，必须显式地以return语句结束，除非函数明确不会走完整个执行流程（如在函数中抛出宕机异常或者函数体内存在一个没有break退出条件的无限for循环）
    - 函数体
- **函数的类型****称为****函数的签名**。当两个函数**拥有相同的形参列表参数类型（和名字无关）和返回列表参数类型（和名字无关）**时，认为这**两个函数的类型或签名是相同的**。
    ```go
    func name(parameter-list) (result-list) {   // 和变量命名一样，相同类型可合并写
        body
    }  
    func add(x int, y int) int   {return x + y}  // 函数类型：func(int, int) int
    ```
- **Go没有默认参数值的概念**，也不能指定实参名，所以除了用于文档说明之外，**形参和返回值的命名不会对调用方有任何影响。**
- 函数形参以及命名返回值同属于**函数最外层作用域的局部变量**。
- **实参是按值传递的**，所以函数接收到的是每个实参的副本；修改函数的形参变量并不会影响到调用者提供的实参。然而，如果提供的实参包含引用类型（如指针、slice、map、函数或者channel)，那么当函数使用形参变量时就有可能会间接地修改实参变量。
- 有些函数的声明没有函数体，那说明这个函数使用除了Go以外的语言实现。
    ```go
    func Sin(x float64) float //implemented in **assembly language(汇编语言)**
    ```
## **递归**

- 函数可以递归调用（**可以直接或间接的调用自己**）。递归是一种实用的技术，**可以处理许多带有递归特性的数据结构**。
- [golang.org/x/](http://golang.org/x/)... （比如网络、国际化语言处理、移动平台、图片处理、加密功能以及开发者工具）都**由Go团队负责设计和维护，但**并不属于标准库，原因是**它们还在开发当中，或者很少被Go程序员使用。**
    ```shell
    // Findlinks1 prints the links in an HTML document read from standard input.
    // See page 122.
    package main
    import (
    	"fmt"
    	"os"
    	"golang.org/x/net/html"
    )
    func main() {
    	doc, err := html.Parse(os.Stdin)
    	if err != nil {
    		fmt.Fprintf(os.Stderr, "findlinks1: %v\n", err)
    		os.Exit(1)
    	}
    	for _, link := range visit(nil, doc) {
    		fmt.Println(link)
    	}
    }
    // visit函数递归调用，遍历HTML的节点树，从每一个anchor元素的href属性获得link,将这些links存入字符串数组中，并返回这个字符串数组。
    // visit appends to links each link found in n and returns the result.
    func visit(links []string, n *html.Node) []string {
    	if n.Type == html.ElementNode && n.Data == "a" {
    		for _, a := range n.Attr {
    			if a.Key == "href" {
    				links = append(links, a.Val)
    			}
    		}
    	}
    	for c := n.FirstChild; c != nil; c = c.NextSibling {
    		// 为了遍历结点n的所有后代结点，每次遇到n的孩子结点时，visit递归的调用自身（逻辑完全一样）。这些孩子结点存放在FirstChild链表中。
    		links = visit(links, c)
    	}
    	return links
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
    //!+html
    package html
    type Node struct {
    	Type                    NodeType
    	Data                    string
    	Attr                    []Attribute
    	FirstChild, NextSibling *Node
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
    //!-html
    */
    ```
- 许多编程语言使用固定长度的函数调用栈（大小在64KB到2MB之间）。递归的深度会受限于固定长度的栈大小，所以当进行深度递归调用时必须谨防栈溢出。固定长度的栈甚至会造成一定的安全隐患。相比固定长的栈，**Go的实现使用了可变长度的栈，栈的大小会随着使用而增长，可达到1GB左右的上限。这使得我们可以安全地使用递归而不用担心溢出问题。**
## **多值返回**

- 返回一个计算结果和一个错误值或是否调用正确的布尔值
- Go与众不同的特性之一就是函数和方法可返回多个值。
    - 这可改善C中一些笨拙的习惯：将错误值返回（例如用 `-1` 表示 `EOF`）和修改通过地址传入的实参。在C中，写入操作发生的错误会用一个负数标记，而错误码会隐藏在某个不确定的位置。
    - 而在Go中，`Write` 会返回写入的字节数**以及**一个错误： “是的，您写入了一些字节，但并未全部写入，因为设备已满”。 在 `os` 包中，`File.Write` 的签名为：
        ```go
        func (file *File) Write(b []byte) (n int, err error)
        ```
    - 正如文档所述，它返回写入的字节数，并在`n != len(b)` 时返回一个非 `nil` 的 `error` 错误值。 这是一种常见的编码风格
    - 我们可以采用一种简单的方法。来避免为模拟引用参数而传入指针。 以下简单的函数可从字节数组中的特定位置获取其值，并返回该数值和下一个位置。
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
    - 你可以像下面这样，通过它扫描输入的切片 `b` 来获取数字。
        ```go
        	for i := 0; i < len(b); {
        		x, i = nextInt(b, i)
        		fmt.Println(x)
        	}
        ```
- Go的GC机制将**回收未使用的内存**，但不能指望它会释放未使用的**操作系统资源（**如打开的文件、网络连接），必须显式地关闭它们。
    ```go
    resp.Body.Close()
    ```
- 良好的名称可以使得返回值更加有意义。尤其在一个函数返回多个结果且类型相同时，名字的选择更加重要
- 可命名的结果形参，起到文档的作用，使代码更加简短清晰：如nexPos一看就知道返回的 `int` 就值如其意了。
    ```go
    func Size(rect image.Rectangle) (width, height int)
    func Split(path string) (dir, file string)
    func HourMinSec(t time.Time) (hour, minute, second int)
    ```
    ```go
    func nextInt(b []byte, pos int) (value, nextPos int) {
    ```
- **按照惯例，函数的最后一个bool类型的返回值表示函数是否运行成功，error类型的返回值代表函数的错误信息，它们都无需再使用变量名解释。**
- **bare return （**裸返回）/ber/：如果返回值列表均为命名返回值，那么该函数的return语句可以省略操作数，代码更简洁。默认按照返回值列表的次序，返回所有的返回值。**但是使得代码可读性很差**。
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
# **可变参数**

- 可变参数函数：**参数数量可变的函数**。声明时需要在参数列表的最后一个参数类型之前加上省略符号“...”，表示该函数会接收任意数量的该类型参数。常被用于格式化字符串
    - Printf：首先接收一个必备的参数format string，之后接收任意个数的后续参数a ...anys。
        ```go
        func Printf(format string, a ...any) (n int, err error) {
        	return Fprintf(os.Stdout, format, a...)
        }
        ```
    - 函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串
    - errorf：构造了一个以行号开头的，经过格式化的错误信息
        - **interface{}表示函数的最后一个参数可以接收任意类型**
        ```go
        func errorf(linenum int, format string, args ...interface{}) {
        	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
        	fmt.Fprintf(os.Stderr, format, args...)
        	fmt.Fprintln(os.Stderr)
        }
        ```
- **在函数体中，vals****被看作是类型为[] int的切片 （所以也是语法糖？）。**
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
    - **调用者****隐式的创建一个数组****，并将原始参数复制到数组中。****再把数组的一个切片作为参数传给被调用函数****。**
        ```go
        fmt.Println(sum(1, 2, 3, 4))
        ```
    - **可变参数函数和以切片作为参数的函数****是不同的函数类型**
        ```go
        func([]int)
        func(...int)
        ```
    - 如果原始参数已经是切片类型，只需在最后一个参数后加上省略符，即可将切片的元素进行传递sum函数
        ```go
        values := []int{1, 2, 3, 4}
        fmt.Println(sum(values...)) // "10"
        ```


