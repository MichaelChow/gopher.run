// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 130.

// The wait program waits for an HTTP server to start responding.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// !+
// WaitForServer attempts to contact the server of a URL.
// It tries for one minute using exponential back-off.
// It reports an error if all attempts fail.
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	deadline := time.Now().Add(timeout)
	// 处理错误的第二种策略：如果错误的发生是偶然性的，或由不可预知的问题导致的，重试
	// 需要限制重试的时间间隔或重试的次数，防止无限制的重试
	for tries := 0; time.Now().Before(deadline); tries++ {
		_, err := http.Head(url)
		if err == nil {
			return nil // success
		}
		log.Printf("server not responding (%s); retrying...", err)
		time.Sleep(time.Second << uint(tries)) // exponential back-off
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}

//!-

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: wait url\n")
		os.Exit(1)
	}
	url := os.Args[1]
	//!+main
	// (In function main.)
	// 如果错误发生后，程序无法继续运行，我们就可以采用处理错误的第三种策略：输出错误信息并结束程序。
	// 需要注意的是，这种策略只应在main中执行。
	// 对库函数而言，应仅向上传播错误，除非该错误意味着程序内部包含不一致性，即遇到了bug，才能在库函数中结束程序。
	if err := WaitForServer(url); err != nil {
		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		os.Exit(1)
	}
	// 调用log.Fatalf可以更简洁的代码达到与上文相同的效果。log中的所有函数，都默认会在错误信息之前输出时间信息。
	if err := WaitForServer(url); err != nil { // 等价代码
		// 长时间运行的服务器常采用默认的时间格式，而交互式工具很少采用包含如此多信息的格式
		log.Fatalf("Site is down: %v\n", err)
		// 2006/01/02 15:04:05 Site is down: no such domain:
		// bad.gopl.io

		// 我们可以设置log的前缀信息屏蔽时间信息，一般而言，前缀信息会被设置成命令名
		log.SetPrefix("wait: ")
		log.SetFlags(0)
	}

	// 第四种策略：有时，我们只需要输出错误信息就足够了，不需要中断程序的运行。我们可以通过log包提供函数，或者标准错误流输出错误信息。
	// log包中的所有函数会为没有换行符的字符串增加换行符
	if err := WaitForServer(url); err != nil {
		log.Printf("ping failed: %v; networking disabled", err)
		// fmt.Fprintf(os.Stderr, "ping failed: %v; networking disabled\n", err)
	}

	// 第五种策略：我们可以直接忽略掉错误
	dir, err := os.MkdirTemp("", "scratch")
	if err != nil {
		// return fmt.Errorf("failed to create temp dir: %v", err)
	}

	// ...use temp dir…
	os.RemoveAll(dir) // ignore errors; $TMPDIR is cleaned periodically
	// 尽管os.RemoveAll会失败，但上面的例子并没有做错误处理。这是因为操作系统会定期的清理临时目录。正因如此，虽然程序没有处理错误，但程序的逻辑不会因此受到影响。

}
