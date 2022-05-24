package core

import "math"

type Vector struct {
	X, Y, Z float64
}

func NewVector(x, y, z float64) *Vector {
	return &Vector{x, y, z}
}
func (this *Vector) LengthSq() float64 {

	return this.X*this.X + this.Y*this.Y + this.Z*this.Z

}
func (v *Vector) DistanceTo(d *Vector) float64 {
	return math.Sqrt(v.DistanceToSquared(d))
}
func (v *Vector) CrossVectors(a, b *Vector) *Vector {
	ax := a.X
	ay := a.Y
	az := a.Z

	bx := b.X
	by := b.Y
	bz := b.Z

	v.X = ay*bz - az*by
	v.Y = az*bx - ax*bz
	v.Z = ax*by - ay*bx

	return v
}
func (v *Vector) Copy(a *Vector) *Vector {
	v.X = a.X
	v.Y = a.Y
	v.Z = a.Z
	return v
}
func (v *Vector) Set(x, y, z float64) *Vector {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}
func (this *Vector) DistanceToSquared(v *Vector) float64 {
	dx := this.X - v.X
	dy := this.Y - v.Y
	dz := this.Z - v.Z
	return dx*dx + dy*dy + dz*dz
}

func (this *Vector) TransformDirection(m *Matrix4) *Vector {

	// input: THREE.Matrix4 affine matrix
	// vector interpreted as a direction

	x := this.X
	y := this.Y
	z := this.Z
	e := m.elements

	this.X = e[0]*x + e[4]*y + e[8]*z
	this.Y = e[1]*x + e[5]*y + e[9]*z
	this.Z = e[2]*x + e[6]*y + e[10]*z

	return this.Normalize()

}
func (v *Vector) SubVectors(a, b *Vector) *Vector {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}
func (v *Vector) MultiplyScalar(f float64) *Vector {
	v.X *= f
	v.Y *= f
	v.Z *= f

	return v
}
func (v *Vector) Max(a *Vector) *Vector {
	v.X = math.Max(v.X, a.X)
	v.Y = math.Max(v.Y, a.Y)
	v.Z = math.Max(v.Z, a.Z)
	return v
}
func (v *Vector) Min(a *Vector) *Vector {
	v.X = math.Min(v.X, a.X)
	v.Y = math.Min(v.Y, a.Y)
	v.Z = math.Min(v.Z, a.Z)
	return v
}
func (v *Vector) Clone() *Vector {
	return &Vector{v.X, v.Y, v.Z}
}

func (v *Vector) Negate() *Vector {
	v.X *= -1
	v.Y *= -1
	v.Z *= -1
	return v
}
func (v *Vector) Add(x *Vector) *Vector {
	v.X += x.X
	v.Y += x.Y
	v.Z += x.Z
	return v
}

func (v *Vector) Sub(x *Vector) *Vector {
	v.X -= x.X
	v.Y -= x.Y
	v.Z -= x.Z
	return v
}

func (v *Vector) Times(x float64) *Vector {
	v.X *= x
	v.Y *= x
	v.Z *= x
	return v
}

func (v *Vector) DivideBy(x float64) *Vector {
	v.X /= x
	v.Y /= x
	v.Z /= x
	return v
}

func (v *Vector) Lerp(a *Vector, t float64) *Vector {
	x := NewVector(a.X, a.Y, a.Z)
	return v.Add(x.Sub(v).Times(t))
}

func (v *Vector) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}
func (v *Vector) Unit() *Vector {
	return v.DivideBy(v.Length())
}
func (v *Vector) Normalize() *Vector {
	return v.Unit()
}
func (v *Vector) Cross(b *Vector) *Vector {
	a := v.Clone()

	v.X = a.Y*b.Z - a.Z*b.Y
	v.Y = a.Z*b.X - a.X*b.Z
	v.Z = a.X*b.Y - a.Y*b.X
	return v
}
func (v *Vector) Dot(b *Vector) float64 {
	return v.X*b.X + v.Y*b.Y + v.Z*b.Z
}
func (v *Vector) SetFromMatrix3Column(m *Matrix3, index int) *Vector {
	return v.FromArray(m.elements, index*3)
}
func (this *Vector) ApplyMatrix3(m *Matrix4) *Vector {
	x := this.X
	y := this.Y
	z := this.Z

	e := m.elements

	this.X = e[0]*x + e[3]*y + e[6]*z
	this.Y = e[1]*x + e[4]*y + e[7]*z
	this.Z = e[2]*x + e[5]*y + e[8]*z

	return this
}
func (v *Vector) ApplyMatrix4(m *Matrix4) *Vector {
	x := v.X
	y := v.Y
	z := v.Z

	e := m.elements
	w := 1 / (e[3]*x + e[7]*y + e[11]*z + e[15])
	v.X = (e[0]*x + e[4]*y + e[8]*z + e[12]) * w
	v.Y = (e[1]*x + e[5]*y + e[9]*z + e[13]) * w
	v.Z = (e[2]*x + e[6]*y + e[10]*z + e[14]) * w
	return v
}
func (v *Vector) FromArray(arr []float64, offset int) *Vector {
	v.X = arr[offset]
	v.Y = arr[offset+1]
	v.Z = arr[offset+2]

	return v
}
