// The trace program uses defer to add entry/exit diagnostics to a function.
// See page 146.
package main

import (
	"fmt"
	"log"
	"time"
)

// !+main
// 调试复杂程序时，defer机制也常被用于记录何时进入和退出函数
func bigSlowOperation() {
	// 函数值会在bigSlowOperation退出时被调用
	// 只通过一条语句就能控制函数的入口和所有的出口，甚至可以记录函数的运行时间
	// 注意（很微妙）：defer语句后有圆括号，否则不会执行函数值（本该在退出时执行的，永远不会被执行），而本该在进入时执行的操作会在退出时执行 // 2025/01/02 19:34:25 enter bigSlowOperation
	defer trace("bigSlowOperation")() // don't forget the extra parentheses
	// ...lots of work...
	time.Sleep(3 * time.Second) // simulate slow operation by sleeping
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg) // 2025/01/02 19:33:45 enter bigSlowOperation
	// 返回一个函数值
	return func() { log.Printf("exit %s (%s)", msg, time.Since(start)) } // 2025/01/02 19:33:55 exit bigSlowOperation (10.00141025s)
}

func triple(x int) (result int) {
	// 1. 调用匿名函数/函数值（有圆括号），先计算表达式 result = 4
	defer func() { result += x }()
	// 2. 先计算表达式（调用double(x)）
	// 5. 执行triple函数的defer：result = 8 + 4
	return double(x)
}

func double(x int) (result int) {
	// 3. 调用匿名函数/函数值（有圆括号），先计算表达式 x = 4
	defer func() { fmt.Printf("double(%d) = %d\n", x, result) }()
	// 4. 先计算表达式 result = 8
	// 5. 执行double函数的defer：double(4) = 8
	return x + x
}

//!-main

func main() {
	bigSlowOperation()
	_ = double(4)          // double(4) = 8
	fmt.Println(triple(4)) // double(4) = 8  12
}

/*
!+output
$ go build gopl.io/ch5/trace
$ ./trace
2015/11/18 09:53:26 enter bigSlowOperation
2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
!-output
*/
