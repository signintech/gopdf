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
	pdf.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 })  
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
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4 })  
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

### Links

```go

package main

import (
	"log"
	"github.com/signintech/gopdf"
)

func main()  {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 }) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("times", "./test/res/times.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("times", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.SetX(30)
	pdf.SetY(40)
	pdf.Text("Link to example.com")
	pdf.AddExternalLink("http://example.com/", 27.5, 28, 125, 15)

	pdf.SetX(30)
	pdf.SetY(70)
	pdf.Text("Link to second page")
	pdf.AddInternalLink("anchor", 27.5, 58, 120, 15)

	pdf.AddPage()
	pdf.SetX(30)
	pdf.SetY(100)
	pdf.SetAnchor("anchor")
	pdf.Text("Anchor position")

	pdf.WritePdf("hello.tmp.pdf")

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

### Rotation text or image
```go
pdf.SetX(100)
pdf.SetY(100)
pdf.Rotate(270.0, 100.0, 100.0)
pdf.Text("Hello...")
pdf.RotateReset() //reset
```

### Set transparency for text
Read about [transparency in pdf](https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf) `(page 320, section 11)`
```go
pdf.SetAlpha(alpha float64, blendMode string)
```
- `alpha` - value of transparency, can be between `0` and `1` 
- `blendMode` - default value is `Normal` - read about [blendMode and kinds of its](https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf) `(page 325, section 11.3.5)`  

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
		PageSize: *gopdf.PageSizeA4, //595.28, 841.89 = A4
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
### Import existing PDF 
Import existing PDF power by package [gofpdi](https://github.com/phpdave11/gofpdi) created by @phpdave11 (thank you :smile:)
```go
package main

import (
        "github.com/signintech/gopdf"
        "io"
        "net/http"
        "os"
)

func main() {
        var err error

        // Download a Font
        fontUrl := "https://github.com/google/fonts/raw/master/ofl/daysone/DaysOne-Regular.ttf"
        if err = DownloadFile("example-font.ttf", fontUrl); err != nil {
            panic(err)
        }

        // Download a PDF
        fileUrl := "https://tcpdf.org/files/examples/example_012.pdf"
        if err = DownloadFile("example-pdf.pdf", fileUrl); err != nil {
            panic(err)
        }

        pdf := gopdf.GoPdf{}
        pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4

        pdf.AddPage()

        err = pdf.AddTTFFont("daysone", "example-font.ttf")
        if err != nil {
            panic(err)
        }

        err = pdf.SetFont("daysone", "", 20)
        if err != nil {
            panic(err)
        }

        // Color the page
        pdf.SetLineWidth(0.1)
        pdf.SetFillColor(124, 252, 0) //setup fill color
        pdf.RectFromUpperLeftWithStyle(50, 100, 400, 600, "FD")
        pdf.SetFillColor(0, 0, 0)

        pdf.SetX(50)
        pdf.SetY(50)
        pdf.Cell(nil, "Import existing PDF into GoPDF Document")

        // Import page 1
        tpl1 := pdf.ImportPage("example-pdf.pdf", 1, "/MediaBox")

        // Draw pdf onto page
        pdf.UseImportedTemplate(tpl1, 50, 100, 400, 0)

        pdf.WritePdf("example.pdf")

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
        // Get the data
        resp, err := http.Get(url)
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        // Create the file
        out, err := os.Create(filepath)
        if err != nil {
            return err
        }
        defer out.Close()

        // Write the body to file
        _, err = io.Copy(out, resp.Body)
        return err
}

```

  
visit https://github.com/oneplus1000/gopdfsample for more samples.
