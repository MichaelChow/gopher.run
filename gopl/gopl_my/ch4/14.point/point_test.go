package point_test

import (
	"fmt"
	"testing"

	point "gopher.run/go/src/ch4/14.point"
)

func Scale(p point.Point, factor int) point.Point {
	return point.Point{p.X * factor, p.Y * factor}
}

// 如果考虑效率的话，较大的结构体通常会用指针的方式传入和返回 （避免分配一份大结构体的内存再拷贝其数据的操作）
// 在Go语言中，所有的函数参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量
func Bonus(p *point.Point, factor int) {
	p.X *= factor
	p.Y *= factor
}

func TestPoint(t *testing.T) {
	// 第一种写法：要求以结构体成员定义的顺序为每个结构体成员指定一个字面值。
	// 它要求写代码和读代码的人要记住结构体的每个成员的类型和顺序，不过结构体成员有细微的调整就可能导致上述代码不能编译
	// 因此，上述的语法一般只在定义结构体的包内部使用，或者是在较小的结构体中使用，这些结构体的成员排列比较规则
	// 如image.Point{x, y}或color.RGBA{red, green, blue, alpha}
	p1 := point.Point{1, 2}
	fmt.Println(p1)

	// 第二种写法(更常用)：以成员名字和相应的值来初始化，可以包含部分或全部的成员。如果成员被忽略的话将默认用零值。
	// 如Lissajous程序的写法：
	// anim := gif.GIF{LoopCount: nframes}
	p2 := point.Point{X: 1, Y: 2}

	fmt.Println(p2)

	// 两种不同形式的写法不能混合使用
	// 而且，你不能企图在外部包中用第一种顺序赋值的技巧来偷偷地初始化结构体中未导出的成员（小写字母开头的）
	// p2 := point.Point{X: 1, Y: 2, z: 3}  // cannot refer to unexported field z in struct literal of type point.Point

	fmt.Println(Scale(p1, 5)) // "{5 10}"
	Bonus(&p1, 5)
	fmt.Println(p1) // "{5 10}"

	pp1 := &point.Point{1, 2} // 可以直接在表达式中使用，如一个函数调用
	pp2 := new(point.Point)
	*pp2 = point.Point{1, 2} // 这两句等价上述pp1一句的写法
	fmt.Println(pp1, pp2)
}
