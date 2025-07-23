// incr study
package main

import "fmt"

func main() {
	v := 1
	incr(&v)              // v = 2,return 2
	fmt.Println(incr(&v)) // v = 3 return 3
}

func incr(p *int) int {
	*p++ // 非常重要：只是增加p指向的变量的值，并不改变p指针！！！
	return *p
}
