package core

import (
	"bytes"
	"fmt"
	"math"
)

type Geometry struct {
	//Attributes map[string]interface{}
	Index                []int32
	Normal               []float64
	Position             []float64
	UV                   []float64
	Color                []float64
	BoundingSphere       *Sphere
	BoundingBox          *Box3
	MorphAttribute       map[string][][]float64
	MorphTargetsRelative bool
	DrawRange            DrawRange
	Groups               []Group
}

type DrawRange struct {
	Start int
	Count float64
}

func (g *Geometry) Clone() *Geometry {
	newGeom := &Geometry{}
	newGeom.MorphTargetsRelative = false
	newGeom.DrawRange = DrawRange{0, math.Inf(1)}
	newGeom.Color = make([]float64, len(g.Color))
	for i := range newGeom.Color {
		newGeom.Color[i] = g.Color[i]
	}
	newGeom.Index = make([]int32, len(g.Index))
	for i := range newGeom.Index {
		newGeom.Index[i] = g.Index[i]
	}
	newGeom.Position = make([]float64, len(g.Position))
	for i := range newGeom.Position {
		newGeom.Position[i] = g.Position[i]
	}
	newGeom.UV = make([]float64, len(g.UV))
	for i := range newGeom.UV {
		newGeom.UV[i] = g.UV[i]
	}
	newGeom.Normal = make([]float64, len(g.Normal))
	for i := range newGeom.Normal {
		newGeom.Normal[i] = g.Normal[i]
	}
	newGeom.Groups = make([]Group, len(g.Groups))
	for i := range newGeom.Groups {
		newGeom.Groups[i] = g.Groups[i]
	}
	newGeom.BoundingSphere = nil
	newGeom.BoundingBox = nil
	newGeom.MorphAttribute = map[string][][]float64{}
	return newGeom
}

type Group struct {
	Start         int32
	Count         int32
	MaterialIndex int32
}

func (g *Group) Clone() *Group {
	newGroup := &Group{}

	newGroup.Start = g.Start
	newGroup.Count = g.Count
	newGroup.MaterialIndex = g.MaterialIndex

	return newGroup
}

