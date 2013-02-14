package obj

import (
	"bytes"
)

type FontObj struct { //impl IObj
	buffer bytes.Buffer
}

func (me *FontObj) Init() {
}

func (me *FontObj) Build() {
	me.buffer.WriteString("\t/Type /" + me.GetType() + "\n")
	me.buffer.WriteString("\t/Subtype /Type1\n")
	me.buffer.WriteString("\t/BaseFont /Times-Roman\n")
}

func (me *FontObj) GetType() string {
	return "Font"
}

func (me *FontObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}
