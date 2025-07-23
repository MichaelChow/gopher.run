// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 179.

// The sleep program sleeps for a specified period of time.
package main

import (
	"flag"
	"fmt"
	"time"
)

// flag.Duration函数创建一个time.Duration类型的标记变量，并且允许用户通过多种用户友好的方式来设置这个变量的大小，这种方式还包括和String方法相同的符号排版形式。这种对称设计使得用户交互良好。
var period = flag.Duration("period", 1*time.Second, "sleep period")

func main() {
	flag.Parse()
	// fmt包调用time.Duration的String方法打印这个时间周期是以用户友好的注解方式，而不是一个纳秒数字
	fmt.Printf("Sleeping for %v...", *period)
	time.Sleep(*period)
	fmt.Println()
}

//$ go run sleep.go -period 1m3s
//$ Sleeping for 1m3s...
