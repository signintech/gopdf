package gopdf

import (
	"fmt"
	"io"
)

type cacheContentOval struct {
	pageHeight float64
	x1         float64
	y1         float64
	x2         float64
	y2         float64
}

func (c *cacheContentOval) write(w io.Writer, protection *PDFProtection) error {

	h := c.pageHeight
	x1 := c.x1
	y1 := c.y1
	x2 := c.x2
	y2 := c.y2

	cp := 0.55228                              // Magnification of the control point
	v1 := [2]float64{x1 + (x2-x1)/2, h - y2}   // Vertex of the lower
	v2 := [2]float64{x2, h - (y1 + (y2-y1)/2)} // .. Right
	v3 := [2]float64{x1 + (x2-x1)/2, h - y1}   // .. Upper
	v4 := [2]float64{x1, h - (y1 + (y2-y1)/2)} // .. Left

	fmt.Fprintf(w, "%0.2f %0.2f m\n", v1[0], v1[1])
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v1[0]+(x2-x1)/2*cp, v1[1], v2[0], v2[1]-(y2-y1)/2*cp, v2[0], v2[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v2[0], v2[1]+(y2-y1)/2*cp, v3[0]+(x2-x1)/2*cp, v3[1], v3[0], v3[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c\n",
		v3[0]-(x2-x1)/2*cp, v3[1], v4[0], v4[1]+(y2-y1)/2*cp, v4[0], v4[1],
	)
	fmt.Fprintf(w,
		"%0.2f %0.2f %0.2f %0.2f %0.2f %0.2f c S\n",
		v4[0], v4[1]-(y2-y1)/2*cp, v1[0]-(x2-x1)/2*cp, v1[1], v1[0], v1[1],
	)

	return nil
}
