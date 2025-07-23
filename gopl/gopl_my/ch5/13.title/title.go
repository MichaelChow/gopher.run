// Title1 prints the title of an HTML document specified by a URL.
// See page 144.
package main

/*
//!+output
$ go build gopl.io/ch5/title1
$ ./title1 http://gopl.io
The Go Programming Language
$ ./title1 https://golang.org/doc/effective_go.html
Effective Go - The Go Programming Language
$ ./title1 https://golang.org/doc/gopher/frontpage.png
title: https://golang.org/doc/gopher/frontpage.png
    has type image/png, not text/html
//!-output
*/

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
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

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	// 而Go语言独有的defer机制可以让代码变得简单, 只需要在调用普通函数或方法前加上关键字defer
	// 当执行到该条语句时，函数和参数表达式先得到计算
	// 但直到包含该defer语句所在的函数执行完毕时，defer后的函数才会被执行，包括通过return正常结束、由于panic导致的异常结束等所有函数出口
	// 一个函数中的多条defer语句的执行顺序与声明顺序相反
	// defer语句常被用于成对的操作语句写在一块，如打开/关闭、连接/断开连接、加锁/释放锁
	defer resp.Body.Close()

	// Check Content-Type is HTML (e.g., "text/html; charset=utf-8").
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		// 重复代码1
		// resp.Body.Close()
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	// 虽然Go的垃圾回收机制会回收不被使用的内存，但是这不包括操作系统层面的资源，比如打开的文件、网络连接。
	// 因此我们必须显式的释放这些资源
	// 重复代码2
	// 为了确保title在所有执行路径下（即使函数运行失败）都关闭了网络连接
	// 随着函数变得复杂，需要处理的错误也变多，维护清理逻辑变得越来越困难
	// resp.Body.Close()
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			fmt.Println(n.FirstChild.Data)
		}
	}
	forEachNode(doc, visitNode, nil)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}
