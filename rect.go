package gopdf

type Rect struct {
	W            float64
	H            float64
	unitOverride int
}

// PointsToUnits converts the rectanlges width and height to Units. When this is called it is assumed the values of the rectangle are in Points
func (rect *Rect) PointsToUnits(t int) (r *Rect) {
	if rect == nil {
		return
	}

	if rect.unitOverride != Unit_Unset {
		t = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	PointsToUnitsVar(t, &r.W, &r.H)
	return
}

// UnitsToPoints converts the rectanlges width and height to Points. When this is called it is assumed the values of the rectangle are in Units
func (rect *Rect) UnitsToPoints(t int) (r *Rect) {
	if rect == nil {
		return
	}

	if rect.unitOverride != Unit_Unset {
		t = rect.unitOverride
	}

	r = &Rect{W: rect.W, H: rect.H}
	UnitsToPointsVar(t, &r.W, &r.H)
	return
}
