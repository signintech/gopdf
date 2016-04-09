package gopdf

//CURRENT_FONT_TYPE_IFONT this font add by Gopdf.AddFont(...)
const CURRENT_FONT_TYPE_IFONT = 0

//CURRENT_FONT_TYPE_SUBSET this font add by Gopdf.AddTTFFont(...)
const CURRENT_FONT_TYPE_SUBSET = 1

//Current current state
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

	Font_IFont   IFont          // Font_Type == CURRENT_FONT_TYPE_IFONT Deprecated
	Font_ISubset *SubsetFontObj // Font_Type == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches []ImageCache

	//text color
	txtColor Rgb

	//text grayscale
	grayFill float64
	//draw grayscale
	grayStroke float64
}

func (c *Current) setTextColor(rgb Rgb) {
	c.txtColor = rgb
}

func (c *Current) textColor() Rgb {
	return c.txtColor
}

type ImageCache struct {
	Path  string
	Index int
}

type Rgb struct {
	r uint8
	g uint8
	b uint8
}
