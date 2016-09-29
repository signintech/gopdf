package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentLine struct {
	pageHeight float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
}

func (c *cacheContentLine) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	h := c.pageHeight
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
	return &buff, nil
}
