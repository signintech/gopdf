gopdf
====

gopdf is a simple library for generating PDF document written in Go lang.


####Changelogs

**2015-08-07**

- Add support for Unicode subfont embedding. (Chinese, Korean and Japanese fonts are now supported.)
- No longer need to create font maps.


##Installation
 ```
 go get -u github.com/signintech/gopdf
 ```

##Sample code

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
    err := pdf.AddTTFFont("HDZB_5", "../ttf/wts11.ttf")
    if err != nil {
        log.Print(err.Error())
        return
    }
    
    err = pdf.SetFont("HDZB_5", "", 14)
    if err != nil {
        log.Print(err.Error())
        return
    }
    pdf.Cell(nil, "您好")
    pdf.WritePdf("hello.pdf")

  }

  ```
  
visit https://github.com/oneplus1000/gopdfsample for more samples.
