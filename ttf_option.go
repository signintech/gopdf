package gopdf

//TtfOption  font option
type TtfOption struct {
	UseKerning        bool
	Style             int  // Regular|Bold|Italic
	UseOpenTypeLayout bool //https://mpdf.github.io/fonts-languages/opentype-layout-otl.html
}

func defaultTtfFontOption() TtfOption {
	var defa TtfOption
	defa.UseKerning = false
	defa.Style = Regular
	return defa
}
