package core

type Matrix4 struct {
	elements []float64
}

func NewMatrix4() *Matrix4 {
	mat4 := &Matrix4{}
	mat4.elements = []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return mat4
}

func (m *Matrix4) Identity() *Matrix4 {
	m.elements = []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
	return m
}
func (m *Matrix4) Clone() *Matrix4 {
	mat4 := &Matrix4{}
	mat4.elements = make([]float64, 16)
	for i := range m.elements {
		mat4.elements[i] = m.elements[i]
	}
	return mat4
}

func (m *Matrix4) Copy(otherM *Matrix4) *Matrix4 {
	for i := range otherM.elements {
		m.elements[i] = otherM.elements[i]
	}
	return m
}
func (m *Matrix4) CopyPosition(otherM *Matrix4) *Matrix4 {
	m.elements[12] = otherM.elements[12]
	m.elements[13] = otherM.elements[13]
	m.elements[14] = otherM.elements[14]
	return m
}
func (m *Matrix4) Decompose(positino *Vector, quaterion Quaterion, scale *Vector) {

}
func (m *Matrix4) Compose(position *Vector, quaternion *Quaterion, scale *Vector) *Matrix4 {
	te := m.elements
	x := quaternion._x
	y := quaternion._y
	z := quaternion._z
	w := quaternion._w

	x2 := x + x
	y2 := y + y
	z2 := z + z

	xx := x * x2
	xy := x * y2
	xz := x * z2

	yy := y * y2
	yz := y * z2
	zz := z * z2

	wx := w * x2
	wy := w * y2
	wz := w * z2

	sx := scale.X
	sy := scale.Y
	sz := scale.Z

	te[0] = (1 - (yy + zz)) * sx
	te[1] = (xy + wz) * sx
	te[2] = (xz - wy) * sx
	te[3] = 0

	te[4] = (xy - wz) * sy
	te[5] = (1 - (xx + zz)) * sy
	te[6] = (yz + wx) * sy
	te[7] = 0

	te[8] = (xz + wy) * sz
	te[9] = (yz - wx) * sz
	te[10] = (1 - (xx + yy)) * sz
	te[11] = 0

	te[12] = position.X
	te[13] = position.Y
	te[14] = position.Z
	te[15] = 1

	m.elements = te

	return m
}
func (m *Matrix4) MakeRotationFromQuaternion(q *Quaterion) *Matrix4 {
	return m.Compose(_ZERO, q, _ONES)
}
func (m *Matrix4) Set(n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44 float64) *Matrix4 {
	te := m.elements

	te[0] = n11
	te[4] = n12
	te[8] = n13
	te[12] = n14
	te[1] = n21
	te[5] = n22
	te[9] = n23
	te[13] = n24
	te[2] = n31
	te[6] = n32
	te[10] = n33
	te[14] = n34
	te[3] = n41
	te[7] = n42
	te[11] = n43
	te[15] = n44

	m.elements = te
	return m
}
func (m *Matrix4) Invert() *Matrix4 {
	te := m.elements

	n11 := te[0]
	n21 := te[1]
	n31 := te[2]
	n41 := te[3]
	n12 := te[4]
	n22 := te[5]
	n32 := te[6]
	n42 := te[7]
	n13 := te[8]
	n23 := te[9]
	n33 := te[10]
	n43 := te[11]
	n14 := te[12]
	n24 := te[13]
	n34 := te[14]
	n44 := te[15]

	t11 := n23*n34*n42 - n24*n33*n42 + n24*n32*n43 - n22*n34*n43 - n23*n32*n44 + n22*n33*n44
	t12 := n14*n33*n42 - n13*n34*n42 - n14*n32*n43 + n12*n34*n43 + n13*n32*n44 - n12*n33*n44
	t13 := n13*n24*n42 - n14*n23*n42 + n14*n22*n43 - n12*n24*n43 - n13*n22*n44 + n12*n23*n44
	t14 := n14*n23*n32 - n13*n24*n32 - n14*n22*n33 + n12*n24*n33 + n13*n22*n34 - n12*n23*n34

	det := n11*t11 + n21*t12 + n31*t13 + n41*t14

	if det == 0 {
		return m.Set(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
	}

	detInv := 1 / det

	te[0] = t11 * detInv
	te[1] = (n24*n33*n41 - n23*n34*n41 - n24*n31*n43 + n21*n34*n43 + n23*n31*n44 - n21*n33*n44) * detInv
	te[2] = (n22*n34*n41 - n24*n32*n41 + n24*n31*n42 - n21*n34*n42 - n22*n31*n44 + n21*n32*n44) * detInv
	te[3] = (n23*n32*n41 - n22*n33*n41 - n23*n31*n42 + n21*n33*n42 + n22*n31*n43 - n21*n32*n43) * detInv

	te[4] = t12 * detInv
	te[5] = (n13*n34*n41 - n14*n33*n41 + n14*n31*n43 - n11*n34*n43 - n13*n31*n44 + n11*n33*n44) * detInv
	te[6] = (n14*n32*n41 - n12*n34*n41 - n14*n31*n42 + n11*n34*n42 + n12*n31*n44 - n11*n32*n44) * detInv
	te[7] = (n12*n33*n41 - n13*n32*n41 + n13*n31*n42 - n11*n33*n42 - n12*n31*n43 + n11*n32*n43) * detInv

	te[8] = t13 * detInv
	te[9] = (n14*n23*n41 - n13*n24*n41 - n14*n21*n43 + n11*n24*n43 + n13*n21*n44 - n11*n23*n44) * detInv
	te[10] = (n12*n24*n41 - n14*n22*n41 + n14*n21*n42 - n11*n24*n42 - n12*n21*n44 + n11*n22*n44) * detInv
	te[11] = (n13*n22*n41 - n12*n23*n41 - n13*n21*n42 + n11*n23*n42 + n12*n21*n43 - n11*n22*n43) * detInv

	te[12] = t14 * detInv
	te[13] = (n13*n24*n31 - n14*n23*n31 + n14*n21*n33 - n11*n24*n33 - n13*n21*n34 + n11*n23*n34) * detInv
	te[14] = (n14*n22*n31 - n12*n24*n31 - n14*n21*n32 + n11*n24*n32 + n12*n21*n34 - n11*n22*n34) * detInv
	te[15] = (n12*n23*n31 - n13*n22*n31 + n13*n21*n32 - n11*n23*n32 - n12*n21*n33 + n11*n22*n33) * detInv

	m.elements = te
	return m
}
func (m *Matrix4) MultiplyMatrices(a, b *Matrix4) *Matrix4 {
	ae := a.elements
	be := b.elements
	te := m.elements

	a11 := ae[0]
	a12 := ae[4]
	a13 := ae[8]
	a14 := ae[12]

	a21 := ae[1]
	a22 := ae[5]
	a23 := ae[9]
	a24 := ae[13]

	a31 := ae[2]
	a32 := ae[6]
	a33 := ae[10]
	a34 := ae[14]

	a41 := ae[3]
	a42 := ae[7]
	a43 := ae[11]
	a44 := ae[15]

	b11 := be[0]
	b12 := be[4]
	b13 := be[8]
	b14 := be[12]

	b21 := be[1]
	b22 := be[5]
	b23 := be[9]
	b24 := be[13]

	b31 := be[2]
	b32 := be[6]
	b33 := be[10]
	b34 := be[14]

	b41 := be[3]
	b42 := be[7]
	b43 := be[11]
	b44 := be[15]

	te[0] = a11*b11 + a12*b21 + a13*b31 + a14*b41
	te[4] = a11*b12 + a12*b22 + a13*b32 + a14*b42
	te[8] = a11*b13 + a12*b23 + a13*b33 + a14*b43
	te[12] = a11*b14 + a12*b24 + a13*b34 + a14*b44

	te[1] = a21*b11 + a22*b21 + a23*b31 + a24*b41
	te[5] = a21*b12 + a22*b22 + a23*b32 + a24*b42
	te[9] = a21*b13 + a22*b23 + a23*b33 + a24*b43
	te[13] = a21*b14 + a22*b24 + a23*b34 + a24*b44

	te[2] = a31*b11 + a32*b21 + a33*b31 + a34*b41
	te[6] = a31*b12 + a32*b22 + a33*b32 + a34*b42
	te[10] = a31*b13 + a32*b23 + a33*b33 + a34*b43
	te[14] = a31*b14 + a32*b24 + a33*b34 + a34*b44

	te[3] = a41*b11 + a42*b21 + a43*b31 + a44*b41
	te[7] = a41*b12 + a42*b22 + a43*b32 + a44*b42
	te[11] = a41*b13 + a42*b23 + a43*b33 + a44*b43
	te[15] = a41*b14 + a42*b24 + a43*b34 + a44*b44

	m.elements = te
	return m
}

var _ZERO = NewVector(0, 0, 0)

var _ONES = NewVector(1, 1, 1)
