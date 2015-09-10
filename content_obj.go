package gopdf

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
)

type ContentObj struct { //impl IObj
	buffer bytes.Buffer
	stream bytes.Buffer

	//text bytes.Buffer
	getRoot func() *GoPdf
}

func (me *ContentObj) Init(funcGetRoot func() *GoPdf) {
	me.getRoot = funcGetRoot
}

func (me *ContentObj) Build() error {
	streamlen := me.stream.Len()
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Length " + strconv.Itoa(streamlen) + "\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.Write(me.stream.Bytes())
	me.buffer.WriteString("endstream\n")
	return nil
}

func (me *ContentObj) GetType() string {
	return "Content"
}

func (me *ContentObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *ContentObj) AppendStreamSubsetFont(rectangle *Rect, text string) {

	sumWidth := uint64(0)
	var buff bytes.Buffer
	for _, r := range text {
		index, err := me.getRoot().Curr.Font_ISubset.CharIndex(r)
		if err != nil {
			log.Fatalf("err:%s", err.Error())
		}
		buff.WriteString(fmt.Sprintf("%04X", index))
		width, err := me.getRoot().Curr.Font_ISubset.CharWidth(r)
		if err != nil {
			log.Fatalf("err:%s", err.Error())
		}
		sumWidth += width
	}

	fontSize := me.getRoot().Curr.Font_Size
	x := fmt.Sprintf("%0.2f", me.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f", me.getRoot().config.PageSize.H-me.getRoot().Curr.Y-(float64(fontSize)*0.7))

	me.stream.WriteString("BT\n")
	me.stream.WriteString(x + " " + y + " TD\n")
	me.stream.WriteString("/F" + strconv.Itoa(me.getRoot().Curr.Font_FontCount+1) + " " + strconv.Itoa(fontSize) + " Tf\n")
	me.stream.WriteString("<" + buff.String() + "> Tj\n")
	me.stream.WriteString("ET\n")
	if rectangle == nil {
		fontSize := me.getRoot().Curr.Font_Size
		me.getRoot().Curr.X += float64(sumWidth) * (float64(fontSize) / 1000.0)
	} else {
		me.getRoot().Curr.X += rectangle.W
	}
}

func (me *ContentObj) AppendStream(rectangle *Rect, text string) {

	fontSize := me.getRoot().Curr.Font_Size

	x := fmt.Sprintf("%0.2f", me.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f", me.getRoot().config.PageSize.H-me.getRoot().Curr.Y-(float64(fontSize)*0.7))

	me.stream.WriteString("BT\n")
	me.stream.WriteString(x + " " + y + " TD\n")
	me.stream.WriteString("/F" + strconv.Itoa(me.getRoot().Curr.Font_FontCount+1) + " " + strconv.Itoa(fontSize) + " Tf\n")
	me.stream.WriteString("(" + text + ") Tj\n")
	me.stream.WriteString("ET\n")
	if rectangle == nil {
		me.getRoot().Curr.X += StrHelperGetStringWidth(text, fontSize, me.getRoot().Curr.Font_IFont)
	} else {
		me.getRoot().Curr.X += rectangle.W
	}

}

func (me *ContentObj) AppendStreamLine(x1 float64, y1 float64, x2 float64, y2 float64) {

	h := me.getRoot().config.PageSize.H
	me.stream.WriteString(fmt.Sprintf("%0.2f %0.2f m %0.2f %0.2f l s\n", x1, h-y1, x2, h-y2))
}

func (me *ContentObj) AppendUnderline(startX float64, y float64, endX float64, endY float64, text string) {

	h := me.getRoot().config.PageSize.H
	ut := int(0)
	if me.getRoot().Curr.Font_IFont != nil {
		ut = me.getRoot().Curr.Font_IFont.GetUt()
	} else if me.getRoot().Curr.Font_ISubset != nil {
		ut = int(me.getRoot().Curr.Font_ISubset.GetUt())
	} else {
		log.Fatal("error AppendUnderline not found font")
	}

	textH := ContentObj_CalTextHeight(me.getRoot().Curr.Font_Size)
	arg3 := float64(h) - float64(y) - textH - textH*0.07
	arg4 := (float64(ut) / 1000.00) * float64(me.getRoot().Curr.Font_Size)
	me.stream.WriteString(fmt.Sprintf("%0.2f %0.2f %0.2f -%0.2f re f\n", startX, arg3, endX-startX, arg4))
}

func (me *ContentObj) AppendStreamSetLineWidth(w float64) {

	me.stream.WriteString(fmt.Sprintf("%.2f w\n", w))

}

//  Set the grayscale fills
func (me *ContentObj) AppendStreamSetGrayFill(w float64) {
	me.stream.WriteString(fmt.Sprintf("%.2f g\n", w))
}

//  Set the grayscale stroke
func (me *ContentObj) AppendStreamSetGrayStroke(w float64) {
	me.stream.WriteString(fmt.Sprintf("%.2f G\n", w))
}

func (me *ContentObj) AppendStreamImage(index int, x float64, y float64, rect *Rect) {
	//fmt.Printf("index = %d",index)
	h := me.getRoot().config.PageSize.H
	me.stream.WriteString(fmt.Sprintf("q %0.2f 0 0 %0.2f %0.2f %0.2f cm /I%d Do Q\n", rect.W, rect.H, x, h-(y+rect.H), index+1))
}

//cal text height
func ContentObj_CalTextHeight(fontsize int) float64 {
	return (float64(fontsize) * 0.7)
}
