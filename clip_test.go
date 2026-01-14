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
