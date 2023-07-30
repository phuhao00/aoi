package aoi

type HubNode struct {
	owner    INode
	XPre     INode
	XNext    INode
	ZPre     INode
	ZNext    INode
	Category NodeCategory
}

func (h *HubNode) X() CoordinateVal {
	//TODO implement me
	panic("implement me")
}

func (h *HubNode) Z() CoordinateVal {
	//TODO implement me
	panic("implement me")
}

func (h *HubNode) moveToPrevX() {
	if h.XNext != nil {
		h.XNext.SetXPre(h.XPre)
	}
	h.XPre.SetXNext(h.XNext)
	h.XNext = h.XPre
	h.XPre = h.XPre.GetXPre()

	if h.XPre != nil {
		h.XPre.SetXNext(h.owner)
	}
	h.XNext.SetXPre(h.owner)
}

func (h *HubNode) moveToNextX() {

	if h.XPre != nil {
		h.XPre.SetXNext(h.XNext)
	}
	h.XNext.SetXPre(h.XPre)

	h.XPre.SetXPre(h.XNext)
	h.XNext = h.XNext.GetXNext()

	if h.XNext != nil {
		h.XNext.SetXPre(h.owner)
	}
	h.XPre.SetXNext(h.owner)
}

func (h *HubNode) moveToPrevZ() {
	if h.ZNext != nil {
		h.ZNext.SetZPre(h.ZPre)
	}
	h.ZPre.SetZNext(h.ZNext)
	h.ZNext = h.ZPre
	h.ZPre = h.ZPre.GetZPre()

	if h.ZPre != nil {
		h.ZPre.SetZNext(h.owner)
	}
	h.ZNext.SetZPre(h.owner)
}

func (h *HubNode) moveToNextZ() {
	if h.ZPre != nil {
		h.ZPre.SetZNext(h.ZNext)
	}
	h.ZNext.SetZPre(h.ZPre)

	h.ZPre.SetZPre(h.ZNext)
	h.ZNext = h.ZNext.GetZNext()

	if h.ZNext != nil {
		h.ZNext.SetZPre(h.owner)
	}
	h.ZPre.SetZNext(h.owner)
}

func (h *HubNode) insertBeforeX(newNode INode) {
	if h.XPre != nil {
		h.XPre.SetXNext(newNode)
	}

	newNode.SetXPre(h.XPre)

	h.XPre = newNode
	newNode.SetXNext(h.owner)
}

func (h *HubNode) insertBeforeZ(newNode INode) {
	if h.ZPre != nil {
		h.ZPre.SetZNext(newNode)
	}

	newNode.SetZPre(h.ZPre)

	h.ZPre = newNode
	newNode.SetZNext(h.owner)
}
