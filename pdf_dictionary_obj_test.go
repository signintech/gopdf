package gopdf

import "testing"

func TestSubsetIncludesComponentsOfNestedCompositeGlyphs(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	// A subset of DejaVu Sans with only U+01D6, a composite of macron and udieresis,
	// which is itself a composite of u and dieresis.
	if err := pdf.AddTTFFont("DejaVuSans-Subset", "./test/res/DejaVuSans-Subset.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("DejaVuSans-Subset", "", 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("\u01D6"); err != nil {
		t.Fatal(err)
	}

	var dict *PdfDictionaryObj
	for _, obj := range pdf.pdfObjs {
		if d, ok := obj.(*PdfDictionaryObj); ok {
			dict = d
		}
	}
	if dict == nil {
		t.Fatal("no font dictionary")
	}

	_, loca, err := dict.makeGlyfAndLocaTable()
	if err != nil {
		t.Fatal(err)
	}

	// Every glyph after .notdef is part of U+01D6, so each one needs its data in the embedded font.
	for glyph := 1; glyph < len(loca)-1; glyph++ {
		if loca[glyph+1] == loca[glyph] {
			t.Errorf("glyph %d has no data in the embedded font", glyph)
		}
	}
}
