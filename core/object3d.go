package core

type Object3D struct {
	Matrix      *Matrix4
	MatrixWorld *Matrix4
	Parent      *Object3D
	Children    []*Object3D
	Quaterion   *Quaterion
	Position    *Vector
	Rotation    *Euler
	Scale       *Vector
}

func (o *Object3D) UpdateMatrixWorld(force bool) {
	if o.Parent == nil {
		o.MatrixWorld.Copy(o.Matrix)
	} else {
		o.MatrixWorld.MultiplyMatrices(o.Parent.MatrixWorld, o.Matrix)
	}
	for i := range o.Children {
		o.Children[i].UpdateMatrixWorld(force)
	}
}

func (o *Object3D) Clone() *Object3D {
	newObj := &Object3D{}
	newObj.Children = make([]*Object3D, len(o.Children))
	for idx, k := range o.Children {
		newObj.Children[idx] = k.Clone()
	}

	return newObj
}
