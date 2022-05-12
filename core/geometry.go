package core

type Geometry struct {
	//Attributes map[string]interface{}
	Index    []int32
	Normal   []float64
	Position []float64
	UV       []float64
	Color    []float64
	Groups   []Group
}

type Group struct {
	Start         int32
	Count         int32
	MaterialIndex int32
}

func NewGeometry() *Geometry {
	geo := &Geometry{}
	geo.Index = nil
	geo.Position = nil
	return geo
}
func (g *Geometry) AddGroup(start, count, materialIndex int) {
	newGroup := Group{}
	g.Groups = append(g.Groups, newGroup)
}
