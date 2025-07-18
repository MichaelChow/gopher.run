---
title: "2.8 struct、组合嵌套、json"
date: 2025-04-03T08:02:00Z
draft: false
weight: 2008
---

# 2.8 struct、组合嵌套、json

- struct是**将零个或多个任意类型的命名变量组合在一起的聚合数据类型**，每个变量都叫做**结构体的成员**；
- 数据处理领域用结构体的经典案例是员工信息记录 **Employee /ɪmˈplɔɪiː/**
    - 结构体成员的不同的输入顺序，定义了不同的结构体类型；
    - **结构体的成员是变量（可寻址的值），可以取地址**；
    ```go
    type Employee struct {
    	ID int                // 工号
    	Name, Address string  // 姓名、地址，相同的成员类型且相关的成员，合并到一行
    	DoB       time.Time   // 出生日期 
    	Position string      // 职位
    	Salary    int        // 薪水
    	ManagerID int        // 直属领导
    }
    var dilbert Employee
    dilbert.Salary -= 5000
    position := &dilbert.Position
    *position = "Senior " + *position 
    // 点操作符也可以和指向结构体的指针一起工作：
    var employeeOfTheMonth *Employee = &dilbert
    employeeOfTheMonth.Position += "(proactive team player)"
    (*employeeOfTheMonth).Position += " (proactive team player)"  // 等价写法
    ```
- 结构体类型的值可以通过**结构体字面量**来设置，即通过设置结构体的成员变量来设置；
    ```go
    type Point struct { X, Y int }
    p := Point{1,2}  // 初始化写法1：要求完全按顺序写；仅在定义结构体的包内部使用，或者是在较小的结构体中使用
    p := point.Point{X: 1, Y: 2}  // 初始化写法2：name:value，与顺序无关，默认用成员变量的零值
    ```
- **较大的结构体通常会用指针的方式传入和返回，（避免分配一份大结构体的内存再拷贝其数据的操作）；**（在Go语言中，所有的函数参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量）
    ```go
    func Scale(p point.Point, factor int) point.Point {  // 传结构体，赋值给函数的新的局部变量，再返回结构体再赋值
    	return point.Point{p.X * factor, p.Y * factor}
    }
    func Bonus(p *point.Point, factor int) {   // 传结构体指针，无需返回
        p.X *= factor
    	  p.Y *= factor
    }
    ```
- **如果结构体的全部成员都是可以比较的，那么结构体也是可以比较的**，可以使用==或!=运算符进行比较两个结构体的每个成员
    ```go
    type Point struct{ X, Y int }
    p := Point{1, 2}
    q := Point{2, 1}
    fmt.Println(p.X == q.X && p.Y == q.Y) // "false"
    fmt.Println(p == q)                   // "false"
    ```
    ```go
    // 可比较的结构体类型（和其他可比较的类型一样）可用于map的key类型
    type address struct {
        hostname string
        port     int
    }
    hits := make(map[address]int)
    hits[address{"golang.org", 443}]++
    ```
- Go中不同寻常的**结构体嵌套机制**，将一个命名结构体当做另一个结构体类型的匿名成员使用；并可以通过简单的表达式就可以访问嵌套的成员（如x.d.e.f）
    ```go
    type Circle struct {   
        X, Y, Radius int  // 圆形类型，圆心的X、Y坐标、半径
    }
    type Wheel struct {
        X, Y, Radius, Spokes int // 轮形类型，增加了Spokes 径向辐条数
    }
    type Point struct {      // 独立相同的属性
        X, Y int
    }
    type Circle struct {   // 改进后
        Center Point
        Radius int
    }
    type Wheel struct {   // 改进后
        Circle Circle
        Spokes int
    }
    var w Wheel
    w.Circle.Center.X = 8     // 但**访问每个成员变得繁琐**
    w.Circle.Center.Y = 8
    w.Circle.Radius = 5
    w.Spokes = 20
    ```
- Go允许定义不带名称的结构体成员（即**匿名成员**），只需要指定类型即可；**匿名成员实际拥有隐式的名字（即对应类型的名字，如Point、Circle）**，只是在点操作符时可以省略，仅仅是点操作符的语法糖。
    ```go
    type Circle struct {
        Point            // Point类型被嵌入到了Circle结构体
        Radius int
    }
    type Wheel struct {
        Circle           // Circle类型被嵌入到了Wheel结构体
        Spokes int
    }
    var w Wheel
    w.X = 8            // 等价写法：**w.Circle.Point.X = 8**
    w.Radius = 5       // 等价写法： w.Circle.Radius = 5
    ```
    - 匿名成员的类型不能是匿名类型; **匿名成员的类型必须是 任何一个命名类型或指向命名类型的指针(任何，不仅仅是结构体类型)；**
    - **嵌套一个非结构体类型（没有子成员）有什么用？以该快捷方式访问匿名成员的内部变量同样适用于访问匿名成员的方法。这个机制就是从简单类型对象组合成复杂的组合类型的主要方式，Go中组合是面向对象编程方式的核心；**
    - 遗憾的是结构体字面值并没有简短表示匿名成员的语法，结构体字面值必须遵循 形状类型声明时的结构；
        ```go
        	w = Wheel{Circle{Point{8, 8}, 5}, 20}
        ```




### treesort 

```go
// 使用一个二叉树来实现一个插入排序
package treesort

// 结构体类型的零值是每个成员都是零值。通常会将零值作为最合理的默认值
// 如，对于bytes.Buffer类型，结构体初始值就是一个随时可用的空缓存
// sync.Mutex的零值也是有效的未锁定状态
// 有时候这种零值可用的特性是自然获得的，但是也有些类型需要一些额外的工作

// 一个命名为S的结构体类型将不能再包含S类型的成员：因为一个聚合的值不能包含它自身。（该限制同样适用于数组。）
// 但S类型的结构体可以包含*S指针类型的成员，这可以让我们创建递归的数据结构，比如链表和树结构等。
// 声明一个struct，二叉树
type tree struct {
	value       int
	left, right *tree
}

// Sort 函数对整数切片进行排序
func Sort(values []int) {
	// 初始化根节点为空
	var root *tree
	// 遍历切片中的每个元素
	for _, v := range values {
		// 将元素插入到二叉树中
		root = add(root, v)
	}
	// 将排序后的元素追加到原始切片的前缀
	appendValues(values[:0], root)
}

// add 函数向二叉树 t 中插入一个值为 value 的节点，并返回插入后的二叉树
func add(t *tree, value int) *tree {
	// 如果树为空，则创建一个新的树节点
	if t == nil {
		// 等价于返回 &tree{value: value}
		t = new(tree)
		t.value = value
		return t
	}
	// 如果值小于当前节点的值，则将其插入到左子树中
	if value < t.value {
		t.left = add(t.left, value)
		// 如果值大于等于当前节点的值，则将其插入到右子树中
	} else {
		t.right = add(t.right, value)
	}
	// 返回插入后的二叉树
	return t
}

// appendValues 函数将二叉树 t 中的元素按顺序追加到 values 切片中，并返回结果切片
func appendValues(values []int, t *tree) []int {
	// 如果树不为空
	if t != nil {
		// 递归地将左子树中的元素追加到 values 中
		values = appendValues(values, t.left)
		// 将当前节点的值追加到 values 中
		values = append(values, t.value)
		// 递归地将右子树中的元素追加到 values 中
		values = appendValues(values, t.right)
	}
	// 返回最终的排序结果
	return values
}
```



```go
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
// 用单元测试上述代码

package treesort_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	treesort "gopher.run/go/src/ch4/13.treesort"
)

func TestSort(t *testing.T) {
	data := make([]int, 50)
	for i := range data {
		data[i] = rand.Int() % 50
	}
	fmt.Println(data)
	treesort.Sort(data)
	fmt.Println(data)
	if !sort.IntsAreSorted(data) {
		t.Errorf("not sorted: %v", data)
	}
}

=== RUN   TestSort
[39 38 30 39 23 1 18 26 33 22 6 26 1 48 32 14 35 38 27 36 35 49 6 44 21 38 32 9 30 8 1 48 29 23 45 27 20 0 37 1 22 11 24 44 24 38 10 21 24 42]
[0 1 1 1 1 6 6 8 9 10 11 14 18 20 21 21 22 22 23 23 24 24 24 26 26 27 27 29 30 30 32 32 33 35 35 36 37 38 38 38 38 39 39 42 44 44 45 48 48 49]
--- PASS: TestSort (0.00s)
PASS
```

```go
// 如果结构体没有任何成员的话就是空结构体，写作struct{}。它的大小为0，也不包含任何信息，但是有时候依然是有价值的。
// 有些Go语言程序员用map来模拟set数据结构时，用它来代替map中布尔类型的value，只是强调key的重要性，但是因为节约的空间有限，而且语法比较复杂，所以我们通常会避免这样的用法。
seen := make(map[string]struct{}) // set of strings
// ...
if _, ok := seen[s]; !ok {
    seen[s] = struct{}{}
    // ...first time seeing s...
}
```

# JSON

- 在类似的协议中，JSON并不是唯一的一个标准协议。 XML（§7.14）、ASN.1和Google的Protocol Buffers都是类似的协议，并且有各自的特色，但是由于简洁性、可读性和流行程度等原因，JSON是应用最广泛的一个。
- Go语言对于这些标准格式的编码和解码都有良好的支持，由标准库中的**encoding/json**、encoding/xml、encoding/asn1等包提供支持（译注：Protocol Buffers的支持由 [github.com/golang/protobuf](http://github.com/golang/protobuf) 包提供），并且这类包都有着相似的API接口。
- JavaScript对象表示法（JSON）是一种用于发送和接收结构化信息的标准协议。JSON是对JavaScript中各种类型的值——字符串、数字、布尔值和对象——Unicode本文编码。它可以用有效可读的方式表示基础数据类型和聚合数据类型(数组、slice、结构体和map)。
    - 基本的JSON类型有数字（十进制或科学记数法）、布尔值（true或false）、字符串，其中字符串是以双引号包含的Unicode字符序列，支持和Go语言类似的反斜杠转义特性，不过**JSON使用的是**`**\Uhhhh**`**转义数字来表示一****个UTF-16编码，而不是Go语言的rune类型**。（译注：UTF-16和UTF-8一样是一种变长的编码，有些Unicode码点较大的字符需要用4个字节表示；而且UTF-16还有大端和小端的问题）
    - 这些基础类型可以通过JSON的数组和对象类型进行递归组合。一个JSON数组是一个有序的值序列，写在一个方括号中并以逗号分隔；一个JSON数组可以用于编码Go语言的数组和slice。
    - 一个JSON对象是一个字符串到值的映射，写成一系列的name:value对形式，用花括号包含并以逗号分隔；
    - **JSON的对象类型可以用于编码 Go语言的map类型（key类型是字符串，Go中Map的key有意处理为无序的）和Struct**。如：
    - 将一个Go语言中类似movies的 结构体slice 转为 JSON 的过程叫编组/编码（marshaling） 美 /ˈmɑːrʃl/
        - 编组通过调用json.Marshal函数完成，返回一个编码后的字节slice，包含很长的字符串，并且没有空白缩进以紧凑的表示。(类似Java中的序列化为Json串)
        - 注意：JavaScript中Map的所有键均按插入顺序排列；
```go
boolean         true
number          -273.15
string          "She said \"Hello, BF\""
array           ["gold", "silver", "bronze"]
object          {"year": 1980,
                 "event": "archery",
                 "medals": ["gold", "silver", "bronze"]}
```



### movie.go

```go
// Movie prints Movies as JSON.
// See page 108.
// 程序负责收集各种电影评论并提供反馈功能
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Movie struct {
	Title string
	// 在结构体声明中，Year和Color成员后面的字符串面值是结构体成员Tag
	// 结构体的成员Tag可以是任意的字符串面值，但是通常是一系列用空格分隔的key:"value"键值对序列
	// 因为值中含有双引号字符，因此成员Tag一般用原生字符串面值（即反引号`包裹）的形式书写，这样不再需要转义字符来表示双引号
	// "json"开头键名对应的值，用于控制encoding/json包的编码和解码的行为，并且encoding/...下面其它的包也遵循这个约定
	// 值的第一部分用于指定JSON对象的名字，如将Go语言中的TotalCount成员对应到JSON中的total_count对象
	// Year名字的成员在编码后变成了released，还有Color成员编码后变成了小写字母开头的color
	// 一个结构体成员Tag是和在编译阶段关联到该成员的元信息字符串
	Year int `json:"released"`
	// 额外的omitempty选项，表示当Go语言结构体成员为空或零值时不生成该JSON对象（这里false为零值）。
	// 即Marshal的json串中"Casablanca"由于其为零值fasle，所以没有Color字段
	Color  bool `json:"color,omitempty"`
	Actors []string
}

var movies = []Movie{
	{Title: "Casablanca", Year: 1942, Color: false,
		Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
	{Title: "Cool Hand Luke", Year: 1967, Color: true,
		Actors: []string{"Paul Newman"}},
	{Title: "Bullitt", Year: 1968, Color: true,
		Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
	// ...
}

func main() {
	{
		data, err := json.Marshal(movies)
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("%s\n", data)
	}

	{
		// 格式化json输出，产生整齐缩进的输出
		// MarshalIndent()有两个额外的字符串参数，用于表示每一行输出的前缀和每一个层级的缩进：
		// 译注：在最后一个成员或元素后面并没有逗号分隔符

		// 在编码时，默认使用Go语言结构体的成员名字作为JSON的对象（通过reflect反射技术，我们将在12.6节讨论）。
		// 只有导出的结构体成员才会被编码，这也就是我们为什么选择用大写字母开头的成员名称。
		data, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("%s\n", data)

		// 编码的逆操作是解码，对应将JSON数据解码为Go语言的数据结构，Go语言中一般叫unmarshaling
		// 类似Java中的将json串反序列化为对象
		// 代码将JSON格式的电影数据data解码为一个结构体slice，结构体中只有Title成员
		var titles []struct{ Title string }
		// 通过定义合适的Go语言数据结构，我们可以选择性地解码JSON中感兴趣的成员
		if err := json.Unmarshal(data, &titles); err != nil {
			log.Fatalf("JSON unmarshaling failed: %s", err)
		}
		// 这里的slice将被只含有Title信息的值填充，其它JSON成员将被忽略
		fmt.Println(titles) //  [{Casablanca} {Cool Hand Luke} {Bullitt}]
	}
}

/*
//!+output
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingr
id Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Ac
tors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"
Actors":["Steve McQueen","Jacqueline Bisset"]}]
//!-output
*/

