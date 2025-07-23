// point study
package main

import (
	"fmt"
	"math"
	"os"
)

func main() {
	// 在局部变量地址&v被返回之后依然有效，因为指针p依然引用这个变量。
	//var p = f()
	//fmt.Println(p)
	//fmt.Println(f() == f()) // "false"
	i := 1
	fmt.Println(&i)
	urls := []string{"", ""}
	fmt.Println(&urls[0], &urls[1])
	f1()
	fmt.Println(math.Pi)
}

func f1() *int {
	v := 1
	err := g()
	if err != nil {
		return nil
	}
	return &v
}

func g() error {
	fname := "abc.txt"
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	f1()
	return nil
}

type Circle struct {
	radius float64
	Point
}
type Point struct {
	x, y int
}
