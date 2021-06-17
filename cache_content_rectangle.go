package gopdf

import (
	"fmt"
	"io"
)

type cacheContentRectangle struct {
	pageHeight       float64
	x                float64
	y                float64
	width            float64
	height           float64
	style            PaintStyle
	extGStateIndexes []int
}

func NewCacheContentRectangle(pageHeight float64, rectOpts DrawableRectOptions) ICacheContent {
	if rectOpts.PaintStyle == "" {
		rectOpts.PaintStyle = DrawPaintStyle
	}

	return cacheContentRectangle{
		x:                rectOpts.X,
		y:                rectOpts.Y,
		width:            rectOpts.W,
		height:           rectOpts.H,
		pageHeight:       pageHeight,
		style:            rectOpts.PaintStyle,
		extGStateIndexes: rectOpts.extGStateIndexes,
	}
}

func (c cacheContentRectangle) write(w io.Writer, protection *PDFProtection) error {
	stream := "q\n"

	for _, extGStateIndex := range c.extGStateIndexes {
		stream += fmt.Sprintf("/GS%d gs\n", extGStateIndex)
	}

	stream += fmt.Sprintf("%0.2f %0.2f %0.2f %0.2f re %s\n", c.x, c.pageHeight-c.y, c.width, c.height, c.style)

	stream += "Q\n"

	if _, err := io.WriteString(w, stream); err != nil {
		return err
	}

	return nil
}
