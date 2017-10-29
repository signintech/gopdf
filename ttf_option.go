package gopdf

//TtfOption  font option
type TtfOption struct {
	UseKerning bool
	Style      int // Regular|Bold|Italic
}

func defaultTtfFontOption() TtfOption {
	var defa TtfOption
	defa.UseKerning = false
	defa.Style = Regular
	return defa
}
