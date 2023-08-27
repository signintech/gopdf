package gopdf

import (
	"fmt"
	"io"
	"os"
)

// EmbedFontObj is an embedded font object.
type EmbedFontObj struct {
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

func (e *EmbedFontObj) write(w io.Writer, objID int) error {
	b, err := os.ReadFile(e.zfontpath)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "<</Length %d\n", len(b))
	io.WriteString(w, "/Filter /FlateDecode\n")
	fmt.Fprintf(w, "/Length1 %d\n", e.font.GetOriginalsize())
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if e.protection() != nil {
		tmp, err := rc4Cip(e.protection().objectkey(objID), b)
		if err != nil {
			return err
		}
		w.Write(tmp)
		io.WriteString(w, "\n")
	} else {
		w.Write(b)
	}
	io.WriteString(w, "\nendstream\n")
	return nil
}

func (e *EmbedFontObj) getType() string {
	return "EmbedFont"
}

// SetFont sets the font of an embedded font object.
func (e *EmbedFontObj) SetFont(font IFont, zfontpath string) {
	e.font = font
	e.zfontpath = zfontpath
}
