// Findlinks3 crawls the web, starting with the URLs on the command line.
// See page 139.
package main

import (
	"fmt"
	"log"
	"os"

	links "gopher.run/go/src/ch5/9.links"
)

// !+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
// 下面的函数实现了广度优先算法。
// 调用者需要输入一个函数值: f func(item string) []string 和一个初始的待访问列表。待访问列表中的每个元素被定义为string类型。
// 广度优先算法会为每个元素调用一次f。每次f执行完毕后，会返回一组待访问元素。这些元素会被加入到待访问列表中。当待访问列表中的所有元素都被访问后，breadthFirst函数运行结束。
// 为了避免同一个元素被访问两次，代码中维护了一个map。
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				// append的参数“f(item)...”，会将f返回的一组元素[]string 一个个添加到worklist中
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//!-breadthFirst

// !+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-crawl

// !+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	// go run findlinks.go https://www.taobao.com  爬淘宝
	// 将crawl作为参数传递给breadthFirst
	// 当所有发现的链接都已经被访问 或 电脑的内存耗尽时，程序运行结束
	breadthFirst(crawl, os.Args[1:])
}

//!-main
