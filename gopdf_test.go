package gopdf

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkPdfWithImageHolder(b *testing.B) {

	err := initTesting()
	if err != nil {
		b.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	bytesOfImg, err := ioutil.ReadFile("./test/res/chilli.jpg")
	if err != nil {
		b.Error(err)
		return
	}

	imgH, err := ImageHolderByBytes(bytesOfImg)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		pdf.ImageByHolder(imgH, 20.0, float64(i)*2.0, nil)
	}

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_bench.pdf")
}

func initTesting() error {
	err := os.MkdirAll("./test/out", 0777)
	if err != nil {
		return err
	}
	return nil
}

func TestPdfWithImageHolder(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	bytesOfImg, err := ioutil.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.ImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.ImageByHolder(imgH, 20.0, 200, nil)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test.pdf")
}

func TestRetrievingNumberOfPdfPage(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}

	pdf.AddPage()
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	bytesOfImg, err := ioutil.ReadFile("./test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.ImageByHolder(imgH, 20.0, 20, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if pdf.GetNumberOfPages() != 1 {
		t.Error(err)
		return
	}

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/number_of_pages_test.pdf")
}

func TestImageCrop(t *testing.T) {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}

	pdf.AddPage()
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	bytesOfImg, err := ioutil.ReadFile("./test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		return
	}

	imgH, err := ImageHolderByBytes(bytesOfImg)
	if err != nil {
		t.Error(err)
		return
	}

	//err = pdf.ImageByHolder(imgH, 20.0, 20, nil)
	err = pdf.ImageByHolderWithOptions(imgH, ImageOptions{
		//VerticalFlip: true,
		//HorizontalFlip: true,
		Rect: &Rect{
			W: 100,
			H: 100,
		},
		Crop: &CropOptions{
			X:      0,
			Y:      0,
			Width:  10,
			Height: 100,
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	if pdf.GetNumberOfPages() != 1 {
		t.Error(err)
		return
	}

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetX(250)
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/image_crop.pdf")
}

/*
func TestBuffer(t *testing.T) {
	b := bytes.NewReader([]byte("ssssssss"))

	b1, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("->%s\n", string(b1))
	b.Seek(0, 0)
	b2, err := ioutil.ReadAll(b)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("+>%s\n", string(b2))
}*/

func BenchmarkAddTTFFontByReader(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := ioutil.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &GoPdf{}
		pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
		if err := pdf.AddTTFFontByReader("LiberationSerif-Regular", bytes.NewReader(fontData)); err != nil {
			return
		}
	}
}

func BenchmarkAddTTFFontData(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := ioutil.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &GoPdf{}
		pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
		if err := pdf.AddTTFFontData("LiberationSerif-Regular", fontData); err != nil {
			return
		}
	}
}

