gopdf
=====

A simple library for generating PDF written in Go lang.

<strike>Use [fpdfGo](https://github.com/signintech/fpdfGo) to generate fonts.</strike><br />
<strike>Use fontmaker to generate fonts.</strike><br />
Sample code [here](https://github.com/oneplus1000/gopdfsample)

####Installation
 ```
 go get github.com/signintech/gopdf
 ```

####Example
  ```go
  package main
  import (
	"fmt"
	"github.com/signintech/gopdf"
  )

  func main() {

    pdf := gopdf.GoPdf{}
    pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
    pdf.AddPage()
    var err error
    err = pdf.AddTTFFont("HDZB_5", "../ttf/wts11.ttf")
    if err != nil {
        log.Printf("%s", err.Error())
        return
    }
    err = pdf.SetFont("HDZB_5", "", 14)
    if err != nil {
        log.Printf("ERROR:%s\n", err.Error())
        return
    }
    pdf.Cell(nil, "您好")
    pdf.WritePdf("hello.pdf")

  }


  ```
