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
	angle := (cc.angle * 22.0) / (180.0 * 7.0)
	c := math.Cos(angle)
	s := math.Sin(angle)
	cy := cc.pageHeight - cc.y

	contentStream := fmt.Sprintf("q %.5f %.5f %.5f %.5f %.2f %.2f cm 1 0 0 1 %.2f %.2f cm\n", c, s, -s, c, cc.x, cy, -cc.x, -cy)
	if _, err := io.WriteString(w, contentStream); err != nil {
		return err
	}

	return nil
}
