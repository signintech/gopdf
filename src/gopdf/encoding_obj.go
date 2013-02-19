package gopdf

import (
	"bytes"
	"gopdf/fonts"
)

type EncodingObj struct{
	buffer    bytes.Buffer
	font fonts.IFont
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

func (me * EncodingObj) SetFont(font fonts.IFont){
	me.font  = font
}

func (me * EncodingObj) GetFont() fonts.IFont{
	return me.font
}