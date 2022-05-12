package core

type Polygon struct {
	Plane    *Plane
	Shared   interface{}
	Vertices []*Vertex
}

func NewPolygon(vertices []*Vertex, shared interface{}) *Polygon {
	p := &Polygon{}
	p.Shared = shared
	p.Vertices = vertices
	p.Plane = FromPoints(vertices[0].Pos, vertices[1].Pos, vertices[2].Pos)
	return p
}
func (p *Polygon) Clone() *Polygon {
	newVertices := []*Vertex{}
	for _, v := range p.Vertices {
		newVertices = append(newVertices, v)
	}
	newPolygon := &Polygon{}
	newPolygon.Shared = p.Shared
	newPolygon.Vertices = newVertices
	newPolygon.Plane = p.Plane
	return newPolygon
}
func ReverseVertexArr(s []*Vertex) []*Vertex {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
func (p *Polygon) Flip() {
	p.Vertices = ReverseVertexArr(p.Vertices)
	for idx := range p.Vertices {
		p.Vertices[idx].Flip()
	}
	p.Plane.Flip()
}
