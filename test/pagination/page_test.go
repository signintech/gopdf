package pagination

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-text/render"
	"github.com/go-text/typesetting/font"
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
		//var textH float64 = 25 // if text height is 25px.
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
		//pdf.SetX(x)
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

func TestLineBreak(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "../res/LiberationSerif-Regular.ttf")
	if err != nil {
		log.Fatalln(err)
	}
	err = pdf.SetFont("Ubuntu-L", "", 28)
	if err != nil {
		log.Fatalln(err)
	}
	pdf.SetMargins(0, 20, 0, 10)
	pdf.AddPage()

	w := 500.0

	var breakOptionTests = []*gopdf.BreakOption{
		&gopdf.DefaultBreakOption,
		{
			Mode:           gopdf.BreakModeIndicatorSensitive,
			BreakIndicator: ' ',
		},
	}

	y := (gopdf.PageSizeA4.H/2 + 100.0*float64(len(breakOptionTests))) / 2
	linebreakText := strings.Repeat("MultiCell* methods don't respect linebreaking rules.", 2)
	for i, opt := range breakOptionTests {
		pdf.SetXY(gopdf.PageSizeA4.W/2-w/2, y+100.0*float64(i))
		err = pdf.MultiCellWithOption(&gopdf.Rect{
			W: w,
			H: 1000,
		}, linebreakText, gopdf.CellOption{
			BreakOption: opt,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = pdf.WritePdf("page_linebreak.pdf")
	if err != nil {
		log.Fatalln(err)
	}
}
	
func TestHindiRendering(t *testing.T) {
	var err error
	pdf := &gopdf.GoPdf{}
	pageSize := *gopdf.PageSizeA4
	pdf.Start(gopdf.Config{PageSize: pageSize})
	pdf.AddPage()

	err = pdf.AddTTFFont("mangal", "../res/MangalRegular.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("mangal", "", 150)
	if err != nil {
		log.Print(err.Error())
		return
	}

	r := &render.Renderer{
		FontSize: 30,
		Color:    color.Black,
	}

	text := "नमस्ते"

	tw, err := pdf.MeasureTextWidth(text)
	if err != nil {
		log.Print(err.Error())
		return
	}

	th, err := pdf.MeasureCellHeightByText(text)
	if err != nil {
		log.Print(err.Error())
		return
	}

	img := image.NewNRGBA(image.Rect(0, 0, int(tw), int(th)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	data, _ := os.Open("../res/MangalRegular.ttf")
	face, _ := font.ParseTTF(data)

	r.FontSize = 150
	r.DrawString(text, img, face)

	mid := pageSize.H/2 - th/2

	var x float64 = 100
	var y = mid - 100.0
	imgRect := &gopdf.Rect{
		W: tw,
		H: th,
	}
	err = pdf.ImageFrom(img, x, y, imgRect)
	if err != nil {
		log.Fatal(err)
	}

	pdf.SetX(pageSize.W/2 - tw/2)
	pdf.SetY(mid + 100.0)
	pdf.Cell(nil, text)
	pdf.WritePdf("page_hindi_rendering.pdf")
}
