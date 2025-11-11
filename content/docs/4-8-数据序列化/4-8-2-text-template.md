---
title: "4.8.2 text/template"
date: 2025-08-05T00:51:00Z
draft: false
weight: 4008
---

# 4.8.2 text/template



> [https://pkg.go.dev/text/template](https://pkg.go.dev/text/template)



package template implements data-driven templates for generating textual output.

text/template 是 Go 标准库提供的数据驱动的文本模板引擎，用于生成动态文本内容。



模板是通过将其应用于数据结构来执行的。模板中的注释引用数据结构中的元素（通常是结构体的字段或映射的键），以控制执行并导出要显示的值。模板的执行会遍历结构，并将光标（用点'.'表示，称为"dot"）设置为执行过程中当前结构位置的值。



模板的输入文本是任何格式的 UTF-8 编码文本。"Actions"（数据评估或控制结构）由"{{"and"}}"分隔；所有动作外的文本都会原样复制到输出中。

- 默认情况下，当模板执行时，所有在操作符之间的文本都会原封不动地复制。
- 修剪空白字符：如果操作符的左定界符（默认为"{{"）紧跟着一个减号和空格，那么会从紧邻的前一个文本中删除所有尾随的空格。如果右定界符（"}}"）前面有空格和减号，会从紧邻的下一个文本中删除所有开头的空格。
- 空白字符的定义与 Go 中的定义相同：空格、水平制表符、回车和换行。


### **基本结构**

**1. Template 对象**

```go
type Template struct {
    name string
    *parse.Tree
    *common
    leftDelim  string
    rightDelim string
}
```

**2. 模板语法**

- **分隔符**: {{ 和 }}
- **变量**: {{.VariableName}}
- **管道**: {{.Value | function}}
- **控制结构**: {{if}}, {{range}}, {{with}}


为什么设计成{{.v}}，而不是更简洁的{v}? Go 语言"明确优于简洁"的设计哲学

```go
// 问题示例
template := `
用户信息:
姓名: {name}
邮箱: {email}
地址: {address}
价格: {price} 元
`

// 问题：
// 1. 普通文本中的大括号会被误解析
// 2. 难以区分变量和文本
// 3. 解析器复杂度增加
```

```go
// 优势示例
template := `
用户信息:
姓名: {{.name}}
邮箱: {{.email}}
地址: {{.address}}
价格: {{.price}} 元
`

// 优势：
// 1. 明确的语法边界
// 2. 易于解析和识别
// 3. 与主流模板引擎一致
// 4. 支持复杂语法扩展
```

**主流模板引擎对比：**

| 模板引擎 | 语言 | 变量语法 | 条件语法 | 循环语法 | 函数调用 | 
| --- | --- | --- | --- | --- | --- | 
| **Go text/template** | go | {{.variable}} | {{if .condition}} | {{range .items}} | {{.value \| function}} | 
| **Mustache** |   | {{variable}} | {{#condition}} | {{#items}} | 不支持 | 
| **Handlebars** | JavaScript | {{variable}} | {{#if condition}} | {{#each items}} | {{helper value}} | 
| **Django/Jinja2** | python | {{ variable }} | {% if condition %} | {% for item in items %} | {{ function(value) }} | 
| **EJS** | JavaScript | <%= variable %> | <% if (condition) { %> | <% for (item of items) { %> | <%= function(value) %> | 
| **Thymeleaf** |   | ${variable} | th:if="${condition}" | th:each="item : items" | ${#functions.function(value)} | 







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



