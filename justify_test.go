package gopdf

import (
	"bytes"
	"testing"
)

func TestJustifyAdjustment(t *testing.T) {
	tests := []struct {
		name     string
		slack    float64
		gapCount int
		fontSize float64
		want     int
	}{
		{"zero gaps", 30, 0, 10, 0},
		{"zero slack", 0, 3, 10, 0},
		{"negative slack", -5, 3, 10, 0},
		{"zero font size", 30, 3, 0, 0},
		{"even distribution", 30, 3, 10, -1000},   // 10pt/gap -> -(10*1000/10)
		{"rounded distribution", 10, 3, 12, -278}, // 3.3333pt/gap -> -(277.78) rounded
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := justifyAdjustment(tt.slack, tt.gapCount, tt.fontSize)
			if got != tt.want {
				t.Fatalf("justifyAdjustment(%v,%d,%v) = %d, want %d",
					tt.slack, tt.gapCount, tt.fontSize, got, tt.want)
			}
		})
	}
}

func TestNonSpaceBounds(t *testing.T) {
	first, last := nonSpaceBounds("  hi  ")
	if first != 2 || last != 3 {
		t.Fatalf("nonSpaceBounds = (%d,%d), want (2,3)", first, last)
	}
	if f, l := nonSpaceBounds("   "); f != -1 || l != -1 {
		t.Fatalf("all-spaces bounds = (%d,%d), want (-1,-1)", f, l)
	}
	if f, l := nonSpaceBounds(""); f != -1 || l != -1 {
		t.Fatalf("empty bounds = (%d,%d), want (-1,-1)", f, l)
	}
}

