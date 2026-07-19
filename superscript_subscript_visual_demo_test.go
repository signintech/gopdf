package gopdf

import (
	"os"
	"testing"
)

// TestSuperscriptSubscriptVisualDemo renders mixed normal/script runs so the
// behaviour can be checked visually: script glyphs are smaller, superscripts
// sit above the baseline, subscripts below, and the surrounding normal text
// keeps a single baseline grid.
//
// It writes ./test/out/superscript_subscript_demo.pdf. Run it on its own with:
//
//	go test -run TestSuperscriptSubscriptVisualDemo
func TestSuperscriptSubscriptVisualDemo(t *testing.T) {
	if err := os.MkdirAll("./test/out", 0755); err != nil {
		t.Fatal(err)
	}

	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	const family = "LiberationSerif-Regular"
	if err := pdf.AddTTFFont(family, "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont(family, "", 14); err != nil {
		t.Fatal(err)
	}

	// seg draws one run in the given style at the current position;
	// Cell(nil, ...) advances x, so runs chain into one visual line.
	seg := func(style int, s string) {
		if err := pdf.SetFontWithStyle(family, style, 14); err != nil {
			t.Fatal(err)
		}
		if err := pdf.Cell(nil, s); err != nil {
			t.Fatal(err)
		}
	}

	y := 40.0
	row := func(draw func()) {
		pdf.SetXY(40, y)
		draw()
		y += 30
	}

	row(func() {
		seg(Regular, "E = mc")
		seg(Superscript, "2")
		seg(Regular, ", said Einstein")
	})
	row(func() {
		seg(Regular, "H")
		seg(Subscript, "2")
		seg(Regular, "O, CO")
		seg(Subscript, "2")
		seg(Regular, " and C")
		seg(Subscript, "6")
		seg(Regular, "H")
		seg(Subscript, "12")
		seg(Regular, "O")
		seg(Subscript, "6")
	})
	row(func() {
		seg(Regular, "x")
		seg(Subscript, "1")
		seg(Regular, " + y")
		seg(Superscript, "2")
		seg(Regular, " = z")
		seg(Superscript, "n")
	})
	row(func() {
		seg(Regular, "a footnote marker")
		seg(Superscript, "[3]")
		seg(Regular, " and the baseline stays put")
	})

	out := "./test/out/superscript_subscript_demo.pdf"
	if err := pdf.WritePdf(out); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %s", out)
}
