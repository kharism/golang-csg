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
func (m *Matrix4) Decompose(positino *Vector, quaterion Quaterion, scale *Vector)
