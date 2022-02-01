package gopdf

import (
	"fmt"
	"io"
	"log"
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

func (c *cacheContentImage) write(writer io.Writer, protection *PDFProtection) error {
	width := c.rect.W
	height := c.rect.H

	// проблема в том когда пишется стрим маски также пишется и поворот второй раз
	var angle float64
	if c.withMask {
		angle = 0
	} else {
		angle = c.imageAngle
	}

	if angle != 0 {
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
			angle:      angle,
		}
		if err := cacheRotate.write(writer, protection); err != nil {
			return err
		}

		defer func() {
			resetCacheRotate := cacheContentRotate{isReset: true}

			if err := resetCacheRotate.write(writer, protection); err != nil {
				log.Println(err)
			}
		}()
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

		var rotateMat string
		if c.imageAngle != 0 && c.maskAngle != 0 {
			x := c.x + width/2
			y := c.y + height/2

			rotateMat = computeRotateTransformationMatrix(x, y, c.imageAngle, c.pageHeight)
		}

		contentStream += fmt.Sprintf("q\n %s %0.2f 0 0\n %0.2f %0.2f %0.2f cm /I%d Do \nQ\n", rotateMat, width, height, startX, startY, c.index+1)
	} else {
		x := c.x
		y := c.pageHeight - (c.y + height)

		if c.horizontalFlip {
			x = -x - width
		}

		if c.verticalFlip {
			y = -y - height
		}

		var rotateMat string
		if c.imageAngle != 0 && c.maskAngle != 0 {
			rotatedX := c.x + width/2
			rotatedY := c.y + height/2

			rotateMat = computeRotateTransformationMatrix(rotatedX, rotatedY, c.maskAngle+c.imageAngle, c.pageHeight)
		}

		contentStream += fmt.Sprintf("q\n %s %0.2f 0 0\n %0.2f %0.2f %0.2f cm\n /I%d Do \nQ\n", rotateMat, width, height, x, y, c.index+1)
	}

	contentStream += "Q\n"

	if _, err := io.WriteString(writer, contentStream); err != nil {
		return err
	}

	return nil
}
