package gopdf

//PageOption option of page
type PageOption struct {
	PageSize *Rect
}

func (p PageOption) isEmpty() bool {
	if p.PageSize == nil {
		return true
	}
	return false
}
