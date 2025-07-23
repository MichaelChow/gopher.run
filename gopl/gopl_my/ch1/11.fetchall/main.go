// Fetchall fetches URLs in parallel and reports their times and sizes.
// See page 17.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()

	ch := make(chan string) // 用make函数 创建了一个传递string类型参数的channel /ˈtʃænl/；
	for _, url := range os.Args[1:] {
		// 用go关键字来创建一个goroutine, 并发的执行fetch函数
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // 从channel里接收一个值
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		// <- 是用于通道（channel）操作的箭头运算符，用于发送和接收通道数据。
		ch <- fmt.Sprint(err) //  往channel里发送一个值
		return
	}
	// Body内容拷贝到ioutil.Discard输出流中，即丢弃。因为这里我们不需要Body内容。
	nbytes, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url) // 往channel里发送一个值
}
