package aoi

type UnitCategory uint16

type IUint interface {
	ID() uint64
	X() CoordinateVal
	Z() CoordinateVal
	Range() CoordinateVal
	Category() UnitCategory
	OnEnterRange(u IUint)
	OnLeaveRange(u IUint)
}
