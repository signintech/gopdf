package gopdf

import (
	"bytes"
	"strings"
	"testing"
)

func TestClipPolygonWrite(t *testing.T) {
	cache := &cacheContentClipPolygon{
		pageHeight: 792,
		points: []Point{
			{X: 150, Y: 100},
			{X: 250, Y: 250},
			{X: 50, Y: 250},
		},
	}

	var buf bytes.Buffer
	err := cache.write(&buf, nil)
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	got := buf.String()
	// Y coordinates flipped: 792 - y
	expected := "150.00 692.00 m 250.00 542.00 l 50.00 542.00 l h W n\n"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}

func TestClipPolygonIntegration(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.SetNoCompression() // Disable compression for content inspection
	pdf.AddPage()

	// Triangle clipping path
	pdf.ClipPolygon([]Point{
		{X: 150, Y: 50},
		{X: 250, Y: 200},
		{X: 50, Y: 200},
	})

	pdf.SetFillColor(255, 0, 0)
	pdf.RectFromUpperLeftWithStyle(50, 50, 200, 200, "F")

	pdfBytes := pdf.GetBytesPdf()
	if len(pdfBytes) == 0 {
		t.Fatal("Expected non-empty PDF")
	}

	// Verify PDF content contains clipping operators
	content := string(pdfBytes)

	// Check for move operator (first point of triangle: 150, 50 -> Y flipped to 792)
	if !strings.Contains(content, "150.00 792.00 m") {
		t.Error("Expected move to first point '150.00 792.00 m' in PDF")
	}

	// Check for line operators (l) used in polygon path
	if !strings.Contains(content, " l ") {
		t.Error("Expected line operator 'l' in PDF")
	}

	// Check for clip operator sequence: h (close), W (clip), n (end path)
	if !strings.Contains(content, "h W n") {
		t.Error("Expected clipping operators 'h W n' in PDF")
	}
}

func TestGraphicsStateSaveRestore(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.SetNoCompression()
	pdf.AddPage()

	// Save state, clip, draw clipped rect, restore
	pdf.SaveGraphicsState()
	pdf.ClipPolygon([]Point{{X: 100, Y: 100}, {X: 200, Y: 100}, {X: 150, Y: 200}})
	pdf.SetFillColor(255, 0, 0)
	pdf.RectFromUpperLeftWithStyle(50, 50, 200, 200, "F")
	pdf.RestoreGraphicsState()

	// This rect should NOT be clipped (state restored)
	pdf.SetFillColor(0, 0, 255)
	pdf.RectFromUpperLeftWithStyle(250, 50, 100, 100, "F")

	pdfBytes := pdf.GetBytesPdf()
	content := string(pdfBytes)

	// Verify q (save graphics state) operator present
	if !strings.Contains(content, "q\n") {
		t.Error("Expected 'q' (save graphics state) operator")
	}
	// Verify Q (restore graphics state) operator present
	if !strings.Contains(content, "Q\n") {
		t.Error("Expected 'Q' (restore graphics state) operator")
	}
}

func TestGraphicsStateWriteMethods(t *testing.T) {
	// Test save graphics state write
	saveBuf := bytes.Buffer{}
	save := &cacheContentSaveGraphicsState{}
	err := save.write(&saveBuf, nil)
	if err != nil {
		t.Fatalf("save write failed: %v", err)
	}
	if saveBuf.String() != "q\n" {
		t.Errorf("save: got %q, want %q", saveBuf.String(), "q\n")
	}

	// Test restore graphics state write
	restoreBuf := bytes.Buffer{}
	restore := &cacheContentRestoreGraphicsState{}
	err = restore.write(&restoreBuf, nil)
	if err != nil {
		t.Fatalf("restore write failed: %v", err)
	}
	if restoreBuf.String() != "Q\n" {
		t.Errorf("restore: got %q, want %q", restoreBuf.String(), "Q\n")
	}
}

func TestGraphicsStateClipScoping(t *testing.T) {
	// This test verifies that RestoreGraphicsState properly ends clipping.
	// We draw: blue rect, clipped red rect, unclipped green rect.
	// The green rect must appear AFTER the Q (restore) operator.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.SetNoCompression()
	pdf.AddPage()

	// 1. Blue rectangle (before clip)
	pdf.SetFillColor(0, 0, 255) // blue
	pdf.RectFromUpperLeftWithStyle(50, 50, 100, 100, "F")

	// 2. Red rectangle (clipped)
	pdf.SaveGraphicsState()
	pdf.ClipPolygon([]Point{{X: 250, Y: 50}, {X: 350, Y: 200}, {X: 150, Y: 200}})
	pdf.SetFillColor(255, 0, 0) // red
	pdf.RectFromUpperLeftWithStyle(150, 50, 200, 200, "F")
	pdf.RestoreGraphicsState()

	// 3. Green rectangle (after restore - must NOT be clipped)
	pdf.SetFillColor(0, 255, 0) // green
	pdf.RectFromUpperLeftWithStyle(400, 50, 100, 100, "F")

	pdfBytes := pdf.GetBytesPdf()
	content := string(pdfBytes)

	// Find the clip position
	clipPos := strings.Index(content, "h W n")
	if clipPos == -1 {
		t.Fatal("Missing clip operators 'h W n'")
	}

	// Find q immediately before clip (our SaveGraphicsState)
	beforeClip := content[:clipPos]
	qPos := strings.LastIndex(beforeClip, "q\n")
	if qPos == -1 {
		t.Fatal("Missing 'q' before clip")
	}

	// Find Q after clip (our RestoreGraphicsState)
	afterClip := content[clipPos:]
	qRelPos := strings.Index(afterClip, "Q\n")
	if qRelPos == -1 {
		t.Fatal("Missing 'Q' after clip")
	}
	bigQPos := clipPos + qRelPos

	// Green rect must come after the Q
	greenPos := strings.LastIndex(content, "0.000 1.000 0.000 rg")
	if greenPos == -1 {
		t.Fatal("Missing green fill '0.000 1.000 0.000 rg'")
	}

	// Verify order: q < clip < Q < green
	if !(qPos < clipPos && clipPos < bigQPos && bigQPos < greenPos) {
		t.Errorf("Wrong order: q=%d, clip=%d, Q=%d, green=%d", qPos, clipPos, bigQPos, greenPos)
	}
}
