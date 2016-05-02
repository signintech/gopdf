package gopdf

//Current current state
type Current struct {
	setXCount int //many times we go func SetX()
	X         float64
	Y         float64

	//font
	IndexOfFontObj int
	CountOfFont    int
	CountOfL       int

	Font_Size      int
	Font_Style     string
	Font_FontCount int
	Font_Type      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

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

	lineWidth float64
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

func (rgb Rgb) equal(obj Rgb) bool {
	if rgb.r == obj.r && rgb.g == obj.g && rgb.b == obj.b {
		return true
	}
	return false
}
