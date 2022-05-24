package core

import (
	"math"
)

type Mesh struct {
	*Object3D
	Geometry *Geometry
	//material
}

func NewMesh(geometry *Geometry) *Mesh {
	pp := NewObject3D()
	mesh := &Mesh{}
	mesh.Object3D = pp

	mesh.Geometry = geometry

	return mesh
}
func (m *Mesh) Raycast(r *Raycaster, intersects *[]Intersection) {
	if m.Geometry.BoundingSphere == nil {
		m.Geometry.ComputeBoundingSphere()
	}
	// material := m.Material
	sphere := Sphere{}
	sphere.Copy(m.Geometry.BoundingSphere)
	sphere.Center.ApplyMatrix4(m.MatrixWorld)
	if r.Ray.IntersectsSphere(&sphere) == false {
		return
	}
	_inverseMatrix := NewMatrix4()
	_inverseMatrix.Copy(m.MatrixWorld).Invert()
	// fmt.Println("_inverseMatrix", _inverseMatrix)
	_ray := Ray{}
	_ray.Copy(r.Ray).ApplyMatrix4(_inverseMatrix)

	// Check boundingBox before continuing

	if m.Geometry.BoundingBox != nil {

		if _ray.IntersectsBox(m.Geometry.BoundingBox) == false {
			return
		}
	}
	// intersection := []Intersection{}
	index := m.Geometry.Index
	position := m.Geometry.Position
	morphPosition := m.Geometry.MorphAttribute["position"]  //MorphAttributes.position
	morphTargetsRelative := m.Geometry.MorphTargetsRelative //morphTargetsRelative
	uv := m.Geometry.UV
	// uv2 := geometry.attributes.uv2
	// groups := m.Geometry.Groups
	drawRange := m.Geometry.DrawRange
	if index != nil {
		start := math.Max(0, float64(drawRange.Start))
		end := math.Min(float64(len(index)), (float64(drawRange.Start) + drawRange.Count))
		il := end
		// fmt.Println("Start", start, "end", end, int(il), len(index))
		for i := int(start); i < int(il); i += 3 {

			a := index[i] //.getX(i)
			b := index[(i + 1)]
			c := index[(i + 2)]

			intersection := checkBufferGeometryIntersection(m, r, &_ray, position, morphPosition, morphTargetsRelative, uv, a, b, c)

			// fmt.Println(index, i, a, b, c, intersection)
			if intersection != nil {

				intersection.FaceIndex = int(math.Floor(float64(i) / 3)) // triangle number in indexed buffer semantics
				*intersects = append(*intersects, *intersection)

			}

		}
	}
}
func CheckIntersection(object *Mesh, raycaster *Raycaster, ray *Ray, pA, pB, pC, point *Vector) *Intersection {

	intersect := &Vector{}

	intersect = ray.intersectTriangle(pA, pB, pC, false, point)

	if intersect == nil {
		return nil
	}
	_intersectionPointWorld := Vector{}
	_intersectionPointWorld.Copy(point)
	_intersectionPointWorld.ApplyMatrix4(object.MatrixWorld)

	distance := raycaster.Ray.Origin.DistanceTo(&_intersectionPointWorld)

	// if distance < raycaster.near || distance > raycaster.far {
	// 	return nil
	// }

	return &Intersection{
		Distance: distance,
		Point:    _intersectionPointWorld,
	}

}
func checkBufferGeometryIntersection(object *Mesh, raycaster *Raycaster, ray *Ray, position []float64, morphPosition [][]float64, morphTargetsRelative bool, uv []float64, a, b, c int32) *Intersection {
	_vA := Vector{}
	_vB := Vector{}
	_vC := Vector{}
	_vA.Set(position[a*3], position[a*3+1], position[a*3+2])
	_vB.Set(position[b*3], position[b*3+1], position[b*3+2])
	_vC.Set(position[c*3], position[c*3+1], position[c*3+2])

	// morphInfluences := object.MorphTargetInfluences
	_intersectionPoint := Vector{}
	intersection := CheckIntersection(object, raycaster, ray, &_vA, &_vB, &_vC, &_intersectionPoint)

	if intersection != nil {

		// if ( uv ) {

		// 	_uvA.fromBufferAttribute( uv, a );
		// 	_uvB.fromBufferAttribute( uv, b );
		// 	_uvC.fromBufferAttribute( uv, c );

		// 	intersection.uv = Triangle.getUV( _intersectionPoint, _vA, _vB, _vC, _uvA, _uvB, _uvC, new Vector2() );

		// }

		// face := {
		// 	a: a,
		// 	b: b,
		// 	c: c,
		// 	normal: new Vector3(),
		// 	materialIndex: 0
		// };
		face := Face{a, b, c, &Vector{}, 0}

		GetNormal(&_vA, &_vB, &_vC, face.Normal)

		intersection.Face = &face

	}
	return intersection
}
func (m *Mesh) Clone() *Mesh {
	jj := m.Object3D.Clone()
	newMesh := &Mesh{Object3D: jj}
	newMesh.Geometry = m.Geometry.Clone()
	return m
}
