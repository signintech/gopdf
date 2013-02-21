package gopdf

import (
	"gopdf/fonts"
)

type Current struct{
	X float64
	Y float64
	
	//font
	IndexOfFontObj int
	CountOfFont int
	
	Font_Size int
	Font_Style string
	Font_IFont fonts.IFont
	Font_FontCount int
	
	//page
	IndexOfPageObj int
}