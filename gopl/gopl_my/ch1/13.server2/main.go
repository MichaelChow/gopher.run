// Server2 is an "echo" server that displays request parameters.
// See page 21.
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("start webserver http://localhost:8000 ... ")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the HTTP request.
func handler(w http.ResponseWriter, r *http.Request) {
	getHttpRequest(w, r)
}

// getHttpRequest 解析http请求
func getHttpRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s %s %s\n", r.Method, r.URL, r.Proto) // URI
	fmt.Fprintf(w, "Host: %s\n", r.Host)                   // Host
	for k, v := range r.Header {                           // Header
		fmt.Fprintf(w, "%s: %s\n", k, strings.Join(v, ""))
	}

	// fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	// 用if和ParseForm结合可以让代码更简洁，并且可以限制err这个变量的作用域
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	fmt.Fprintf(w, "\n")                           // 空行
	if PostForm := r.PostForm; len(PostForm) > 0 { // post form 表单

		seq := ""
		for k, v := range PostForm { //  POST_DATA
			fmt.Fprintf(w, "%s%s=%s", seq, k, strings.Join(v, ""))
			seq = "&"
		}
	} else { // 非post json body
		b, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Print(err)
		}
		fmt.Fprintf(w, "%s", string(b))
	}
}
