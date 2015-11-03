package gopdftest

import (
	"testing"

	"github.com/signintech/gopdf"
)

func TestHello(t *testing.T) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("HDZB_5", "./ttf/wts11.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.SetFont("HDZB_5", "", 14)
	if err != nil {
		t.Error(err)
		return
	}

	pdf.SetGrayFill(0.5)
	pdf.Cell(nil, "您好")
	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Error(err)
		return
	}

	if len(data) <= 0 {
		t.Error(err)
		return
	}

}

func TestHello2(t *testing.T) {

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	var err error
	err = pdf.AddTTFFont("HDZB_5", "./ttf/wts11.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.AddTTFFont("TakaoPGothic", "./ttf/TakaoPGothic.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.AddTTFFont("loma", "./ttf/Loma.ttf")
	if err != nil {
		t.Error(err)
		return
	}

	err = pdf.AddTTFFont("namum", "./ttf/NanumBarunGothic.ttf")

	//china
	err = pdf.SetFont("HDZB_5", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.Cell(nil, "Hello")
	pdf.Br(20)
	pdf.Cell(nil, "您好")
	pdf.Br(20)

	//japan
	err = pdf.SetFont("TakaoPGothic", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.Cell(nil, "こんにちは")
	pdf.Br(20)

	//thai
	err = pdf.SetFont("loma", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.Cell(nil, "สวัสดี")
	pdf.Br(20)

	//korean
	err = pdf.SetFont("namum", "", 14)
	if err != nil {
		t.Error(err)
		return
	}
	pdf.Cell(nil, "안녕하세요")

	data, err := pdf.GetBytesPdfReturnErr()
	if err != nil {
		t.Error(err)
		return
	}

	if len(data) <= 0 {
		t.Error(err)
		return
	}
}

/*
func BenchmarkHello(b *testing.B) {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()
	err := pdf.AddTTFFont("HDZB_5", "./ttf/wts11.ttf")
	if err != nil {
		b.Error(err)
		return
	}

	err = pdf.SetFont("HDZB_5", "", 14)
	if err != nil {
		b.Error(err)
		return
	}

	pdf.SetGrayFill(0.5)
	pdf.Cell(nil, "您好")
	_, err = pdf.GetBytesPdfReturnErr()
	if err != nil {
		b.Error(err)
		return
	}
}*/
