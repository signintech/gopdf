package gopdf

import (
	"bytes"
	"fmt"
)

type CIDFontObj struct {
	buffer                    bytes.Buffer
	PtrToSubsetFontObj        *SubsetFontObj
	indexObjSubfontDescriptor int
}

func (me *CIDFontObj) Init(funcGetRoot func() *GoPdf) {
}

func (me *CIDFontObj) Build() {

	me.buffer.WriteString("<<\n")
	me.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(me.PtrToSubsetFontObj.GetFamily())))
	me.buffer.WriteString("/CIDSystemInfo\n")
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Ordering (Identity)\n")
	me.buffer.WriteString("  /Registry (Adobe)\n")
	me.buffer.WriteString("  /Supplement 0\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString(fmt.Sprintf("/FontDescriptor %d 0 R\n", me.indexObjSubfontDescriptor+1)) //TODO fix
	me.buffer.WriteString("/Subtype /CIDFontType2\n")
	me.buffer.WriteString("/Type /Font\n")
	characterToGlyphIndex := me.PtrToSubsetFontObj.CharacterToGlyphIndex
	me.buffer.WriteString("/W [")
	for _, v := range characterToGlyphIndex {
		width := me.PtrToSubsetFontObj.GlyphIndexToPdfWidth(v)
		me.buffer.WriteString(fmt.Sprintf("%d[%d]", v, width))
	}
	me.buffer.WriteString("]\n")
	me.buffer.WriteString(">>\n")
}

func (me *CIDFontObj) SetIndexObjSubfontDescriptor(index int) {
	me.indexObjSubfontDescriptor = index
}

func (me *CIDFontObj) GetType() string {
	return "CIDFont"
}

func (me *CIDFontObj) GetObjBuff() *bytes.Buffer {
	//fmt.Printf("%s\n", me.buffer.String())
	return &me.buffer
}

func (me *CIDFontObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}
