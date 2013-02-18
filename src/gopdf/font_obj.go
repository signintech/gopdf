package gopdf


import (
	"bytes"
)

type FontObj struct { //impl IObj
	buffer bytes.Buffer
}

func (me *FontObj) Init(funcGetRoot func()(*GoPdf)) {
}

func (me *FontObj) Build() {
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Type /" + me.GetType() + "\n")
	me.buffer.WriteString("  /Subtype /Type1\n")
	me.buffer.WriteString("  /BaseFont /Times-Roman\n")
	me.buffer.WriteString(">>\n")
	
}

func (me *FontObj) GetType() string {
	return "Font"
}

func (me *FontObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

