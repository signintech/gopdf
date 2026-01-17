package gopdf

import (
	"fmt"
	"io"
)

type cacheContentClipPolygon struct {
	pageHeight float64
	points     []Point
}

func (c *cacheContentClipPolygon) write(w io.Writer, protection *PDFProtection) error {
	for i, p := range c.points {
		fmt.Fprintf(w, "%.2f %.2f", p.X, c.pageHeight-p.Y)
		if i == 0 { // first point
			fmt.Fprint(w, " m ") // moveto: start new path
		} else {
			fmt.Fprint(w, " l ") // lineto: draw line from current point
		}
	}
	fmt.Fprint(w, "h W n\n") // h=close path, W=clip, n=end path
	return nil
}
