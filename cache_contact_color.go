package gopdf

import (
	"bytes"
	"fmt"
)

const colorTypeStroke = "RG"

const colorTypeFill = "rg"

type cacheContentColor struct {
	colorType string
	r, g, b   uint8
}

func (c *cacheContentColor) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	rFloat := float64(c.r) * 0.00392156862745
	gFloat := float64(c.g) * 0.00392156862745
	bFloat := float64(c.b) * 0.00392156862745
	buff.WriteString(fmt.Sprintf("%.2f %.2f %.2f %s\n", rFloat, gFloat, bFloat, c.colorType))
	return &buff, nil
}
