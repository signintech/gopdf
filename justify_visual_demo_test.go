package gopdf

import (
	"os"
	"testing"
)

// TestJustifyVisualDemo renders the alignment options into a single PDF so the
// justify behaviour can be checked visually. Each sample sits in a bordered box:
// a justified line is flush against both inner edges, a justified paragraph has
// every line flush except the last (which stays left-aligned), and the graceful
// fallbacks (a line with no spaces, or one that overflows) stay left-aligned.
//
// It writes ./test/out/justify_demo.pdf. Run it on its own with:
//
//	go test -run TestJustifyVisualDemo
func TestJustifyVisualDemo(t *testing.T) {
	if err := os.MkdirAll("./test/out", 0755); err != nil {
		t.Fatal(err)
	}

	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		t.Fatal(err)
	}

	const (
		labelX = 40.0
		boxX   = 200.0
		boxW   = 340.0
		lineH  = 20.0
	)
	// Wrap paragraphs on spaces (word boundaries) instead of the mid-word default.
	wordBreak := &BreakOption{Mode: BreakModeIndicatorSensitive, BreakIndicator: ' '}

	label := func(y float64, s string) {
		pdf.SetXY(labelX, y)
		pdf.Cell(nil, s)
	}
	// line draws a single line of text in a bordered box.
	line := func(y float64, text string, align int) {
		pdf.SetXY(boxX, y)
		if err := pdf.CellWithOption(&Rect{W: boxW, H: lineH}, text,
			CellOption{Align: align | Top, Border: AllBorders}); err != nil {
			t.Fatal(err)
		}
	}
	// paragraph wraps text over multiple lines and draws ONE border around the
	// whole block. MultiCell would redraw a full box for every line (stacking
	// overlapping borders), so we render without a border and box it ourselves.
	paragraph := func(y float64, text string, align int) {
		pdf.SetXY(boxX, y)
		if err := pdf.MultiCellWithOption(&Rect{W: boxW, H: lineH}, text,
			CellOption{Align: align, BreakOption: wordBreak}); err != nil {
			t.Fatal(err)
		}
		pdf.RectFromUpperLeft(boxX, y, boxW, pdf.GetY()-y)
	}
	row := func(y *float64, name, text string, align int) {
		label(*y, name)
		line(*y, text, align)
		*y += 34
	}

	const sample = "the quick brown fox jumps over"
	const para = "the quick brown fox jumps over the lazy dog while the sun sets " +
		"slowly over the quiet western hills as a gentle breeze drifts across the " +
		"meadow carrying the faint scent of pine and distant rain this evening"

	y := 40.0
	row(&y, "Left", sample, Left)
	row(&y, "Center", sample, Center)
	row(&y, "Right", sample, Right)
	row(&y, "Justify (single line)", sample, Justify)
	y += 10

	label(y, "Justify (paragraph)")
	paragraph(y, para, Justify)
	y = pdf.GetY() + 26

	row(&y, "Fallback: no spaces", "supercalifragilistic", Justify)
	label(y, "Fallback: overflow")
	line(y, "this single line of text is far too long to ever fit inside the width of this box", Justify)

	out := "./test/out/justify_demo.pdf"
	if err := pdf.WritePdf(out); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %s", out)
}
