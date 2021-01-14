package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImage struct {
	flipVertical   bool
	flipHorizontal bool
	index          int
	x              float64
	y              float64
	h              float64
	rect           Rect
	transparency   Transparency
}

func (c *cacheContentImage) write(w io.Writer, protection *PDFProtection) error {
	contentStream := "q\n"

	if c.transparency.IndexOfExtGState != 0 {
		contentStream += fmt.Sprintf("/GS%d gs\n", c.transparency.IndexOfExtGState)
	}

	if c.flipHorizontal || c.flipVertical {
		fh := "1"
		if c.flipHorizontal {
			fh += "-1"
		}

		fv := "1"
		if c.flipVertical {
			fv = "-1"
		}

		contentStream += fmt.Sprintf("%s 0 0 %s 0 0 cm \n", fh, fv)
	}

	contentStream += fmt.Sprintf("%0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", c.rect.W, c.rect.H, c.x, c.h-(c.y+c.rect.H), c.index+1)

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
