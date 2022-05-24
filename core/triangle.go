package core

import "math"

type Triangle struct {
	a *Vector
	b *Vector
	c *Vector
}

func NewTriangle(a *Vector, b *Vector, c *Vector) *Triangle {
	newT := &Triangle{a, b, c}
	return newT
}

func GetNormal(a, b, c, target *Vector) *Vector {
	target.SubVectors(c, b)
	_v0 := Vector{}
	_v0.SubVectors(a, b)
	target.Cross(&_v0)

	targetLengthSq := target.LengthSq()
	if targetLengthSq > 0 {

		return target.MultiplyScalar(1 / math.Sqrt(targetLengthSq))

	}

	return target.Set(0, 0, 0)
}
