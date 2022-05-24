package core

type Ray struct {
	Origin    *Vector
	Direction *Vector
}

func NewRay(origin, direction *Vector) *Ray {
	r := &Ray{}
	r.Origin = origin
	r.Direction = direction
	return r
}
func (r *Ray) At(t float64, target *Vector) *Vector {
	return target.Copy(r.Direction).MultiplyScalar(t).Add(r.Origin)
}

func (r *Ray) Copy(g *Ray) *Ray {
	r.Origin = g.Origin.Clone()
	r.Direction = g.Direction.Clone()
	return r
}

func (r *Ray) IntersectsBox(box *Box3) bool {
	_vector := NewVector(0, 0, 0)
	return r.IntersectsBoxVector(box, _vector) != nil
}
func (r *Ray) IntersectsBoxVector(box *Box3, target *Vector) *Vector {
	var tmin float64
	var tmax float64
	var tymin float64
	var tymax float64
	var tzmin float64
	var tzmax float64

	invdirx := 1 / r.Direction.X
	invdiry := 1 / r.Direction.Y
	invdirz := 1 / r.Direction.Z

	origin := r.Origin

	if invdirx >= 0 {
		tmin = (box.Min.X - origin.X) * invdirx
		tmax = (box.Max.X - origin.X) * invdirx
	} else {
		tmin = (box.Max.X - origin.X) * invdirx
		tmax = (box.Min.X - origin.X) * invdirx
	}

	if invdiry >= 0 {
		tymin = (box.Min.Y - origin.Y) * invdiry
		tymax = (box.Max.Y - origin.Y) * invdiry
	} else {
		tymin = (box.Max.Y - origin.Y) * invdiry
		tymax = (box.Min.Y - origin.Y) * invdiry
	}
	if (tmin > tymax) || (tymin > tmax) {
		return nil
	}

	// These lines also handle the case where tmin or tmax is NaN
	// (result of 0 * Infinity). x !== x returns true if x is NaN

	if tymin > tmin || tmin != tmin {
		tmin = tymin
	}
	if tymax < tmax || tmax != tmax {
		tmax = tymax
	}
	if invdirz >= 0 {
		tzmin = (box.Min.Z - origin.Z) * invdirz
		tzmax = (box.Max.Z - origin.Z) * invdirz
	} else {
		tzmin = (box.Max.Z - origin.Z) * invdirz
		tzmax = (box.Min.Z - origin.Z) * invdirz
	}

	if (tmin > tzmax) || (tzmin > tmax) {
		return nil
	}

	if tzmin > tmin || tmin != tmin {
		tmin = tzmin
	}

	if tzmax < tmax || tmax != tmax {
		tmax = tzmax
	}

	//return point closest to the ray (positive side)

	if tmax < 0 {
		return nil
	}
	if tmin >= 0 {
		return r.At(tmin, target)
	} else {
		return r.At(tmax, target)
	}
	// return this.at(  ? tmin : tmax, target );
}

func (r *Ray) IntersectsSphere(s *Sphere) bool {
	return r.DistanceSqToPoint(s.Center) <= (s.Radius * s.Radius)
}
func (r *Ray) ApplyMatrix4(matrix4 *Matrix4) *Ray {

	r.Origin.ApplyMatrix4(matrix4)
	r.Direction.TransformDirection(matrix4)

	return r

}
func (r *Ray) intersectTriangle(a, b, c *Vector, backfaceCulling bool, target *Vector) *Vector {
	// Compute the offset origin, edges, and normal.

	// from http://www.geometrictools.com/GTEngine/Include/Mathematics/GteIntrRay3Triangle3.h
	_edge1 := Vector{}
	_edge2 := Vector{}
	_normal := Vector{}
	_edge1.SubVectors(b, a)
	_edge2.SubVectors(c, a)
	_normal.CrossVectors(&_edge1, &_edge2)

	// Solve Q + t*D = b1*E1 + b2*E2 (Q = kDiff, D = ray direction,
	// E1 = kEdge1, E2 = kEdge2, N = Cross(E1,E2)) by
	//   |Dot(D,N)|*b1 = sign(Dot(D,N))*Dot(D,Cross(Q,E2))
	//   |Dot(D,N)|*b2 = sign(Dot(D,N))*Dot(D,Cross(E1,Q))
	//   |Dot(D,N)|*t = -sign(Dot(D,N))*Dot(Q,N)
	DdN := r.Direction.Dot(&_normal)

	// fmt.Println("_normal", _normal, "DdN", DdN)
	// let sign;
	sign := float64(0)

	if DdN > 0 {
		if backfaceCulling {
			return nil
		}
		sign = 1
	} else if DdN < 0 {
		sign = -1
		DdN = -DdN
	} else {
		return nil
	}
	_diff := Vector{}
	_diff.SubVectors(r.Origin, a)
	DdQxE2 := sign * r.Direction.Dot(_edge2.CrossVectors(&_diff, &_edge2))

	// b1 < 0, no intersection
	if DdQxE2 < 0 {
		return nil
	}

	DdE1xQ := sign * r.Direction.Dot(_edge1.Cross(&_diff))

	// b2 < 0, no intersection
	if DdE1xQ < 0 {
		return nil
	}

	// b1+b2 > 1, no intersection
	if DdQxE2+DdE1xQ > DdN {
		return nil
	}

	// Line intersects triangle, check if ray does.
	QdN := -sign * _diff.Dot(&_normal)

	// t < 0, no intersection
	if QdN < 0 {
		return nil
	}

	// Ray intersects triangle.
	return r.At(QdN/DdN, target)
}
func (r *Ray) DistanceSqToPoint(point *Vector) float64 {
	_vector := Vector{}
	directionDistance := _vector.SubVectors(point, r.Origin).Dot(r.Direction)
	// point behind the ray

	if directionDistance < 0 {

		return r.Origin.DistanceToSquared(point)

	}

	_vector.Copy(r.Direction).MultiplyScalar(directionDistance).Add(r.Origin)

	return _vector.DistanceToSquared(point)
}
