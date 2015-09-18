package gopdf

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

type EmbedFontObj struct {
	buffer    bytes.Buffer
	Data      string
	zfontpath string
	font      IFont
}

func (e *EmbedFontObj) Init(funcGetRoot func() *GoPdf) {
}

func (e *EmbedFontObj) Build() error {
	b, err := ioutil.ReadFile(e.zfontpath)
	if err != nil {
		return err
	}
	e.buffer.WriteString("<</Length " + strconv.Itoa(len(b)) + "\n")
	e.buffer.WriteString("/Filter /FlateDecode\n")
	e.buffer.WriteString("/Length1 " + strconv.Itoa(e.font.GetOriginalsize()) + "\n")
	e.buffer.WriteString(">>\n")
	e.buffer.WriteString("stream\n")
	e.buffer.Write(b)
	e.buffer.WriteString("\nendstream\n")
	return nil
}

func (e *EmbedFontObj) GetType() string {
	return "EmbedFont"
}

func (e *EmbedFontObj) GetObjBuff() *bytes.Buffer {
	return &(e.buffer)
}

func (e *EmbedFontObj) SetFont(font IFont, zfontpath string) {
	e.font = font
	e.zfontpath = zfontpath
}
