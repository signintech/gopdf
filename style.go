package gopdf

type PaintStyle string

const (
	DrawPaintStyle     PaintStyle = "S"
	FillPaintStyle     PaintStyle = "f"
	DrawFillPaintStyle PaintStyle = "B"
)

func parseStyle(style string) PaintStyle {
	op := DrawPaintStyle
	if style == "F" {
		op = FillPaintStyle
	} else if style == "FD" || style == "DF" {
		op = DrawFillPaintStyle
	}

	return op
}
