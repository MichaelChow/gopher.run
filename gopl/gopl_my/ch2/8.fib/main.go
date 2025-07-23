// fib,Fibonacci，计算斐波纳契数列的第N个数。

package main

import "fmt"

func main() {
	fmt.Println(fib(10))
}

// 参数列表（入参）、返回值列表（出参）
func fib(n int) (int, []int) {
	fib := []int{}
	x, y := 0, 1
	// 初始值第一个、第二个数都是1，后续每个数都是前两个数之和
	for i := 0; i < n; i++ {
		x, y = y, x+y
		fib = append(fib, x)
	}
	return x, fib
}
