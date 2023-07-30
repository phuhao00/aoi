package aoi

import "math"

type ListNode struct {
	HubNode
	unit *Unit
	PosX CoordinateVal
	PosZ CoordinateVal
}

func newListNode(entAOINode *Unit, x CoordinateVal, z CoordinateVal) *ListNode {
	eln := new(ListNode)
	eln.unit = entAOINode
	eln.PosX = x
	eln.PosZ = z
	eln.HubNode.Category = UNIT
	return eln
}

func (l *ListNode) SetXPre(node INode) {
	l.XPre = node
}

func (l *ListNode) SetXNext(node INode) {
	l.XNext = node
}

func (l *ListNode) SetZPre(node INode) {
	l.ZPre = node
}

func (l *ListNode) SetZNext(node INode) {
	l.ZNext = node
}

func (l *ListNode) GetXPre() INode {
	return l.HubNode.XPre
}

func (l *ListNode) GetXNext() INode {
	return l.HubNode.XNext
}

func (l *ListNode) GetZPre() INode {
	return l.HubNode.ZPre
}

func (l *ListNode) GetZNext() INode {
	return l.HubNode.ZNext
}

func (l *ListNode) Category() NodeCategory {
	return l.HubNode.Category
}

func (l *ListNode) getEntityID() uint64 {
	return l.unit.ID
}

func (l *ListNode) removeMyself(oldZ CoordinateVal) {
	l.PosZ = math.MaxFloat32
	shuffleZ(l, l.PosX, oldZ)
	l.removeFromRangeList()
}

func (l *ListNode) moveToPos(tgtX CoordinateVal, tgtZ CoordinateVal) {
	oldX, oldZ := l.PosX, l.PosZ
	l.PosX, l.PosZ = tgtX, tgtZ
	shuffleXThenZ(l, oldX, oldZ)
}

func (l *ListNode) removeFromRangeList() {
	if l.XPre != nil {
		l.XPre.SetXNext(l.XNext)
	}

	if l.XNext != nil {
		l.XNext.SetXPre(l.XPre)
	}

	if l.ZPre != nil {
		l.ZPre.SetZNext(l.ZNext)
	}

	if l.ZNext != nil {
		l.ZNext.SetZPre(l.ZPre)
	}

	l.XPre = nil
	l.XNext = nil
	l.ZPre = nil
	l.ZNext = nil
}

func (l *ListNode) X() CoordinateVal {
	return l.PosX
}

func (l ListNode) Z() CoordinateVal {
	return l.PosZ
}
