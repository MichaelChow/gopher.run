// Embed demonstrates basic struct embedding.
// See page 106.
package main

import "fmt"

type Point struct{ X, Y int }

type Circle struct {
	Point
	Radius int
}

type Wheel struct {
	Circle
	Spokes int
}

func main() {
	var w Wheel
	w = Wheel{Circle{Point{8, 8}, 5}, 20}
	// %#v，#副词表示用和Go语言类似的语法打印值。对于结构体类型来说，将包含每个成员的名字和值。
	fmt.Printf("%#v\n", w) // main.Wheel{Circle:main.Circle{Point:main.Point{X:8, Y:8}, Radius:5}, Spokes:20}

	w = Wheel{
		Circle: Circle{
			Point:  Point{X: 8, Y: 8},
			Radius: 5,
		},
		Spokes: 20, // NOTE: trailing comma necessary here (and at Radius)
	}

	fmt.Printf("%#v\n", w) // main.Wheel{Circle:main.Circle{Point:main.Point{X:8, Y:8}, Radius:5}, Spokes:20}

	w.X = 42

	fmt.Printf("%#v\n", w) // main.Wheel{Circle:main.Circle{Point:main.Point{X:42, Y:8}, Radius:5}, Spokes:20}
}
