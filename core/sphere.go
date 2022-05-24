package core

type Sphere struct {
	Center *Vector
	Radius float64
}

func NewSphere(center *Vector, radius float64) *Sphere {
	s := &Sphere{}
	s.Center = center
	s.Radius = radius
	return s
}

func (s *Sphere) Copy(s2 *Sphere) *Sphere {
	if s.Center == nil {
		s.Center = NewVector(0, 0, 0)
	}
	s.Center.Copy(s2.Center)
	s.Radius = s2.Radius
	return s
}
