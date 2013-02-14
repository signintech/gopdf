package main

import (
	"fmt"
	"gopdf"
)

func main() {
	fmt.Println("start...")
	pdf := gopdf.GoPdf{}
	pdf.Start()
	pdf.AddPage()
	pdf.SetFont("Arial","B",16)
	pdf.Cell(gopdf.Rect{H: 10, W: 10}, "")
	//pdf.AddPage()
	//pdf.Cell(gopdf.Rect{H: 10, W: 10}, "xxxx")
	pdf.WritePdf("/home/oneplus/x.pdf")
	fmt.Println("end...")
}
