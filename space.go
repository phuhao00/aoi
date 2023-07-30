package aoi

import (
	"errors"
	"math"
)

type EventCategory uint16

const (
	EVENT_ENTER EventCategory = 1 << iota
	EVENT_LEAVE               = 1 << iota
	EVENT_ALL                 = EVENT_ENTER | EVENT_LEAVE
)

type Space struct {
	tailNode     *TailNode
	unitNodesMap map[uint64]*Unit

	oldAoiUnitsCache map[uint64]IUint

	newAoiUnitsCache map[uint64]IUint

	diffAoiUnitsCache map[uint64]IUint
}

func NewSpace() *Space {
	space := new(Space)
	space.tailNode = newCLNodeTail()
	space.unitNodesMap = make(map[uint64]*Unit, 64)
	space.oldAoiUnitsCache = make(map[uint64]IUint, 32)
	space.newAoiUnitsCache = make(map[uint64]IUint, 32)
	space.diffAoiUnitsCache = make(map[uint64]IUint, 32)
	return space
}

func (s *Space) clearOldAoiUnitsCache() {
	if len(s.oldAoiUnitsCache) > 0 {
		for k := range s.oldAoiUnitsCache {
			delete(s.oldAoiUnitsCache, k)
		}
	}
}

func (s *Space) clearNewAoiUnitsCache() {

	if len(s.newAoiUnitsCache) > 0 {
		for k := range s.newAoiUnitsCache {
			delete(s.newAoiUnitsCache, k)
		}
	}
}

func (s *Space) clearDiffAoiUnitsCache() {
	if len(s.diffAoiUnitsCache) > 0 {
		for k := range s.diffAoiUnitsCache {
			delete(s.diffAoiUnitsCache, k)
		}
	}
}

func (s *Space) onUnitEnterRange(whoID uint64, enteringID uint64, rangeID CoordinateVal) {
	who, ok := s.unitNodesMap[whoID]
	if !ok {
		return
	}
	enteringEntity, ok := s.unitNodesMap[enteringID]
	if !ok {
		return
	}
	who.unitInstance.OnEnterRange(enteringEntity.unitInstance)
}

func (s *Space) onUnitLeaveRange(whoID uint64, leavingID uint64, rangeID CoordinateVal) {
	who, ok := s.unitNodesMap[whoID]
	if !ok {
		return
	}
	leaving, ok := s.unitNodesMap[leavingID]
	if !ok {
		return
	}
	who.unitInstance.OnLeaveRange(leaving.unitInstance)
}

// AddUnit : add an Entity to this AOI space
func (s *Space) AddUnit(aoiUnit IUint) (map[uint64]IUint, error) {
	entID := aoiUnit.ID()
	_, ok := s.unitNodesMap[entID]
	if ok {
		return nil, errors.New("entity already in aoi")
	}

	x := aoiUnit.X()
	z := aoiUnit.Z()

	if !IsValidAoiCLPosXZ(x, z) {
		return nil, errors.New("entity pos not valid")
	}

	aoiNode := newUnit(s, aoiUnit, x, z)
	s.unitNodesMap[entID] = aoiNode

	s.tailNode.insertBeforeX(aoiNode.ListNode)
	s.tailNode.insertBeforeZ(aoiNode.ListNode)
	shuffleXThenZ(aoiNode.ListNode, math.MaxFloat32, math.MaxFloat32)

	aoiEntities := make(map[uint64]IUint)

	err := s.getUnitsInRange(aoiEntities, aoiUnit, aoiUnit.Range(), true)
	if err != nil {
		return nil, err
	}

	return aoiEntities, nil
}

// RemoveUnit : remove an Entity from this AOI space
func (s *Space) RemoveUnit(aoiUnit IUint) error {
	entID := aoiUnit.ID()
	aoiNode, ok := s.unitNodesMap[entID]
	if !ok {
		return errors.New("u not in aoi")
	}

	s.clearNewAoiUnitsCache()
	err := s.getUnitsInRange(s.newAoiUnitsCache, aoiUnit, aoiUnit.Range(), true)
	if err != nil {
		return err
	}
	for _, u := range s.newAoiUnitsCache {
		if aoiUnit.ID() == u.ID() {
			continue
		}
		u.OnLeaveRange(aoiUnit)
	}

	aoiNode.removeMyself()
	delete(s.unitNodesMap, entID)

	return nil
}

