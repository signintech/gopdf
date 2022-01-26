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
	rotateMat := computeRotateMat(c.radianAngle)

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
		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm %s /I%d Do Q\n", width, height, startX, startY, rotateMat, c.index+1)
	} else {
		x := c.x
		y := c.pageHeight - (c.y + height)

		if c.horizontalFlip {
			x = -x - width
		}

		if c.verticalFlip {
			y = -y - height
		}

		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm %s /I%d Do Q\n", width, height, x, y, rotateMat, c.index+1)
	}

	contentStream += "Q\n"

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}

func computeRotateMat(radianAngle float64) string {
	if radianAngle == 0 {
		return ""
	}

	cos := math.Cos(radianAngle)
	sin := math.Sin(radianAngle)

	degreeAngle := int(math.Round(radianAngle / math.Pi * 180))
	if math.Abs(float64(degreeAngle)) > 360 {
		degreeAngle = degreeAngle % 360
	}

	translateX := 0
	translateY := 0

	if degreeAngle == 180 || degreeAngle == 360 || degreeAngle == -180 {
		translateX, translateY = 1, 1
	} else if degreeAngle == 90 || degreeAngle == -270 {
		translateX, translateY = 1, 0
	} else if degreeAngle == 0 {
		translateX, translateY = 0, 0
	} else {
		translateX, translateY = 0, 1
	}

	return fmt.Sprintf("%.5f %.5f %.5f %.5f %d %d cm\n", cos, sin, -sin, cos, translateX, translateY)
}
