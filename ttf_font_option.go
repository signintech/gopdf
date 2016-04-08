package gopdf

//TtfFontOption  font option
type TtfFontOption struct {
	UseKerning bool
}

func defaultTtfFontOption() TtfFontOption {
	var defa TtfFontOption
	defa.UseKerning = false
	return defa
}
