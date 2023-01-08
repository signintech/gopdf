package gopdf

import (
	"fmt"
	"io"
)

// UnicodeMap unicode map
type UnicodeMap struct {
	PtrToSubsetFontObj *SubsetFontObj
	//getRoot            func() *GoPdf
	pdfProtection *PDFProtection
}

func (u *UnicodeMap) init(funcGetRoot func() *GoPdf) {
	//u.getRoot = funcGetRoot
}

func (u *UnicodeMap) setProtection(p *PDFProtection) {
	u.pdfProtection = p
}

func (u *UnicodeMap) protection() *PDFProtection {
	return u.pdfProtection
}

// SetPtrToSubsetFontObj set pointer to SubsetFontObj
func (u *UnicodeMap) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	u.PtrToSubsetFontObj = ptr
}

func (u *UnicodeMap) getType() string {
	return "Unicode"
}

func (u *UnicodeMap) write(w io.Writer, objID int) error {
	//stream
	//characterToGlyphIndex := u.PtrToSubsetFontObj.CharacterToGlyphIndex
	prefix :=
		"/CIDInit /ProcSet findresource begin\n" +
			"12 dict begin\n" +
			"begincmap\n" +
			"/CIDSystemInfo << /Registry (Adobe)/Ordering (UCS)/Supplement 0>> def\n" +
			"/CMapName /Adobe-Identity-UCS def /CMapType 2 def\n"
	suffix := "endcmap CMapName currentdict /CMap defineresource pop end end"

	glyphIndexToCharacter := newMapGlyphIndexToCharacter() //make(map[int]rune)
	lowIndex := 65536
	hiIndex := -1

	keys := u.PtrToSubsetFontObj.CharacterToGlyphIndex.AllKeys()
	for _, k := range keys {
		v, _ := u.PtrToSubsetFontObj.CharacterToGlyphIndex.Val(k)
		index := int(v)
		if index < lowIndex {
			lowIndex = index
		}
		if index > hiIndex {
			hiIndex = index
		}
		//glyphIndexToCharacter[index] = k
		glyphIndexToCharacter.set(index, k)
	}

	buff := GetBuffer()
	defer PutBuffer(buff)

	buff.WriteString(prefix)
	buff.WriteString("1 begincodespacerange\n")
	fmt.Fprintf(buff, "<%04X><%04X>\n", lowIndex, hiIndex)
	buff.WriteString("endcodespacerange\n")
	fmt.Fprintf(buff, "%d beginbfrange\n", glyphIndexToCharacter.size())
	indexs := glyphIndexToCharacter.allIndexs()
	for _, k := range indexs {
		v, _ := glyphIndexToCharacter.runeByIndex(k)
		fmt.Fprintf(buff, "<%04X><%04X><%04X>\n", k, k, v)
	}
	buff.WriteString("endbfrange\n")
	buff.WriteString(suffix)
	buff.WriteString("\n")

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/Length %d\n", buff.Len())
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if u.protection() != nil {
		tmp, err := rc4Cip(u.protection().objectkey(objID), buff.Bytes())
		if err != nil {
			return err
		}
		w.Write(tmp)
		//streambuff.WriteString("\n")
	} else {
		buff.WriteTo(w)
	}
	io.WriteString(w, "endstream\n")

	return nil
}

type mapGlyphIndexToCharacter struct {
	runes  []rune
	indexs []int
}

func newMapGlyphIndexToCharacter() *mapGlyphIndexToCharacter {
	var m mapGlyphIndexToCharacter
	return &m
}

func (m *mapGlyphIndexToCharacter) set(index int, r rune) {
	m.runes = append(m.runes, r)
	m.indexs = append(m.indexs, index)
}

func (m *mapGlyphIndexToCharacter) size() int {
	return len(m.indexs)
}

func (m *mapGlyphIndexToCharacter) allIndexs() []int {
	return m.indexs
}

func (m *mapGlyphIndexToCharacter) runeByIndex(index int) (rune, bool) {
	var r rune
	ok := false
	for i, idx := range m.indexs {
		if idx == index {
			r = m.runes[i]
			ok = true
			break
		}
	}
	return r, ok
}
