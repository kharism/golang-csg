package core

const EPSILON = 1e-5

type Plane struct {
	Normal *Vector
	W      float64
}

func NewPlane(Normal *Vector, W float64) *Plane {
	return &Plane{Normal: Normal, W: W}
}

func FromPoints(a, b, c *Vector) *Plane {
	n := b.Sub(a).Cross(c.Sub(a)).Normalize()
	return NewPlane(n, n.Dot(a))
}

func (p *Plane) Clone() *Plane {
	return &Plane{p.Normal, p.W}
}

func (p *Plane) Flip() {
	p.Normal.Negate()
	p.W = -p.W
}

// Split `polygon` by this plane if needed, then put the polygon or polygon
// fragments in the appropriate lists. Coplanar polygons go into either
// `coplanarFront` or `coplanarBack` depending on their orientation with
// respect to this plane. Polygons in front or in back of this plane go into
// either `front` or `back`.
func (p *Plane) SplitPolygon(polygon *Polygon, coplanarFront *[]*Polygon, coplanarBack, front, back *[]*Polygon) {
	const COPLANAR = 0
	const FRONT = 1
	const BACK = 2
	const SPANNING = 3

	polygonType := 0
	vertTypes := []int{}

	for i := 0; i < len(polygon.Vertices); i++ {
		t := p.Normal.Dot(polygon.Vertices[i].Pos) - p.W
		vertType := 0
		if t < EPSILON {
			vertType = BACK
		} else if t > EPSILON {
			vertType = FRONT
		} else {
			vertType = COPLANAR
		}
		polygonType |= vertType
		vertTypes = append(vertTypes, vertType)
	}
	switch polygonType {
	case COPLANAR:
		if p.Normal.Dot(polygon.Plane.Normal) > 0 {
			*coplanarFront = append(*coplanarFront, polygon)
		} else {
			*coplanarBack = append(*coplanarBack, polygon)
		}
	case FRONT:
		*front = append(*front, polygon)
	case BACK:
		*back = append(*back, polygon)
	case SPANNING:
		f := []*Vertex{}
		b := []*Vertex{}
		for i := 0; i < len(polygon.Vertices); i++ {
			j := (i + 1) % len(polygon.Vertices)
			ti := vertTypes[i]
			tj := vertTypes[j]
			vi := polygon.Vertices[i]
			vj := polygon.Vertices[j]
			if ti != BACK {
				f = append(f, vi)
			}
			if ti != FRONT {
				vb := &Vertex{}
				if ti != BACK {
					vb = vi.Clone()
				} else {
					vb = vi
				}
				b = append(b, vb)
			}
			if (ti | tj) == SPANNING {
				t := (p.W - p.Normal.Dot(vi.Pos)) / (p.Normal.Dot(vj.Pos.Clone().Sub(vi.Pos)))
				v := vi.interpolate(vj, t)
				f = append(f, v)
				b = append(b, v.Clone())
			}
		}
		if len(f) >= 3 {
			*front = append(*front, NewPolygon(f, polygon.Shared))
		}
		if len(b) >= 3 {
			*back = append(*back, NewPolygon(b, polygon.Shared))
		}
	}
}
