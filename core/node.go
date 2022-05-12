package core

type Node struct {
	Polygons []*Polygon
	Plane    *Plane
	Front    *Node
	Back     *Node
}

func NewNode(polygons []*Polygon) *Node {
	newNode := &Node{}
	newNode.Plane = nil
	newNode.Front = nil
	newNode.Back = nil
	newNode.Polygons = []*Polygon{}
	if polygons != nil {
		newNode.Build(polygons)
	}
	return newNode
}

// Build a BSP tree out of `polygons`. When called on an existing tree, the
// new polygons are filtered down to the bottom of the tree and become new
// nodes there. Each set of polygons is partitioned using the first polygon
// (no heuristic is used to pick a good split).
func (n *Node) Build(polygons []*Polygon) {
	if len(polygons) == 0 {
		return
	}
	if n.Plane == nil {
		n.Plane = polygons[0].Plane.Clone()
	}
	FrontPolygon := []*Polygon{}
	BackPolygon := []*Polygon{}

	for i := 0; i < len(polygons); i++ {
		n.Plane.SplitPolygon(polygons[i], &(n.Polygons), &(n.Polygons), &FrontPolygon, &BackPolygon)
	}
	if len(FrontPolygon) > 0 {
		if n.Front == nil {
			n.Front = NewNode(nil)
		}
		n.Front.Build(FrontPolygon)
	}
	if len(BackPolygon) > 0 {
		if n.Back == nil {
			n.Back = NewNode(nil)
		}
		n.Back.Build(BackPolygon)
	}
}

func (n *Node) Clone() *Node {
	node := NewNode(nil)
	node.Plane = n.Plane.Clone()
	node.Back = n.Back.Clone()
	node.Front = n.Front.Clone()
	node.Polygons = []*Polygon{}
	for idx := range n.Polygons {
		node.Polygons = append(node.Polygons, n.Polygons[idx].Clone())
	}
	return node
}
