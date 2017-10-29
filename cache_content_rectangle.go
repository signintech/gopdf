package gopdf

import (
	"fmt"
	"io"
)

type cacheContentRectangle struct {
	pageHeight float64
	x          float64
	y          float64
	width      float64
	height     float64
	style      string
}

func (c *cacheContentRectangle) write(w io.Writer, protection *PDFProtection) error {

	h := c.pageHeight
	x := c.x
	y := c.y
	width := c.width
	height := c.height

	op := parseStyle(c.style)
	fmt.Fprintf(w, "%0.2f %0.2f %0.2f %0.2f re %s\n", x, h-y, width, height, op)
	return nil
}
