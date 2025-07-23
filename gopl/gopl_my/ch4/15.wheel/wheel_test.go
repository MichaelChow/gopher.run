package wheel_test

import (
	"fmt"
	"testing"

	wheel "gopher.run/go/src/ch4/15.wheel"
)

func TestWheal(t *testing.T) {
	w := wheel.Wheel{
		Cicle: wheel.Cicle{
			Point: wheel.Point{
				X: 1,
				Y: 2,
			},
			Radius: 3,
		},
		Spokes: 4,
	}
	fmt.Println(w, w.X, w.Y, w.Radius, w.Spokes) // {{{1 2} 3} 4} 1 2 3 4
}
