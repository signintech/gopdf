package gopdf

import (
	"bytes"
)

type EncodingObj struct {
	buffer bytes.Buffer
	font   IFont
}

func (e *EncodingObj) init(funcGetRoot func() *GoPdf) {

}
func (e *EncodingObj) getType() string {
	return "Encoding"
}
func (e *EncodingObj) getObjBuff() *bytes.Buffer {
	return &e.buffer
}
func (e *EncodingObj) build(objID int) error {
	e.buffer.WriteString("<</Type /Encoding /BaseEncoding /WinAnsiEncoding /Differences [")
	e.buffer.WriteString(e.font.GetDiff())
	e.buffer.WriteString("]>>\n")
	return nil
}

func (e *EncodingObj) SetFont(font IFont) {
	e.font = font
}

func (e *EncodingObj) GetFont() IFont {
	return e.font
}
