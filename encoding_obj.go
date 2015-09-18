package gopdf

import (
	"bytes"
)

type EncodingObj struct {
	buffer bytes.Buffer
	font   IFont
}

func (e *EncodingObj) Init(funcGetRoot func() *GoPdf) {

}
func (e *EncodingObj) GetType() string {
	return "Encoding"
}
func (e *EncodingObj) GetObjBuff() *bytes.Buffer {
	return &e.buffer
}
func (e *EncodingObj) Build() error {
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
