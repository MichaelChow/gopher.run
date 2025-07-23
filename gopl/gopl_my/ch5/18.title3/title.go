// Title3 prints the title of an HTML document specified by a URL.
// See page 153.
// title函数的变形
package main

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

// !+
// soleTitle returns the text of the first non-empty title element
// in doc, and an error if there was not exactly one.
func soleTitle(doc *html.Node) (title string, err error) {
	type bailout struct{}
	// deferred函数调用recover，并检查panic value
	defer func() {
		switch p := recover(); p {
		case nil:
			// no panic
		// 当panic value是bailout{}类型时（预期发生的错误），deferred函数生成一个error返回给调用者
		// 请注意：对可预期的错误采用了panic（应该使用error），不符合Go语言风格，这里仅是为了向读者演示这种机制。
		case bailout{}:
			// "expected" panic
			err = fmt.Errorf("multiple title elements")
		// 当panic value是其他non-nil值时，表示发生了未知的panic异常，deferred函数将调用panic函数并将当前的panic value作为参数传入。等同于recover没有做任何操作
		default:
			panic(p) // unexpected panic; carry on panicking
		}
	}()

	// Bail out of recursion if we find more than one non-empty title.
	forEachNode(doc, func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" &&
			n.FirstChild != nil {
			// 如果检测到有多个<title>：调用panic，阻止函数继续递归，并将特殊类型bailout作为panic的参数
			if title != "" {
				panic(bailout{}) // multiple title elements
			}
			title = n.FirstChild.Data
		}
	}, nil)
	if title == "" {
		return "", fmt.Errorf("no title element")
	}
	return title, nil
}

//!-

func title(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// Check Content-Type is HTML (e.g., "text/html; charset=utf-8").
	ct := resp.Header.Get("Content-Type")
	if ct != "text/html" && !strings.HasPrefix(ct, "text/html;") {
		resp.Body.Close()
		return fmt.Errorf("%s has type %s, not text/html", url, ct)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	title, err := soleTitle(doc)
	if err != nil {
		return err
	}
	fmt.Println(title)
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if err := title(arg); err != nil {
			fmt.Fprintf(os.Stderr, "title: %v\n", err)
		}
	}
}
