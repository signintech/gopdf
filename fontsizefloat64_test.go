package gopdf

import (
	"fmt"
	"testing"
)

func TestContentObjCalTextHeight(t *testing.T) {
	intfontsize := 7
	have := ContentObjCalTextHeight(intfontsize)
	want := float64(intfontsize) * 0.7
	if have != want {
		t.Errorf("ContentObjCalTextHeight(%d) = %f; want %f\n", intfontsize, have, want)
	}

	floatfontsize := 7.2
	have = ContentObjCalTextHeightPrecise(floatfontsize)
	want = float64(floatfontsize) * 0.7
	if have != want {
		t.Errorf("ContentObjCalTextHeight(%d) = %f; want %f\n", intfontsize, have, want)
	}

}

func TestSetFontCheckGetX(t *testing.T) {
	prefix := "Afont"
	font := "test/res/LiberationSerif-Regular.ttf"
	pdf := GoPdf{}
	pdf.Start(Config{Unit: UnitPT, PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFontWithOption(prefix, font, TtfOption{UseKerning: true}); err != nil {
		t.Error(err)
	}

	// Baseline: SetFont(int) + AddText
	if err := pdf.SetFont(prefix, "", int(50)); err != nil {
		t.Error(err)
	}

	pdf.SetXY(30.0, 30.0)
	pdf.Text(prefix)
	wantx, wanty := pdf.GetX(), pdf.GetY()

	// ensure that GetX, GetY work as expected
	pdf.SetXY(30.0, 30.0)
	pdf.Text(prefix + prefix)
	havex, havey := pdf.GetX(), pdf.GetY()
	if havex <= wantx || havey != wanty {
		t.Errorf("wanted (x,y) => (x', y') with x<x' and y==y', got (%f, %f) => (%f, %f)\n", wantx, wanty, havex, havey)
	}
}

func moveAndAdd(pdf *GoPdf, prefix string) (float64, float64) {
	pdf.SetXY(30.0, 30.0)
	pdf.Text(prefix)
	return pdf.GetX(), pdf.GetY()
}
func setup(t *testing.T) (string, *GoPdf, float64, float64) {
	prefix := "Afont"
	font := "test/res/LiberationSerif-Regular.ttf"
	pdf := GoPdf{}
	pdf.Start(Config{Unit: UnitPT, PageSize: *PageSizeA4})
	pdf.AddPage()
	if err := pdf.AddTTFFontWithOption(prefix, font, TtfOption{UseKerning: true}); err != nil {
		t.Error(err)
	}

	// Baseline: SetFont(int) + AddText
	if err := pdf.SetFont(prefix, "", int(50)); err != nil {
		t.Error(err)
	}
	hx, hy := moveAndAdd(&pdf, prefix)
	return prefix, &pdf, hx, hy
}
func TestSetFontWithFloat(t *testing.T) {
	prefix, pdf, wantx, wanty := setup(t)
	// try it with fontsize = float64(50)
	if err := pdf.SetFont(prefix, "", float64(50)); err != nil {
		t.Error(err)
	}
	havex, havey := moveAndAdd(pdf, prefix)
	if (havex != wantx) || (havey != wanty) {
		t.Errorf("SetFont(float64) + '%s' => \nhave = %f, %f;\nwant = %f,%f\n",
			prefix, havex, havey, wantx, wanty)
	}
}

func TestSetFontWithUint(t *testing.T) {
	prefix, pdf, wantx, wanty := setup(t)
	// try it with fontsize = uint8(50)
	if err := pdf.SetFont(prefix, "", uint8(50)); err != nil {
		t.Error(err)
	}
	havex, havey := moveAndAdd(pdf, prefix)
	if (havex != wantx) || (havey != wanty) {
		t.Errorf("SetFont(uint8) + AddText => %f, %f; want = %f,%f\n",
			havex, havey, wantx, wanty)
	}
}
func TestSetFontWithString(t *testing.T) {
	prefix, pdf, _, _ := setup(t)
	// Try with a string
	err := pdf.SetFont(prefix, "", string("50"))
	if err == nil {
		t.Errorf("SetFont(string) + AddText: Should have gotten an error!\n")
	}
}

func ExampleFormatFloatTrim() {
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10.0))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10.01))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10.001))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10.0001))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(10.00001))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(9.99999))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(9.9999))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(9.999))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(9.99))
	fmt.Printf("/F1 %s Tf\n", FormatFloatTrim(9.9))
	// Output: /F1 10 Tf
	// /F1 10 Tf
	// /F1 10.01 Tf
	// /F1 10.001 Tf
	// /F1 10 Tf
	// /F1 10 Tf
	// /F1 10 Tf
	// /F1 10 Tf
	// /F1 9.999 Tf
	// /F1 9.99 Tf
	// /F1 9.9 Tf
}
