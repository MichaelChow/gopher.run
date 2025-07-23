// Package tempconv performs Celsius and Fahrenheit temperature computations.
package tempconv

import "fmt"

func Example_one() {
	// 计算摄氏温度差

	// 打印 BoilingC（沸点）和 FreezingC（冰点）之间的摄氏温度差
	fmt.Printf("%g\n", BoilingC-FreezingC) // "100" °C
	// 将 BoilingC 转换为华氏温度
	boilingF := CToF(BoilingC)
	// 打印 boilingF（沸点的华氏温度）和 CToF(FreezingC)（冰点的华氏温度）之间的华氏温度差
	fmt.Printf("%g\n", boilingF-CToF(FreezingC)) // "180" °F

	// 尝试打印 boilingF（沸点的华氏温度）和 FreezingC（冰点的摄氏温度）之间的温度差，这会导致类型不匹配错误
	// fmt.Printf("%g\n", boilingF-FreezingC) // compile error: type mismatch
}

func Example_two() {
	//!+printf
	c := FToC(212.0)
	fmt.Println(c.String()) // 打印摄氏温度的字符串表示，"100°C"
	fmt.Printf("%v\n", c)   // %v，自动调用 String 方法。"100°C"; no need to call String explicitly
	fmt.Printf("%s\n", c)   // %s，自动调用 String 方法。"100°C"
	fmt.Println(c)          // 直接打印，自动调用 String 方法。"100°C"
	fmt.Printf("%g\n", c)   // %g，不会调用 String 方法。"100"; does not call String
	fmt.Println(float64(c)) // "100"; does not call String
}
