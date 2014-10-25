gopdf
=====

A simple library for generating PDF written in Go lang.

Use [fpdfGo](https://github.com/signintech/fpdfGo) to generate fonts.<br />
Sample code [here](https://github.com/oneplus1000/gopdfusecase) 


Sample
======

	```go
	import (
		"fmt"
		 iconv "github.com/djimenez/iconv-go"
		 "github.com/signintech/gopdf"
	 "github.com/signintech/gopdf/fonts"
	)
	
	func main() {
	
		pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddFont("THSarabunPSK",new(fonts.THSarabun),"THSarabun.z")
	pdf.AddFont("Loma",new(fonts.Loma),"Loma.z")
	pdf.AddPage()
	pdf.SetFont("THSarabunPSK","B",14)
	pdf.Cell(nil,   ToCp874("Hello world  = สวัสดี โลก in thai"))
		pdf.WritePdf("x.pdf")
		fmt.Println("Done...")
	}
	
	func ToCp874(str string) string{
		str, _ = iconv.ConvertString( str, "utf-8", "cp874") 
		return  str
	}
	```
	
fontmaker
======

###Build fontmaker

$ cd {GOPATH}/src/github.com/signintech/gopdf/fontmaker/fontmaker

$ go build

###Usage:
	fontmaker encoding map_folder font_file output_folder

###Example:
	fontmaker cp874 ../map  ../ttf/Loma.ttf ./tmp 


