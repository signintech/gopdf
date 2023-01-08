package gopdf

// TtfOption  font option
type TtfOption struct {
	UseKerning                bool
	Style                     int               //Regular|Bold|Italic
	OnGlyphNotFound           func(r rune)      //Called when a glyph cannot be found, just for debugging
	OnGlyphNotFoundSubstitute func(r rune) rune //Called when a glyph cannot be found, we can return a new rune to replace it.
}

func defaultTtfFontOption() TtfOption {
	var defa TtfOption
	defa.UseKerning = false
	defa.Style = Regular
	defa.OnGlyphNotFoundSubstitute = DefaultOnGlyphNotFoundSubstitute
	return defa
}

func DefaultOnGlyphNotFoundSubstitute(r rune) rune {
	return rune('\u0020')
}
