package outline_example

import (
	"fmt"
	"os"
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

func OutlineWithPositionExample() (err error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "./Ubuntu-L.ttf")
	if err != nil {
		return err
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		return err
	}
	pdf.AddPage()

	pdf.SetXY(150, 400)
	pdf.AddOutlineWithPosition("first page")
	err = pdf.Text("1.Hello World")
	if err != nil {
		return err
	}

	pdf.AddPage()
	pdf.SetXY(150, 700)
	pdf.AddOutlineWithPosition("second page")
	err = pdf.Text("2.Hello World")
	if err != nil {
		return err
	}

	pdf.AddPage()
	pdf.SetXY(150, 700)
	pdf.AddOutlineWithPosition("third page")
	err = pdf.Text("3.Hello World")
	if err != nil {
		return err
	}

	pdf.AddPage()
	pdf.SetXY(150, 200)
	pdf.AddOutlineWithPosition("forth page")
	err = pdf.Text("4.Hello World")
	if err != nil {
		return err
	}

	pdf.AddPage()

	err = pdf.WritePdf(fmt.Sprintf("./outline_demo_%d.pdf", time.Now().Unix()-1637580000))
	if err != nil {
		return err
	}
	return err
}

// OutlineWithLevelExample outlines with multiple level
func OutlineWithLevelExample() (err error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	err = GetFont(pdf, "./Ubuntu-L.ttf")
	if err != nil {
		return err
	}
	err = pdf.SetFont("Ubuntu-L", "", 14)
	if err != nil {
		return err
	}

	/*
		expectation:

			-- first page
			-- second page
			   -- level 2-1
			   -- level 2-2
		    	   -- level 3-1
		    	   -- level 3-2
		    -- third page
	*/

	// outline nodes
	var outlineNodes gopdf.OutlineNodes
	var first = new(gopdf.OutlineNode)
	var second = new(gopdf.OutlineNode)
	var third = new(gopdf.OutlineNode)
	var nodeObjs = make([]*gopdf.OutlineNode, 0)
	nodeObjs = append(nodeObjs, first, second, third)
	outlineNodes = nodeObjs

	pdf.AddPage()
	pdf.SetXY(150, 200)
	first.Obj = pdf.AddOutlineWithPosition("first page")
	err = pdf.Text("first page")
	if err != nil {
		return err
	}

	pdf.AddPage()
	pdf.SetXY(150, 200)
	second.Obj = pdf.AddOutlineWithPosition("second page")
	err = pdf.Text("second page")
	if err != nil {
		return err
	}

	pdf.SetY(250)
	lv21 := pdf.AddOutlineWithPosition("level 2-1")
	var node21 = new(gopdf.OutlineNode)
	node21.Obj = lv21
	second.Children = append(second.Children, node21)
	err = pdf.Text("level 2-1...")
	if err != nil {
		return err
	}

	pdf.SetY(350)
	lv22 := pdf.AddOutlineWithPosition("level 2-2")
	var node22 = new(gopdf.OutlineNode)
	node22.Obj = lv22
	second.Children = append(second.Children, node22)
	err = pdf.Text("level 2-2...")
	if err != nil {
		return err
	}

	pdf.SetY(500)
	lv31 := pdf.AddOutlineWithPosition("level 3-1")
	var node31 = new(gopdf.OutlineNode)
	node31.Obj = lv31
	node22.Children = append(node22.Children, node31)
	err = pdf.Text("level 3-1...")
	if err != nil {
		return err
	}

	pdf.SetY(600)
	lv32 := pdf.AddOutlineWithPosition("level 3-2")
	var node32 = new(gopdf.OutlineNode)
	node32.Obj = lv32
	node22.Children = append(node22.Children, node32)
	err = pdf.Text("level 3-2...")
	if err != nil {
		return err
	}

	pdf.AddPage()
	pdf.SetXY(150, 200)
	third.Obj = pdf.AddOutlineWithPosition("third page")
	err = pdf.Text("third page")
	if err != nil {
		return err
	}

	// parse outline nodes
	outlineNodes.Parse()

	err = pdf.WritePdf(fmt.Sprintf("./outline_demo.pdf"))
	if err != nil {
		return err
	}
	return err
}
