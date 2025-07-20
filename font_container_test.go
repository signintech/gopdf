package gopdf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func BenchmarkAddTTFFontFromFontContainer(b *testing.B) {
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
	container := &FontContainer{}
	container.AddTTFFontData("LiberationSerif-Regular", fontData)

	for n := 0; n < b.N; n++ {
		pdf := &GoPdf{}
		pdf.Start(Config{PageSize: *PageSizeA4})
		if err := pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", container); err != nil {
			return
		}
	}
}

func TestFontContainer_ReuseFontData(t *testing.T) {
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

	container := &FontContainer{}
	err = container.AddTTFFontData("LiberationSerif-Regular", fontData)
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &GoPdf{}
	rst1, err := generatePDFBytesByAddTTFFromFontContainer(pdf1, container)
	if err != nil {
		t.Error(err)
		return
	}

	// Reuse the parsed font data.
	pdf2 := &GoPdf{}
	rst2, err := generatePDFBytesByAddTTFFromFontContainer(pdf2, container)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Error(errors.New("the generated files must be exactly the same"))
		return
	}

	if err := writeFile("./test/out/result1_by_parsed_ttf_font_with_font_container.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_by_parsed_ttf_font_with_font_container.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func TestFontContainer_CompareAddTTFFontData(t *testing.T) {
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

	container := &FontContainer{}
	err = container.AddTTFFontData("LiberationSerif-Regular", fontData)
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &GoPdf{}
	rst1, err := generatePDFBytesByAddTTFFromFontContainer(pdf1, container)
	if err != nil {
		t.Error(err)
		return
	}

	pdf2 := &GoPdf{}
	rst2, err := generatePDFBytesByAddTTFFontData(pdf2, fontData)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Error(errors.New("the generated files must be exactly the same"))
		return
	}

	if err := writeFile("./test/out/result1_font_container_compare_add_ttf_font_data.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_font_container_compare_add_ttf_font_data.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func TestFontContainer_CompareAddTTFFont(t *testing.T) {
	container := &FontContainer{}
	err := container.AddTTFFont("LiberationSerif-Regular", "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &GoPdf{}
	rst1, err := generatePDFBytesByAddTTFFromFontContainer(pdf1, container)
	if err != nil {
		t.Error(err)
		return
	}

	pdf2 := &GoPdf{}
	rst2, err := generatePDFBytesByAddTTFFont(pdf2, "test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Error(errors.New("the generated files must be exactly the same"))
		return
	}

	if err := writeFile("./test/out/result1_font_container_compare_add_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_font_container_compare_add_ttf_font.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func TestFontContainer_CompareAddTTFFontByReader(t *testing.T) {
	ttf, err := os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	defer ttf.Close()

	container := &FontContainer{}
	err = container.AddTTFFontByReader("LiberationSerif-Regular", ttf)
	if err != nil {
		t.Error(err)
		return
	}

	pdf1 := &GoPdf{}
	rst1, err := generatePDFBytesByAddTTFFromFontContainer(pdf1, container)
	if err != nil {
		t.Error(err)
		return
	}

	ttf, err = os.Open("test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}
	defer ttf.Close()

	pdf2 := &GoPdf{}
	rst2, err := generatePDFBytesByAddTTFFontByReader(pdf2, ttf)
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare(rst1, rst2) != 0 {
		t.Error(errors.New("the generated files must be exactly the same"))
		return
	}

	if err := writeFile("./test/out/result1_font_container_compare_add_ttf_font_by_reader.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
	if err := writeFile("./test/out/result2_font_container_compare_add_ttf_font_by_reader.pdf", rst1, 0644); err != nil {
		t.Error(err)
		return
	}
}

func generatePDFBytesByAddTTFFromFontContainer(pdf *GoPdf, container *FontContainer) ([]byte, error) {
	pdf.Start(Config{PageSize: *PageSizeA4})
	if pdf.GetNumberOfPages() != 0 {
		return nil, errors.New("invalid starting number of pages, should be 0")
	}

	if err := pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", container); err != nil {
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

func generatePDFBytesByAddTTFFont(pdf *GoPdf, ttfpath string) ([]byte, error) {
	pdf.Start(Config{PageSize: *PageSizeA4})
	if pdf.GetNumberOfPages() != 0 {
		return nil, errors.New("invalid starting number of pages, should be 0")
	}

	if err := pdf.AddTTFFont("LiberationSerif-Regular", ttfpath); err != nil {
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

func generatePDFBytesByAddTTFFontByReader(pdf *GoPdf, rd io.Reader) ([]byte, error) {
	pdf.Start(Config{PageSize: *PageSizeA4})
	if pdf.GetNumberOfPages() != 0 {
		return nil, errors.New("invalid starting number of pages, should be 0")
	}

	if err := pdf.AddTTFFontByReader("LiberationSerif-Regular", rd); err != nil {
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

func ExampleFontContainer_AddTTFFont() {
	fontContainer := &FontContainer{}
	err := fontContainer.AddTTFFont("LiberationSerif-Regular", "path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontWithOption() {
	fontContainer := &FontContainer{}
	err := fontContainer.AddTTFFontWithOption(
		"LiberationSerif-Regular",
		"path/to/LiberationSerif-Regular.ttf",
		TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontByReader() {
	ttf, err := os.Open("path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	defer ttf.Close()

	fontContainer := &FontContainer{}
	err = fontContainer.AddTTFFontByReader("LiberationSerif-Regular", ttf)
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontByReaderWithOption() {
	ttf, err := os.Open("path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	defer ttf.Close()

	fontContainer := &FontContainer{}
	err = fontContainer.AddTTFFontByReaderWithOption("LiberationSerif-Regular", ttf, TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontData() {
	ttf, err := os.Open("path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		// handle error
	}

	fontContainer := &FontContainer{}
	err = fontContainer.AddTTFFontData("LiberationSerif-Regular", fontData)
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontDataWithOption() {
	ttf, err := os.Open("path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	defer ttf.Close()

	fontData, err := io.ReadAll(ttf)
	if err != nil {
		// handle error
	}

	fontContainer := &FontContainer{}
	err = fontContainer.AddTTFFontDataWithOption("LiberationSerif-Regular", fontData, TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}
