// Defer2 demonstrates a deferred call to runtime.Stack during a panic.
// See page 151.
package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// 通过在main函数中延迟调用printStack输出堆栈信息
	defer printStack()
	f(3)
}

func printStack() {
	var buf [4096]byte
	// 为了方便诊断问题，runtime包允许程序员输出堆栈信息
	n := runtime.Stack(buf[:], false)
	os.Stdout.Write(buf[:n])
}

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x)      // f(3)、f(2)、f(1)、f(0) panics
	defer fmt.Printf("defer %d\n", x) // defer 3、defer 2、defer 1，多条defer语句的执行顺序与声明顺序相反
	f(x - 1)
}

/*
f(3)
f(2)
f(1)
defer 1
defer 2
defer 3
goroutine 1 [running]:
main.printStack()
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:18 +0x38
panic({0x104e55440?, 0x104ee5960?})
        /usr/local/go/src/runtime/panic.go:785 +0x124
main.f(0x104e69078?)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:23 +0xec
将panic机制类比其他语言异常机制的读者可能会惊讶，runtime.Stack为何能输出已经被释放函数的信息？在Go的panic机制中，延迟函数的调用在释放堆栈信息之前。
main.f(0x1)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.f(0x2)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.f(0x3)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.main()
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:13 +0x3c
panic: runtime error: integer divide by zero

goroutine 1 [running]:
main.f(0x104e69078?)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:23 +0xec
main.f(0x1)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.f(0x2)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.f(0x3)
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:25 +0xcc
main.main()
        /Users/xxx/3. go/gopher.run/src/ch5/17.defer2/defer.go:13 +0x3c
exit status 2
*/
