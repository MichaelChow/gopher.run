// Defer1 demonstrates a deferred call being invoked during a panic.
// See page 150.
package main

import "fmt"

// !+f
func main() {
	f(3)
}

func f(x int) {
	fmt.Printf("f(%d)\n", x+0/x)      // f(3)、f(2)、f(1)、f(0) panics
	defer fmt.Printf("defer %d\n", x) // defer 3、defer 2、defer 1，多条defer语句的执行顺序与声明顺序相反
	f(x - 1)
}

/*
//!+stdout
f(3)
f(2)
f(1)
defer 1
defer 2
defer 3
panic: runtime error: integer divide by zero

goroutine 1 [running]:
main.f(0x102734fe8?)
        /Users/xxx/3. go/gopher.run/src/ch5/16.defer1/defer.go:13 +0xec
main.f(0x1)
        /Users/xxx/3. go/gopher.run/src/ch5/16.defer1/defer.go:15 +0xcc
main.f(0x2)
        /Users/xxx/3. go/gopher.run/src/ch5/16.defer1/defer.go:15 +0xcc
main.f(0x3)
        /Users/xxx/3. go/gopher.run/src/ch5/16.defer1/defer.go:15 +0xcc
main.main()
        /Users/xxx/3. go/gopher.run/src/ch5/16.defer1/defer.go:9 +0x20
exit status 2
*/
