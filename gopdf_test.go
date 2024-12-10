package gopdf

import (
	"bytes"
	"errors"
	"io"
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
	pdf.Start(Config{PageSize: *PageSizeA4})
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

	bytesOfImg, err := os.ReadFile("./test/res/chilli.jpg")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_bench.pdf")
}

func TestPdfWithImageHolder(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test.pdf")
}

func TestPdfWithImageHolderGif(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher03.gif")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.WritePdf("./test/out/image_test_gif.pdf")
}

func TestRetrievingNumberOfPdfPage(t *testing.T) {
	pdf := setupDefaultA4PDF(t)
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}
	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher01.jpg")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/number_of_pages_test.pdf")
}

func TestImageCrop(t *testing.T) {
	pdf := setupDefaultA4PDF(t)
	if pdf.GetNumberOfPages() != 0 {
		t.Error("Invalid starting number of pages, should be 0")
		return
	}

	pdf.AddPage()

	bytesOfImg, err := os.ReadFile("./test/res/gopher01.jpg")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")

	pdf.AddPage()

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher again")

	if pdf.GetNumberOfPages() != 2 {
		t.Error(err)
		return
	}

	pdf.WritePdf("./test/out/image_crop.pdf")
}

func BenchmarkAddTTFFontByReader(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &GoPdf{}
		pdf.Start(Config{PageSize: *PageSizeA4})
		if err := pdf.AddTTFFontByReader("LiberationSerif-Regular", bytes.NewReader(fontData)); err != nil {
			return
		}
	}
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

func BenchmarkAddTTFFontData(b *testing.B) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		b.Error(err)
		return
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		b.Error(err)
		return
	}

	for n := 0; n < b.N; n++ {
		pdf := &GoPdf{}
		pdf.Start(Config{PageSize: *PageSizeA4})
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

	fontData, err := io.ReadAll(ttf)
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

	if err := writeFile("./test/out/result1_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_by_parsed_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func writeFile(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err1 := f.Close(); err1 != nil && err == nil {
		err = err1
	}
	return err
}

func generatePDFBytesByAddTTFFontData(pdf *GoPdf, fontData []byte) ([]byte, error) {
	pdf.Start(Config{PageSize: *PageSizeA4})
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
	pdf.SetXY(350, 50)
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
	pdf.SetXY(350, 50)
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

func TestClearValue(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4, Protection: PDFProtectionConfig{
		UseProtection: true,
		OwnerPass:     []byte("123456"),
		UserPass:      []byte("123456"),
	}})
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

	bytesOfImg, err := os.ReadFile("./test/res/PNG_transparency_demonstration_1.png")
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

	pdf.SetXY(250, 200)
	pdf.Cell(nil, "gopher and gopher")
	pdf.SetInfo(PdfInfo{
		Title: "xx",
	})
	pdf.WritePdf("./test/out/test_clear_value.pdf")

	//reset
	pdf.Start(Config{PageSize: *PageSizeA4})

	pdf2 := GoPdf{}
	pdf2.Start(Config{PageSize: *PageSizeA4})

	//check
	if pdf.margins != pdf2.margins {
		t.Fatal("pdf.margins != pdf2.margins")
	}

	if len(pdf2.pdfObjs) != len(pdf.pdfObjs) {
		t.Fatalf("len(pdf2.pdfObjs) != len(pdf.pdfObjs)")
	}

	if len(pdf.anchors) > 0 {
		t.Fatalf("len( pdf.anchors) = %d", len(pdf.anchors))
	}

	if len(pdf.indexEncodingObjFonts) != len(pdf2.indexEncodingObjFonts) {
		t.Fatalf("len(pdf.indexEncodingObjFonts) != len(pdf2.indexEncodingObjFonts)")
	}

	if pdf.indexOfContent != pdf2.indexOfContent {
		t.Fatalf("pdf.indexOfContent != pdf2.indexOfContent")
	}

	if pdf.buf.Len() > 0 {
		t.Fatalf("pdf.buf.Len() > 0")
	}

	if pdf.pdfProtection != nil {
		t.Fatalf("pdf.pdfProtection is not nil")
	}
	if pdf.encryptionObjID != 0 {
		t.Fatalf("encryptionObjID %d", pdf.encryptionObjID)
	}

	if pdf.info != nil {
		t.Fatalf("pdf.info %v", pdf.info)
	}
}

