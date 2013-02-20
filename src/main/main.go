package main

import (
	"fmt"
	"gopdf"
	 iconv "github.com/djimenez/iconv-go"
	 "gopdf/fonts"
)

func main() {
	fmt.Println("start...")
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 A4
	pdf.AddFont("THSarabunPSK",new(fonts.THSarabun),"res/fonts/THSarabun.z")
	pdf.AddFont("Loma",new(fonts.Loma),"res/fonts/Loma.z")
	pdf.AddPage()
	pdf.SetFont("THSarabunPSK", "B", 12)
	output , _ := iconv.ConvertString( "กAโจตลาด  2 ล้อ พุ่ง 20% รับปีใหม่ คาดเอที toyota ยังแรงกุ้งตั้ว ", "utf-8", "cp874") 
	pdf.Cell(gopdf.Rect{H: 100, W: 100},  output)
	pdf.SetFont("Loma", "B", 12)
	output , _ = iconv.ConvertString( "การบ้านx", "utf-8", "cp874") 
	pdf.Cell(gopdf.Rect{H: 100, W: 100}, output)
	//pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xzxzxzxzx")
	//pdf.AddPage()
	//pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xxxx")
	pdf.WritePdf("/home/oneplus/pdf/x.pdf")
	fmt.Println("end...")
}
