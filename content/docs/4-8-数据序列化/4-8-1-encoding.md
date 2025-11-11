---
title: "4.8.1 encoding"
date: 2025-08-16T10:09:00Z
draft: false
weight: 4008
---

# 4.8.1 encoding

# encoding/json

**JavaScript对象表示法（JSON）**是一种用于发送和接收结构化信息的标准协议，由于其简洁性、可读性和流行程度等原因，JSON是应用最广泛的一个。类似协议还有、Google的Protocol Buffers、XML、ASN.1。

Go语言对于这些标准格式的编码和解码都有良好的支持，由标准库中的**encoding/json**、encoding/xml、encoding/asn1等包提供支持（译注：Protocol Buffers的支持由 [github.com/golang/protobuf](http://github.com/golang/protobuf) 包提供），并且这类包都有着相似的API接口。

JSON是对JavaScript中各种类型的值——字符串、数字、布尔值和对象——Unicode本文编码。它可以用有效可读的方式表示基础数据类型和聚合数据类型(数组、slice、结构体和map)。

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
type Movie struct {
	Title string
	Year int `json:"released"`
	// 额外的omitempty选项，表示当Go语言结构体成员为空或零值时不生成该JSON对象（这里false为零值）。
	// 即Marshal的json串中"Casablanca"由于其为零值fasle，所以没有Color字段
	Color  bool `json:"color,omitempty"`
	Actors []string
}
```

**结构体成员Tag**：在结构体声明中，Year和Color成员后面的字符串面值。**通过reflect.StructTag.Get实现**

- 结构体的成员Tag 可以是任意的字符串面值，但是通常是一系列用空格分隔的key:"value"键值对序列
- 由于值中含有双引号字符，因此成员Tag一般用原生字符串面值（即反引号`包裹）的形式书写，这样不再需要转义字符来表示双引号
- "json"开头键名对应的值，用于控制**json.Marshal json.UnMarshal**的行为，并且encoding/...下面其它的包也遵循这个约定
- 值的第一部分用于指定JSON对象的名字，如将Go语言中的**TotalCount**成员对应到JSON中的**total_count对象，**Year名字的成员在编码后变成了released，还有Color成员编码后变成了小写字母开头的color
- 在编码时，默认使用Go语言结构体的成员名字作为JSON的对象（通过reflect反射技术）。**只有导出的结构体成员才会被编码**，这也就是我们为什么选择用大写字母开头的成员名称。
- 一个结构体成员Tag是和在**编译阶段关联**到该成员的**元信息字符串**
- **json.Marshal**的逆操作是**json.UnMarshal**，对应将JSON数据解码为Go语言的数据结构，类似Java中的将json串反序列化为对象。


```go
// Movie prints Movies as JSON.
// See page 108.
// 程序负责收集各种电影评论并提供反馈功
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

		data, err := json.MarshalIndent(movies, "", "    ")
		if err != nil {
			log.Fatalf("JSON marshaling failed: %s", err)
		}
		fmt.Printf("%s\n", data)

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



