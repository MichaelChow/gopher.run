// The sum program demonstrates a variadic function.
// See page 142.
package main

import (
	"fmt"
	"os"
)

// sum可以接收任意数量的int型参数
func sum(vals ...int) int {
	// 在函数体中，vals被看作是类型为[] int的切片
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func main() {
	//!+main
	fmt.Println(sum())           //  "0"
	fmt.Println(sum(3))          //  "3"
	fmt.Println(sum(1, 2, 3, 4)) //  "10"
	// 在上面的代码中，调用者隐式的创建一个数组，并将原始参数复制到数组中。再把数组的一个切片作为参数传给被调用函数。
	//!-main

	//!+slice
	values := []int{1, 2, 3, 4}
	// 如果原始参数已经是切片类型，只需在最后一个参数后加上省略符，即可将切片的元素进行传递sum函数
	fmt.Println(sum(values...)) // "10"
	//!-slice

	// 但实际上，可变参数函数和以切片作为参数的函数是不同的函数类型
	fmt.Printf("%T\n", f) // func(...int)
	fmt.Printf("%T\n", g) // func([]int)

	linenum, name := 12, "count"
	errorf(linenum, "undefined: %s", name) // "Line 12: undefined: count"

}

func f(...int) {}
func g([]int)  {}

// 可变参数函数经常被用于格式化字符串
// 下面的errorf函数构造了一个以行号开头的，经过格式化的错误信息
// 函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串
func errorf(linenum int, format string, args ...interface{}) {
	// interface{}表示函数的最后一个参数可以接收任意类型，我们会在第7章详细介绍。
	fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}
