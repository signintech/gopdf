package obj

import (
	"bytes"
)

type ContentObj struct { //impl IObj
	buffer bytes.Buffer
}

func (me *ContentObj) Init() {
}

func (me *ContentObj) Build() {
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Length 44\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.WriteString("BT\n")
	me.buffer.WriteString("70 50 TD\n")
	me.buffer.WriteString("/F1 14 Tf\n")
	me.buffer.WriteString("(Hello, world!) Tj\n")
	me.buffer.WriteString("ET\n")
	me.buffer.WriteString("endstream\n")
}

func (me *ContentObj) GetType() string {
	return "Content"
}

func (me *ContentObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}