func TestSplitTextWithOptions(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := setupDefaultA4PDF(t)

	var splitTextTests = []struct {
		name string
		in   string
		opts *BreakOption
		exp  []string
	}{
		{
			"strict breaks no separator",
			"Lorem ipsum dolor sit amet, consetetur",
			&DefaultBreakOption,
			[]string{"Lorem ipsum dol", "or sit amet, conse", "tetur"},
		},
		{
			"no options given",
			"Lorem ipsum dolor sit amet, consetetur",
			nil,
			[]string{"Lorem ipsum dol", "or sit amet, conse", "tetur"},
		},
		{
			"strict breaks with separator",
			"Lorem ipsum dolor sit amet, consetetur",
			&BreakOption{
				Separator: "-",
				Mode:      BreakModeStrict,
			},
			[]string{"Lorem ipsum d-", "olor sit amet, c-", "onsetetur"},
		},
		{
			"text with possible word-wrap",
			"Lorem ipsum dolor sit amet, consetetur",
			&BreakOption{
				BreakIndicator: ' ',
				Mode:           BreakModeIndicatorSensitive,
			},
			[]string{"Lorem ipsum", "dolor sit amet,", "consetetur"},
		},
		{
			"text without possible word-wrap",
			"Loremipsumdolorsitamet,consetetur",
			&BreakOption{
				BreakIndicator: ' ',
				Mode:           BreakModeIndicatorSensitive,
			},
			[]string{"Loremipsumdolo", "rsitamet,consetet", "ur"},
		},
		{
			"text with only empty spaces",
			"                                                ",
			&BreakOption{
				BreakIndicator: ' ',
				Mode:           BreakModeIndicatorSensitive,
			},
			[]string{"                           ", "                    "},
		},
	}

	for _, tt := range splitTextTests {
		t.Run(tt.name, func(t *testing.T) {
			lines, err := pdf.SplitTextWithOption(tt.in, 100, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if len(lines) != len(tt.exp) {
				t.Fatalf("amount of expected and split lines invalid. Expected: %d, result: %d", len(tt.exp), len(lines))
			}
			for i, e := range tt.exp {
				if e != lines[i] {
					t.Fatalf("split text invalid. Expected: '%s', result: '%s'", e, lines[i])
				}
			}
		})
	}
}

func TestTextColor(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	// create pdf.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	pdf.AddPage()
	err = pdf.AddTTFFont("LiberationSerif", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("LiberationSerif", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetTextColor(255, 0, 2)
	pdf.Br(20)
	pdf.Cell(nil, "a")

	pdf.SetTextColorCMYK(0, 6, 14, 0)
	pdf.Br(20)
	pdf.Cell(nil, "b")

	err = pdf.WritePdf("./test/out/colored_text.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAddHeaderFooter(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	// create pdf.
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})

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

	pdf.AddHeader(func() {
		pdf.SetY(5)
		pdf.Cell(nil, "header")
	})
	pdf.AddFooter(func() {
		pdf.SetY(825)
		pdf.Cell(nil, "footer")
	})

	pdf.AddPage()
	pdf.SetY(400)
	pdf.Text("page 1 content")
	pdf.AddPage()
	pdf.SetY(400)
	pdf.Text("page 2 content")

	err = pdf.WritePdf("./test/out/header_footer.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}

func initTesting() error {
	err := os.MkdirAll("./test/out", 0777)
	if err != nil {
		return err
	}
	return nil
}

// setupDefaultA4PDF creates an A4 sized pdf with a plain configuration adding and setting the required fonts for
// further processing. Tests will fail in case adding or setting the font fails.
func setupDefaultA4PDF(t *testing.T) *GoPdf {
	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err := pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Fatal(err)
	}

	err = pdf.SetFont("LiberationSerif-Regular", "", 14)
	if err != nil {
		log.Fatal(err)
	}
	return &pdf
}

func TestImportPagesFromFile(t *testing.T) {
	err := os.MkdirAll("./test/out", 0777)
	if err != nil {
		t.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})

	err = pdf.ImportPagesFromSource("./examples/outline_example/outline_demo.pdf", "/MediaBox")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetFont("LiberationSerif-Regular", "", 14)

	err = pdf.SetPage(1)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetPage(2)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}
	err = pdf.SetPage(3)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.SetXY(350, 50)
	err = pdf.Cell(&Rect{W: 20, H: 30}, "Hello World")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.WritePdf("./test/out/open-existing-pdf.pdf")
	if err != nil {
		t.Error(err)
		return
	}
}
