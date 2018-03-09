package gopdf

type anchorOption struct {
	page int
	y    float64
}

type linkOption struct {
	x, y, w, h float64
	url        string
	anchor     string
}
