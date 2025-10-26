package shaping

import (
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/signintech/gopdf"
)

// scriptForRune returns a coarse script identifier used to pick a font family.
func scriptForRune(r rune) string {
	switch {
	case unicode.In(r, unicode.Arabic):
		return "Arabic"
	case unicode.In(r, unicode.Devanagari):
		return "Devanagari"
	case unicode.In(r, unicode.Myanmar):
		return "Myanmar"
	case unicode.In(r, unicode.Thai):
		return "Thai"
	case unicode.In(r, unicode.Hangul):
		return "Korean"
	case unicode.In(r, unicode.Hiragana, unicode.Katakana):
		return "Japanese"
	case unicode.In(r, unicode.Han): // CJK Han; route to SC by default
		return "CJK"
	default:
		return "Latin"
	}
}

func familyForScript(script string) string {
	switch script {
	case "Arabic":
		return "NotoSansArabic"
	case "Devanagari":
		return "NotoSansDevanagari"
	case "Myanmar":
		return "NotoSansMyanmar"
	case "Thai":
		return "NotoSansThai"
	case "Korean":
		return "NotoSansKR"
	case "Japanese":
		return "NotoSansJP"
	case "CJK":
		return "NotoSansSC"
	default:
		return "NotoSans"
	}
}

type textRun struct {
	text   string
	family string
}

func splitRunsByScript(s string) []textRun {
	runes := []rune(s)
	if len(runes) == 0 {
		return nil
	}
	curScript := scriptForRune(runes[0])
	curFamily := familyForScript(curScript)
	var sb strings.Builder
	runs := make([]textRun, 0, 8)
	for _, r := range runes {
		// keep whitespace with the current run to avoid jitter
		if unicode.IsSpace(r) {
			sb.WriteRune(r)
			continue
		}
		s := scriptForRune(r)
		fam := familyForScript(s)
		if fam == curFamily {
			sb.WriteRune(r)
			continue
		}
		// flush
		if sb.Len() > 0 {
			runs = append(runs, textRun{text: sb.String(), family: curFamily})
			sb.Reset()
		}
		curFamily = fam
		sb.WriteRune(r)
	}
	if sb.Len() > 0 {
		runs = append(runs, textRun{text: sb.String(), family: curFamily})
	}
	return runs
}

func TestShapingSingleStringMultiScript(t *testing.T) {
	// Ensure output directory exists
	if err := os.MkdirAll("./out", 0o777); err != nil {
		t.Fatal(err)
	}

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Register fonts present in ./test directory (this file resides in the same package)
	if err := pdf.AddTTFFont("NotoSans", "./NotoSans-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansArabic", "./NotoSansArabic-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansJP", "./NotoSansJP-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansKR", "./NotoSansKR-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansSC", "./NotoSansSC-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansThai", "./NotoSansThai-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansDevanagari", "./NotoSansDevanagari-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansMyanmar", "./NotoSansMyanmar-Regular.ttf"); err != nil {
		t.Fatal(err)
	}

	pdf.AddPage()

	// One single string containing multiple languages (LTR and RTL mixed)
	line := "Hello مرحبا こんにちは 안녕하세요 你好 สวัสดี नमस्ते မင်္ဂလာပါ​​ ဒီနေ့တော့ ရေမချိုးတော့ဘူး"

	// Start position
	pdf.SetXY(50, 120)

	// Render sequentially, font-per-run. X/Y advance is handled by the library.
	for _, run := range splitRunsByScript(line) {
		if err := pdf.SetFont(run.family, "", 14); err != nil {
			t.Fatal(err)
		}
		if err := pdf.Text(run.text); err != nil {
			t.Fatal(err)
		}
	}

	out := "./out/shaping_single_string.pdf"
	if err := pdf.WritePdf(out); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(out)
	if err != nil {
		t.Fatal(err)
	}
	if st.Size() <= 0 {
		t.Fatalf("generated PDF is empty: %s", out)
	}
}
