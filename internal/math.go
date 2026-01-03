package internal

type Vector2 struct {
	X int
	Y int
}

func NewVector2(x int, y int) Vector2 {
	return Vector2{X: x, Y: y}
}

func (v *Vector2) Add(v2 Vector2) Vector2 {
	return NewVector2(v.X+v2.X, v.Y+v2.Y)
}
