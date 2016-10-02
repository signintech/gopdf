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
	getRoot   func() *GoPdf
}

func (e *EmbedFontObj) init(funcGetRoot func() *GoPdf) {
	e.getRoot = funcGetRoot
}

func (e *EmbedFontObj) protection() *PDFProtection {
	return e.getRoot().protection()
}

func (e *EmbedFontObj) build(objID int) error {
	b, err := ioutil.ReadFile(e.zfontpath)
	if err != nil {
		return err
	}
	e.buffer.WriteString("<</Length " + strconv.Itoa(len(b)) + "\n")
	e.buffer.WriteString("/Filter /FlateDecode\n")
	e.buffer.WriteString("/Length1 " + strconv.Itoa(e.font.GetOriginalsize()) + "\n")
	e.buffer.WriteString(">>\n")
	e.buffer.WriteString("stream\n")
	if e.protection() != nil {
		tmp, err := rc4Cip(e.protection().objectkey(objID), b)
		if err != nil {
			return err
		}
		e.buffer.Write(tmp)
		e.buffer.WriteString("\n")
	} else {
		e.buffer.Write(b)
	}
	e.buffer.WriteString("\nendstream\n")
	return nil
}

func (e *EmbedFontObj) getType() string {
	return "EmbedFont"
}

func (e *EmbedFontObj) getObjBuff() *bytes.Buffer {
	return &(e.buffer)
}

func (e *EmbedFontObj) SetFont(font IFont, zfontpath string) {
	e.font = font
	e.zfontpath = zfontpath
}
