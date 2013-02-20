package gopdf

import (
	"bytes"
	"fmt"
)

type PageObj struct { //impl IObj
	buffer   bytes.Buffer
	Contents string
	Realtes RelateFonts
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
	max := len(me.Realtes)
	for i < max {
		//me.buffer.WriteString("      "+me.realtes[i]+"\n") //example: /F1 8 0 R 
		realte := me.Realtes[i]
		me.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n",realte.IndexOfFontInPage + 1, realte.IndexOfObj + 1))
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

type  RelateFonts []RelateFont

func (me * RelateFonts ) IsContainsFamily(family string) bool{
	i := 0
	max := len(*me)
	for i < max {
		if (*me)[i].Family == family {
			return true
		}
		i++;
	}
	return false
}


type RelateFont struct{
	
	Family string
	//เช่น /F1
	IndexOfFontInPage int 
	//เช่น  5 0 R
	IndexOfObj int
	
}
