package gopdf

import (
	"bytes"
	"fmt"
)

type cacheContentLineType struct {
	lineType string
}

func (c *cacheContentLineType) toStream(protection *PDFProtection) (*bytes.Buffer, error) {
	var buff bytes.Buffer
	switch c.lineType {
	case "dashed":
		buff.WriteString(fmt.Sprint("[5] 2 d\n"))
	case "dotted":
		buff.WriteString(fmt.Sprint("[2 3] 11 d\n"))
	default:
		buff.WriteString(fmt.Sprint("[] 0 d\n"))
	}
	return &buff, nil
}
