package gopdf

import (
	"fmt"
	"io"
)

const colorTypeStrokeCMYK = "K"

const colorTypeFillCMYK = "k"

type cacheContentColorCMYK struct {
	colorType  string
	c, m, y, k uint8
}

func (c *cacheContentColorCMYK) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.2f %.2f %.2f %.2f %s\n", float64(c.c)/100, float64(c.m)/100, float64(c.y)/100, float64(c.k)/100, c.colorType)
	return nil
}
