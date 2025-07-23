package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Printf("int8: %d \nint16: %d\nint32: %d\nint64: %d\n", math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64)
	fmt.Printf("float32: %f \nfloat64: %f\n", math.MaxFloat32, math.MaxFloat64)

	fmt.Printf("%.20f\n", 0.1+0.2) // 0.29999999999999998890
	var f float32 = 1 << 24        // 16777216
	fmt.Println(f == f+1)          // "true"!
	fmt.Printf("%g", f)

	// 浮点数的字面值可以直接写小数部分：
	const e = 2.71828 // (approximately)

	for x := 0; x < 8; x++ {
		fmt.Printf("x = %d e^x = %8.3f\n", x, math.Exp(float64(x)))
	}
}