// MoveUnit : move a entity to a new position, and auto recalc the aoi
func (s *Space) MoveUnit(aoiUnit IUint, tgtX CoordinateVal, tgtZ CoordinateVal) error {
	entID := aoiUnit.ID()
	aoiNode, ok := s.unitNodesMap[entID]
	if !ok {
		return errors.New("entity not in aoi")
	}

	if !IsValidAoiCLPosXZ(tgtX, tgtZ) {
		return errors.New("entity pos not valid")
	}

	s.clearOldAoiUnitsCache()

	err := s.getUnitsInRange(s.oldAoiUnitsCache, aoiUnit, aoiUnit.Range(), true)
	if err != nil {
		return err
	}

	aoiNode.moveToPos(tgtX, tgtZ)

	s.clearNewAoiUnitsCache()

	err = s.getUnitsInRange(s.newAoiUnitsCache, aoiUnit, aoiUnit.Range(), true)
	if err != nil {
		return err
	}

	s.clearDiffAoiUnitsCache()
	difference(s.newAoiUnitsCache, s.oldAoiUnitsCache, s.diffAoiUnitsCache)
	for _, entity := range s.diffAoiUnitsCache {
		if aoiUnit.ID() == entity.ID() {
			continue
		}
		aoiUnit.OnEnterRange(entity)
	}

	s.clearDiffAoiUnitsCache()
	difference(s.oldAoiUnitsCache, s.newAoiUnitsCache, s.diffAoiUnitsCache)
	for _, entity := range s.diffAoiUnitsCache {
		if aoiUnit.ID() == entity.ID() {
			continue
		}
		aoiUnit.OnLeaveRange(entity)
		entity.OnLeaveRange(aoiUnit)
	}

	return nil
}

// UnitsInRange : get entities in specified range of this AOI space
func (s *Space) UnitsInRange(aoiUnit IUint, r CoordinateVal, includeThis bool) ([]uint64, error) {
	entID := aoiUnit.ID()
	aoiNode, ok := s.unitNodesMap[entID]
	if !ok {
		return nil, errors.New("entity not in aoi")
	}

	if r <= 0 {
		return nil, errors.New("r should be greater than 0")
	}

	listNode := aoiNode.ListNode
	centerPosX, centerPoxZ := listNode.X, listNode.Z

	res := []uint64{}
	if includeThis {
		res = append(res, entID)
	}

	var xDist, zDist CoordinateVal

	xDist, zDist = 0, 0
	cursor := listNode.XNext
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.Category() == UNIT {
			xDist, zDist = Abs(cursor.X()-centerPosX), Abs(cursor.Z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				res = append(res, cursor.(*ListNode).getEntityID())
			}
		}
		cursor = cursor.GetXNext()
	}

	xDist, zDist = 0, 0
	cursor = listNode.XPre
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.Category() == UNIT {
			xDist, zDist = Abs(cursor.X()-centerPosX), Abs(cursor.Z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				res = append(res, cursor.(*ListNode).getEntityID())
			}
		}
		cursor = cursor.GetXPre()
	}

	return res, nil
}

// getUnitsInRange : get entities in specified range of this AOI space
func (s *Space) getUnitsInRange(retMap map[uint64]IUint, aoiUnit IUint, r CoordinateVal, includeThis bool) error {
	entID := aoiUnit.ID()
	aoiNode, ok := s.unitNodesMap[entID]
	if !ok {
		return errors.New("entity not in aoi")
	}

	if r <= 0 {
		return errors.New("r should be greater than 0")
	}

	listNode := aoiNode.ListNode
	centerPosX, centerPoxZ := listNode.X, listNode.Z

	if includeThis {
		retMap[entID] = aoiUnit
	}

	var xDist, zDist CoordinateVal

	xDist, zDist = 0, 0
	cursor := listNode.XNext
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.Category() == UNIT {
			xDist, zDist = Abs(cursor.X()-centerPosX), Abs(cursor.Z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				entity, ok := s.unitNodesMap[cursor.(*ListNode).getEntityID()]
				if ok {
					retMap[entity.ID] = entity.unitInstance
				}
			}
		}
		cursor = cursor.GetXNext()
	}

	xDist, zDist = 0, 0
	cursor = listNode.XPre
	for {
		if cursor == nil || xDist > r {
			break
		}
		if cursor.Category() == UNIT {
			xDist, zDist = Abs(cursor.X()-centerPosX), Abs(cursor.Z()-centerPoxZ)
			if xDist <= r && zDist <= r {
				entity, ok := s.unitNodesMap[cursor.(*ListNode).getEntityID()]
				if ok {
					retMap[entity.ID] = entity.unitInstance
				}
			}
		}
		cursor = cursor.GetXPre()
	}

	return nil
}
