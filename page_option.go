package gopdf

// PageOption option of page
type PageOption struct {
	TrimBox  *Box
	PageSize *Rect
}

func (p PageOption) isEmpty() bool {
	return p.PageSize == nil
}

func (p PageOption) isTrimBoxSet() bool {
	if p.TrimBox == nil {
		return false
	}
	if p.TrimBox.Top == 0 && p.TrimBox.Left == 0 && p.TrimBox.Bottom == 0 && p.TrimBox.Right == 0 {
		return false
	}

	return true
}
