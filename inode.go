package aoi

type INode interface {
	NodeGetter
	NodeSetter
	Category() NodeCategory
	NodeAction
	NodeVector
}

type NodeAction interface {
	moveToPrevX()
	moveToNextX()
	moveToPrevZ()
	moveToNextZ()
}

type NodeGetter interface {
	GetXPre() INode
	GetXNext() INode
	GetZPre() INode
	GetZNext() INode
}

type NodeSetter interface {
	SetXPre(INode)
	SetXNext(INode)
	SetZPre(INode)
	SetZNext(INode)
}

type NodeVector interface {
	X() CoordinateVal
	Z() CoordinateVal
}

type NodeCategory uint8

const (
	TAIL NodeCategory = iota
	UNIT
	TRIGGER
)
