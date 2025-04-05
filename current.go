package gopdf

// Current current state
type Current struct {
	setXCount int //many times we go func SetX()
	X         float64
	Y         float64

	//font
	IndexOfFontObj int
	CountOfFont    int
	CountOfL       int

	FontSize      float64
	FontStyle     int // Regular|Bold|Italic|Underline
	FontFontCount int
	FontType      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

	IndexOfColorSpaceObj int
	CountOfColorSpace    int

	CharSpacing float64

	FontISubset *SubsetFontObj // FontType == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches map[int]ImageCache

	//text color mode
	txtColorMode string //color, gray

	//text color
	txtColor ICacheColorText

	//text grayscale
	grayFill float64
	//draw grayscale
	grayStroke float64

	lineWidth float64

	//current page size
	pageSize *Rect

	//current trim box
	trimBox *Box

	sMasksMap       SMaskMap
	extGStatesMap   ExtGStatesMap
	transparency    *Transparency
	transparencyMap TransparencyMap
}

func (c *Current) setTextColor(color ICacheColorText) {
	c.txtColor = color
}

func (c *Current) textColor() ICacheColorText {
	return c.txtColor
}

// ImageCache is metadata for caching images.
type ImageCache struct {
	Path  string //ID or Path
	Index int
	Rect  *Rect
}
