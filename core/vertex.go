package core

type Vertex struct {
	Pos, Normal, UV, Color *Vector
}

func NewVertex(Pos, Normal, UV, Color *Vector) Vertex {
	return Vertex{Pos, Normal, UV, Color}
}

func (v *Vertex) Clone() *Vertex {
	return &Vertex{v.Pos, v.Normal, v.UV, v.Color}
}

// Invert all orientation-specific data (e.g. vertex normal). Called when the
// orientation of a polygon is flipped.
func (v *Vertex) Flip() {
	v.Normal.Negate()
}

// Create a new vertex between this vertex and `other` by linearly
// interpolating all properties using a parameter of `t`. Subclasses should
// override this to interpolate additional properties.
func (v *Vertex) interpolate(other *Vertex, t float64) *Vertex {
	if v.Color != nil {
		return &Vertex{
			v.Pos.Clone().Lerp(other.Pos, t),
			v.Normal.Clone().Lerp(other.Normal, t),
			v.UV.Clone().Lerp(other.UV, t),
			v.Color.Clone().Lerp(other.Color, t),
		}
	} else {
		return &Vertex{
			v.Pos.Clone().Lerp(other.Pos, t),
			v.Normal.Clone().Lerp(other.Normal, t),
			v.UV.Clone().Lerp(other.UV, t),
			nil,
			//v.Color.Clone().Lerp(other.Color, t),
		}
	}

}