/*
//!+indented
[
    {
        "Title": "Casablanca",
        "released": 1942,
        "Actors": [
            "Humphrey Bogart",
            "Ingrid Bergman"
        ]
    },
    {
        "Title": "Cool Hand Luke",
        "released": 1967,
        "color": true,
        "Actors": [
            "Paul Newman"
        ]
    },
    {
        "Title": "Bullitt",
        "released": 1968,
        "color": true,
        "Actors": [
            "Steve McQueen",
            "Jacqueline Bisset"
        ]
    }
]
//!-indented
*/

```

### github.go

```go
// Package github provides a Go API for the GitHub issue tracker.
// See page 110.
// See https://developer.github.com/v3/search/#search-issues.
// 通过Github的issue查询服务来演示
// 定义合适的类型和常量
package github

import "time"

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number  int
	HTMLURL string `json:"html_url"`
	Title   string
	State   string
	User    *User
	// 有些JSON成员名字和Go结构体成员名字并不相同，因此需要Go语言结构体成员Tag来指定对应的JSON名字
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
```



### search.go

```go
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SearchIssues函数发出一个HTTP请求，然后解码返回的JSON格式的结果
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	// 用户提供的查询条件可能包含类似?和&之类的特殊字符，为了避免对URL造成冲突，用url.QueryEscape来对查询中的特殊字符进行转义操作
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	//!-
	// For long-term stability, instead of http.Get, use the
	// variant below which adds an HTTP request header indicating
	// that only version 3 of the GitHub API is acceptable.
	//
	//   req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	//   if err != nil {
	//       return nil, err
	//   }
	//   req.Header.Set(
	//       "Accept", "application/vnd.github.v3.text-match+json")
	//   resp, err := http.DefaultClient.Do(req)
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	// movie.go中使用了json.Unmarshal函数来将JSON格式的字符串解码为字节slice
	// 这里使用了基于流式的解码器json.Decoder，它可以从一个输入流解码JSON数据，尽管这不是必须的
	// 相对应的，还有一个针对输出流的json.Encoder编码对象
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

