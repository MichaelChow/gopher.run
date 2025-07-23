// Package tempconv performs Celsius and Fahrenheit temperature computations.
// See page 180.
package tempconv

import (
	"flag"
	"fmt"
)

type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

/*
//!+flagvalue
package flag

// Value is the interface to the value stored in a flag.
type Value interface {
	String() string
	Set(string) error
}
//!-flagvalue
*/

// !+celsiusFlag
// *celsiusFlag satisfies the flag.Value interface.
type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

// CelsiusFlag defines a Celsius flag with the specified name,
// default value, and usage, and returns the address of the flag variable.
// The flag argument must have a quantity and a unit, e.g., "100C".
// CelsiusFlag函数将所有逻辑都封装在一起
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value} // Celsius字段是一个会通过Set方法在标记处理的过程中更新的变量
	// 调用Var方法将标记加入应用的命令行标记集合中，有异常复杂命令行接口的全局变量flag.CommandLine.Programs可能有几个这个类型的变量
	// 调用Var方法将一个*celsiusFlag参数赋值给一个flag.Value参数，导致编译器去检查*celsiusFlag是否有必须的方法
	flag.CommandLine.Var(&f, name, usage)
	// 它返回一个内嵌在celsiusFlag变量f中的Celsius指针给调用者
	return &f.Celsius
}
