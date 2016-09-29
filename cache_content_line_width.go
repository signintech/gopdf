package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentLineWidth struct {
	width float64
}

func (c *cacheContentLineWidth) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("%.2f w\n", c.width))
	return &buff, nil
}
