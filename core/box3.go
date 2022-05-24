package core

import (
	"math"
)

type Box3 struct {
	Min *Vector
	Max *Vector
}

func NewBox(min, max *Vector) *Box3 {
	b := &Box3{}
	b.Max = max
	b.Min = min
	return b
}
func NewBoxAutoFill() *Box3 {
	b := &Box3{}
	b.Max = NewVector(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)
	b.Min = NewVector(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	return b
}
func (this *Box3) MakeEmpty() *Box3 {
	this.Min.X = math.MaxFloat64
	this.Min.Y = math.MaxFloat64
	this.Min.Z = math.MaxFloat64
	this.Max.X = -math.MaxFloat64
	this.Max.Y = -math.MaxFloat64
	this.Max.Z = -math.MaxFloat64

	return this
}
func (b *Box3) ExpandByPoint(v *Vector) *Box3 {
	b.Max.Max(v)
	b.Min.Min(v)
	return b
}
func (this *Box3) IsEmpty() bool {

	// this is a more robust check for empty than ( volume <= 0 ) because volume can get positive with two negative axes

	return (this.Max.X < this.Min.X) || (this.Max.Y < this.Min.Y) || (this.Max.Z < this.Min.Z)

}
func (b *Box3) GetCenter(target *Vector) *Vector {
	if b.IsEmpty() {
		return target.Set(0, 0, 0)
	} else {
		return target.Set(b.Min.X, b.Min.Y, b.Min.Z).Add(b.Max).MultiplyScalar(0.5)
	}
}
func (b *Box3) SetFromBufferAttribute(attribute []float64) *Box3 {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	minZ := math.MaxFloat64

	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64
	maxZ := -math.MaxFloat64

	l := len(attribute) / 3
	// fmt.Println("attr len", len(attribute)/3)
	for i := 0; i < l; i += 1 {
		x := attribute[i*3]
		y := attribute[i*3+1]
		z := attribute[i*3+2]
		// fmt.Println("SetFromAttr", x, y, z)
		minX = math.Min(x, minX)
		minY = math.Min(y, minY)
		minZ = math.Min(z, minZ)

		maxX = math.Max(x, maxX)
		maxY = math.Max(y, maxY)
		maxZ = math.Max(z, maxZ)
	}
	b.Max.X = maxX
	b.Max.Y = maxY
	b.Max.Z = maxZ

	b.Min.X = minX
	b.Min.Y = minY
	b.Min.Z = minZ

	return b
}
