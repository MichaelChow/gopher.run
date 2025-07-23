// Dup3 prints the count and text of lines that
// See page 12.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	// range 是GO的25个关键字之一，不同于Python中的range()内置函数形式
	for _, filename := range os.Args[1:] {
		// ioutil.ReadFile is deprecated: As of Go 1.16, this function simply calls [os.ReadFile]
		// 一次性把全部输入数据读到内存中
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		// 把字符串按换行符切割成行的切片
		// 与字节切片（byte slice，类似java的byte[]）转成string ，然后拼接成字符串。
		// 相反： strings.Join(os.Args[1:], " ")
		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%s: %d\n", line, n)
			}
		}
	}
}
