package pagination

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/signintech/gopdf"
)

func GetFont(pdf *gopdf.GoPdf, fontPath string) (err error) {
	b, err := os.Open(fontPath)
	if err != nil {
		return err
	}
	err = pdf.AddTTFFontByReader("Ubuntu-L", b)
	if err != nil {
		return err
	}
	return err
}

func TestSetY(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		// var textH float64 = 25 // if text height is 25px.
		pdf.SetXY(x, y)
		err = pdf.Text(text)
		if err != nil {
			log.Fatalln(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("page_sety-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestSetNewY(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		pdf.SetX(x)
		pdf.SetNewY(y, textH)
		y = pdf.GetY()
		err = pdf.Text(text)
		if err != nil {
			log.Fatalln(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("page_setnewy-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestSetNewXY(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		// pdf.SetX(x)
		pdf.SetNewXY(y, x, textH)
		y = pdf.GetY()
		err = pdf.Text(text)
		if err != nil {
			log.Fatalln(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("page_setnewxy-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestSetNewYX(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	var x float64 = 100
	var y float64 = 10
	for i := 0; i < 200; i++ {
		text := fmt.Sprintf("---------line no: %d -----------", i)
		var textH float64 = 25 // if text height is 25px.
		pdf.SetNewY(y, textH)
		y = pdf.GetY()
		pdf.SetX(x) // must after pdf.SetNewY() called.
		err = pdf.Text(text)
		if err != nil {
			log.Fatalln(err)
		}
		y += 20
	}

	err = pdf.WritePdf(fmt.Sprintf("page_setnewyx-%s.pdf", time.Now().Format("01-02-15-04-05")))
	if err != nil {
		log.Fatalln(err)
	}
}

func TestSetNewYCheckHeight(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	y := 10.0
	pdf.SetNewY(y, 0)
	if y != pdf.GetY() {
		log.Fatalln(" y != pdf.GetY()")
	}

	y = 1000.0
	pdf.SetNewY(y, 0)
	if y != pdf.GetY() {
		log.Fatalln(" y != pdf.GetY()")
	}
}
