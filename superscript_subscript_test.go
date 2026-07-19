package gopdf

import (
	"testing"
)

// newScriptTestPDF returns a started A4 PDF with a font loaded and content
// compression disabled so the content stream can be inspected as raw bytes.
func newScriptTestPDF(t *testing.T) *GoPdf {
	t.Helper()
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		t.Fatal(err)
	}
	pdf.SetNoCompression()
	return pdf
}

// Liberation Serif carries designer-provided script metrics; all four parsed
// values must be positive (OS/2 stores even the subscript offset as a
// positive "distance below the baseline").
func TestParseOS2ScriptMetrics(t *testing.T) {
	pdf := newScriptTestPDF(t)
	ttfp := pdf.curr.FontISubset.ttfp
	if v := ttfp.SubscriptYSize(); v <= 0 {
		t.Fatalf("SubscriptYSize = %d, want > 0", v)
	}
	if v := ttfp.SubscriptYOffset(); v <= 0 {
		t.Fatalf("SubscriptYOffset = %d, want > 0", v)
	}
	if v := ttfp.SuperscriptYSize(); v <= 0 {
		t.Fatalf("SuperscriptYSize = %d, want > 0", v)
	}
	if v := ttfp.SuperscriptYOffset(); v <= 0 {
		t.Fatalf("SuperscriptYOffset = %d, want > 0", v)
	}
}
