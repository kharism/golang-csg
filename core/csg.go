package core

import (
	"fmt"
	"math"

	"github.com/eaciit/toolkit"
)

type CSG struct {
	Polygons []*Polygon
}

func FromPolygons(polygons []*Polygon) *CSG {
	csg := &CSG{}
	csg.Polygons = polygons
	return csg
}

func FromGeometry(geom *Geometry, objectIndex interface{}) *CSG {
	polygons := []*Polygon{}
	position := geom.Position
	normal := geom.Normal
	groups := geom.Groups
	colors := geom.Color
	uv := geom.UV

	var index []int32
	// fmt.Println("DDDDD", len(geom.Index), len(position), 3)
	if geom.Index != nil {
		index = geom.Index
	} else {
		indexSize := len(position) / 3
		index = make([]int32, indexSize)
		for i := int32(0); i < int32(indexSize); i++ {
			index[i] = i
		}
	}
	triCount := len(index) / 3
	polygons = make([]*Polygon, triCount)

	pli := 0
	l := len(index)
	// fmt.Println(len(polygons), l, position)
	for i := 0; i < l; i += 3 {
		Vertices := make([]*Vertex, 3)
		for j := 0; j < 3; j++ {
			vi := index[i+j]
			vp := vi * 3
			vt := vi * 2
			x := position[vp]
			y := position[vp+1]
			z := position[vp+2]
			// Logger.Println("VP", vp, position[vp], vp+1, position[vp+1], vp+2, position[vp+2])
			nx := normal[vp]
			ny := normal[vp+1]
			nz := normal[vp+2]
			u := uv[vt]
			v := uv[vt+1]
			pp := &Vertex{}
			pp.Pos = NewVector(x, y, z)
			pp.Normal = NewVector(nx, ny, nz)
			pp.UV = NewVector(u, v, 0)
			if colors != nil {
				pp.Color = NewVector(colors[vt], colors[vt+1], colors[vt+2])
			}
			Vertices[j] = pp
			// fmt.Println("Vertices[j]", i+j, vi, toolkit.JsonString(Vertices[j]), pli)
		}
		// fmt.Println("Vertices DONE", toolkit.JsonString(Vertices))
		if groups != nil && len(groups) > 0 {
			for _, grp := range groups {
				if index[i] >= grp.Start && index[i] < grp.Start+grp.Count {
					polygons[pli] = NewPolygon(Vertices, grp.MaterialIndex)
					// fmt.Println("polygons[pli]XY", pli, toolkit.JsonString(Vertices), toolkit.JsonString(polygons[pli]))
				}

			}
		} else {
			polygons[pli] = NewPolygon(Vertices, objectIndex)
		}
		pli++
	}
	realPolygon := []*Polygon{}
	fmt.Println(polygons)
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
	grps := map[int32][]int32{}
	dgrp := []int32{}
	for idx, p := range ps {
		pvs := p.Vertices
		pvlen := len(p.Vertices)
		fmt.Println("idx", idx, len(grps), len(dgrp))

		// fmt.Println("grps")
		// fmt.Println(toolkit.JsonString(grps), )
		// fmt.Println("dgrp")
		// fmt.Println(toolkit.JsonString(grps), )
		if p.Shared != -1 {
			if _, ok := grps[p.Shared]; !ok {
				grps[p.Shared] = []int32{}
			}
		}

		if pvlen > 0 && pvs[0].Color != nil {
			if colors == nil {
				colors = NewNBuf3(triCount * 3 * 3)
			}
		}
		if idx == 26 || idx == 27 {
			fmt.Println("grps", grps)
			fmt.Println("grps", dgrp)
			fmt.Println("p.Shared", p.Shared)
			fmt.Println(p.Vertices[0].Pos, p.Vertices[1].Pos, p.Vertices[2].Pos)
			fmt.Println(">>>>")
		}
		for j := 3; j <= pvlen; j++ {
			//var grp map[interface{}]interface{}
			if p.Shared == -1 {
				dgrp = append(dgrp, int32(vertices.Top/3))
				dgrp = append(dgrp, int32(vertices.Top/3+1))
				dgrp = append(dgrp, int32(vertices.Top/3+2))
			} else {
				ll := grps[p.Shared]
				ll = append(ll, int32(vertices.Top/3))
				ll = append(ll, int32(vertices.Top/3+1))
				ll = append(ll, int32(vertices.Top/3+2))
				grps[p.Shared] = ll
			}
			if idx == 26 {
				fmt.Println(toolkit.JsonString(dgrp), len(dgrp))
				fmt.Println(toolkit.JsonString(grps), len(grps))
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
		fmt.Println("=======")
	}
	geom.Position = vertices.Arr
	geom.Normal = normals.Arr
	if uvs != nil {
		geom.UV = uvs.Arr
	}
	if colors != nil {
		geom.Color = colors.Arr
	}
	// fmt.Println(len(grps))
	for gi := 0; gi < len(grps); gi++ {
		if _, ok := grps[int32(gi)]; !ok {
			grps[int32(gi)] = []int32{} //.Add(int32(gi), []interface{}{})
		}
	}
	// fmt.Println(grps)
	if len(grps) > 0 {
		index := []int32{}
		gbase := 0
		for gi := 0; gi < len(grps); gi++ {
			geom.AddGroup(int32(gbase), int32(len(grps[int32(gi)])), int32(gi))
			gbase += len(grps[int32(gi)])
			kk := grps[int32(gi)]
			for _, val := range kk {
				index = append(index, val)
			}
		}
		geom.AddGroup(int32(gbase), int32(len(dgrp)), int32(len(grps)))
		index = append(index, dgrp...)
		geom.Index = index

	}
	return geom
}

func FromMesh(mesh *Mesh, objectIndex interface{}) *CSG {
	csg := FromGeometry(mesh.Geometry, objectIndex)
	ttvv0 := NewVector(0, 0, 0)
	tmpm3 := NewMatrix3()
	// fmt.Println("mesh.Matrix", mesh.Matrix)
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
	//fmt.Println(a.Polygons, a.Front.Polygons, a.Back.Polygons)
	//fmt.Println(b.Polygons, b.Front.Polygons, b.Back.Polygons)
	// fmt.Println("Before invert")
	// fmt.Println(toolkit.JsonString(m.Polygons))
	a.Invert()
	// fmt.Println(toolkit.JsonString(m.Polygons))
	//fmt.Println(a)
	// os.Exit(-1)
	a.ClipTo(b)
	// fmt.Println(len(a.Polygons))
	// fmt.Println(toolkit.JsonString(a.Polygons))
	b.ClipTo(a)
	// fmt.Println(len(b.Polygons))
	b.Invert()
	b.ClipTo(a)
	b.Invert()
	a.Build(b.AllPolygons())
	a.Invert()
	fmt.Println("a.polygon", len(a.AllPolygons()))
	// fmt.Println(toolkit.JsonString(a.Polygons))
	return FromPolygons(a.AllPolygons())
}
