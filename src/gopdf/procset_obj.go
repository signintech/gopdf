package gopdf

import (
	"bytes"
)

type ProcSet struct{
	buffer bytes.Buffer
	//Font
}


func (me *ProcSet) Init(funcGetRoot func()(*GoPdf)) {
	//me.getRoot = funcGetRoot
}

func (me *ProcSet) Build() {
	
	
}

func (me *ProcSet) GetType() string {
	return "ProcSet"
}

func (me *ProcSet) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}