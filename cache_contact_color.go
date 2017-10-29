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
	rFloat := float64(c.r) * 0.00392156862745
	gFloat := float64(c.g) * 0.00392156862745
	bFloat := float64(c.b) * 0.00392156862745
	fmt.Fprintf(w, "%.2f %.2f %.2f %s\n", rFloat, gFloat, bFloat, c.colorType)
	return nil
}
