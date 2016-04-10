package gopdf

//TtfOption  font option
type TtfOption struct {
	UseKerning bool
}

func defaultTtfFontOption() TtfOption {
	var defa TtfOption
	defa.UseKerning = false
	return defa
}