func NewGeometry() *Geometry {
	geo := &Geometry{}
	geo.Index = nil
	geo.Position = nil
	geo.DrawRange = DrawRange{0, math.Inf(1)}
	geo.MorphAttribute = map[string][][]float64{}
	return geo
}
func (g *Geometry) ToObj() string {
	output := []byte{}
	buffer := bytes.NewBuffer(output)
	for i := 0; i < len(g.Position); i += 3 {
		buffer.WriteString(fmt.Sprintf("v %f %f %f\n", g.Position[i], g.Position[i+1], g.Position[i+2]))
	}
	for i := 0; i < len(g.Index); i += 3 {
		buffer.WriteString(fmt.Sprintf("f %d %d %d\n", g.Index[i]+1, g.Index[i+1]+1, g.Index[i+2]+1))
	}
	return buffer.String()
}
func (g *Geometry) SetDrawRange(start int, count float64) *Geometry {
	g.DrawRange.Count = count
	g.DrawRange.Start = start
	return g
}
func (g *Geometry) ComputeBoundingSphere() {
	if g.BoundingSphere == nil {
		g.BoundingSphere = NewSphere(NewVector(0, 0, 0), 0)
	}
	position := g.Position
	morphAttributesPosition, morphAttributesPosition_ok := g.MorphAttribute["position"]
	// fmt.Println("position", position)
	if position != nil {
		center := g.BoundingSphere.Center
		_box := NewBox(
			NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64),
			NewVector(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64))
		_box.SetFromBufferAttribute(position)
		// fmt.Println(_box.Max, _box.Min, morphAttributesPosition_ok, morphAttributesPosition, g.BoundingSphere.Center)
		if morphAttributesPosition_ok {
			il := len(morphAttributesPosition)
			for i := 0; i < il; i++ {
				morphAttribute := morphAttributesPosition[i]
				_boxMorphTargets := NewBoxAutoFill()
				_boxMorphTargets.SetFromBufferAttribute(morphAttribute)
				if g.MorphTargetsRelative {
					_vector := Vector{}
					_vector.Add(_box.Min).Add(_boxMorphTargets.Min)
					_box.ExpandByPoint(&_vector)
					_vector = Vector{}
					_vector.Add(_box.Max).Add(_boxMorphTargets.Max)
					_box.ExpandByPoint(&_vector)
				} else {
					_box.ExpandByPoint(_boxMorphTargets.Min)
					_box.ExpandByPoint(_boxMorphTargets.Max)
				}
			}
		}
		// fmt.Println("center-pre", center, g.BoundingSphere.Center)
		_box.GetCenter(center)
		// fmt.Println("center", g.BoundingSphere.Center)
		// g.BoundingSphere.Center = center
		maxRadiusSq := float64(0)
		// second, try to find a boundingSphere with a radius smaller than the
		// boundingSphere of the boundingBox: sqrt(3) smaller in the best case
		il := len(position) / 3
		for i := 0; i < il; i += 1 {
			_vector := Vector{}
			_vector.Set(position[i], position[i+1], position[i+2])
			maxRadiusSq = math.Max(maxRadiusSq, center.DistanceToSquared(&_vector))
		}
		// process morph attributes if present

		if morphAttributesPosition != nil {
			il := len(morphAttributesPosition)
			for i := 0; i < il; i++ {
				// TODO: fix this bracket
				morphAttribute := morphAttributesPosition[i]
				morphTargetsRelative := g.MorphTargetsRelative
				jl := len(morphAttribute) / 3
				for j := 0; j < jl; j++ {
					_vector := Vector{}
					_vector.Set(morphAttribute[j*3], morphAttribute[j*3+1], morphAttribute[j*3+2])

					if morphTargetsRelative {
						_offset := Vector{}
						_offset.Set(position[j*3], position[j*3+1], position[j*3+2])
						_vector.Add(&_offset)

					}

					maxRadiusSq = math.Max(maxRadiusSq, center.DistanceToSquared(&_vector))

				}

			}
		}
		g.BoundingSphere.Radius = math.Sqrt(maxRadiusSq)
	}
}
func ApplyMatrix4(attribute []float64, m *Matrix4) []float64 {
	l := len(attribute) / 3
	for i := 0; i < l; i += 3 {
		v := Vector{}
		v.X = attribute[i]
		v.Y = attribute[i+1]
		v.Z = attribute[i+2]
		v.ApplyMatrix4(m)
		attribute[i] = v.X
		attribute[i+1] = v.Y
		attribute[i+2] = v.Z
	}
	return attribute
}
func (g *Geometry) ComputeVertexNormals() {
	index := g.Index
	positionAttribute := g.Position
	normalAttribute := make([]float64, len(positionAttribute))
	pA := NewVector(0, 0, 0)
	pB := NewVector(0, 0, 0)
	pC := NewVector(0, 0, 0)

	nA := NewVector(0, 0, 0)
	nB := NewVector(0, 0, 0)
	nC := NewVector(0, 0, 0)

	cb := NewVector(0, 0, 0)
	ab := NewVector(0, 0, 0)
	il := len(index)
	for i := 0; i < il; i += 3 {
		vA := index[i]
		vB := index[(i + 1)]
		vC := index[(i + 2)]

		pA.Set(positionAttribute[vA*3], positionAttribute[vA*3+1], positionAttribute[vA*3+2])
		pB.Set(positionAttribute[vB*3], positionAttribute[vB*3+1], positionAttribute[vB*3+2])
		pC.Set(positionAttribute[vC*3], positionAttribute[vC*3+1], positionAttribute[vC*3+2])

		// fmt.Println("pA, pB, pC", pA, pB, pC)

		cb.SubVectors(pC, pB)
		ab.SubVectors(pA, pB)
		cb.Cross(ab)

		nA.Set(normalAttribute[vA*3], normalAttribute[vA*3+1], normalAttribute[vA*3+2])
		nB.Set(normalAttribute[vB*3], normalAttribute[vB*3+1], normalAttribute[vB*3+2])
		nC.Set(normalAttribute[vC*3], normalAttribute[vC*3+1], normalAttribute[vC*3+2])

		// fmt.Println("nA,nB,nC", nA, nB, nC)

		nA.Add(cb)
		nB.Add(cb)
		nC.Add(cb)

		normalAttribute[vA*3] = nA.X
		normalAttribute[vA*3+1] = nA.Y
		normalAttribute[vA*3+2] = nA.Z

		normalAttribute[vB*3] = nB.X
		normalAttribute[vB*3+1] = nB.Y
		normalAttribute[vB*3+2] = nB.Z

		normalAttribute[vC*3] = nC.X
		normalAttribute[vC*3+1] = nC.Y
		normalAttribute[vC*3+2] = nC.Z

	}
	g.Normal = normalAttribute
	g.NormalizeNormal()

}
func (g *Geometry) NormalizeNormal() {
	normal := g.Normal
	il := len(normal) / 3
	for i := 0; i < il; i++ {
		_vector := Vector{normal[i*3], normal[i*3+1], normal[i*3+2]}
		_vector.Normalize()
		normal[i*3] = _vector.X
		normal[i*3+1] = _vector.Y
		normal[i*3+2] = _vector.Z
	}
	g.Normal = normal
}
func (g *Geometry) ComputeBoundingBox() {
	if g.BoundingBox == nil {
		g.BoundingBox = &Box3{}
	}
	position := g.Position
	morphAttributesPosition := g.MorphAttribute["position"]
	if position != nil {
		g.BoundingBox.SetFromBufferAttribute(position)

		// process morph attributes if present

		if morphAttributesPosition != nil {
			il := len(morphAttributesPosition)
			for i := 0; i < il; i++ {
				_box := Box3{}
				morphAttribute := morphAttributesPosition[i]
				_box.SetFromBufferAttribute(morphAttribute)
				if g.MorphTargetsRelative {
					_vector := Vector{}
					_vector.Copy(g.BoundingBox.Min).Add(_box.Min)
					g.BoundingBox.ExpandByPoint(&_vector)

					_vector = Vector{}
					_vector.Copy(g.BoundingBox.Max).Add(_box.Max)
					g.BoundingBox.ExpandByPoint(&_vector)
				} else {
					g.BoundingBox.ExpandByPoint(_box.Min)
					g.BoundingBox.ExpandByPoint(_box.Max)
				}
			}
		} else {
			g.BoundingBox.MakeEmpty()
		}
	}
}
func (g *Geometry) ApplyMatrix4(m *Matrix4) *Geometry {
	//position := g.Position
	g.Position = ApplyMatrix4(g.Position, m)
	if g.BoundingBox != nil {
		g.ComputeBoundingBox()
	}

	if g.BoundingSphere != nil {
		g.ComputeBoundingSphere()
	}
	return g
}
func (g *Geometry) AddGroup(start, count, materialIndex int32) {
	newGroup := Group{}
	newGroup.Start = start
	newGroup.Count = count
	newGroup.MaterialIndex = materialIndex
	g.Groups = append(g.Groups, newGroup)
}
