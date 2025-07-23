// The squares program demonstrates a function value with state.
// See page 135.
package main

import "fmt"

// squares返回一个匿名函数，其函数的类型（即函数的签名）为: func() int
func squares() func() int {
	// 创建一个局部变量x
	var x int
	// 并返回另一个匿名函数，其函数的类型（即函数的签名）同样为: func() int
	return func() int {
		// 先将x的值加1
		x++
		// 再返回x的平方
		return x * x
	}
}

// squares的例子证明，函数值不仅仅是一串代码，还记录了局部变量x的状态
// 在squares中定义的匿名内部函数可以访问和更新squares中的局部变量x，这意味着匿名函数和squares中，存在变量引用
// 这就是函数值属于引用类型和函数值不可比较的原因
// Go使用闭包（closures）技术实现函数值，Go程序员也把函数值叫做闭包。

func main() {
	f := squares()
	fmt.Printf("%T\n", squares)     // func() func() int
	fmt.Printf("%T\n", squares())   // func() int
	fmt.Printf("%T\n", squares()()) // int
	fmt.Println(f())                // x=1 1
	fmt.Println(f())                // x=2 4
	fmt.Println(f())                // x=3 9
	fmt.Println(f())                // x=4 16
	// 通过这个例子，我们看到变量的生命周期不由它的作用域决定：squares返回后，变量x仍然隐式的存在于f中。
}
