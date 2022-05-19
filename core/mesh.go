package core

type Mesh struct {
	*Object3D
	Geometry *Geometry
	//material
}

func NewMesh(geometry *Geometry) *Mesh {
	pp := NewObject3D()
	mesh := &Mesh{}
	mesh.Object3D = pp

	mesh.Geometry = geometry

	return mesh
}

func (m *Mesh) Clone() *Mesh {
	jj := m.Object3D.Clone()
	newMesh := &Mesh{Object3D: jj}
	newMesh.Geometry = m.Geometry.Clone()
	return m
}
