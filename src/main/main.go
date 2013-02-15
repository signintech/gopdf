package main

import (
	"fmt"
	"gopdf"
)

func main() {
	fmt.Println("start...")
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{Unit: "pt", PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 A4
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 120)
	pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xxxx")
	pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xzxzxzxzx")
	//pdf.AddPage()
	//pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xxxx")
	pdf.WritePdf("/home/oneplus/x.pdf")
	fmt.Println("end...")
}
