package gopdf

import (
	"fmt"
	"io"
)

const colorTypeStroke = "RG"

const colorTypeFill = "rg"

type cacheContentColor struct {
	colorType string
	r, g, b   uint8
}

func (c *cacheContentColor) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.3f %.3f %.3f %s\n", float64(c.r)/255, float64(c.g)/255, float64(c.b)/255, c.colorType)
	return nil
}
