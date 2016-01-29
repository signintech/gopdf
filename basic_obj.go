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
func (b *BasicObj) init(funcGetRoot func() *GoPdf) {
}

//Build : build buff
func (b *BasicObj) build() error {
	b.buffer.WriteString(b.Data)
	return nil
}

//GetType : type of object
func (b *BasicObj) getType() string {
	return "Basic"
}

//GetObjBuff : get buffer
func (b *BasicObj) getObjBuff() *bytes.Buffer {
	return &(b.buffer)
}
