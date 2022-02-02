package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImage struct {
	withMask         bool
	maskAngle        float64
	imageAngle       float64
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

func (c *cacheContentImage) openImageRotateTrMt(writer io.Writer, protection *PDFProtection) error {
	w := c.rect.W
	h := c.rect.H

	if c.crop != nil {
		w = c.crop.Width
		h = c.crop.Height
	}

	x := c.x + w/2
	y := c.y + h/2

	cacheRotate := cacheContentRotate{
		x:          x,
		y:          y,
		pageHeight: c.pageHeight,
		angle:      c.imageAngle,
	}
	if err := cacheRotate.write(writer, protection); err != nil {
		return err
	}

	return nil
}

func (c *cacheContentImage) closeImageRotateTrMt(writer io.Writer, protection *PDFProtection) error {
	resetCacheRotate := cacheContentRotate{isReset: true}

	return resetCacheRotate.write(writer, protection)
}

func (c *cacheContentImage) computeMaskImageRotateTrMt() string {
	angle := c.maskAngle + c.imageAngle
	if angle == 0 {
		return ""
	}

	x := c.x + c.rect.W/2
	y := c.y + c.rect.H/2

	rotateMat := computeRotateTransformationMatrix(x, y, angle, c.pageHeight)

	return rotateMat
}

func (c *cacheContentImage) write(writer io.Writer, protection *PDFProtection) error {
	width := c.rect.W
	height := c.rect.H

	if !c.withMask {
		if err := c.openImageRotateTrMt(writer, protection); err != nil {
			return err
		}

		defer c.closeImageRotateTrMt(writer, protection)
	}

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

	x := c.x
	y := c.pageHeight - c.y

	if c.crop != nil {
		clippingX := x
		if c.horizontalFlip {
			clippingX = -clippingX - c.crop.Width
		}

		clippingY := y - c.crop.Height
		if c.verticalFlip {
			clippingY = -clippingY - c.crop.Height
		}

		contentStream += fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re W* n\n", clippingX, clippingY, c.crop.Width, c.crop.Height)

		x -= c.crop.X
		if c.horizontalFlip {
			x = -x - width
		}

		y += c.crop.Y - height
		if c.verticalFlip {
			y = -y - height
		}
	} else {
		y -= height

		if c.horizontalFlip {
			x = -x - width
		}

		if c.verticalFlip {
			y = -y - height
		}
	}

	var maskImageRotateMat string
	if c.withMask {
		maskImageRotateMat = c.computeMaskImageRotateTrMt()
	}

	contentStream += fmt.Sprintf("q\n %s %0.2f 0 0\n %0.2f %0.2f %0.2f cm /I%d Do \nQ\n", maskImageRotateMat, width, height, x, y, c.index+1)

	contentStream += "Q\n"

	if _, err := io.WriteString(writer, contentStream); err != nil {
		return err
	}

	return nil
}
