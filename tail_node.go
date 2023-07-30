package aoi

import "math"

type TailNode struct {
	HubNode
}

func newCLNodeTail() *TailNode {
	tail := new(TailNode)
	tail.Category = TAIL
	return tail
}

func (n *TailNode) isTail() bool {
	return true
}

func (n *TailNode) X() CoordinateVal {
	return math.MaxFloat32
}

func (n *TailNode) Z() CoordinateVal {
	return math.MaxFloat32
}
