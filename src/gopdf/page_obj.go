package gopdf

import (
	"bytes"
)

type PageObj struct { //impl IObj
	buffer   bytes.Buffer
	Contents string
}

func (me *PageObj) Init(funcGetRoot func()(*GoPdf)) {
}

func (me *PageObj) Build() {
	me.buffer.WriteString("\t/Type /" + me.GetType() + "\n")
	me.buffer.WriteString("\t/Parent 2 0 R\n")
	me.buffer.WriteString("\t/Resources <<\n")
	me.buffer.WriteString("\t\t/Font <<\n")
	me.buffer.WriteString("\t\t\t/F1 4 0 R \n")
	me.buffer.WriteString("\t\t>>\n")
	me.buffer.WriteString("\t>>\n")
	me.buffer.WriteString("\t/Contents " + me.Contents + "\n") //sample  Contents 8 0 R
}

func (me *PageObj) GetType() string {
	return "Page"
}

func (me *PageObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

