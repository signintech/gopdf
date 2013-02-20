package gopdf

import (
	"bytes"
)

type PageObj struct { //impl IObj
	buffer   bytes.Buffer
	Contents string
	realtes []string
}

func (me *PageObj) Init(funcGetRoot func()(*GoPdf)) {

}

func (me *PageObj) Build(){

	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Type /" + me.GetType() + "\n")
	me.buffer.WriteString("  /Parent 2 0 R\n")
	me.buffer.WriteString("  /Resources <<\n")
	me.buffer.WriteString("    /Font <<\n")
	i := 0
	max := len(me.realtes)
	for i < max {
		me.buffer.WriteString("      "+me.realtes[i]+"\n") //example: /F1 8 0 R 
		i++
	}
	me.buffer.WriteString("    >>\n")
	me.buffer.WriteString("  >>\n")
	me.buffer.WriteString("  /Contents " + me.Contents + "\n") //sample  Contents 8 0 R
	me.buffer.WriteString(">>\n")
}

func (me *PageObj) GetType() string {
	return "Page"
}

func (me *PageObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me * PageObj) SetResFontRelates(realtes []string){
	
}

