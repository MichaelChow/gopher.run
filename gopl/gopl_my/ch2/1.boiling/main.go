// Boiling prints the boiling point of water.
// See page 29.

package main

import "fmt"

// 包一级范围声明语句
const boilingF = 212.0

func main() {
	// 名字的作用域比较小，生命周期也比较短的函数内部声明的变量，使用简短的名字
	var f = boilingF
	var c = (f - 32) * 5 / 9
	fmt.Printf("boiling point = %g°F or %g°C\n", f, c)
}
