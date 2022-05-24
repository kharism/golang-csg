package core

type Raycaster struct {
	Origin    Vector
	Direction Vector
	Ray       *Ray
}
type Face struct {
	a             int32
	b             int32
	c             int32
	Normal        *Vector
	MaterialIndex int
}
type Intersection struct {
	Distance  float64
	Point     Vector
	FaceIndex int
	Face      *Face
}

func NewRaycaster(origin, direction *Vector) *Raycaster {
	r := &Raycaster{
		Origin:    *origin,
		Direction: *direction,
		Ray:       NewRay(origin, direction),
	}
	return r
}

//
type IRaycaster interface {
	Raycast(*Raycaster, *[]Intersection)
}

func (r *Raycaster) IntersectObject(object IRaycaster, recursive bool) []Intersection {
	intersection := []Intersection{}
	object.Raycast(r, &intersection)
	return intersection
}