```



### issues.go

```go
// Issues prints a table of GitHub issues matching the search terms.
// See page 112.
package main

import (
	"fmt"
	"log"
	"os"

	github "gopher.run/go/src/ch4/18.github"
)

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s %s %s\n",
			item.Number, item.User.Login, item.Title, item.HTMLURL, item.CreatedAt)
	}
}

/*
//!+textoutput
$ go build gopl.io/ch4/issues
$ ./issues repo:golang/go is:open json decoder
92 issues:
#48298     dsnet encoding/json: add Decoder.DisallowDuplicateFields https://github.com/golang/go/issues/48298 2021-09-09 19:39:33 +0000 UTC
#69449   gazerro encoding/json: Decoder.Token does not return an error f https://github.com/golang/go/issues/69449 2024-09-13 15:50:14 +0000 UTC
#5901        rsc encoding/json: allow per-Encoder/per-Decoder registrati https://github.com/golang/go/issues/5901 2013-07-17 16:39:15 +0000 UTC
#56733 rolandsho encoding/json: add (*Decoder).SetLimit https://github.com/golang/go/issues/56733 2022-11-14 18:51:33 +0000 UTC
#6647    btracey x/pkgsite: display type kind of each named type https://github.com/golang/go/issues/6647 2013-10-23 17:19:48 +0000 UTC
#42571     dsnet encoding/json: clarify Decoder.InputOffset semantics https://github.com/golang/go/issues/42571 2020-11-13 00:09:09 +0000 UTC
#11046     kurin encoding/json: Decoder internally buffers full input https://github.com/golang/go/issues/11046 2015-06-03 19:25:08 +0000 UTC
#67525 mateusz83 encoding/json: don't silently ignore errors from (*Deco https://github.com/golang/go/pull/67525 2024-05-20 14:10:55 +0000 UTC
#58649 nabokihms encoding/json: show nested fields path if DisallowUnkno https://github.com/golang/go/issues/58649 2023-02-22 23:20:53 +0000 UTC
#43716 ggaaooppe encoding/json: increment byte counter when using decode https://github.com/golang/go/pull/43716 2021-01-15 08:58:39 +0000 UTC
#36225     dsnet encoding/json: the Decoder.Decode API lends itself to m https://github.com/golang/go/issues/36225 2019-12-19 22:26:12 +0000 UTC
#26946    deuill encoding/json: clarify what happens when unmarshaling i https://github.com/golang/go/issues/26946 2018-08-12 18:19:01 +0000 UTC
#29035    jaswdr proposal: encoding/json: add error var to compare  the  https://github.com/golang/go/issues/29035 2018-11-30 11:21:31 +0000 UTC
#61627    nabice x/tools/gopls: feature: CLI syntax for renaming by iden https://github.com/golang/go/issues/61627 2023-07-28 06:40:34 +0000 UTC
#34543  maxatome encoding/json: Unmarshal & json.(*Decoder).Token report https://github.com/golang/go/issues/34543 2019-09-25 22:13:24 +0000 UTC
#32779       rsc encoding/json: memoize strings during decode https://github.com/golang/go/issues/32779 2019-06-25 21:08:30 +0000 UTC
#40128  rogpeppe proposal: encoding/json: garbage-free reading of tokens https://github.com/golang/go/issues/40128 2020-07-09 07:58:19 +0000 UTC
#40982   Segflow encoding/json: use different error type for unknown fie https://github.com/golang/go/issues/40982 2020-08-22 21:07:03 +0000 UTC
#59053   joerdav proposal: encoding/json: add a generic Decode function https://github.com/golang/go/issues/59053 2023-03-15 16:20:31 +0000 UTC
#65691  Merovius encoding/xml: Decoder does not reject xml-ProcInst prec https://github.com/golang/go/issues/65691 2024-02-13 10:33:20 +0000 UTC
#14750 cyberphon encoding/json: parser ignores the case of member names https://github.com/golang/go/issues/14750 2016-03-10 13:04:44 +0000 UTC
#40127  rogpeppe encoding/json: add Encoder.EncodeToken method https://github.com/golang/go/issues/40127 2020-07-09 07:52:47 +0000 UTC
#16212 josharian encoding/json: do all reflect work before decoding https://github.com/golang/go/issues/16212 2016-06-29 16:07:36 +0000 UTC
#41144 alvaroale encoding/json: Unmarshaler breaks DisallowUnknownFields https://github.com/golang/go/issues/41144 2020-08-31 14:27:19 +0000 UTC
#64847 zephyrtro encoding/json: UnmarshalJSON methods of embedded fields https://github.com/golang/go/issues/64847 2023-12-22 17:08:52 +0000 UTC
#56332    gansvv encoding/json: clearer error message for boolean like p https://github.com/golang/go/issues/56332 2022-10-19 19:30:20 +0000 UTC
#43513 Alexander encoding/json: add line number to SyntaxError https://github.com/golang/go/issues/43513 2021-01-05 10:59:27 +0000 UTC
#22752  buyology proposal: encoding/json: add access to the underlying d https://github.com/golang/go/issues/22752 2017-11-15 23:46:13 +0000 UTC
#33835     Qhesz encoding/json: unmarshalling null into non-nullable gol https://github.com/golang/go/issues/33835 2019-08-26 10:27:12 +0000 UTC
#33854     Qhesz encoding/json: unmarshal option to treat omitted fields https://github.com/golang/go/issues/33854 2019-08-27 00:20:25 +0000 UTC
//!-textoutput
*/
```

# **文本模板和HTML模板**

最简单的格式化，使用Printf是完全足够的
但是有时候会需要复杂的打印格式，这时候一般需要将格式化代码分离出来以便更安全地修改
这些功能是由text/template和html/template等模板包提供的，它们提供了一个将变量值填充到一个文本或HTML格式的模板的机制

### issuesreport.go

```go
// Issuesreport prints a report of issues matching the search terms.
// See page 113.
package main

