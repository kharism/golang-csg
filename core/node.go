package core

import (
	"log"
	"os"
)

type Node struct {
	Polygons []*Polygon
	Plane    *Plane
	Front    *Node
	Back     *Node
}

var Logger = log.New(os.Stdout, ">>", log.Llongfile)

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

func (n *Node) Invert() {
	for i := range n.Polygons {
		n.Polygons[i].Flip()
	}
	if n.Plane != nil {
		n.Plane.Flip()
	}
	if n.Front != nil {
		n.Front.Invert()
	}
	if n.Back != nil {
		n.Back.Invert()
	}
	n.Front, n.Back = n.Back, n.Front
}

// Recursively remove all polygons in `polygons` that are inside this BSP
// tree
func (n *Node) ClipPolygon(polygons []*Polygon) []*Polygon {
	if n.Plane == nil {
		return polygons
	}
	front := []*Polygon{}
	back := []*Polygon{}

	for i := range polygons {
		n.Plane.SplitPolygon(polygons[i], &front, &back, &front, &back)
	}

	if n.Front != nil {
		front = n.Front.ClipPolygon(front)
	}
	if n.Back != nil {
		back = n.Back.ClipPolygon(back)
	} else {
		back = []*Polygon{}
	}
	front = append(front, back...)
	return front
}
func (n *Node) AllPolygons() []*Polygon {
	polygons := n.Polygons
	if n.Front != nil {
		polygons = append(polygons, n.Front.AllPolygons()...)
	}
	if n.Back != nil {
		polygons = append(polygons, n.Back.AllPolygons()...)
	}
	return polygons
}
func (n *Node) ClipTo(bsp *Node) {
	n.Polygons = bsp.ClipPolygon(n.Polygons)
	if n.Front != nil {
		n.Front.ClipTo(bsp)
	}
	if n.Back != nil {
		n.Back.ClipTo(bsp)
	}
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
	// fmt.Println(len(polygons))
	// os.Exit(-1)
	// Logger.Println(n.Plane.Normal, n.Plane.W)
	for i := 0; i < len(polygons); i++ {
		n.Plane.SplitPolygon(polygons[i], &(n.Polygons), &(n.Polygons), &FrontPolygon, &BackPolygon)
		// Logger.Printf("%s\n", tk.JsonString(polygons[i]))
		// Logger.Println("Node Build Split polygon", i, len(n.Polygons), len(FrontPolygon), len(BackPolygon))
	}
	// Logger.Println("FrontPolygon Length", len(FrontPolygon))
	// Logger.Println("BackPolygon Length", len(BackPolygon))
	// os.Exit(-1)
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
