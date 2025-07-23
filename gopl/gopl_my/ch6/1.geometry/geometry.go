// Package geometry defines simple types for plane geometry.
// See page 156.
package geometry

import "math"

type Point struct{ X, Y float64 }

// Distance 函数，为包级别的函数 geometry.Distance
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// Distance 方法，附加到Point类型上，为Point类下声明的Point.Distance方法
// 在函数声明时，在函数名前附加一个参数（即是一个方法，相当于为这种类型定义了一个独占的方法），将该函数附加到这种类型上
// 附加的参数叫方法的接收器（receiver），Go中可任意的选择接收器的名字（不限制为this或者self），通常统一简写为类型的首字母（由于高频使用）
// Go中能给任意命名类型（数值、字符串、slice、map等）定义方法，只要这个命名类型的底层类型不是指针或者interface
func (p Point) Distance(q Point) float64 {
	// p.Distance的表达式叫做选择器，同选择一个struct类型的字段时
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}
