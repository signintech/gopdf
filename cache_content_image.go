package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentImage struct {
	index int
	x     float64
	y     float64
	h     float64
	rect  Rect
}

func (c *cacheContentImage) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", c.rect.W, c.rect.H, c.x, c.h-(c.y+c.rect.H), c.index+1))
	return &buff, nil
}
