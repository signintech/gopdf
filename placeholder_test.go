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
	pdf.Start(Config{PageSize: *PageSizeA4})
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

	err = pdf.FillInPlaceHoldText("totalnumber", fmt.Sprintf("%d", 5), Left) //<-- fillin text to PlaceHolder
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.WritePdf("./test/out/placeholder_text.pdf")
}

func TestPlaceHolderText2(t *testing.T) {
	err := initTesting()
	if err != nil {
		t.Error(err)
		return
	}

	pdf := GoPdf{}
	pdf.Start(Config{PageSize: *PageSizeA4})
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
		pdf.Text("page")
		pagenumberPH := fmt.Sprintf("pagenumber_%d", i)
		err = pdf.PlaceHolderText(pagenumberPH, 20) //<-- create PlaceHolder
		if err != nil {
			log.Print(err.Error())
			return
		}

		err := pdf.Text("of")
		if err != nil {
			log.Print(err.Error())
			return
		}
		err = pdf.PlaceHolderText("totalnumber", 20) //<-- create PlaceHolder
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

		err = pdf.FillInPlaceHoldText(pagenumberPH, fmt.Sprintf("%d", i+1), Center) //<-- fillin text to PlaceHolder
		if err != nil {
			log.Print(err.Error())
			return
		}

	}

	err = pdf.FillInPlaceHoldText("totalnumber", fmt.Sprintf("%d", 5), Center) //<-- fillin text to PlaceHolder
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.WritePdf("./test/out/placeholder_text_2.pdf")
}
