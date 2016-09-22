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


##Sample code : Print text 

  ```go
  
  package main
  import (
	"log"
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
  
##Sample code : Image
  
```go

package main
import (
	"log"
	"github.com/signintech/gopdf"
)

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	var err error
	err = pdf.AddTTFFont("loma", "../ttf/Loma.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}
	
	pdf.Image("../imgs/gopher.jpg", 200, 50, nil) //print image
	err = pdf.SetFont("loma", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.SetX(250) //move current location
	pdf.SetY(200)
	pdf.Cell(nil, "gopher and gopher") //print text

	pdf.WritePdf("image.pdf")
}
  
```
  
visit https://github.com/oneplus1000/gopdfsample for more samples.
