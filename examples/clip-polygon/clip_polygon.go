package clip_polygon

import (
	"github.com/signintech/gopdf"
)

func ClipPolygonExample() error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	// 1. Blue rectangle (unclipped, for reference)
	pdf.SetFillColor(200, 200, 255)
	pdf.RectFromUpperLeftWithStyle(50, 50, 100, 100, "F")

	// 2. Red rectangle clipped to triangle
	pdf.SaveGraphicsState() // Save state before clipping
	pdf.ClipPolygon([]gopdf.Point{
		{X: 250, Y: 50},  // top
		{X: 350, Y: 200}, // bottom-right
		{X: 150, Y: 200}, // bottom-left
	})
	pdf.SetFillColor(255, 0, 0)
	pdf.RectFromUpperLeftWithStyle(150, 50, 200, 200, "F")
	pdf.RestoreGraphicsState() // Restore state - clipping no longer active

	// 3. Green rectangle (NOT clipped - proves RestoreGraphicsState works)
	pdf.SetFillColor(0, 200, 0)
	pdf.RectFromUpperLeftWithStyle(400, 50, 100, 100, "F")

	return pdf.WritePdf("clip-polygon.pdf")
}
