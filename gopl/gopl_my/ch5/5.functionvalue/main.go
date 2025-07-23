package main

import (
	"fmt"
	"strings"
)

func square(n int) int     { return n * n }
func negative(n int) int   { return -n }
func product(m, n int) int { return m * n }

func add1(r rune) rune { return r + 1 }

func main() {
	f := square
	fmt.Println(f(3)) // "9"

	f = negative
	fmt.Println(f(3)) // "-3"

	// f = product // cannot use product (value of type func(m int, n int) int) as func(n int) int value in assignmen

	// 函数类型的零值是nil。调用值为nil的函数值会引起panic错误：
	var f1 func(int) int
	// f1(3) // f1为nil，引发panic错误。 n. 恐慌 panic: runtime error: invalid memory address or nil pointer dereference

	// 函数值可以与nil比较：
	// 但是函数值之间是不可比较的，也不能用函数值作为map的key。
	if f1 != nil {
		f1(3)
	}

	// 函数值使得我们不仅仅可以通过数据来参数化函数，亦可通过行为
	// strings.Map对字符串中的每个字符调用add1函数，并将每个add1函数的返回值组成一个新的字符串返回给调用者
	fmt.Println(strings.Map(add1, "HAL-9000")) // "IBM.:111"
	fmt.Println(strings.Map(add1, "VMS"))      // "WNT"
	fmt.Println(strings.Map(add1, "Admix"))    // "Benjy"

}
