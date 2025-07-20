package gopdf_test

import (
	"github.com/signintech/gopdf"
	"io"
	"os"
)

func ExampleFontContainer_AddTTFFont() {
	fontContainer := &gopdf.FontContainer{}
	err := fontContainer.AddTTFFont("LiberationSerif-Regular", "path/to/LiberationSerif-Regular.ttf")
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleFontContainer_AddTTFFontWithOption() {
	fontContainer := &gopdf.FontContainer{}
	err := fontContainer.AddTTFFontWithOption(
		"LiberationSerif-Regular",
		"path/to/LiberationSerif-Regular.ttf",
		gopdf.TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
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

	fontContainer := &gopdf.FontContainer{}
	err = fontContainer.AddTTFFontByReader("LiberationSerif-Regular", ttf)
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
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

	fontContainer := &gopdf.FontContainer{}
	err = fontContainer.AddTTFFontByReaderWithOption("LiberationSerif-Regular", ttf, gopdf.TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
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

	fontContainer := &gopdf.FontContainer{}
	err = fontContainer.AddTTFFontData("LiberationSerif-Regular", fontData)
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
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

	fontContainer := &gopdf.FontContainer{}
	err = fontContainer.AddTTFFontDataWithOption("LiberationSerif-Regular", fontData, gopdf.TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}

func ExampleGoPdf_AddTTFFontFromFontContainer() {
	fontContainer := &gopdf.FontContainer{}
	err := fontContainer.AddTTFFontWithOption(
		"LiberationSerif-Regular",
		"path/to/LiberationSerif-Regular.ttf",
		gopdf.TtfOption{})
	if err != nil {
		// handle error
	}
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = pdf.AddTTFFontFromFontContainer("LiberationSerif-Regular", fontContainer)
	if err != nil {
		// handle error
	}
}
