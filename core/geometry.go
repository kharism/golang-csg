package core

import (
	"bytes"
	"fmt"
)

type Geometry struct {
	//Attributes map[string]interface{}
	Index    []int32
	Normal   []float64
	Position []float64
	UV       []float64
	Color    []float64
	Groups   []Group
}

func (g *Geometry) Clone() *Geometry {
	newGeom := &Geometry{}
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
func (g *Geometry) AddGroup(start, count, materialIndex int32) {
	newGroup := Group{}
	newGroup.Start = start
	newGroup.Count = count
	newGroup.MaterialIndex = materialIndex
	g.Groups = append(g.Groups, newGroup)
}
