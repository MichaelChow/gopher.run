// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 123.

// Outline prints the outline of an HTML document tree.
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

// !+
func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outline: %v\n", err)
		os.Exit(1)
	}
	outline(nil, doc)
}

// 在函数outline中，我们通过递归的方式遍历整个HTML结点树，并输出树的结构。
// 在outline内部，每遇到一个HTML元素标签，就将其入栈，并输出
func outline(stack []string, n *html.Node) {
	if n.Type == html.ElementNode {
		// 有一点值得注意：这里的outline有入栈操作，但没有相对应的出栈操作。
		stack = append(stack, n.Data) // push tag
		fmt.Println(stack)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// 当outline调用自身时，被调用者接收的是stack的拷贝。
		// 被调用者对stack的元素追加操作，修改的是stack的拷贝，其可能会修改slice底层的数组甚至是申请一块新的内存空间进行扩容；
		// 但这个过程并不会修改调用方的stack。因此当函数返回时，调用方的stack与其调用自身之前完全一致。（此处反倒利用了这个特性）
		outline(stack, c)
	}
}

// 大部分HTML页面只需几层递归就能被处理，但仍然有些页面需要深层次的递归。
// 大部分编程语言使用固定大小的函数调用栈，常见的大小从64KB到2MB不等。固定大小栈会限制递归的深度，当你用递归处理大量数据时，需要避免栈溢出；除此之外，还会导致安全性问题。
// 与此相反，Go语言使用可变栈，栈的大小按需增加（初始时很小）。这使得我们使用递归时不必考虑溢出和安全问题。

//!-
// $ ../../ch1/10.fetch/fetch https://www.taobao.com | ./outline
// [html]
// [html head]
// [html head meta]
// [html head title]
// [html head meta]
// [html head link]
// [html head script]
// [html head script]
// [html body div div div div div div style]
// ...
