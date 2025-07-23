// Fetch prints the content found at each specified URL.
// See page 16.

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		res, err := http.Get(arg)
		// 检查http.Get的错误
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		// res.Body 可读的服务器响应流是一个 io.Reader 类型，使用io.ReadAll读取到字节切片
		b, err := io.ReadAll(res.Body)
		// 关闭Body，释放相关资源，防止资源泄露
		res.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s", b)
	}
}