func TestInteriorSpaceCount(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"hello world", 1},
		{"a b c", 2},
		{"  hello   world  ", 3}, // 3 interior spaces; leading/trailing excluded
		{"word", 0},
		{"", 0},
		{"   ", 0},
	}
	for _, tt := range tests {
		if got := interiorSpaceCount(tt.text); got != tt.want {
			t.Fatalf("interiorSpaceCount(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

func TestLineAlign(t *testing.T) {
	if got := lineAlign(Justify, false); got != Justify {
		t.Fatalf("non-last justify line = %d, want %d (Justify)", got, Justify)
	}
	if got := lineAlign(Justify, true); got != Left {
		t.Fatalf("last justify line = %d, want %d (Left)", got, Left)
	}
	if got := lineAlign(Justify|Top, true); got != Left|Top {
		t.Fatalf("last justify+top line = %d, want %d (Left|Top)", got, Left|Top)
	}
	if got := lineAlign(Center, true); got != Center {
		t.Fatalf("last center line = %d, want %d (unchanged)", got, Center)
	}
	if got := lineAlign(Left, false); got != Left {
		t.Fatalf("non-last left line = %d, want %d (unchanged)", got, Left)
	}
}

// newJustifyTestPDF returns a started A4 PDF with a font loaded and content
// compression disabled so the content stream can be inspected as raw bytes.
func newJustifyTestPDF(t *testing.T) *GoPdf {
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

// countAdjust counts justify/kerning adjustment markers ">-" in a content
// stream. Justify numbers are always negative, so ">-" marks each injected gap.
func countAdjust(b []byte) int {
	return bytes.Count(b, []byte(">-"))
}

func TestJustifyCellInjectsAdjustments(t *testing.T) {
	const text = "the quick brown fox" // 3 interior spaces

	justified := newJustifyTestPDF(t)
	justified.SetXY(20, 40)
	if err := justified.CellWithOption(&Rect{W: 400, H: 20}, text,
		CellOption{Align: Justify | Top}); err != nil {
		t.Fatal(err)
	}
	jb, err := justified.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	left := newJustifyTestPDF(t)
	left.SetXY(20, 40)
	if err := left.CellWithOption(&Rect{W: 400, H: 20}, text,
		CellOption{Align: Left | Top}); err != nil {
		t.Fatal(err)
	}
	lb, err := left.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	// The only systematic difference between the two documents is the injected
	// justify adjustments — one per interior space.
	if got := countAdjust(jb) - countAdjust(lb); got != 3 {
		t.Fatalf("justify adjustment count delta = %d, want 3", got)
	}
}

// A paragraph that wraps to multiple lines should have justify adjustments
// (on all but the last line).
func TestJustifyMultiCellWraps(t *testing.T) {
	const text = "the quick brown fox jumps over the lazy dog again and again"

	justified := newJustifyTestPDF(t)
	justified.SetXY(20, 40)
	if err := justified.MultiCellWithOption(&Rect{W: 150, H: 400}, text,
		CellOption{Align: Justify}); err != nil {
		t.Fatal(err)
	}
	jb, err := justified.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	left := newJustifyTestPDF(t)
	left.SetXY(20, 40)
	if err := left.MultiCellWithOption(&Rect{W: 150, H: 400}, text,
		CellOption{Align: Left}); err != nil {
		t.Fatal(err)
	}
	lb, err := left.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	if countAdjust(jb)-countAdjust(lb) <= 0 {
		t.Fatalf("expected justify adjustments in a wrapped paragraph, got delta %d",
			countAdjust(jb)-countAdjust(lb))
	}
}

// The last line of a justified paragraph is left-aligned. When the whole text
// fits on ONE line, that sole line is the last line, so MultiCell must NOT
// justify it (unlike CellWithOption, which does).
func TestJustifyMultiCellLastLineNotStretched(t *testing.T) {
	const text = "quick brown" // fits on one line in a 400-wide rect

	justified := newJustifyTestPDF(t)
	justified.SetXY(20, 40)
	if err := justified.MultiCellWithOption(&Rect{W: 400, H: 100}, text,
		CellOption{Align: Justify}); err != nil {
		t.Fatal(err)
	}
	jb, err := justified.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	left := newJustifyTestPDF(t)
	left.SetXY(20, 40)
	if err := left.MultiCellWithOption(&Rect{W: 400, H: 100}, text,
		CellOption{Align: Left}); err != nil {
		t.Fatal(err)
	}
	lb, err := left.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}

	if got := countAdjust(jb) - countAdjust(lb); got != 0 {
		t.Fatalf("last (only) line must not be justified; adjustment delta = %d, want 0", got)
	}
}

// renderKernedCell renders a single cell using the Ubuntu font (which has a
// legacy `kern` table, unlike LiberationSerif) with kerning optionally enabled,
// and returns the raw, uncompressed content bytes.
func renderKernedCell(t *testing.T, text string, align int, useKerning bool) []byte {
	t.Helper()
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFontWithOption("Ubuntu-L",
		"./examples/outline_example/Ubuntu-L.ttf", TtfOption{UseKerning: useKerning}); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("Ubuntu-L", "", 14); err != nil {
		t.Fatal(err)
	}
	pdf.SetNoCompression()
	pdf.SetXY(20, 40)
	if err := pdf.CellWithOption(&Rect{W: 500, H: 20}, text,
		CellOption{Align: align | Top}); err != nil {
		t.Fatal(err)
	}
	b, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Justify must work when kerning is enabled: both mechanisms emit position
// numbers into the same TJ array. This verifies (1) kerning genuinely alters
// the output for this font/text, and (2) enabling justify on top of kerning
// adds exactly one adjustment per interior space — i.e. justify coexists with
// kerning and is unaffected by it.
func TestJustifyCoexistsWithKerning(t *testing.T) {
	const text = "AVA WAV YAV" // capital pairs that kern in Ubuntu; 2 interior spaces

	leftNoKern := renderKernedCell(t, text, Left, false)
	leftKern := renderKernedCell(t, text, Left, true)
	justifyKern := renderKernedCell(t, text, Justify, true)

	// (1) Kerning must actually be active for this font/text.
	if bytes.Equal(leftNoKern, leftKern) {
		t.Fatal("kerning enabled produced identical output to kerning disabled; " +
			"the font/text no longer exercises kerning, so this test is vacuous")
	}

	// (2) Justify adds exactly its interior-space adjustments on top of kerning.
	if got, want := countAdjust(justifyKern)-countAdjust(leftKern), interiorSpaceCount(text); got != want {
		t.Fatalf("with kerning enabled, justify adjustment delta = %d, want %d", got, want)
	}
}
