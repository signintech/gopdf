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
}

func (c *cacheContentRectangle) toStream() (*bytes.Buffer, error) {

	h := c.pageHeight
	x := c.x
	y := c.y
	width := c.width
	height := c.height

	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re s\n", x, h-y, width, height))
	return &buff, nil
}
