package gopdf

import (
	"bytes"
	"io/ioutil"
	"strconv"
)


type EmbedFontObj struct{
	buffer    bytes.Buffer
	Data string
	zfontpath string
	font IFont
}

func (me *EmbedFontObj) Init(funcGetRoot func()(*GoPdf)) {
}

func (me *EmbedFontObj) Build() {
	b, err := ioutil.ReadFile(me.zfontpath)
	if err != nil {
		return
	}
	me.buffer.WriteString("<</Length "+ strconv.Itoa(len(b)) +"\n")
	me.buffer.WriteString("/Filter /FlateDecode\n")
	me.buffer.WriteString("/Length1 "+strconv.Itoa(me.font.GetOriginalsize())+"\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.Write(b)
	me.buffer.WriteString("\nendstream\n")
}

func (me *EmbedFontObj) GetType() string {
	return "EmbedFont"
}

func (me *EmbedFontObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *EmbedFontObj) SetFont(font IFont,zfontpath string){
	me.font = font
	me.zfontpath = zfontpath
}