import (
	"log"
	"os"
	"text/template"
	"time"

	github "gopher.run/go/src/ch4/18.github"
)

// !+template

// 一个模板是一个字符串或一个文件，里面包含了一个或多个由双花括号包含的{{action}}对象。大部分的字符串只是按字面值（字符串“”包裹的内容）打印。
// 每个actions都包含了一个用模板语言书写的表达式，一个action虽然简短但是可以输出复杂的打印值
// 模板语言包含通过选择结构体的成员、调用函数或方法、表达式控制流if-else语句和range循环语句，还有其它实例化模板等诸多特性。

// 这个模板先打印匹配到的issue总数，然后打印每个issue的编号、创建用户、标题还有存在的时间。
// 对于每一个action，都有一个当前值的概念，对应点操作符，写作“.”。最初被初始化为调用模板时的参数 （对应github.IssuesSearchResult类型的变量）

// {{.TotalCount}}：对应action将展开为结构体中TotalCount成员以默认的方式打印的值
// {{range .Items}}和{{end}}：对应一个循环action，因此它们之间的内容可能会被展开多次，循环每次迭代的当前值对应当前的Items元素的值
const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number: {{.Number}}
User:   {{.User.Login}}
Title:  {{.Title | printf "%.64s"}}
Age:    {{.CreatedAt | daysAgo}} days
URL:    {{.HTMLURL}}
{{end}}`

// 在一个action中，|操作符表示将前一个表达式的结果作为后一个函数的输入，类似于UNIX中管道的概念。
// 在Title这一行的action中，第二个操作是一个printf函数，是一个基于fmt.Sprintf实现的内置函数，所有模板都可以直接使用。
// 对于Age部分，第二个动作是一个叫daysAgo的函数，通过time.Since函数将CreatedAt成员转换为过去的时间长度

//!-template

// 需要注意的是CreatedAt的参数类型是time.Time，并不是字符串
func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

// 创建并分析上面定义的模板templ
// 注意方法调用链的顺序：template.New先创建并返回一个模板；Funcs方法将daysAgo等自定义函数注册到模板中，并返回模板；最后调用Parse函数分析模板。

// 因为模板通常在编译时就测试好了，如果模板解析失败将是一个致命的错误。
// template.Must辅助函数可以简化这个致命错误的处理：它接受一个模板和一个error类型的参数，检测error是否为nil（如果不是nil则发出panic异常），然后返回传入的模板。我们将在5.9节再讨论这个话题。
var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	// 使用github.IssuesSearchResult作为输入源、os.Stdout作为输出源来执行模板
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

func noMust() {
	//!+parse
	report, err := template.New("report").
		Funcs(template.FuncMap{"daysAgo": daysAgo}).
		Parse(templ)
	if err != nil {
		log.Fatal(err)
	}
	//!-parse
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := report.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

/*
//!+output
$ go build gopl.io/ch4/issuesreport
$ ./issuesreport repo:golang/go is:open json decoder
13 issues:
----------------------------------------
Number: 5680
User:   eaigner
Title:  encoding/json: set key converter on en/decoder
Age:    750 days
----------------------------------------
Number: 6050
User:   gopherbot
Title:  encoding/json: provide tokenizer
Age:    695 days
----------------------------------------
...
//!-output
*/

```



