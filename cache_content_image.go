package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImage struct {
	VerticalFlip   bool
	HorizontalFlip bool
	index          int
	x              float64
	y              float64
	h              float64
	rect           Rect
	transparency   Transparency
}

func (c *cacheContentImage) write(w io.Writer, protection *PDFProtection) error {
	x := c.x
	width := c.rect.W
	height := c.rect.H
	y := c.h-(c.y+c.rect.H)

	contentStream := "q\n"

	if c.transparency.IndexOfExtGState != 0 {
		contentStream += fmt.Sprintf("/GS%d gs\n", c.transparency.IndexOfExtGState)
	}

	if c.HorizontalFlip || c.VerticalFlip {
		fh := "1"
		if c.HorizontalFlip {
			fh = "-1"
			x = -1 * x - width
		}

		fv := "1"
		if c.VerticalFlip {
			fv = "-1"
			y = -1 * y - height
		}

		contentStream += fmt.Sprintf("%s 0 0 %s 0 0 cm \n", fh, fv)
	}

	contentStream += fmt.Sprintf("%0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, x, y, c.index+1)

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
