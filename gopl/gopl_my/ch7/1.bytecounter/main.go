// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// See page 173.
// Byte counter demonstrates an implementation of io.Writer that counts bytes.
package main

import (
	"fmt"
)

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}

func main() {
	var c ByteCounter
	_, _ = c.Write([]byte("hello"))
	fmt.Println(c) // "5", = len("hello")

	c = 0 // reset the counter
	var name = "Dolly"
	// 因为*ByteCounter满足io.Writer接口的约定，所以可以把它传入Fprintf函数中
	// Fprintf函数执行字符串格式化的过程不会去关注ByteCounter正确的累加结果的长度
	fmt.Fprintf(&c, "hello, %s", name)
	fmt.Println(c) // "12", = len("hello, Dolly")
}
