package main

import "fmt"

var q [3]int = [3]int{1, 2, 3}

func main() {
	p := [...]int{1, 2, 3}
	fmt.Println(q, p)

}
