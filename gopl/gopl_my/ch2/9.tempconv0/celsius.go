// Package tempconv performs Celsius and Fahrenheit temperature computations.
// See page 39.
package tempconv

import "fmt"

type Celsius float64    // 摄氏温度类型 /ˈselsiəs/
type Fahrenheit float64 // 华氏温度类型 /'færən'haɪt/

const (
	AbsoluteZeroC Celsius = -273.15 // 绝对零度
	FreezingC     Celsius = 0       // 冰点
	BoilingC      Celsius = 100     // 沸点
)

// CToF 函数将摄氏温度转换为华氏温度
func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) }

// FToC 函数将华氏温度转换为摄氏温度
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) }

// String 方法将 Celsius 类型的值转换为字符串表示
func (c Celsius) String() string {
	// 使用 fmt.Sprintf 函数将 Celsius 值格式化为字符串，并添加 °C 后缀
	return fmt.Sprintf("%g°C", c)
}
