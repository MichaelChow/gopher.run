// Server1 is a minimal "echo" server.
// See page 19.
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler) // 每一个访问/根路径的请求，调用绑定的handler函数
	fmt.Println("start web serve  http://localhost:8000 ...")
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // 监听8000端口，nil 表示服务器使用默认的路由器。log.Fatal 函数会在服务器启动失败时输出错误信息并结束程序的执行。
}

// handler echoes the Path component of the requested URL.
func handler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf 将格式化的字符串输出到指定的 io.Writer接口实现的对象
	fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path) // 向客户端发送响应,%q 是一个动词，用于将字符串值用双引号包围.
}
