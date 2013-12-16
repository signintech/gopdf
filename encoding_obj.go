package gopdf

import (
	"bytes"
)

type EncodingObj struct{
	buffer    bytes.Buffer
	font IFont
}

func (me * EncodingObj )Init( funcGetRoot func()(*GoPdf)){
	
}
func (me * EncodingObj )GetType() string {
	return "Encoding"
}
func (me * EncodingObj )GetObjBuff() *bytes.Buffer {
	return &me.buffer
}
func (me * EncodingObj ) Build(){
	me.buffer.WriteString("<</Type /Encoding /BaseEncoding /WinAnsiEncoding /Differences [")
	me.buffer.WriteString(me.font.GetDiff())
	me.buffer.WriteString("]>>\n");
}

func (me * EncodingObj) SetFont(font IFont){
	me.font  = font
}

func (me * EncodingObj) GetFont() IFont{
	return me.font
}