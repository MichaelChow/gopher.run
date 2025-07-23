// new study
package main

import "fmt"

func main() {
	// 特殊情况：如果类型的大小都为0，地址可能会相同
	p, q := new([0]int), new([0]int)
	fmt.Println(p == q) // false
}
