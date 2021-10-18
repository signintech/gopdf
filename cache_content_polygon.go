package gopdf

import (
	"fmt"
	"io"
)

type cacheContentPolygon struct {
	pageHeight float64
	style      string
	points     []Point
	opts       polygonOptions
}

func (c *cacheContentPolygon) write(w io.Writer, protection *PDFProtection) error {

	fmt.Fprintf(w, "q\n")
	for _, extGStateIndex := range c.opts.extGStateIndexes {
		fmt.Fprintf(w, "/GS%d gs\n", extGStateIndex)
	}

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

	fmt.Fprintf(w, "Q\n")
	return nil
}
