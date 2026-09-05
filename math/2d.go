package math

import (
	"gioui.org/f32"
)

func Abs[T ~int | ~float32](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

func Distance(p1, p2 f32.Point) (dx, dy float32) {
	return Abs(p2.X - p1.X), Abs(p2.Y - p1.Y)
}
