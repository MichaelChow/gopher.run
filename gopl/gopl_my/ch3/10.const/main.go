package main

import (
	"fmt"
	"math"
)

// iota 是go语言的常量计数器，只能在常量的表达式中使用。
// iota 表示从0开始自动加1，所以Sunday=0，Monday=1，以此类推。
type Weekday int

// 它首先定义了一个Weekday命名类型，然后为一周的每天定义了一个常量，从周日0开始。
// 在其它编程语言中，这种类型一般被称为枚举类型。

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Flags uint

// 给一个无符号整数的最低5bit的每个bit指定一个名字
// 使用这些常量可以用于测试、设置或清除对应的bit位的值：
const (
	FlagUp           Flags = 1 << iota // 1 << 0，1 * 2^0
	FlagBroadcast                      // 1 << 1，1 * 2^1
	FlagLoopback                       // 1 << 2，1 * 2^2
	FlagPointToPoint                   // 1 << 3，1 * 2^3
	FlagMulticast                      // 1 << 4，1 * 2^4
)

// 每个常量都是1024的幂
const (
	_   = 1 << (10 * iota)
	KiB // 1024   , 1 << (10 * 1)
	MiB // 1048576 , 1 << (10 * 2)
	GiB // 1073741824, 1 << (10 * 3)
	TiB // 1099511627776, 1 << (10 * 4)   (exceeds 1 << 32)
	PiB // 1125899906842624, 1 << (10 * 5)
	EiB // 1152921504606846976, 1 << (10 * 6)
	ZiB // 1180591620717411303424, 1 << (10 * 7)    (exceeds 1 << 64) 溢出
	YiB // 1208925819614629174706176, 1 << (10 * 8) 溢出
)

var x float32 = math.Pi
var y float64 = math.Pi
var z complex128 = math.Pi

func main() {
	fmt.Println(Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday)       // 0 1 2 3 4 5 6
	fmt.Println(FlagUp, FlagBroadcast, FlagLoopback, FlagPointToPoint, FlagMulticast) // 1 2 4 8 16
	var i int64
	j := i << 63
	fmt.Println("j: ", j)
	// fmt.Println(KiB, MiB, GiB, TiB, PiB, EiB, ZiB, YiB)
}
