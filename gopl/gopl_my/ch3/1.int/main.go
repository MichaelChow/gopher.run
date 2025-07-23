// type test
package main

import "fmt"

func main() {
	bytes := []byte("hello")
	s := "hello"
	fmt.Println(bytes)        // [104 101 108 108 111]  h的ascii码的十进制表示为104，e的ascii码的十进制表示为101，以此类推
	fmt.Printf("%b\n", bytes) // [1101000 1100101 1101100 1101100 1101111]
	fmt.Printf("%b\n", s)     // %!b(string=hello) wrong type

	var u uint8 = 255
	fmt.Println(u, u+1, u+2, u+3, u+4)                              // "255 0 1"
	fmt.Printf("%08b %08b %08b %08b %08b\n", u, u+1, u+2, u+3, u+4) // 11111111 10000000(截断存后8位 0) 10000001(截断存后8位 1) 10000010(截断存后8位 10)  10000011(截断存后8位 11)

	// 最高位为 0 表示正数，最高位为 1 表示负数。
	var i int8 = 127
	// fmt.Println(i, i+1, i+2, i+3, i+4)       // "127 -128 1"
	fmt.Printf("%08b %[1]v\n", i+1) // 01111111 10000000（先计算，后截断存后8位。最高位符号位为1，变成负数） -1111111 -1111110 -1111101

	// 128的二进制表示是10000000
	// -128的二进制表示法：取补码：取反后得到01111111，再加1得到10000000。
}
