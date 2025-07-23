// Echo4 prints its command-line arguments.
// See page 33.
package main

import (
	"flag"
	"fmt"
	"strings"
)

// 定义命令行参数名、默认值、和描述信息。
var n = flag.Bool("n", false, "omit trailing newline") // *bool类型的指针，默认为false不换行（终端会打印出一个%作为无换行符的标识）
var sep = flag.String("s", " ", "separator")           // *string类型的指针，默认为空格

func main() {
	// 解析命令行参数，更新每个标志参数对应变量的值（之前是默认值）。
	// 解析命令行参数时遇到错误，默认将打印相关的提示信息，然后调用os.Exit(2)终止程序。
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), *sep)) // 打印命令行参数，*sep 是一个字符串指针，它的值是通过命令行参数 -s 指定的分隔符。

	if !*n { // 如果命令行参数 -n 没有指定，就打印一个换行符
		fmt.Println()
	}
}
