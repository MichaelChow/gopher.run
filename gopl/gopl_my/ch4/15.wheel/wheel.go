package wheel

type Point struct {
	X, Y int
}

type Cicle struct {
	Point
	Radius int
}

type Wheel struct {
	Cicle
	Spokes int
}
