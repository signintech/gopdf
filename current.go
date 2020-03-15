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

	FontSize      int
	FontStyle     int // Regular|Bold|Italic|Underline
	FontFontCount int
	FontType      int // CURRENT_FONT_TYPE_IFONT or  CURRENT_FONT_TYPE_SUBSET

	FontISubset *SubsetFontObj // FontType == CURRENT_FONT_TYPE_SUBSET

	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//cache of image in pdf file
	ImgCaches []ImageCache

	//text color mode
	txtColorMode string //color, gray

	//text color
	txtColor Rgb

	//text grayscale
	grayFill float64
	//draw grayscale
	grayStroke float64

	lineWidth float64

	//current page size
	pageSize *Rect

	transparency    Transparency
	transparencyMap map[string]Transparency
}

func (c *Current) setTextColor(rgb Rgb) {
	c.txtColor = rgb
}

func (c *Current) textColor() Rgb {
	return c.txtColor
}

// ImageCache is metadata for caching images.
type ImageCache struct {
	Path  string //ID or Path
	Index int
	Rect  *Rect
}

//Rgb  rgb color
type Rgb struct {
	r uint8
	g uint8
	b uint8
}

//SetR set red
func (rgb *Rgb) SetR(r uint8) {
	rgb.r = r
}

//SetG set green
func (rgb *Rgb) SetG(g uint8) {
	rgb.g = g
}

//SetB set blue
func (rgb *Rgb) SetB(b uint8) {
	rgb.b = b
}

func (rgb Rgb) equal(obj Rgb) bool {
	if rgb.r == obj.r && rgb.g == obj.g && rgb.b == obj.b {
		return true
	}
	return false
}

// Transparency defines an object alpha.
type Transparency struct {
	IndexOfExtGState int
}
