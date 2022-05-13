package core

import (
	"math"
)

type CSG struct {
	Polygons []*Polygon
}

func FromPolygons(polygons []*Polygon) *CSG {
	csg := &CSG{}
	csg.Polygons = polygons
	return csg
}

func FromGeometry(geom *Geometry, objectIndex []float32) *CSG {
	polygons := []*Polygon{}
	position := geom.Position
	normal := geom.Normal
	groups := geom.Groups
	colors := geom.Color
	uv := geom.UV

	var index []int32
	if geom.Index != nil {
		index = geom.Index
	} else {
		indexSize := len(position) / 3
		index = make([]int32, indexSize)
	}
	triCount := len(index) / 3
	polygons = make([]*Polygon, triCount)
	pli := 0
	l := len(index)
	for i := 0; i < l; i += 3 {
		Vertices := []*Vertex{}
		for j := 0; j < 3; j++ {
			vi := index[i+j]
			vp := vi * 3
			vt := vi * 2
			x := position[vp]
			y := position[vp+1]
			z := position[vp+2]
			nx := normal[vp]
			ny := normal[vp+1]
			nz := normal[vp+2]
			u := uv[vt]
			v := uv[vt+1]
			pp := NewVertex(
				NewVector(x, y, z),
				NewVector(nx, ny, nz),
				NewVector(u, v, 0),
				NewVector(colors[vt], colors[vt+1], colors[vt+2]),
			)
			Vertices[j] = &pp
		}
		if groups != nil && len(groups) > 0 {
			for _, grp := range groups {
				if index[i] >= grp.Start && index[i] < grp.Start+grp.Count {
					polygons[pli] = NewPolygon(Vertices, grp.MaterialIndex)
				}

			}
		} else {
			polygons[pli] = NewPolygon(Vertices, objectIndex)
		}
		pli++
	}
	realPolygon := []*Polygon{}
	for _, p := range polygons {
		if !math.IsNaN(p.Plane.Normal.X) {
			realPolygon = append(realPolygon, p)
		}
	}
	return FromPolygons(realPolygon)
}

type JSMap struct {
	Storage map[interface{}]interface{}
	LastIdx int
}

func (s *JSMap) Add(index int32, item interface{}) {
	s.Storage[index] = item
}
func (s *JSMap) Push(item interface{}) {
	s.Storage[s.LastIdx] = item
	s.LastIdx++
}

func ToGeometry(csg *CSG, toMatrix *Matrix4) *Geometry {
	triCount := 0
	ps := csg.Polygons
	for _, p := range ps {
		triCount += len(p.Vertices) - 2
	}
	geom := &Geometry{}

	vertices := NewNBuf3(triCount * 3 * 3)
	normals := NewNBuf3(triCount * 3 * 3)
	uvs := NewNBuf2(triCount * 3 * 3)
	var colors *NBuf3
	grps := JSMap{}
	dgrp := JSMap{}
	for _, p := range ps {
		pvs := p.Vertices
		pvlen := len(p.Vertices)
		if grps.Storage[p.Shared] == nil {
			grps.Storage[p.Shared] = []interface{}{}
		}

		if pvlen > 0 && pvs[0].Color != nil {
			if colors == nil {
				colors = NewNBuf3(triCount * 3 * 3)
			}
		}
		for j := 3; j <= pvlen; j++ {
			//var grp map[interface{}]interface{}
			if p.Shared == nil {
				dgrp.Push(vertices.Top / 3)
				dgrp.Push(vertices.Top/3 + 1)
				dgrp.Push(vertices.Top/3 + 2)
			} else {
				ll := grps.Storage[p.Shared].([]interface{})
				ll = append(ll, vertices.Top/3)
				ll = append(ll, vertices.Top/3+1)
				ll = append(ll, vertices.Top/3+2)
			}
			vertices.Write(pvs[0].Pos)
			vertices.Write(pvs[j-2].Pos)
			vertices.Write(pvs[j-1].Pos)
			normals.Write(pvs[0].Normal)
			normals.Write(pvs[j-2].Normal)
			normals.Write(pvs[j-1].Normal)
			if uvs != nil {
				uvs.Write(pvs[0].UV)
				uvs.Write(pvs[j-2].UV)
				uvs.Write(pvs[j-1].UV)
			}

			if colors != nil {
				colors.Write(pvs[0].Color)
				colors.Write(pvs[j-2].Color)
				colors.Write(pvs[j-1].Color)
			}
		}
	}
	geom.Position = vertices.Arr
	geom.Normal = normals.Arr
	if uvs != nil {
		geom.UV = uvs.Arr
	}
	if colors != nil {
		geom.Color = colors.Arr
	}
	for gi := 0; gi < len(grps.Storage); gi++ {
		if _, ok := grps.Storage[gi]; !ok {
			grps.Add(int32(gi), []interface{}{})
		}
	}
	if len(grps.Storage) > 0 {
		index := []int32{}
		gbase := 0
		for gi := 0; gi < len(grps.Storage); gi++ {
			geom.AddGroup(gbase, len(grps.Storage[gi].([]interface{})), gi)
			gbase += len(grps.Storage[gi].([]interface{}))
			kk := grps.Storage[gi].([]interface{})
			for _, val := range kk {
				index = append(index, val.(int32))
			}
		}
	}
	return geom
}

func FromMesh(mesh *Mesh, objectIndex []float32) *CSG {
	csg := FromGeometry(mesh.Geometry, objectIndex)
	ttvv0 := NewVector(0, 0, 0)
	tmpm3 := NewMatrix3()

	tmpm3.getNormalMatrix(mesh.Matrix)
	for i := 0; i < len(csg.Polygons); i++ {
		p := csg.Polygons[i]
		for j := 0; j < len(p.Vertices); j++ {
			v := p.Vertices[j]
			v.Pos.Copy(ttvv0.Copy(v.Pos)).ApplyMatrix4(mesh.Matrix)
			v.Normal.Copy(ttvv0.Copy(v.Normal)).ApplyMatrix3(mesh.Matrix)
		}
	}
	return csg
}
func ToMesh(csg *CSG, toMatrix *Matrix4) *Mesh {
	geom := ToGeometry(csg, toMatrix)
	m := &Mesh{}
	m.Geometry = geom
	m.Matrix.Copy(toMatrix)
	m.Matrix.Decompose(m.Position, *m.Quaterion, m.Scale)
	m.Rotation.SetFromQuaternion(m.Quaterion, m.Rotation._order, false)
	m.UpdateMatrixWorld(false)
	return m
}
func (m *CSG) Clone() *CSG {
	newcsg := &CSG{}
	newcsg.Polygons = []*Polygon{}
	for i := range m.Polygons {
		if !math.IsInf(m.Polygons[i].Plane.W, 0) {
			newcsg.Polygons = append(newcsg.Polygons, m.Polygons[i].Clone())
		}
	}
	return newcsg
}
func (m *CSG) Subtract(csg *CSG) *CSG {
	a := NewNode(m.Clone().Polygons)
	b := NewNode(csg.Clone().Polygons)
	a.Invert()
	a.ClipTo(b)
	b.ClipTo(a)
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.Build(b.AllPolygons())
	a.Invert()
	return FromPolygons(a.AllPolygons())
}
