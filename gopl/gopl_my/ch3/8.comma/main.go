// Comma prints its argument numbers with a comma at each power of 1000.
// See page 73.
//
// Example:
//
//	$ go build gopl.io/ch3/comma
//	$ ./comma 1 12 123 1234 1234567890
//	1
//	12
//	123
//	1,234
//	1,234,567,890
package main

import (
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", comma(os.Args[i]))
	}
}

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	fmt.Println(s)
	n := len(s)
	if n <= 3 {
		return s
	}
	// 在最后三个字符前的位置插入逗号，拼接前面部分；
	// 前面部分通过递归调用自身comma来得出前面的子串；（很巧妙）
	return comma(s[:n-3]) + "," + s[n-3:]
}
