package gopdf

import (
	"bytes"
	"fmt"
)

const grayTypeFill = "g"
const grayTypeStroke = "G"

type cacheContentGray struct {
	grayType string
	scale    float64
}

func (c *cacheContentGray) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%.2f %s\n", c.scale, c.grayType))
	return &buff, nil
}
