package aoi

type Unit struct {
	ID           uint64
	ListNode     *ListNode
	nextRange    CoordinateVal
	space        *Space
	unitInstance IUint
}

func newUnit(space *Space, unitInstance IUint, x CoordinateVal, z CoordinateVal) *Unit {
	aoiNode := new(Unit)
	aoiNode.ID = unitInstance.ID()
	aoiNode.ListNode = newListNode(aoiNode, x, z)
	aoiNode.space = space
	aoiNode.unitInstance = unitInstance
	return aoiNode
}

func (u *Unit) onEntityEnterRange(entID uint64, rangeID CoordinateVal) {
	u.space.onUnitEnterRange(u.ID, entID, rangeID)
}

func (u *Unit) onEntityLeaveRange(entID uint64, rangeID CoordinateVal) {
	u.space.onUnitLeaveRange(u.ID, entID, rangeID)
}

func (u *Unit) removeMyself() {
	oldZ := u.ListNode.PosZ
	u.ListNode.removeMyself(oldZ)
}

func (u *Unit) moveToPos(tgtX, tgtZ CoordinateVal) {
	u.ListNode.moveToPos(tgtX, tgtZ)
}
