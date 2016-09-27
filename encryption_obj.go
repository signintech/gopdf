package gopdf

import (
	"bytes"
	"fmt"
)

//EncryptionObj  encryption object res
type EncryptionObj struct {
	buffer bytes.Buffer
	uValue []byte //U entry in pdf document
	oValue []byte //O entry in pdf document
	pValue int    //P entry in pdf document
}

func (e *EncryptionObj) init(func() *GoPdf) {

}

func (e *EncryptionObj) getType() string {
	return "Encryption"
}

func (e *EncryptionObj) getObjBuff() *bytes.Buffer {
	return &e.buffer
}

func (e *EncryptionObj) build() error {
	e.buffer.WriteString("/Filter /Standard\n")
	e.buffer.WriteString("/V 1\n")
	e.buffer.WriteString("/R 2\n")
	e.buffer.WriteString(fmt.Sprintf("/O (%s)\n", e.oValue))
	e.buffer.WriteString(fmt.Sprintf("/U (%s)\n", e.uValue))
	e.buffer.WriteString(fmt.Sprintf("/P %d", e.pValue))
	return nil
}
