package gopdf

type Box struct {
	Left, Top, Right, Bottom float64
	unitOverride             defaultUnitConfig
}

// UnitsToPoints converts the box coordinates to Points. When this is called it is assumed the values of the box are in Units
func (box *Box) UnitsToPoints(t int) (b *Box) {
	if box == nil {
		return
	}

	unitCfg := defaultUnitConfig{Unit: t}
	if box.unitOverride.getUnit() != UnitUnset {
		unitCfg = box.unitOverride
	}

	b = &Box{
		Left:   box.Left,
		Top:    box.Top,
		Right:  box.Right,
		Bottom: box.Bottom,
	}
	unitsToPointsVar(unitCfg, &b.Left, &b.Top, &b.Right, &b.Bottom)
	return
}

func (box *Box) unitsToPoints(unitCfg unitConfigurator) (b *Box) {
	if box == nil {
		return
	}

	if box.unitOverride.getUnit() != UnitUnset {
		unitCfg = box.unitOverride
	}

	b = &Box{
		Left:   box.Left,
		Top:    box.Top,
		Right:  box.Right,
		Bottom: box.Bottom,
	}
	unitsToPointsVar(unitCfg, &b.Left, &b.Top, &b.Right, &b.Bottom)
	return
}
