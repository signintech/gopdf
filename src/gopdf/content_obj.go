package gopdf


import (
	"bytes"
	"strconv"
	"fmt"
)

type ContentObj struct { //impl IObj
	buffer bytes.Buffer
	stream  bytes.Buffer
	
	//text bytes.Buffer
	getRoot func()(*GoPdf)
}

func (me *ContentObj) Init(funcGetRoot func()(*GoPdf)) {
	me.getRoot = funcGetRoot
}

func (me *ContentObj) Build() {

	streamlen := me.stream.Len()
		
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Length "+strconv.Itoa(streamlen)+"\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.Write(me.stream.Bytes())
	me.buffer.WriteString("endstream\n")
}

func (me *ContentObj) GetType() string {
	return "Content"
}

func (me *ContentObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}



/*
func (me *ContentObj) buildStream() (*bytes.Buffer,int){
	fontsize := fmt.Sprintf("%d",me.getRoot().Curr.FontSize)
	x := fmt.Sprintf("%0.2f",me.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f",me.getRoot().config.PageSize.H  - me.getRoot().Curr.Y - (float64(me.getRoot().Curr.FontSize) *0.7) )

	temp := new(bytes.Buffer)
	temp.WriteString("BT\n")
	temp.WriteString(x+" "+y+" TD\n")
	temp.WriteString("/F1 "+fontsize+" Tf\n")
	temp.WriteString("("+me.text.String()+") Tj\n")
	temp.WriteString("ET\n")
	
	
	me.getRoot().Curr.X += StrHelper_GetStringWidth(me.text.String(),me.getRoot().Curr.FontSize)

	return temp,temp.Len()
}
*/

func (me *ContentObj) AppendStream(rectangle *Rect,text string){
	
	
	fontobj := me.getRoot().pdfObjs[me.getRoot().Curr.IndexOfFontObj].(*FontObj)
	fontSize := fontobj.Size
	
	x := fmt.Sprintf("%0.2f",me.getRoot().Curr.X)
	y := fmt.Sprintf("%0.2f",me.getRoot().config.PageSize.H  - me.getRoot().Curr.Y - (float64(fontSize) *0.7) )
	
	me.stream.WriteString("BT\n")
	me.stream.WriteString(x+" "+y+" TD\n")
	me.stream.WriteString("/F"+strconv.Itoa(fontobj.CountOfFont + 1)+" "+ strconv.Itoa(fontSize)+" Tf\n")
	me.stream.WriteString("("+text+") Tj\n")
	me.stream.WriteString("ET\n")
	
	if rectangle == nil {
		me.getRoot().Curr.X += StrHelper_GetStringWidth(text,fontSize,fontobj.Font)
	}else{
		me.getRoot().Curr.X += rectangle.W
	}
	
	
}
