package gopdf

import (
	"bytes"
	"fmt"
	"strings"
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

func (e *EncryptionObj) build(objID int) error {
	e.buffer.WriteString("<<\n")
	e.buffer.WriteString("/Filter /Standard\n")
	e.buffer.WriteString("/V 1\n")
	e.buffer.WriteString("/R 2\n")
	e.buffer.WriteString(fmt.Sprintf("/O (%s)\n", e.escape(e.oValue)))
	e.buffer.WriteString(fmt.Sprintf("/U (%s)\n", e.escape(e.uValue)))
	e.buffer.WriteString(fmt.Sprintf("/P %d\n", e.pValue))
	e.buffer.WriteString(">>\n")
	return nil
}

func (e *EncryptionObj) escape(b []byte) string {
	s := string(b)
	s = strings.Replace(s, "\\", "\\\\", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "\r", "\\r", -1)
	return s
}

//GetObjBuff get buffer
func (e *EncryptionObj) GetObjBuff() *bytes.Buffer {
	return e.getObjBuff()
}

//Build build buffer
func (e *EncryptionObj) Build(objID int) error {
	return e.build(objID)
}
