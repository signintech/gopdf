package main

import (
	"log"

	"github.com/signintech/gopdf"
)

func main() {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Amiri-Regular", "./examples/arabic/Amiri-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.SetFont("Amiri-Regular", "", 14)
	if err != nil {
		log.Fatal(err)
	}
	pdf.Cell(nil, gopdf.ToArabic("ولكن لا السلام عليكم ورحمة الله وبركاته"))

	pdf.WritePdf("arabic.pdf")
}
