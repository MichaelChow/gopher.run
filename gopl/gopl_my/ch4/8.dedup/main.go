// Dedup prints only one instance of each line; duplicates are removed.
// See page 96.
package main

import (
	"bufio"
	"fmt"
	"os"
)

// Go语言中并没有提供set类型，但是map类型的键是不重复的，因此我们可以用map实现set的功能：
// Go程序员将这种忽略value的map当做一个字符串集合

// 读取多行输入，但只打印第一次出现的行
// 通过map来表示所有的输入行所对于的set集合，以确保已经出现过的行不会被重复打印
func main() {
	seen := make(map[string]bool) // a set of strings
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if !seen[line] {
			seen[line] = true
			fmt.Println(line)
		}
	}

	if err := input.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "dedup: %v\n", err)
		os.Exit(1)
	}

}

//!-
