package gopdf

import (
	"bytes"
	"fmt"
)

type ProcSetObj struct{
	buffer bytes.Buffer
	//Font
	Realtes RelateFonts
}


func (me *ProcSetObj) Init(funcGetRoot func()(*GoPdf)) {
	//me.getRoot = funcGetRoot
}

func (me *ProcSetObj) Build() {

	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n")
	me.buffer.WriteString("/Font <<\n")
	//me.buffer.WriteString("/F1 9 0 R
	//me.buffer.WriteString("/F2 12 0 R
	//me.buffer.WriteString("/F3 15 0 R
	i := 0
	max := len(me.Realtes)
	for i < max {
		realte := me.Realtes[i]
		me.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n",realte.CountOfFont + 1, realte.IndexOfObj + 1))
		i++
	}
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("/XObject <<\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString(">>\n")
}

func (me *ProcSetObj) GetType() string {
	return "ProcSet"
}

func (me *ProcSetObj) GetObjBuff() *bytes.Buffer {
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
	CountOfFont int 
	//เช่น  5 0 R
	IndexOfObj int
	
}