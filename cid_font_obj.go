package gopdf

import (
	"bytes"
	"fmt"
)

type CIDFontObj struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (me *CIDFontObj) Init(funcGetRoot func() *GoPdf) {
}

func (me *CIDFontObj) Build() {

	me.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(me.PtrToSubsetFontObj.GetFamily())))
	me.buffer.WriteString("/CIDSystemInfo\n")
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Ordering (Identity)\n")
	me.buffer.WriteString("  /Registry (Adobe)\n")
	me.buffer.WriteString("  /Supplement 0\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("/FontDescriptor 7 0 R\n") //TODO fix
	me.buffer.WriteString("/Subtype /CIDFontType2\n")
	me.buffer.WriteString("/Type /Font\n")
	characterToGlyphIndex := me.PtrToSubsetFontObj.CharacterToGlyphIndex
	me.buffer.WriteString("/W [")
	for _, v := range characterToGlyphIndex {
		width := me.PtrToSubsetFontObj.GlyphIndexToPdfWidth(v)
		me.buffer.WriteString(fmt.Sprintf("%d[%d]", v, width))
	}
	me.buffer.WriteString("]\n")

}

func (me *CIDFontObj) GetType() string {
	return "CIDFont"
}

func (me *CIDFontObj) GetObjBuff() *bytes.Buffer {
	fmt.Printf("%s\n", me.buffer.String())
	return &me.buffer
}

func (me *CIDFontObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}
