package gopdf


import (
	"bytes"
	"strconv"
	"fmt"
)

type ContentObj struct { //impl IObj
	buffer bytes.Buffer
	text bytes.Buffer
	getRoot func()(*GoPdf)
}

func (me *ContentObj) Init(funcGetRoot func()(*GoPdf)) {
	me.getRoot = funcGetRoot
}

func (me *ContentObj) Build() {

	stream,streamlen := me.buildStream()
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Length "+strconv.Itoa(streamlen)+"\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.WriteString(stream.String())
	me.buffer.WriteString("endstream\n")
}

func (me *ContentObj) GetType() string {
	return "Content"
}

func (me *ContentObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *ContentObj) AppendText(text string){
	me.text.WriteString(text)
}

func (me *ContentObj) ResetText(){
	me.text.Reset()
}

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
	//x := 1
	//y := 1
	//temp.WriteString("2 J\n");
	//temp.WriteString("0.57 w\n");
	//temp.WriteString("BT /F1 16.00 Tf ET\n");
	//temp.WriteString(fmt.Sprintf("BT 31.19 794.57 Td (%s) Tj ET\n",me.text.String()));
	return temp,temp.Len()
}

