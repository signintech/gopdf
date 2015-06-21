package gopdf

import (
	"bytes"
	"fmt"
)

type UnicodeMap struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (me *UnicodeMap) Init(funcGetRoot func() *GoPdf) {
}

func (me *UnicodeMap) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}

func (me *UnicodeMap) Build() {
	me.buffer.Write(me.pdfToUnicodeMap().Bytes())
}

func (me *UnicodeMap) GetType() string {
	return "Unicode"
}

func (me *UnicodeMap) GetObjBuff() *bytes.Buffer {
	return &me.buffer
}

func (me *UnicodeMap) pdfToUnicodeMap() *bytes.Buffer {

	characterToGlyphIndex := me.PtrToSubsetFontObj.CharacterToGlyphIndex
	var buffer bytes.Buffer
	prefix :=
		"/CIDInit /ProcSet findresource begin\n" +
			"12 dict begin\n" +
			"begincmap\n" +
			"/CIDSystemInfo << /Registry (Adobe)/Ordering (UCS)/Supplement 0>> def\n" +
			"/CMapName /Adobe-Identity-UCS def /CMapType 2 def\n"
	suffix := "endcmap CMapName currentdict /CMap defineresource pop end end"

	glyphIndexToCharacter := make(map[int]rune)
	lowIndex := 65536
	hiIndex := -1
	for k, v := range characterToGlyphIndex {
		index := int(v)
		if index < lowIndex {
			lowIndex = index
		}
		if index > hiIndex {
			hiIndex = index
		}
		glyphIndexToCharacter[index] = k
	}

	buffer.WriteString(prefix)
	buffer.WriteString("1 begincodespacerange\n")
	buffer.WriteString(fmt.Sprintf("<%04X><%04X>\n", lowIndex, hiIndex))
	buffer.WriteString("endcodespacerange\n")
	buffer.WriteString(fmt.Sprintf("%d beginbfrange\n", len(glyphIndexToCharacter)))
	for k, v := range glyphIndexToCharacter {
		buffer.WriteString(fmt.Sprintf("<%04X><%04X><%04X>\n", k, k, v))
	}
	buffer.WriteString("endbfrange\n")
	buffer.WriteString(suffix)
	return &buffer
}
