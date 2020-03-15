package gopdf

// Margins type.
type Margins struct {
	Left, Top, Right, Bottom float64
}

// SetLeftMargin sets left margin.
func (gp *GoPdf) SetLeftMargin(margin float64) {
	gp.UnitsToPointsVar(&margin)
	gp.margins.Left = margin
}

// SetTopMargin sets top margin.
func (gp *GoPdf) SetTopMargin(margin float64) {
	gp.UnitsToPointsVar(&margin)
	gp.margins.Top = margin
}

// SetMargins defines the left, top, right and bottom margins. By default, they equal 1 cm. Call this method to change them.
func (gp *GoPdf) SetMargins(left, top, right, bottom float64) {
	gp.UnitsToPointsVar(&left, &top, &right, &bottom)
	gp.margins = Margins{left, top, right, bottom}
}

// SetMarginLeft sets the left margin
func (gp *GoPdf) SetMarginLeft(margin float64) {
	gp.margins.Left = gp.UnitsToPoints(margin)
}

// SetMarginTop sets the top margin
func (gp *GoPdf) SetMarginTop(margin float64) {
	gp.margins.Top = gp.UnitsToPoints(margin)
}

// SetMarginRight sets the right margin
func (gp *GoPdf) SetMarginRight(margin float64) {
	gp.margins.Right = gp.UnitsToPoints(margin)
}

// SetMarginBottom set the bottom margin
func (gp *GoPdf) SetMarginBottom(margin float64) {
	gp.margins.Bottom = gp.UnitsToPoints(margin)
}

// Margins gets the current margins, The margins will be converted back to the documents units. Returned values will be in the following order Left, Top, Right, Bottom
func (gp *GoPdf) Margins() (float64, float64, float64, float64) {
	return gp.PointsToUnits(gp.margins.Left),
		gp.PointsToUnits(gp.margins.Top),
		gp.PointsToUnits(gp.margins.Right),
		gp.PointsToUnits(gp.margins.Bottom)
}

// MarginLeft returns the left margin.
func (gp *GoPdf) MarginLeft() float64 {
	return gp.PointsToUnits(gp.margins.Left)
}

// MarginTop returns the top margin.
func (gp *GoPdf) MarginTop() float64 {
	return gp.PointsToUnits(gp.margins.Top)
}

// MarginRight returns the right margin.
func (gp *GoPdf) MarginRight() float64 {
	return gp.PointsToUnits(gp.margins.Right)
}

// MarginBottom returns the bottom margin.
func (gp *GoPdf) MarginBottom() float64 {
	return gp.PointsToUnits(gp.margins.Bottom)
}
