package gopdf

//PageOption option of page
type PageOption struct {
	PageSize Rect
}

func (p PageOption) isEmpty() bool {
	if p.PageSize.H == 0 && p.PageSize.W == 0 {
		return true
	}
	return false
}
