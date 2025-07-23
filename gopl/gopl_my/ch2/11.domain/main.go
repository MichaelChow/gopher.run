package main

import (
	"log"
	"os"
)

var cwd string

func init() {
	cwd, err := os.Getwd() // NOTE: wrong!
	if err != nil {
		log.Fatalf("os.Getwd failed: %v\n", err)
	}
	// 由于当前的编译器会检测到局部声明的cwd并没有使用，然后报告这可能是一个错误，但是这种检测并不可靠。
	// 因为一些小的代码变更（如增加一个局部cwd的打印语句），就导致了这种检测失效(包级别词法域的cwd没有被使用，但未被编译器检测到)。
	// 全局的cwd变量依然是没有被正确初始化的，而且看似正常的日志输出更是让这个BUG更加隐晦。
	log.Printf("Working directory = %s", cwd)
}

func main() {
}

// 2024/12/18 19:26:57 Working directory = /Users/xxx/gopher.run/src/ch2/11.domain
