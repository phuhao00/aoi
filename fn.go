package aoi

import "math"

func shuffleXThenZ(thisNode INode, oldX CoordinateVal, oldZ CoordinateVal) {
	shuffleX(thisNode, oldX, oldZ)
	shuffleZ(thisNode, oldX, oldZ)
}

func shuffleX(thisNode INode, oldX CoordinateVal, oldZ CoordinateVal) {
	thisPos := thisNode.X()
	for {
		prevNode := thisNode.GetXPre()
		if prevNode == nil {
			break
		}
		prevPos := prevNode.X()
		if thisPos < prevPos {
			thisNode.moveToPrevX()
		} else {
			break
		}
	}

	for {
		nextNode := thisNode.GetXNext()
		if nextNode == nil {
			break
		}
		nextPos := nextNode.(*ListNode).X
		if thisPos > nextPos {
			thisNode.moveToNextX()
		} else {
			break
		}
	}
}

func shuffleZ(thisNode INode, oldX CoordinateVal, oldZ CoordinateVal) {
	thisPos := thisNode.Z()
	for {
		prevNode := thisNode.GetZPre()
		if prevNode == nil {
			break
		}
		prevPos := prevNode.Z()
		if thisPos < prevPos {
			thisNode.moveToPrevZ()
		} else {
			break
		}
	}

	for {
		nextNode := thisNode.GetZNext()
		if nextNode == nil {
			break
		}
		nextPos := nextNode.(*ListNode).Z
		if thisPos > nextPos {
			thisNode.moveToNextZ()
		} else {
			break
		}
	}
}

func difference(m1, m2, diff map[uint64]IUint) {
	for key, value := range m1 {
		diff[key] = value
	}

	for key := range m2 {
		_, ok := diff[key]
		if ok {
			delete(diff, key)
		}
	}
}

const MAX_VALID_POS_IN_AOI CoordinateVal = math.MaxFloat32

func IsValidAoiCLPosXZ(x, z CoordinateVal) bool {
	return -MAX_VALID_POS_IN_AOI < x && x < MAX_VALID_POS_IN_AOI &&
		-MAX_VALID_POS_IN_AOI < z && z < MAX_VALID_POS_IN_AOI
}

// Abs is a float64 Abs function wrapper
func Abs(f CoordinateVal) CoordinateVal {
	return CoordinateVal(math.Abs(float64(f)))
}

// Max is a float64 function wrapper
func Max(m1 CoordinateVal, m2 CoordinateVal) CoordinateVal {
	return CoordinateVal(math.Max(float64(m1), float64(m2)))
}
