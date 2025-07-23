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
// 最简单的格式化，使用Printf是完全足够的
// 但是有时候会需要复杂的打印格式，这时候一般需要将格式化代码分离出来以便更安全地修改
// 这些功能是由text/template和html/template等模板包提供的，它们提供了一个将变量值填充到一个文本或HTML格式的模板的机制

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
