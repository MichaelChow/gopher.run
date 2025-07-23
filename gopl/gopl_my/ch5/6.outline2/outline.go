// Outline prints the outline of an HTML document tree.]
// See page 133.
package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		outline(url)
	}
}

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

	//!+call
	forEachNode(doc, startElement, endElement)
	//!-call

	return nil
}

// 5.2节的findLinks函数使用了辅助函数visit，遍历和操作了HTML页面的所有结点。
// 使用函数值，我们可以将遍历结点的逻辑和操作结点的逻辑分离，使得我们可以复用遍历的逻辑，从而对结点进行不同的操作。
// 该函数接收2个函数作为参数，分别在结点的孩子被访问前和访问后调用。这样的设计给调用者更大的灵活性。

// forEachNode针对每个结点x，都会调用pre(x)和post(x)。pre和post都是可选的。
// 遍历孩子结点之前，pre被调用；遍历孩子结点之后，post被调用
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

//!-forEachNode

// !+startend
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

//!-startend
//  $ ./outline https://taobao.com
// <html>
//   <head>
//     <meta>
//     </meta>
//     <meta>
//     </meta>
//     <meta>
//     </meta>
// ...
