package gopdf

import (
	"bytes"
	"strings"
	"testing"
)

// Config.ConversionForUnit is documented to replace the built-in factor for
// Config.Unit. Start already honours it for Config.PageSize and Config.TrimBox,
// so the per-call sizes must honour it too.
func TestConversionForUnitAppliesToPageOptionAndImage(t *testing.T) {
	newDoc := func(conversion float64) *GoPdf {
		pdf := new(GoPdf)
		pdf.Start(Config{
			Unit:              UnitPT,
			ConversionForUnit: conversion,
			PageSize:          Rect{W: 100, H: 200},
		})
		pdf.SetCompressLevel(0)
		return pdf
	}

	t.Run("AddPageWithOption page size", func(t *testing.T) {
		for _, conversion := range []float64{0, 2} {
			conversion := conversion
			pdf := newDoc(conversion)
			pdf.AddPageWithOption(PageOption{PageSize: &Rect{W: 100, H: 200}})

			var buf bytes.Buffer
			if _, err := pdf.WriteTo(&buf); err != nil {
				t.Fatalf("conversion %v: %v", conversion, err)
			}
			// The last MediaBox is the page added through PageOption.
			boxes := mediaBoxes(buf.String())
			if len(boxes) == 0 {
				t.Fatalf("conversion %v: no MediaBox written", conversion)
			}
			want := "0 0 100.00 200.00"
			if conversion != 0 {
				want = "0 0 200.00 400.00"
			}
			if got := boxes[len(boxes)-1]; got != want {
				t.Errorf("conversion %v: MediaBox = %q, want %q", conversion, got, want)
			}
		}
	})

	t.Run("Image rect", func(t *testing.T) {
		for _, conversion := range []float64{0, 2} {
			conversion := conversion
			pdf := newDoc(conversion)
			pdf.AddPage()
			if err := pdf.Image("test/res/gopher01.jpg", 10, 10, &Rect{W: 50, H: 50}); err != nil {
				t.Fatalf("conversion %v: %v", conversion, err)
			}

			var buf bytes.Buffer
			if _, err := pdf.WriteTo(&buf); err != nil {
				t.Fatalf("conversion %v: %v", conversion, err)
			}
			// The image matrix carries the rect size and the placement together.
			want := " 50.00 10.00 140.00 cm /I"
			if conversion != 0 {
				want = " 100.00 20.00 280.00 cm /I"
			}
			if !strings.Contains(buf.String(), want) {
				t.Errorf("conversion %v: image matrix %q not found", conversion, want)
			}
		}
	})
}

func mediaBoxes(s string) []string {
	var out []string
	for _, part := range strings.Split(s, "/MediaBox [ ")[1:] {
		if end := strings.Index(part, " ]"); end >= 0 {
			out = append(out, part[:end])
		}
	}
	return out
}

// The same option must reach the other rect-taking entry points.
func TestConversionForUnitAppliesToImageHolders(t *testing.T) {
	for _, tc := range []struct {
		name string
		draw func(*GoPdf)
	}{
		{"ImageByHolder", func(p *GoPdf) {
			h, err := ImageHolderByPath("test/res/gopher01.jpg")
			if err != nil {
				t.Fatal(err)
			}
			if err := p.ImageByHolder(h, 10, 10, &Rect{W: 50, H: 50}); err != nil {
				t.Fatal(err)
			}
		}},
		{"ImageByHolderWithOptions", func(p *GoPdf) {
			h, err := ImageHolderByPath("test/res/gopher01.jpg")
			if err != nil {
				t.Fatal(err)
			}
			if err := p.ImageByHolderWithOptions(h, ImageOptions{X: 10, Y: 10, Rect: &Rect{W: 50, H: 50}}); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, conversion := range []float64{0, 2} {
				conversion := conversion
				pdf := new(GoPdf)
				pdf.Start(Config{Unit: UnitPT, ConversionForUnit: conversion, PageSize: Rect{W: 100, H: 200}})
				pdf.SetCompressLevel(0)
				pdf.AddPage()
				tc.draw(pdf)

				var buf bytes.Buffer
				if _, err := pdf.WriteTo(&buf); err != nil {
					t.Fatal(err)
				}
				want := " 50.00 10.00 140.00 cm /I"
				if conversion != 0 {
					want = " 100.00 20.00 280.00 cm /I"
				}
				if !strings.Contains(buf.String(), want) {
					t.Errorf("conversion %v: image matrix %q not found", conversion, want)
				}
			}
		})
	}
}

// AddPageWithOption must convert the trim box as well as the page size.
func TestConversionForUnitAppliesToPageOptionTrimBox(t *testing.T) {
	for _, conversion := range []float64{0, 2} {
		conversion := conversion
		pdf := new(GoPdf)
		pdf.Start(Config{Unit: UnitPT, ConversionForUnit: conversion, PageSize: Rect{W: 100, H: 200}})
		pdf.SetCompressLevel(0)
		pdf.AddPageWithOption(PageOption{
			PageSize: &Rect{W: 100, H: 200},
			TrimBox:  &Box{Left: 5, Top: 5, Right: 10, Bottom: 10},
		})

		var buf bytes.Buffer
		if _, err := pdf.WriteTo(&buf); err != nil {
			t.Fatal(err)
		}
		want := "/TrimBox [ 5.00 5.00 10.00 10.00 ]"
		if conversion != 0 {
			want = "/TrimBox [ 10.00 10.00 20.00 20.00 ]"
		}
		if !strings.Contains(buf.String(), want) {
			t.Errorf("conversion %v: %q not found", conversion, want)
		}
	}
}
