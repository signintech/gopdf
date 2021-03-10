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
	y := c.h - (c.y + height)

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
			"pageHeight": 842
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

			y = pageHeight - (y + height) = 842 - (191.038 + 223.19) = 427.772
			cropX = x = 272.96
			cropY = y = 427.772
			startX = x - cropOpts.left = 272.96 - 124.21 = 148.75
			startY = y - cropOpts.top = 427.772 - 0 = 427.772
			flipX = -1 * x - width = -272,96 - 173,28 = -446.24
			flipY = -1 * y - height = -191.038 - 223.19 = -414.29
			flipCropX = -1 * cropX - width = -272.96 - 173.28 = -446.24
			flipCropY = -1 * cropY - height = -191.038 - 223.19 = -414.228
			flipStartX = -1 * startX - originalWidth = -148.75 - 297.5 = -446.25
			flipStartY = -1 * startY - originalHeight = -191.038 - 223.19 = -414.228

			Crop
			q
				cropX cropY width height re W* n
				q originalWidth 0 0 originalHeight startX startY cm /I1 Do Q
			Q

			q
				272.96 191.038 173.28 223.19 re W* n
				q 297.5 0 0 223.19 148.75 191.038 cm /I1 Do Q
			Q

			Flip
			q
				-1 0 0 1 0 0 cm
				q width 0 0 height flipX flipY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				q 173.28 0 0 223.19 -446.24 414.29 cm /I1 Do Q
			Q

			Flip+Crop
			q
				-1 0 0 1 0 0 cm
				flipCropX cropY width height re W* n
				q originalWidth 0 0 originalHeight flipStartX startY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				-446.24 427.772 173.28 223.19 re W* n
				q 297.5 0 0 223.19 -446.25 427.772 cm /I1 Do Q
			Q
		*/
		/*
			"pageHeight": 842
			"width" : 129.01,
			"height" : 93.82,
			"y" : 382.58,
			"x" : 239.16,
			"originalWidth" : 575.02,
			"originalHeight" : 383.35,
			"cropOpts": {
	            "left": 360.01,
	            "top": 101.05,
	            "width": 129.01,
	            "height": 93.82
            },

			y = pageHeight - (y + height) = 842 - (382.58 + 93.82) = 365.6
			cropX = x = 382.58
			cropY = y = 365.6
			startX = x - cropOpts.left = 239.16 - 360.01 = -120.85
			startY = y - cropOpts.top = 365.6 - 101.05 = 264.55
			flipX = -1 * x - width = -272,96 - 173,28 = -446.24
			flipY = -1 * y - height = -191.038 - 223.19 = -414.29
			flipCropX = -1 * cropX - width = -382.58 - 129.01 = -511.59
			flipCropY = -1 * cropY - height = -191.038 - 223.19 = -414.228
			flipStartX = -1 * startX - originalWidth = 120.85 - 575.02 = -454.17
			flipStartY = -1 * startY - originalHeight = -191.038 - 223.19 = -414.228

			Crop
			q
				cropX cropY width height re W* n
				q originalWidth 0 0 originalHeight startX startY cm /I1 Do Q
			Q

			q
				272.96 191.038 173.28 223.19 re W* n
				q 297.5 0 0 223.19 148.75 191.038 cm /I1 Do Q
			Q

			Flip
			q
				-1 0 0 1 0 0 cm
				q width 0 0 height flipX flipY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				q 173.28 0 0 223.19 -446.24 414.29 cm /I1 Do Q
			Q

			Flip+Crop
			q
				-1 0 0 1 0 0 cm
				flipCropX cropY width height re W* n
				q originalWidth 0 0 originalHeight flipStartX startY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				-511.59 365.6 129.01 93.82 re W* n
				q 575.02 0 0 383.35 -454.17 264.55 cm /I1 Do Q
			Q
		*/

		/*
			"pageHeight": 842
			"y" : 0,
			"x" : 0,
			"width" : 138.04,
			"height" : 215.4,
			"originalWidth" : 337.30,
			"originalHeight" : 337.30,
			"cropOpts": {
	            "left": 0,
	            "top": 0,
	            "width": 138.04,
	            "height": 215.4
            },

			y = pageHeight - (y + height) = 842 - 0 - 337.30 = 504.7
			cropX = x = 0
			cropY = pageHeight - (cropOpts.top + cropOpts.height) = 842 - 0 - 215.4 = 626.6
			startX = x - cropOpts.left = 0 - 0 = 0
			startY = y - cropOpts.top = 504.7 - 0 = 0
			flipX = -1 * x - width = -138.04
			flipY = -1 * y - height = -215.4
			flipCropX = -1 * cropX - width = -138.04
			flipCropY = -1 * cropY - height = -626.6 - 215.4 = -842
			flipStartX = -1 * startX - originalWidth = -337.30
			flipStartY = -1 * startY - originalHeight = -337.30

			Crop
			q
				cropX cropY width height re W* n
				q originalWidth 0 0 originalHeight startX startY cm /I1 Do Q
			Q

			Crop Calc
			q
				0 504.7 138.04 215.4 re W* n
				q 337.30 0 0 337.30 0 504.7 cm /I1 Do Q
			Q

			True crop
			q
				0.00 627 138.04 215.40 re W* n
				q 337.30 0 0 337.30 0.00 504.70 cm /I1 Do Q
			Q

			Crop Got
			q
				325.43 344.57 47.75 88.35 re W* n
				q 337.30 0 0 337.30 119.63 181.89 cm /I1 Do Q
			Q

			Flip+Crop
			q
				-1 0 0 1 0 0 cm
				flipCropX cropY width height re W* n
				q originalWidth 0 0 originalHeight flipStartX startY cm /I1 Do Q
			Q

			q
				-1 0 0 1 0 0 cm
				-207.88 593.52 47.75 88.35 re W* n
				q 337.30 0 0 337.30 -291.63 430.84 /I1 Do Q
			Q
		*/

		clippingX := c.x
		if c.horizontalFlip {
			clippingX = -clippingX - c.crop.Width
		}

		clippingY := c.h - (c.crop.Y + c.crop.Height)
		if c.verticalFlip {
			clippingY = -(c.crop.Y + c.crop.Height)
		}

		startX := c.x - c.crop.X
		if c.horizontalFlip {
			startX = -1 * startX - width
		}

		startY := y - c.crop.Y
		//if c.verticalFlip {
		//	startY = -1 * startY - height
		//}

		contentStream += fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re W* n\n", clippingX, clippingY, c.crop.Width, c.crop.Height)
		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, startX, startY, c.index+1)
	} else {
		contentStream += fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", width, height, x, y, c.index+1)
	}

	contentStream += "Q\n"

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
