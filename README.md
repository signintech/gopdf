gopdf
=====

A simple library for generating PDF written in Go lang.

<strike>Use [fpdfGo](https://github.com/signintech/fpdfGo) to generate fonts.</strike><br />
Use fontmaker to generate fonts.<br />
Sample code [here](https://github.com/oneplus1000/gopdfusecase) 

####Installation
 ```
 go get github.com/signintech/gopdf
 ```

####Example
  ```go
  package main
  import (
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fonts"
  )

  func main() {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //A4
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
gopdf / fontmaker
======
fontmaker is a font making tool for gopdf.

####Build fontmaker
######open terminal 
  ```
  $ cd {GOPATH}/src/github.com/signintech/gopdf/fontmaker

  $ go build
  ```
####Usage:
  ```
  fontmaker encoding  map_folder  font_file  output_folder
  ```
####Example:
######run command
  ```
  fontmaker  cp874  /gopath/github.com/signintech/gopdf/fontmaker/map   ../ttf/Loma.ttf  ./tmp 
  ```
######result  
  ```
  Save Z file at ./tmp/Loma.z.
  Save GO file at ./tmp/Loma.font.go.
  Finish.
  ```

