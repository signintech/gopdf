package gopdf

import (
	"bytes"
	"fmt"
	"log"

	"github.com/signintech/gopdf/fontmaker/core"
)

//PdfType0 Font
type SubsetFontObj struct {
	ttfp                  core.TTFParser
	familyName            string
	CharacterToGlyphIndex map[rune]uint64
}

func (me *SubsetFontObj) Init(funcGetRoot func() *GoPdf) {
	me.CharacterToGlyphIndex = make(map[rune]uint64)
}

func (me *SubsetFontObj) Build() {
	//print PdfToUnicodeMap
}

func (me *SubsetFontObj) SetFamily(familyname string) {
	me.familyName = familyname
}

func (me *SubsetFontObj) GetFamily() string {
	return me.familyName
}

func (me *SubsetFontObj) SetTTFByPath(ttfpath string) error {
	err := me.ttfp.Parse(ttfpath)
	if err != nil {
		return err
	}
	return nil
}

func (me *SubsetFontObj) AddChars(txt string) {
	for _, runeValue := range txt {
		if _, ok := me.CharacterToGlyphIndex[runeValue]; ok {
			continue
		}
		glyphIndex := me.CharCodeToGlyphIndex(runeValue)
		me.CharacterToGlyphIndex[runeValue] = glyphIndex
	}
}

func (me *SubsetFontObj) GetType() string {
	return "SubsetFont"
}

func (me *SubsetFontObj) GetObjBuff() *bytes.Buffer {

	var buffer bytes.Buffer
	buffer.Write(me.PdfToUnicodeMap().Bytes())
	log.Printf("\n%s\n", buffer.String())
	return &buffer
}

func (me *SubsetFontObj) PdfToUnicodeMap() *bytes.Buffer {

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
	for k, v := range me.CharacterToGlyphIndex {
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

func (me *SubsetFontObj) CharCodeToGlyphIndex(r rune) uint64 {
	seg := uint64(0)
	value := uint64(r)
	segCount := me.ttfp.SegCount
	for seg < segCount {
		if value <= me.ttfp.EndCount[seg] {
			break
		}
		seg++
	}

	if value < me.ttfp.StartCount[seg] {
		return 0
	}

	if me.ttfp.IdRangeOffset[seg] == 0 {
		return (value + me.ttfp.IdDelta[seg]) & 0xFFFF
	}
	//idx := uint64(ttfp.IdRangeOffset[seg]/2 + (value - ttfp.StartCount[seg]) - (segCount - seg))
	log.Panic("unsupport yet!")
	return 0
}
