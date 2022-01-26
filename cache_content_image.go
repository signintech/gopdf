package gopdf

import (
	"fmt"
	"io"
	"math"
)

type cacheContentImage struct {
	radianAngle      float64
	verticalFlip     bool
	horizontalFlip   bool
	index            int
	x                float64
	y                float64
	pageHeight       float64
	rect             Rect
	crop             *CropOptions
	extGStateIndexes []int
}

func (c *cacheContentImage) write(w io.Writer, protection *PDFProtection) error {
	width := c.rect.W
	height := c.rect.H

	contentStream := "q\n"

	for _, extGStateIndex := range c.extGStateIndexes {
		contentStream += fmt.Sprintf("/GS%d gs\n", extGStateIndex)
	}

	if c.horizontalFlip || c.verticalFlip {
		fh := "1"
		if c.horizontalFlip {
			fh = "-1"
		}

		fv := "1"
		if c.verticalFlip {
			fv = "-1"
		}

		contentStream += fmt.Sprintf("%s 0 0 %s 0 0 cm\n", fh, fv)
	}

	contentStream += "q "

	if c.crop != nil {
		clippingX := c.x
		if c.horizontalFlip {
			clippingX = -clippingX - c.crop.Width
		}

		clippingY := c.pageHeight - (c.y + c.crop.Height)
		if c.verticalFlip {
			clippingY = -clippingY - c.crop.Height
		}

		startX := c.x - c.crop.X
		if c.horizontalFlip {
			startX = -startX - width
		}

		startY := c.pageHeight - (c.y - c.crop.Y + c.rect.H)
		if c.verticalFlip {
			startY = -startY - height
		}

		contentStream += fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re W* n\n", clippingX, clippingY, c.crop.Width, c.crop.Height)
		contentStream += fmt.Sprintf("%0.2f 0 0 %0.2f %0.2f %0.2f cm\n", width, height, startX, startY)
	} else {
		x := c.x
		y := c.pageHeight - (c.y + height)

		if c.horizontalFlip {
			x = -x - width
		}

		if c.verticalFlip {
			y = -y - height
		}

		contentStream += fmt.Sprintf("%0.2f 0 0 %0.2f %0.2f %0.2f cm\n", width, height, x, y)
	}

	if c.radianAngle != 0 {
		cos := math.Cos(c.radianAngle)
		sin := math.Sin(c.radianAngle)

		contentStream += fmt.Sprintf("%.5f %.5f %.5f %.5f 1 0 cm\n", cos, sin, -sin, cos)
	}

	contentStream += fmt.Sprintf("/I%d Do Q\n", c.index+1)

	contentStream += "Q\n"

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
