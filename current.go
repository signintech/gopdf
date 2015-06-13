package gopdf

const CURRENT_FONT_TYPE_IFONT = 0
const CURRENT_FONT_TYPE_SUBSET = 1

type Current struct {
	X float64
	Y float64

	//font
	IndexOfFontObj int
	CountOfFont    int
	CountOfL       int

	Font_Size      int
	Font_Style     string
	Font_FontCount int
	Font_Type      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

	Font_IFont   IFont   // depend on Font_Type
	Font_ISubset ISubset // depend on Font_Type

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches []ImageCache
}

type ImageCache struct {
	Path  string
	Index int
}
