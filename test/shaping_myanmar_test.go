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

	pdf.AddPage()
	if err := pdf.SetFont("NotoSans", "", 12); err != nil {
		t.Fatal(err)
	}
	pdf.SetXY(50, 100)
	if err := pdf.Text("Welcome, "); err != nil {
		t.Fatal(err)
	}

	if err := pdf.SetFont("NotoSansMyanmar", "", 12); err != nil {
		t.Fatal(err)
	}
	if err := pdf.Text("ကြိုဆိုပါတယ် ရော ရောင်း သစ္စာ မင်္ဂလာပါ တနင်္ဂနွေ အင်္ဂါ သင်္ကြန်"); err != nil {
		t.Fatal(err)
	}

	out := "./out/shaping_myanmar.pdf"
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
