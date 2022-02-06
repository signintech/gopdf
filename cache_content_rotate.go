package gopdf

import (
	"fmt"
	"io"
	"math"
)

type cacheContentRotate struct {
	isReset     bool
	pageHeight  float64
	angle, x, y float64
}

func (cc *cacheContentRotate) write(w io.Writer, protection *PDFProtection) error {
	if cc.isReset == true {
		if _, err := io.WriteString(w, "Q\n"); err != nil {
			return err
		}

		return nil
	}

	matrix := computeRotateTransformationMatrix(cc.x, cc.y, cc.angle, cc.pageHeight)
	contentStream := fmt.Sprintf("q\n %s", matrix)

	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}

func computeRotateTransformationMatrix(x, y, degreeAngle, pageHeight float64) string {
	radianAngle := degreeAngle * (math.Pi / 180)

	c := math.Cos(radianAngle)
	s := math.Sin(radianAngle)
	cy := pageHeight - y

	return fmt.Sprintf("%.5f %.5f %.5f\n %.5f %.2f %.2f cm\n 1 0 0\n 1 %.2f %.2f cm\n", c, s, -s, c, x, cy, -x, -cy)
}
