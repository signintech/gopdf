package shaping

import (
	"os"
	"testing"

	"github.com/signintech/gopdf"
)

func TestShapingMyanmarSelection(t *testing.T) {
	// Ensure output directory exists under this package
	if err := os.MkdirAll("./out", 0o777); err != nil {
		t.Fatal(err)
	}

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	if err := pdf.AddTTFFont("NotoSans", "./NotoSans-Regular.ttf"); err != nil {
		t.Fatal(err)
	}
	if err := pdf.AddTTFFont("NotoSansMyanmar", "./NotoSansMyanmar-Regular.ttf"); err != nil {
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

	pdf.AddPage()
	if err := pdf.SetFont("NotoSans", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 100)
	if err := pdf.Text("Well, "); err != nil {
		t.Fatal(err)
	}

	if err := pdf.SetFont("NotoSansMyanmar", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Cell(nil, "ကြိုဆိုပါတယ် ရော ရောင်း သစ္စာ မင်္ဂလာပါ တနင်္ဂနွေ အင်္ဂါ သင်္ကြန် ဒီနေ့တော့ ရေမချိုးတော့ဘူး"); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 130)

	if err := pdf.Cell(nil, "ဒီနေ့ နေ့ နွေ နှောင်း နိုး နိုင် နှိုင်း နွှေး လှိုင်း လျှိုး နှိုက် နှီး ဓါတု ဓာတ် နူး နှူး လှိုင်း မှိုင်း နှင့် လှပ် မှား"); err != nil {
		t.Fatal(err)
	}

	// // Arabic
	// if err := pdf.SetFont("NotoSansArabic", "", 12); err != nil {
	// 	t.Fatal(err)
	// }
	// pdf.SetXY(50, 130)
	// if err := pdf.Text("مرحبا بالعالم"); err != nil {
	// 	t.Fatal(err)
	// }

	// Japanese
	if err := pdf.SetFont("NotoSansJP", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 160)
	if err := pdf.Text("こんにちは 世界"); err != nil {
		t.Fatal(err)
	}

	// Korean
	if err := pdf.SetFont("NotoSansKR", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 190)
	if err := pdf.Text("안녕하세요 세계"); err != nil {
		t.Fatal(err)
	}

	// Simplified Chinese
	if err := pdf.SetFont("NotoSansSC", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 220)
	if err := pdf.Text("你好 世界"); err != nil {
		t.Fatal(err)
	}

	// Thai
	if err := pdf.SetFont("NotoSansThai", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 250)
	if err := pdf.Text("สวัสดี โลก"); err != nil {
		t.Fatal(err)
	}

	// Devanagari (Hindi)
	if err := pdf.SetFont("NotoSansDevanagari", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 280)
	if err := pdf.Text("नमस्ते दुनिया"); err != nil {
		t.Fatal(err)
	}

	// Mixed line: Hello in all languages
	pdf.SetXY(50, 310)
	if err := pdf.SetFont("NotoSans", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("Hello "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansArabic", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("مرحبا "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansJP", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("こんにちは "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansKR", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("안녕하세요 "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansSC", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("你好 "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansThai", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("สวัสดี "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansDevanagari", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("नमस्ते "); err != nil {
		t.Fatal(err)
	}
	if err := pdf.SetFont("NotoSansMyanmar", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("မင်္ဂလာပါ"); err != nil {
		t.Fatal(err)
	}

	// Explicit positioning as requested
	// Hello at x:100 (y:200)
	if err := pdf.SetFont("NotoSans", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(100, 200)
	if err := pdf.Text("Hello"); err != nil {
		t.Fatal(err)
	}
	// こんにちは at x:200 (same baseline y:200)
	if err := pdf.SetFont("NotoSansJP", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(200, 200)
	if err := pdf.Text("こんにちは"); err != nil {
		t.Fatal(err)
	}
	// မင်္ဂလာပါ at x:100, y:250
	if err := pdf.SetFont("NotoSansMyanmar", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(100, 250)
	if err := pdf.Text("မင်္ဂလာပါ"); err != nil {
		t.Fatal(err)
	}

	out := "./out/shaping_multiscript.pdf"
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
