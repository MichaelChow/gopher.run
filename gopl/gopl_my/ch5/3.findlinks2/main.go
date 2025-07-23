// Findlinks2 does an HTTP GET on each URL, parses the
// See page 125.
// result as HTML, and prints the links within it.
//
// Usage:
//
//	findlinks url ...
package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

// visit appends to links each link found in n, and returns the result.
func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func main() {
	for _, url := range os.Args[1:] {
		links, err := findLinks(url) // 等价代码1
		log.Println(links, err)      // 等价代码1
		// 将一个返回多参数的函数调用作为该函数的参数
		// 这很少出现在实际生产代码中，但这个特性在debug时很方便，我们只需要一条语句就可以输出所有的返回值
		log.Println(findLinks(url)) // 等价代码2

		if err != nil {
			fmt.Fprintf(os.Stderr, "findlinks2: %v\n", err)
			continue
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}

}

// findLinks performs an HTTP GET request for url, parses the
// response as HTML, and extracts and returns the links.
// 在Go中，一个函数可以返回多个值，这里返回链接列表和错误信息

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	// 虽然Go的垃圾回收机制会回收不被使用的内存，但是这不包括操作系统层面的资源，比如打开的文件、网络连接。
	// 因此我们必须显式的释放这些资源
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	return visit(nil, doc), nil
}

// 准确的变量名可以传达函数返回值的含义。尤其在返回值的类型都相同时：
// 但你也不必为每一个返回值都取一个适当的名字，如按照惯例，函数的最后一个bool类型的返回值表示函数是否运行成功，error类型的返回值代表函数的错误信息，它们都无需再解释
func Size(rect image.Rectangle) (width, height int)
func Split(path string) (dir, file string)
func HourMinSec(t time.Time) (hour, minute, second int)

// func countWordsAndImages(n *html.Node) (words, images int) { /* ... */ }

// CountWordsAndImages does an HTTP GET request for the HTML
// document url and returns the number of words and images in it.

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
		//  return 0,0,err（Go会将返回值 words和images在函数体的开始处，根据它们的类型，将其初始化为0） // 等价代码
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
		//  return 0,0,err（Go会将返回值 words和images在函数体的开始处，根据它们的类型，将其初始化为0） // 等价代码
	}
	fmt.Println(doc)
	// words, images = countWordsAndImages(doc)
	// 如果一个函数所有的返回值都有显式的变量名，那么该函数的return语句可以省略操作数。按照返回值列表的次序，返回所有的返回值
	// 这称之为bare return。/ber/ adj. 光秃秃的, 无遮蔽的 vt. 使赤裸, 使露出, 使暴露
	// 当一个函数有多处return语句以及许多返回值时，bare return 可以减少代码的重复，但是使得代码难以被理解。不宜过度使用bare return。
	return
	// return words, images, err // 等价代码
}
