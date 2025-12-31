package internal

type Vector2 struct {
	X int
	Y int
}

func NewVector2(x int, y int) Vector2 {
	return Vector2{X: x, Y: y}
}
