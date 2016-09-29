package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentRectangle struct {
	pageHeight float64
	x          float64
	y          float64
	width      float64
	height     float64
	style      string
}

func (c *cacheContentRectangle) toStream(protection *PDFProtection) (*bytes.Buffer, error) {

	h := c.pageHeight
	x := c.x
	y := c.y
	width := c.width
	height := c.height

	var buff bytes.Buffer
	op := parseStyle(c.style)
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re %s\n", x, h-y, width, height, op))
	//buff.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re b\n", x, h-y, width, height))
	return &buff, nil
}
