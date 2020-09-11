package gopdf

//PageOption option of page
type PageOption struct {
	TrimBox *Box
	PageSize *Rect
}

func (p PageOption) isEmpty() bool {
	if p.PageSize == nil {
		return true
	}
	return false
}

func (p PageOption) isTrimBoxSet() bool {
	if p.TrimBox == nil {
		return false
	}
	return true
}
