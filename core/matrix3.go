package core

type Matrix3 struct {
	elements []float64
}

func NewMatrix3() *Matrix3 {
	mat3 := &Matrix3{}
	mat3.elements = []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
	return mat3
}

func (m *Matrix3) Identity() *Matrix3 {
	m.elements = []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
	return m
}
func (m *Matrix3) Clone() *Matrix3 {
	mat3 := &Matrix3{}
	mat3.elements = make([]float64, 9)
	for i := range m.elements {
		mat3.elements[i] = m.elements[i]
	}
	return mat3
}

func (m *Matrix3) Copy(otherM *Matrix3) *Matrix3 {
	for i := range otherM.elements {
		m.elements[i] = otherM.elements[i]
	}
	return m
}
func (m *Matrix3) Set(n11, n12, n13, n21, n22, n23, n31, n32, n33 float64) *Matrix3 {
	m.elements[0] = n11
	m.elements[1] = n21
	m.elements[2] = n31
	m.elements[3] = n12
	m.elements[4] = n22
	m.elements[5] = n32
	m.elements[6] = n13
	m.elements[7] = n23
	m.elements[8] = n33

	return m
}
func (m *Matrix3) ExtractBasis(xAxis, yAxis, zAxis Vector) *Matrix3 {
	xAxis.SetFromMatrix3Column(m, 0)
	xAxis.SetFromMatrix3Column(m, 1)
	xAxis.SetFromMatrix3Column(m, 2)

	return m
}
func (m *Matrix3) SetFromMatrix4(m4 *Matrix4) *Matrix3 {
	m.Set(
		m4.elements[0], m4.elements[4], m4.elements[8],
		m4.elements[1], m4.elements[5], m4.elements[9],
		m4.elements[2], m4.elements[6], m4.elements[10],
	)
	return m
}

func (m *Matrix3) Multiply(m2 *Matrix3) *Matrix3 {
	return m.MultiplyMatrix3(m, m2)
}

func (m *Matrix3) PreMultiply(m2 *Matrix3) *Matrix3 {
	return m.MultiplyMatrix3(m2, m)
}

func (m *Matrix3) MultiplyMatrix3(a, b *Matrix3) *Matrix3 {
	ae := a.elements
	be := b.elements

	//result := NewMatrix3()

	a11 := ae[0]
	a12 := ae[3]
	a13 := ae[6]
	a21 := ae[1]
	a22 := ae[4]
	a23 := ae[7]
	a31 := ae[2]
	a32 := ae[5]
	a33 := ae[8]

	b11 := be[0]
	b12 := be[3]
	b13 := be[6]
	b21 := be[1]
	b22 := be[4]
	b23 := be[7]
	b31 := be[2]
	b32 := be[5]
	b33 := be[8]

	m.elements[0] = a11*b11 + a12*b21 + a13*b31
	m.elements[3] = a11*b12 + a12*b22 + a13*b32
	m.elements[6] = a11*b13 + a12*b23 + a13*b33

	m.elements[1] = a21*b11 + a22*b21 + a23*b31
	m.elements[4] = a21*b12 + a22*b22 + a23*b32
	m.elements[7] = a21*b13 + a22*b23 + a23*b33

	m.elements[2] = a31*b11 + a32*b21 + a33*b31
	m.elements[5] = a31*b12 + a32*b22 + a33*b32
	m.elements[8] = a31*b13 + a32*b23 + a33*b33

	return m
}

func (m *Matrix3) MultplyScalar(s float64) *Matrix3 {
	for idx := range m.elements {
		m.elements[idx] *= s
	}
	return m
}
func (m *Matrix3) Determinant() float64 {
	te := m.elements
	a := te[0]
	b := te[1]
	c := te[2]
	d := te[3]
	e := te[4]
	f := te[5]
	g := te[6]
	h := te[7]
	i := te[8]

	return a*e*i - a*f*h - b*d*i + b*f*g + c*d*h - c*e*g
}
func (this *Matrix3) getNormalMatrix(matrix4 *Matrix4) *Matrix3 {
	return this.SetFromMatrix4(matrix4).Invert().Transpose()
}
func (m *Matrix3) Invert() *Matrix3 {
	te := m.elements
	n11 := te[0]
	n21 := te[1]
	n31 := te[2]
	n12 := te[3]
	n22 := te[4]
	n32 := te[5]
	n13 := te[6]
	n23 := te[7]
	n33 := te[8]

	t11 := n33*n22 - n32*n23
	t12 := n32*n13 - n33*n12
	t13 := n23*n12 - n22*n13

	det := n11*t11 + n21*t12 + n31*t13

	if det == 0 {
		return m.Set(0, 0, 0, 0, 0, 0, 0, 0, 0)
	}
	detInv := 1 / det

	m.elements[0] = t11 * detInv
	m.elements[1] = (n31*n23 - n33*n21) * detInv
	m.elements[2] = (n32*n21 - n31*n22) * detInv

	m.elements[3] = t12 * detInv
	m.elements[4] = (n33*n11 - n31*n13) * detInv
	m.elements[5] = (n31*n12 - n32*n11) * detInv

	m.elements[6] = t13 * detInv
	m.elements[7] = (n21*n13 - n23*n11) * detInv
	m.elements[8] = (n22*n11 - n21*n12) * detInv

	return m
}

func (m *Matrix3) Transpose() *Matrix3 {
	var tmp float64
	tmp = m.elements[1]
	m.elements[1] = m.elements[3]
	m.elements[3] = tmp
	tmp = m.elements[2]
	m.elements[2] = m.elements[6]
	m.elements[6] = tmp
	tmp = m.elements[5]
	m.elements[5] = m.elements[7]
	m.elements[7] = tmp

	return m
}
