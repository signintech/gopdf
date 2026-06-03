package gopdf

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func setupFallbackPDF(t *testing.T) *GoPdf {
	t.Helper()
	if err := initTesting(); err != nil {
		t.Fatal(err)
	}
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.SetNoCompression()
	pdf.AddPage()
	if err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("Amiri-Regular", "./examples/arabic/Amiri-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		t.Fatal(err)
	}
	return &pdf
}

func TestFontFallbackConfigurationValidation(t *testing.T) {
	pdf := setupFallbackPDF(t)

	if err := pdf.SetFontFallback("LiberationSerif-Regular", "", FontFallback{Family: "Amiri-Regular"}); err != nil {
		t.Fatalf("set fallback: %v", err)
	}
	key := fontKey{family: "LiberationSerif-Regular", style: 0}
	fallbacks := pdf.fontFallbacks[key]
	if len(fallbacks) != 1 || fallbacks[0].family != "Amiri-Regular" {
		t.Fatalf("unexpected fallback list: %#v", fallbacks)
	}

	if err := pdf.SetFontFallback("LiberationSerif-Regular", ""); err != nil {
		t.Fatalf("clear fallback: %v", err)
	}
	if _, ok := pdf.fontFallbacks[key]; ok {
		t.Fatal("expected empty fallback list to clear configuration")
	}

	if err := pdf.SetFontFallback("Missing", "", FontFallback{Family: "Amiri-Regular"}); !errors.Is(err, ErrFontNotFound) {
		t.Fatalf("expected ErrFontNotFound for missing primary, got %v", err)
	}
	if err := pdf.SetFontFallback("LiberationSerif-Regular", "", FontFallback{Family: "Missing"}); !errors.Is(err, ErrFontNotFound) {
		t.Fatalf("expected ErrFontNotFound for missing fallback, got %v", err)
	}
}

func TestFontFallbackTextUsesFallbackFont(t *testing.T) {
	pdf := setupFallbackPDF(t)
	if err := pdf.SetFontFallback("LiberationSerif-Regular", "", FontFallback{Family: "Amiri-Regular"}); err != nil {
		t.Fatal(err)
	}

	fallbackFont, ok := pdf.subsetFontByKey(fontKey{family: "Amiri-Regular", style: 0})
	if !ok {
		t.Fatal("fallback font not found")
	}
	primaryFont, ok := pdf.subsetFontByKey(fontKey{family: "LiberationSerif-Regular", style: 0})
	if !ok {
		t.Fatal("primary font not found")
	}

	x0 := pdf.GetX()
	width, err := pdf.MeasureTextWidth("AشB")
	if err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("AشB"); err != nil {
		t.Fatal(err)
	}
	if got := pdf.GetX() - x0; math.Abs(got-width) > 0.001 {
		t.Fatalf("Text advance = %f, MeasureTextWidth = %f", got, width)
	}

	if !fallbackFont.CharacterToGlyphIndex.KeyExists('ش') {
		t.Fatal("fallback font did not receive Arabic rune")
	}
	if primaryFont.CharacterToGlyphIndex.KeyExists('ش') {
		t.Fatal("primary font received rune that should use fallback")
	}

	content := string(pdf.GetBytesPdf())
	if !strings.Contains(content, "/F1 14 Tf") || !strings.Contains(content, "/F2 14 Tf") {
		t.Fatalf("expected content to switch between primary and fallback fonts, got:\n%s", content)
	}
}

func TestFontFallbackDisabledKeepsSubstituteBehavior(t *testing.T) {
	pdf := setupFallbackPDF(t)
	missing := 0
	if err := pdf.AddTTFFontWithOption("LiberationSerif-Missing", "./test/res/LiberationSerif-Regular.ttf", TtfOption{
		OnGlyphNotFound: func(r rune) {
			missing++
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Missing", "", 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("ش"); err != nil {
		t.Fatal(err)
	}
	if missing != 1 {
		t.Fatalf("missing glyph hook count = %d, want 1", missing)
	}
	font, ok := pdf.subsetFontByKey(fontKey{family: "Amiri-Regular", style: 0})
	if !ok {
		t.Fatal("fallback font not found")
	}
	if font.CharacterToGlyphIndex.KeyExists('ش') {
		t.Fatal("unconfigured fallback font should not receive missing rune")
	}
}

func TestFontFallbackCloneKeepsFallbackIndependent(t *testing.T) {
	pdf := setupFallbackPDF(t)
	if err := pdf.SetFontFallback("LiberationSerif-Regular", "", FontFallback{Family: "Amiri-Regular"}); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("Aش"); err != nil {
		t.Fatal(err)
	}

	clone := pdf.Clone()
	if err := pdf.Text("س"); err != nil {
		t.Fatal(err)
	}
	if err := clone.Text("ل"); err != nil {
		t.Fatal(err)
	}

	origFallback, ok := pdf.subsetFontByKey(fontKey{family: "Amiri-Regular", style: 0})
	if !ok {
		t.Fatal("original fallback not found")
	}
	cloneFallback, ok := clone.subsetFontByKey(fontKey{family: "Amiri-Regular", style: 0})
	if !ok {
		t.Fatal("clone fallback not found")
	}
	if origFallback == cloneFallback {
		t.Fatal("clone fallback points to original font")
	}
	if !origFallback.CharacterToGlyphIndex.KeyExists('س') {
		t.Fatal("original fallback missing original-only rune")
	}
	if cloneFallback.CharacterToGlyphIndex.KeyExists('س') {
		t.Fatal("clone fallback contains original-only rune")
	}
	if !cloneFallback.CharacterToGlyphIndex.KeyExists('ل') {
		t.Fatal("clone fallback missing clone-only rune")
	}
	if origFallback.CharacterToGlyphIndex.KeyExists('ل') {
		t.Fatal("original fallback contains clone-only rune")
	}

	origContent := string(pdf.GetBytesPdf())
	cloneContent := string(clone.GetBytesPdf())
	if !strings.Contains(origContent, "/F2 14 Tf") || !strings.Contains(cloneContent, "/F2 14 Tf") {
		t.Fatal("expected original and clone PDFs to use fallback font")
	}
}

func TestFontFallbackFromFontContainerPreservesParsedFont(t *testing.T) {
	if err := initTesting(); err != nil {
		t.Fatal(err)
	}
	container := &FontContainer{}
	if err := container.AddTTFFont("Amiri-Regular", "./examples/arabic/Amiri-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFontFromFontContainer("Amiri-Regular", container); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFontFallback("LiberationSerif-Regular", "", FontFallback{Family: "Amiri-Regular"}); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("ش"); err != nil {
		t.Fatal(err)
	}
	fallbackFont, ok := pdf.subsetFontByKey(fontKey{family: "Amiri-Regular", style: 0})
	if !ok {
		t.Fatal("fallback font not found")
	}
	if !fallbackFont.CharacterToGlyphIndex.KeyExists('ش') {
		t.Fatal("font container fallback did not retain parsed glyph data")
	}
}
