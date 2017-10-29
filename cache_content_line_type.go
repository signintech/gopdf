package gopdf

import (
	"fmt"
	"io"
)

type cacheContentLineType struct {
	lineType string
}

func (c *cacheContentLineType) write(w io.Writer, protection *PDFProtection) error {
	switch c.lineType {
	case "dashed":
		fmt.Fprint(w, "[5] 2 d\n")
	case "dotted":
		fmt.Fprint(w, "[2 3] 11 d\n")
	default:
		fmt.Fprint(w, "[] 0 d\n")
	}
	return nil
}
