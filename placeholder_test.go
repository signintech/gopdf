package gopdf

import (
	"fmt"
	"log"
	"testing"
)

func TestPlaceHolderText(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	err = pdf.AddTTFFont("LiberationSerif-Regular", "./test/res/LiberationSerif-Regular.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 5; i++ {
		pdf.AddPage()
		err = pdf.SetFont("LiberationSerif-Regular", "", 14)
		if err != nil {
			log.Print(err.Error())
			return
		}
		pdf.Br(10)
		pdf.SetX(250)
		err := pdf.Text(fmt.Sprintf("%d of ", i+1))
		if err != nil {
			log.Print(err.Error())
			return
		}
		err = pdf.PlaceHolderText("totalnumber", 30) //<-- create PlaceHolder
		if err != nil {
			log.Print(err.Error())
			return
		}
		pdf.Br(20)

		err = pdf.SetFont("LiberationSerif-Regular", "", 11)
		if err != nil {
			log.Print(err.Error())
			return
		}
		pdf.Text("content content content content content contents...")
	}

	err = pdf.FillInPlaceHoldText("totalnumber", fmt.Sprintf("%d", 5)) //<-- fillin text to PlaceHolder
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.WritePdf("./test/out/placeholder_text.pdf")
}