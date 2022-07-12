package gopdf

import (
	"fmt"
	"io"
)

type cacheContentTextColorRGB struct {
	r, g, b uint8
}

func (c cacheContentTextColorRGB) write(w io.Writer, protection *PDFProtection) error {
	fmt.Fprintf(w, "%.3f %.3f %.3f %s\n", float64(c.r)/255, float64(c.g)/255, float64(c.b)/255, colorTypeFillRGB)
	return nil
}

func (c cacheContentTextColorRGB) equal(obj ICacheColorText) bool {
	rgb, ok := obj.(cacheContentTextColorRGB)
	if !ok {
		return false
	}

	return c.r == rgb.r && c.g == rgb.g && c.b == rgb.b
}
