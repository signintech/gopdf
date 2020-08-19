package gopdf

//PageOption option of page
type PageOption struct {
	PageSize *Rect
	TrimSize *Rect
}

func (p PageOption) isEmpty() bool {
	if p.PageSize == nil {
		return true
	}
	return false
}

func (p PageOption) doesTrimSizeSet() bool {
	if p.TrimSize == nil {
		return false
	}
	return true
}
