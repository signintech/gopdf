package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentCurve struct {
	pageHeight float64
	x0         float64
	y0         float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
	x3         float64
	y3         float64
	style      string
}

func (c *cacheContentCurve) toStream(protection *PDFProtection) (*bytes.Buffer, error) {

	h := c.pageHeight
	x0 := c.x0
	y0 := c.y0
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2
	x3 := c.x3
	y3 := c.y3
	style := c.style

	var buff bytes.Buffer
	//cp := 0.55228
	buff.WriteString(fmt.Sprintf("%0.2f %0.2f m\n", x0, h-y0))
	buff.WriteString(fmt.Sprintf(
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c",
		x1, h-y1, x2, h-y2, x3, h-y3,
	))
	op := "S"
	if style == "F" {
		op = "f"
	} else if style == "FD" || style == "DF" {
		op = "B"
	}
	buff.WriteString(fmt.Sprintf(" %s\n", op))
	return &buff, nil
}