### issueshtml.go

使用和text/template包相同的API和模板语言，但是增加了一个将字符串自动转义特性，这可以避免输入字符串和HTML、JavaScript、CSS或URL语法产生冲突的问题。还可以避免通过生成HTML注入的XSS等漏洞。

```go
// Issueshtml prints an HTML table of issues matching the search terms.
// See page 115.
package main

import (
	"html/template"
	"log"
	"os"

	github "gopher.run/go/src/ch4/18.github"
)

// 注意，html/template包已经自动将特殊字符转义(HTML实体转义)，因此我们依然可以看到正确的字面值。
// 如果我们使用text/template包的话，这2个issue将会产生错误，其中“&lt;”四个字符将会被当作小于字符“<”处理，同时“<link>”字符串将会被当作一个链接元素处理，它们都会导致HTML文档结构的改变，从而导致有未知的风险。
var issueList = template.Must(template.New("issuelist").Parse(`
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	if err := issueList.Execute(os.Stdout, result); err != nil {
		log.Fatal(err)
	}
}

/*
//!+output
$ go build gopl.io/ch4/issueshtml
$ ./issueshtml repo:golang/go commenter:gopherbot json encoder >issues.html
...
//!-output
*/

```



### autoescape.go

通过对信任的HTML字符串（非用户可控输入），使用template.HTML类型来抑制这种自动转义的行为：

```go
// Autoescape demonstrates automatic HTML escaping in html/template.
// See page 117.
package main

import (
	"html/template"
	"log"
	"os"
)

func main() {
	const templ = `<p>A: {{.A}}</p><p>B: {{.B}}</p>`
	t := template.Must(template.New("escape").Parse(templ))
	var data struct {
		A string        // untrusted plain text <p>A: &lt;b&gt;Hello!&lt;/b&gt; 污点变量，HTML实体转义
		B template.HTML // trusted HTML  </p><p>B: <b>Hello!</b></p> 用户不可控变量，安全，不转义
	}
	data.A = "<b>Hello!</b>" 
	data.B = "<b>Hello!</b>" 
	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatal(err)
	}
}
```



```shell
# 如果想了解更多的信息，请自己查看包文档
$ go doc text/template
$ go doc html/template
```

