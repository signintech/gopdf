package gopdf_test

import (
	"fmt"
	"gopdf"
	 iconv "github.com/djimenez/iconv-go"
	"testing"
)

func TestNormol01(t *testing.T){

	fmt.Println("start...")
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddFont("THSarabunPSK",new(gopdf.THSarabun),"res/fonts/THSarabun.z")
	pdf.AddFont("Loma",new(gopdf.Loma),"res/fonts/Loma.z")
	pdf.AddPage()
	pdf.SetFont("THSarabunPSK", "B", 14)
	pdf.Cell(nil,  ToCp874("ทดสอบ"))
	pdf.Cell(nil,  ToCp874("Test"))
	pdf.Br(28)
	
	/*pdf.SetFont("THSarabunPSK", "B", 26)
	pdf.Image("/data/CODES/WORKS/GoPdf_ws/GoPdf/test/me.jpg",20,10,nil)
	pdf.Cell(nil,  ToCp874("ด้วยการที่เราไม่ไปทำยังไง"))
	pdf.Cell(nil,  ToCp874("ศาลอาญา ไม่ชี้ชัด ใครยิง “ มานะ อาจราญ” เจ้าหน้าที่สวนสัตว์ดุสิตเสียชีวิต"))
	
	
	pdf.AddPage()
	pdf.SetFont("Loma", "B", 12)
	pdf.Cell(&gopdf.Rect{H: 100, W: 100}, ToCp874("การบ้านx"))
	pdf.SetFont("THSarabunPSK", "B", 12)
	pdf.Cell(&gopdf.Rect{H: 100, W: 100}, ToCp874("การบ้านx") )
	pdf.Br(14)
	pdf.Cell(&gopdf.Rect{H: 100, W: 100}, ToCp874("การบ้านx"))
	pdf.Cell(&gopdf.Rect{H: 100, W: 100}, ToCp874("การบ้านx") )
	pdf.SetLineWidth(1)
	pdf.Line(10,10,500,10)*/

	pdf.WritePdf("/var/www/fpdf17/output/x.pdf")
	fmt.Println("end...")
}

func ToCp874(str string) string{
	str, _ = iconv.ConvertString( str, "utf-8", "cp874") 
	return  str
}