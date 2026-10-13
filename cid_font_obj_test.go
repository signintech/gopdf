package gopdf

import (
	"bytes"
	"testing"
)

func TestCIDFontType2WritesIdentityCIDToGIDMap(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("Hi"); err != nil {
		t.Fatal(err)
	}
	raw, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(raw, []byte("/Subtype /CIDFontType2\n")) {
		t.Fatal("expected a Type 2 CID font")
	}
	if !bytes.Contains(raw, []byte("/CIDToGIDMap /Identity\n")) {
		t.Fatal("CIDFontType2 dictionary is missing /CIDToGIDMap /Identity")
	}
}
