package gopdf

import (
	"fmt"
	"io"
)

type cacheContentImage struct {
	verticalFlip   bool
	horizontalFlip bool
	index          int
	x              float64
	y              float64
	h              float64
	rect           Rect
	transparency   *Transparency
	crop           *CropOptions
}

func (c *cacheContentImage) write(w io.Writer, protection *PDFProtection) error {
	x := c.x
	width := c.rect.W
	height := c.rect.H
	y := c.h - (c.y + c.rect.H)

	contentStream := "q\n"

	if c.transparency != nil && c.transparency.Alpha != 1 {
		contentStream += fmt.Sprintf("/GS%d gs\n", c.transparency.indexOfExtGState)
	}

	if c.horizontalFlip || c.verticalFlip {
		fh := "1"
		if c.horizontalFlip {
			fh = "-1"
			x = -1*x - width
		}

		fv := "1"
		if c.verticalFlip {
			fv = "-1"
			y = -1*y - height
		}

		contentStream += fmt.Sprintf("%s 0 0 %s 0 0 cm\n", fh, fv)
	}

	if c.crop != nil {
		/*
			"width" : 173.28,
			"height" : 223.19,
			"y" : 191.038,
			"x" : 272.96,
			"originalWidth" : 297.5,
			"originalHeight" : 223.19,
			"cropOpts": {
	            "left": 124.21012269938649,
	            "top": 0,
	            "width": 173.2898773006135,
	            "height": 223.18829787234046
            },

			cropX = x = 272.96
			cropY = y = 191.038
			startX = x - cropOpts.left = 272.96 - 124.21 = 148.75
			startY = y - cropOpts.top = 191.038 - 0 = 191.038
			flipX = -1 * x - width = -272,96 - 173,28 = -446.24
			flipY = -1 * y - height = -191.038 - 223.19 = -414.29
			flipCropX = -1 * cropX - width = -272.96 - 173.28 = -446.24
			flipCropY = -1 * cropY - height = -191.038 - 223.19 = -414.228
			flipStartX = -1 * startX - originalWidth = -148.75 - 297.5 = -446.25
			flipStartY = -1 * startY - originalHeight = -191.038 - 223.19 = -414.228
		*/
		/*
			Crop
			q
				cropX cropY width height re W* n
				q originalWidth 0 0 originalHeight startX startY cm /I1 Do Q
			Q

			q
				272.96 191.038 173.28 223.19 re W* n
				q 297.5 0 0 223.19 148.75 191.038 cm /I1 Do Q
			Q
		*/
		/*
			Flip
			q
				-1 0 0 1 0 0 cm
				q width 0 0 height flipX flipY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				q 173.28 0 0 223.19 -446.24 414.29 cm /I1 Do Q
			Q
		*/
		/*
			Flip+Crop
			q
				-1 0 0 1 0 0 cm
				flipCropX cropY width height re W* n
				q originalWidth 0 0 originalHeight flipStartX startY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				-446.24 191.038 173.28 223.19 re W* n
				q 297.5 0 0 223.19 -446.25 191.038 cm /I1 Do Q
			Q
		*/

		clippingX := c.x
		if c.horizontalFlip {
			clippingX = -1 * clippingX - c.crop.Width
		}

		clippingY := c.y
		if c.verticalFlip {
			clippingY = -1 * clippingY - c.crop.Height
		}

		startX := c.x - c.crop.X
		if c.horizontalFlip {
			startX = -1 * startX - width
		}

		startY := c.y - c.crop.Y
		if c.verticalFlip {
			startY = -1 * startY - height
		}

		contentStream += fmt.Sprintf("\t\t%0.2f %0.2f %0.2f %0.2f re W* n\n", clippingX, clippingY, c.crop.Width, c.crop.Height)
		contentStream += fmt.Sprintf("\t\tq %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, startX, startY, c.index+1)
	} else {
		contentStream += fmt.Sprintf("\tq %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, x, y, c.index+1)
	}

	contentStream += "Q\n"

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
