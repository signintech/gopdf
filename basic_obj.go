package gopdf

import (
	"bytes"
)

//BasicObj : basic object in pdf
type BasicObj struct {
	buffer bytes.Buffer
	Data   string
}

//Init : init BasicObj
func (b *BasicObj) Init(funcGetRoot func() *GoPdf) {
}

//Build : build buff
func (b *BasicObj) Build() error {
	b.buffer.WriteString(b.Data)
	return nil
}

//GetType : type of object
func (b *BasicObj) GetType() string {
	return "Basic"
}

//GetObjBuff : get buffer
func (b *BasicObj) GetObjBuff() *bytes.Buffer {
	return &(b.buffer)
}
