package gopdf

import (
	"fmt"
	"io"
)

type cacheContentPolygon struct {
	pageHeight float64
	style      string
	points     []Point
}

func (c *cacheContentPolygon) write(w io.Writer, protection *PDFProtection) error {

	for i, point := range c.points {
		fmt.Fprintf(w, "%.2f %.2f", point.X, c.pageHeight-point.Y)
		if i == 0 {
			fmt.Fprintf(w, " m ")
		} else {
			fmt.Fprintf(w, " l ")
		}

	}

	if c.style == "F" {
		fmt.Fprintf(w, " f\n")
	} else if c.style == "FD" || c.style == "DF" {
		fmt.Fprintf(w, " b\n")
	} else {
		fmt.Fprintf(w, " s\n")
	}

	return nil
}
