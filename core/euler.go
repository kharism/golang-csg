package core

import "math"

type Euler struct {
	_x, _y, _z float64
	_order     string
}

func NewEuler(x, y, z float64, order string) *Euler {
	j := &Euler{}
	j._x = x
	j._y = y
	j._z = z
	if order == "" {
		order = DEFAULT_ORDER
	}
	j._order = order
	return j
}

func (e *Euler) Copy() *Euler {
	p := &Euler{}
	p._x = e._x
	p._y = e._y
	p._z = e._z
	return p
}

func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

func (e *Euler) SetFromRotationMatrix(m *Matrix4, order string, update bool) *Euler {
	if order == "" {
		order = e._order
	}
	te := m.elements
	m11 := te[0]
	m12 := te[4]
	m13 := te[8]

	m21 := te[1]
	m22 := te[5]
	m23 := te[9]

	m31 := te[2]
	m32 := te[6]
	m33 := te[10]

	switch order {
	case "XYZ":
		e._y = math.Asin(Clamp(m13, -1, 1))
		if math.Abs(m13) < 0.9999999 {
			e._x = math.Atan2(-m23, m33)
			e._z = math.Atan2(-m12, m11)
		} else {
			e._x = math.Atan2(m32, m22)
			e._z = 0
		}
	case "YXZ":
		e._x = math.Asin(Clamp(m32, -1, 1))
		if math.Abs(m32) < 0.9999999 {
			e._y = math.Atan2(-m31, m33)
			e._z = math.Atan2(-m12, m22)
		} else {
			e._y = 0
			e._z = math.Atan2(m21, m11)

		}
	case "ZYX":
		e._y = math.Asin(-Clamp(m31, -1, 1))
		if math.Abs(m31) < 0.9999999 {

			e._x = math.Atan2(m32, m33)
			e._z = math.Atan2(m21, m11)

		} else {

			e._x = 0
			e._z = math.Atan2(-m12, m22)

		}
	case "YZX":
		e._z = math.Asin(Clamp(m21, -1, 1))
		if math.Abs(m21) < 0.9999999 {

			e._x = math.Atan2(-m23, m22)
			e._y = math.Atan2(-m31, m11)

		} else {

			e._x = 0
			e._y = math.Atan2(m13, m33)

		}
	case "XZY":
		e._z = math.Asin(-Clamp(m12, -1, 1))

		if math.Abs(m12) < 0.9999999 {

			e._x = math.Atan2(m32, m22)
			e._y = math.Atan2(m13, m11)

		} else {

			e._x = math.Atan2(-m23, m33)
			e._y = 0

		}
	}
	e._order = order
	if update {

	}
	return e
}
func (e *Euler) SetFromQuaternion(q *Quaterion, order string, update bool) *Euler {
	_matrix := NewMatrix4()
	_matrix.MakeRotationFromQuaternion(q)
	return e.SetFromRotationMatrix(_matrix, order, update)
}

const DEFAULT_ORDER = "XYZ"

var ROTATION_ORDER = []string{"XYZ", "YZX", "ZXY", "XZY", "YXZ", "ZYX"}
