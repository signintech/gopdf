package gopdf

import (
	"fmt"
	"io"
)

type ColorSpaceObj struct {
	CountOfSpaceColor int
	Name              string
	Contour           bool
	r                 float64
	g                 float64
	b                 float64
}

func (c *ColorSpaceObj) init(func() *GoPdf) {}

func (c *ColorSpaceObj) getType() string {
	return "ColorSpace"
}

func (c *ColorSpaceObj) write(w io.Writer, objID int) error {

	if c.Contour {
		io.WriteString(w, "[ /Separation /CutContour /DeviceRGB\n")
	} else {
		io.WriteString(w, "[ /Separation /DeviceRGB\n")
	}

	io.WriteString(w, "		<< \n")
	io.WriteString(w, "			/FunctionType 2\n")
	io.WriteString(w, "			/Domain [0 1]\n")
	io.WriteString(w, "			/C0 [0.0 0.0 0.0]\n")
	fmt.Fprintf(w, "			/C1 [%.3f %.3f %.3f]\n", c.r, c.g, c.b)
	io.WriteString(w, "			/N 1\n")
	io.WriteString(w, "		>>\n")
	io.WriteString(w, "]\n")

	return nil
}

func (c *ColorSpaceObj) SetColor(r, g, b uint8) {
	c.r = float64(r) / 255.0
	c.g = float64(g) / 255.0
	c.b = float64(b) / 255.0
}
