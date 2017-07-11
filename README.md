gopdf
====

gopdf is a simple library for generating PDF document written in Go lang.


#### Features

- Unicode subfont embedding. (Chinese, Japanese, Korean, etc.)
- Draw line, oval, rect, curve
- Draw image ( jpg, png )
- Password protection
- Font [kerning](https://en.wikipedia.org/wiki/Kerning)


## Installation
 ```
 go get -u github.com/signintech/gopdf
 ```


### Print text 

```go
  
package main
import (
	"log"
	"github.com/signintech/gopdf"
)

func main() {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("wts11", "../ttf/wts11.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("wts11", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	pdf.Cell(nil, "您好")
	pdf.WritePdf("hello.pdf")

}

```
  
### Image
  
```go

package main
import (
	"log"
	"github.com/signintech/gopdf"
)

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
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

### Draw line
```go
pdf.SetLineWidth(2)
pdf.SetLineType("dashed")
pdf.Line(10, 30, 585, 30)
```

### Draw oval
```go
pdf.SetLineWidth(1)
pdf.Oval(100, 200, 500, 500)
```

### Password protection
```go
package main

import (
	"log"

	"github.com/signintech/gopdf"
)


func main() {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 595.28, H: 841.89}, //595.28, 841.89 = A4
		Protection: gopdf.PDFProtectionConfig{
			UseProtection: true,
			Permissions: gopdf.PermissionsPrint | gopdf.PermissionsCopy | gopdf.PermissionsModify,
			OwnerPass:   []byte("123456"),
			UserPass:    []byte("123456789")},
	})

	pdf.AddPage()
	pdf.AddTTFFont("loma", "../ttf/loma.ttf")
	pdf.Cell(nil,"Hi")
	pdf.WritePdf("protect.pdf")
}

```
  
visit https://github.com/oneplus1000/gopdfsample for more samples.
