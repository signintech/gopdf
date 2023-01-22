package gopdf

import (
	"fmt"
	"io"
)

type cacheContentCustomLineType struct {
	dashArray []float64
	dashPhase float64
}

func (c *cacheContentCustomLineType) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%0.2f %0.2f d\n", c.dashArray, c.dashPhase)
	return nil
}
