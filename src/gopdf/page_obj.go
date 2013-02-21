package gopdf

import (
	"bytes"
	//"fmt"
)

type PageObj struct { //impl IObj
	buffer   bytes.Buffer
	Contents string
	ResourcesRelate string
}

func (me *PageObj) Init(funcGetRoot func()(*GoPdf)) {

}

func (me *PageObj) Build(){

	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Type /" + me.GetType() + "\n")
	me.buffer.WriteString("  /Parent 2 0 R\n")
	me.buffer.WriteString("  /Resources "+me.ResourcesRelate+"\n")
	/*me.buffer.WriteString("    /Font <<\n")
	i := 0
	max := len(me.Realtes)
	for i < max {
		realte := me.Realtes[i]
		me.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n",realte.CountOfFont + 1, realte.IndexOfObj + 1))
		i++
	}
	me.buffer.WriteString("    >>\n")*/
	//me.buffer.WriteString("  >>\n")
	me.buffer.WriteString("  /Contents " + me.Contents + "\n") //sample  Contents 8 0 R
	me.buffer.WriteString(">>\n")
}

func (me *PageObj) GetType() string {
	return "Page"
}

func (me *PageObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}



