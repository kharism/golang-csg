package core

type NBuf3 struct {
	Top int
	Arr []float64
}

func NewNBuf3(ct int) *NBuf3 {
	nbuf3 := &NBuf3{}
	nbuf3.Arr = make([]float64, ct)
	return nbuf3
}
func (b *NBuf3) Write(v *Vector) {
	b.Arr[b.Top] = v.X
	b.Top++
	b.Arr[b.Top] = v.Y
	b.Top++
	b.Arr[b.Top] = v.Z
	b.Top++
}

type NBuf2 struct {
	Top int
	Arr []float64
}

func NewNBuf2(ct int) *NBuf2 {
	nbuf2 := &NBuf2{}
	nbuf2.Arr = make([]float64, ct)
	return nbuf2
}
func (b *NBuf2) Write(v *Vector) {
	b.Arr[b.Top] = v.X
	b.Top++
	b.Arr[b.Top] = v.Y
	b.Top++
}
