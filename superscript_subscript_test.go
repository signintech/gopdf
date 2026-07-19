package gopdf

import (
	"bytes"
	"fmt"
	"math"
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

func TestSetFontWithScriptStyleResolvesFace(t *testing.T) {
	pdf := newScriptTestPDF(t)
	if err := pdf.SetFontWithStyle("LiberationSerif-Regular", Superscript, 14); err != nil {
		t.Fatalf("superscript style should resolve the regular face: %v", err)
	}
	if err := pdf.SetFontWithStyle("LiberationSerif-Regular", Subscript|Underline, 14); err != nil {
		t.Fatalf("subscript|underline should resolve the regular face: %v", err)
	}
}

func TestScriptFontSizeAndRise(t *testing.T) {
	pdf := newScriptTestPDF(t)
	f := pdf.curr.FontISubset

	size, rise := scriptFontSizeAndRise(f, Regular, 14)
	if size != 14 || rise != 0 {
		t.Fatalf("regular = (%v, %v), want (14, 0)", size, rise)
	}

	supSize, supRise := scriptFontSizeAndRise(f, Superscript, 14)
	if supSize <= 0 || supSize >= 14 {
		t.Fatalf("superscript size = %v, want in (0, 14)", supSize)
	}
	if supRise <= 0 {
		t.Fatalf("superscript rise = %v, want > 0", supRise)
	}

	subSize, subRise := scriptFontSizeAndRise(f, Subscript, 14)
	if subSize <= 0 || subSize >= 14 {
		t.Fatalf("subscript size = %v, want in (0, 14)", subSize)
	}
	if subRise >= 0 {
		t.Fatalf("subscript rise = %v, want < 0", subRise)
	}

	// Nonsensical combination: superscript wins.
	bothSize, bothRise := scriptFontSizeAndRise(f, Superscript|Subscript, 14)
	if bothSize != supSize || bothRise != supRise {
		t.Fatalf("both flags = (%v, %v), want superscript values (%v, %v)",
			bothSize, bothRise, supSize, supRise)
	}
}

func TestScriptFontSizeAndRiseFallback(t *testing.T) {
	f := &SubsetFontObj{} // zero-value parser: no usable OS/2 metrics
	const eps = 1e-9

	size, rise := scriptFontSizeAndRise(f, Superscript, 10)
	if math.Abs(size-5.83) > eps || math.Abs(rise-3.33) > eps {
		t.Fatalf("superscript fallback = (%v, %v), want (5.83, 3.33)", size, rise)
	}

	size, rise = scriptFontSizeAndRise(f, Subscript, 10)
	if math.Abs(size-5.83) > eps || math.Abs(rise+1.41) > eps {
		t.Fatalf("subscript fallback = (%v, %v), want (5.83, -1.41)", size, rise)
	}
}

func TestSuperscriptEmitsScaledTfAndTs(t *testing.T) {
	pdf := newScriptTestPDF(t)
	pdf.SetXY(20, 40)
	if err := pdf.Cell(nil, "E = mc"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFontWithStyle("LiberationSerif-Regular", Superscript, 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Cell(nil, "2"); err != nil {
		t.Fatal(err)
	}
	b, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	effSize, rise := scriptFontSizeAndRise(pdf.curr.FontISubset, Superscript, 14)
	wantState := fmt.Sprintf("%s Tf 0 Tc %s Ts", FormatFloatTrim(effSize), FormatFloatTrim(rise))
	if !bytes.Contains(b, []byte(wantState)) {
		t.Fatalf("content stream missing superscript state %q", wantState)
	}
	// The normal-text draw must reset the rise (text state persists across BT/ET).
	if !bytes.Contains(b, []byte("14 Tf 0 Tc 0 Ts")) {
		t.Fatalf("content stream missing 0 Ts reset for normal text")
	}
}

func TestSubscriptEmitsNegativeTs(t *testing.T) {
	pdf := newScriptTestPDF(t)
	pdf.SetXY(20, 40)
	if err := pdf.SetFontWithStyle("LiberationSerif-Regular", Subscript, 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Cell(nil, "2"); err != nil {
		t.Fatal(err)
	}
	b, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	effSize, rise := scriptFontSizeAndRise(pdf.curr.FontISubset, Subscript, 14)
	if rise >= 0 {
		t.Fatalf("subscript rise = %v, want < 0", rise)
	}
	wantState := fmt.Sprintf("%s Tf 0 Tc %s Ts", FormatFloatTrim(effSize), FormatFloatTrim(rise))
	if !bytes.Contains(b, []byte(wantState)) {
		t.Fatalf("content stream missing subscript state %q", wantState)
	}
}

// Justify must keep working for script-styled text; the TJ adjustment is
// computed against the effective (scaled) Tf size.
func TestJustifyWithSuperscript(t *testing.T) {
	pdf := newScriptTestPDF(t)
	if err := pdf.SetFontWithStyle("LiberationSerif-Regular", Superscript, 14); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(20, 40)
	if err := pdf.CellWithOption(&Rect{W: 400, H: 20}, "the quick brown fox",
		CellOption{Align: Justify | Top}); err != nil {
		t.Fatal(err)
	}
	b, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}
	if got := countAdjust(b); got != 3 {
		t.Fatalf("justify adjustments = %d, want 3 (one per interior space)", got)
	}
}
