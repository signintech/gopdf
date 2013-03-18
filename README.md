gopdf
=====

A simple library for generating PDF written in Go lang.

Use [fpdfGo](https://github.com/signintech/fpdfGo) to generate fonts.



Sample
======

	import (
		"fmt"
		"gopdf"
		 iconv "github.com/djimenez/iconv-go"
		 "gopdf/fonts"
	)

	func main() {

		pdf := gopdf.GoPdf{}
		pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
		pdf.AddFont("THSarabunPSK",new(fonts.THSarabun),"res/fonts/THSarabun.z")
		pdf.AddFont("Loma",new(fonts.Loma),"res/fonts/Loma.z")
		pdf.AddPage()
		pdf.SetFont("THSarabunPSK", "B", 14)
		pdf.Cell(nil,  ToCp874("ทดสอบ"))
		pdf.Cell(nil,  ToCp874("Test"))
		pdf.Br(28)
		pdf.WritePdf("/var/www/fpdf17/output/x.pdf")
		fmt.Println("Done...")
	}

	func ToCp874(str string) string{
		str, _ = iconv.ConvertString( str, "utf-8", "cp874") 
		return  str
	}
