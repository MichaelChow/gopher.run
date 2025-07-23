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
