package gopdf

import (
	"bytes"
	"fmt"
)

const drawTypeFill = "g"
const drawTypeStroke = "G"

type cacheContentGray struct {
	drawType string
	scale    float64
}

func (c *cacheContentGray) toStream() (*bytes.Buffer, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%.2f %s\n", c.scale, c.drawType))
	return &buff, nil
}
