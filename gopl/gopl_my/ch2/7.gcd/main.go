// gcd, Greatest Common Divisor 最大公约数，欧几里德的GCD是最早的非平凡算法
package main

import "fmt"

func main() {
	fmt.Println(gcd(12, 18))
}

// 同 var变量申明，类型在后
func gcd(x, y int) int {
	for y != 0 {
		x, y = y, x%y
	}
	return x
}
