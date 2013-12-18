package gopdf

import (
	"bytes"
)

type BasicObj struct{
	buffer    bytes.Buffer
	Data string
}

func (me *BasicObj) Init(funcGetRoot func()(*GoPdf)) {
}

func (me *BasicObj) Build() {
	me.buffer.WriteString(me.Data)
}

func (me *BasicObj) GetType() string {
	return "Basic"
}

func (me *BasicObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}