func TestReuseFontData(t *testing.T) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := ioutil.ReadAll(ttf)
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &GoPdf{}
	rst1, err := generatePDFBytesByAddTTFFontData(pdf1, fontData)
	if err != nil {
		t.Error(err)
		return
	}

	// Reuse the parsed font data.
	pdf2 := &GoPdf{}
	rst2, err := generatePDFBytesByAddTTFFontData(pdf2, fontData)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Error(errors.New("The generated files must be exactly the same."))
		return
	}

	if err := os.WriteFile("./test/out/result1_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := os.WriteFile("./test/out/result2_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func generatePDFBytesByAddTTFFontData(pdf *GoPdf, fontData []byte) ([]byte, error) {
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	if pdf.GetNumberOfPages() != 0 {
		return nil, errors.New("Invalid starting number of pages, should be 0")
	}

	if err := pdf.AddTTFFontData("LiberationSerif-Regular", fontData); err != nil {
		return nil, err
	}

	if err := pdf.SetFont("LiberationSerif-Regular", "", 14); err != nil {
		return nil, err
	}

	pdf.AddPage()
	if err := pdf.Text("Test PDF content."); err != nil {
		return nil, err
	}

	return pdf.GetBytesPdfReturnErr()
}

func TestWhiteTransparent(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	var glyphNotFoundOfLiberationSerif []rune
	err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "test/res/LiberationSerif-Regular.ttf", TtfOption{
		OnGlyphNotFound: func(r rune) { //call when can not find glyph inside ttf file.
			glyphNotFoundOfLiberationSerif = append(glyphNotFoundOfLiberationSerif, r)
			//log.Printf("glyph not found %c", r)
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	// write text.
	op := CellOption{Align: Left | Middle}
	rect := Rect{W: 20, H: 30}
	pdf.SetX(350)
	pdf.SetY(50)
	err = pdf.Cell(&rect, "あい")
	//err = pdf.CellWithOption(&rect, "あい", op)
	//err = pdf.CellWithOption(&rect, "あ", op)
	//err = pdf.CellWithOption(&rect, "a", op)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetY(100)
	err = pdf.CellWithOption(&rect, "abcdef.", op)
	if err != nil {
		t.Error(err)
		return
	}

	//coz あ and い  not contain in "test/res/LiberationSerif-Regular.ttf"
	if len(glyphNotFoundOfLiberationSerif) != 2 {
		t.Error(err)
		return
	}

	//pdf.SetNoCompression()
	err = pdf.WritePdf("./test/out/white_transparent.pdf")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestRectangle(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	pdf.SetStrokeColor(240, 98, 146)
	pdf.SetLineWidth(1)
	pdf.SetFillColor(255, 255, 255)
	// draw rectangle with round radius
	err = pdf.Rectangle(100.6, 150.8, 150.3, 379.3, "DF", 20, 10)
	if err != nil {
		t.Error(err)
		return
	}

	// draw rectangle with round radius but less point number
	err = pdf.Rectangle(200.6, 150.8, 250.3, 379.3, "DF", 20, 2)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetStrokeColor(240, 98, 146)
	pdf.SetLineWidth(1)
	pdf.SetFillColor(255, 255, 255)
	// draw rectangle directly
	err = pdf.Rectangle(100.6, 50.8, 130, 150, "DF", 0, 0)
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/rectangle_with_round_corner.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestWhiteTransparent195(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}
	// create pdf.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()

	var glyphNotFoundOfLiberationSerif []rune
	//err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "/Users/oneplus/Code/Work/gopdf_old/test/res/Meera-Regular.ttf", TtfOption{
	err = pdf.AddTTFFontWithOption("LiberationSerif-Regular", "test/res/LiberationSerif-Regular.ttf", TtfOption{
		OnGlyphNotFound: func(r rune) { //call when can not find glyph inside ttf file.
			glyphNotFoundOfLiberationSerif = append(glyphNotFoundOfLiberationSerif, r)
		},
		OnGlyphNotFoundSubstitute: func(r rune) rune {
			//return r
			return rune('\u20b0') //(U+25A1) = “□”
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	// write text.
	op := CellOption{Align: Left | Middle}
	rect := Rect{W: 20, H: 30}
	pdf.SetX(350)
	pdf.SetY(50)
	//err = pdf.Cell(&rect, "あいうえ") // OK.
	//err = pdf.Cell(&rect, "あうう") // OK.
	err = pdf.CellWithOption(&rect, "あいうえ", op) // NG. "abcdef." is White/Transparent.
	//err = pdf.Cell(&rect, " あいうえ") // NG. "abcdef." is White/Transparent.
	// err = pdf.Cell(&rect, "あいうえ ") // NG. "abcdef." is White/Transparent.
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetY(100)
	err = pdf.CellWithOption(&rect, "abcกdef.", op)
	if err != nil {
		t.Error(err)
		return
	}

	//coz あ い う え  not contain in "test/res/LiberationSerif-Regular.ttf"
	// if len(glyphNotFoundOfLiberationSerif) != 4 {
	// 	t.Error(err)
	// 	return
	// }

	pdf.SetNoCompression()
	err = pdf.WritePdf("./test/out/white_transparent195.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
