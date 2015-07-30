package gopdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf/fontmaker/core"
)

//PdfType0 Font
type SubsetFontObj struct {
	buffer                bytes.Buffer
	ttfp                  core.TTFParser
	Family                string
	CharacterToGlyphIndex map[rune]uint64
	CountOfFont           int
	indexObjCIDFont       int
	indexObjUnicodeMap    int
}

func (me *SubsetFontObj) Init(funcGetRoot func() *GoPdf) {
	me.CharacterToGlyphIndex = make(map[rune]uint64)
}

func (me *SubsetFontObj) Build() error {
	//me.AddChars("จ")
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(me.Family)))
	me.buffer.WriteString(fmt.Sprintf("/DescendantFonts [%d 0 R]\n", me.indexObjCIDFont+1)) //TODO fix
	me.buffer.WriteString("/Encoding /Identity-H\n")
	me.buffer.WriteString("/Subtype /Type0\n")
	me.buffer.WriteString(fmt.Sprintf("/ToUnicode %d 0 R\n", me.indexObjUnicodeMap+1)) //TODO fix
	me.buffer.WriteString("/Type /Font\n")
	me.buffer.WriteString(">>\n")
	return nil
}

func (me *SubsetFontObj) SetIndexObjCIDFont(index int) {
	me.indexObjCIDFont = index
}

func (me *SubsetFontObj) SetIndexObjUnicodeMap(index int) {
	me.indexObjUnicodeMap = index
}

func (me *SubsetFontObj) SetFamily(familyname string) {
	me.Family = familyname
}

func (me *SubsetFontObj) GetFamily() string {
	return me.Family
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

func (me *SubsetFontObj) CharIndex(r rune) (uint64, error) {
	if index, ok := me.CharacterToGlyphIndex[r]; ok {
		return index, nil
	}
	return 0, ErrCharNotFound
}

func (me *SubsetFontObj) CharWidth(r rune) (uint64, error) {
	glyphIndex := me.CharacterToGlyphIndex
	if index, ok := glyphIndex[r]; ok {
		return me.GlyphIndexToPdfWidth(index), nil
	}
	return 0, ErrCharNotFound
}

func (me *SubsetFontObj) GetType() string {
	return "SubsetFont"
}

func (me *SubsetFontObj) GetObjBuff() *bytes.Buffer {
	//fmt.Printf("%s\n", me.buffer.String())
	return &me.buffer
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
	//fmt.Printf("\ncccc--->%#v\n", me.ttfp.Chars())
	if value < me.ttfp.StartCount[seg] {
		return 0
	}

	if me.ttfp.IdRangeOffset[seg] == 0 {
		return (value + me.ttfp.IdDelta[seg]) & 0xFFFF
	}
	//fmt.Printf("IdRangeOffset=%d\n", me.ttfp.IdRangeOffset[seg])
	idx := me.ttfp.IdRangeOffset[seg]/2 + (value - me.ttfp.StartCount[seg]) - (segCount - seg)

	if me.ttfp.GlyphIdArray[int(idx)] == uint64(0) {
		return 0
	}

	return (me.ttfp.GlyphIdArray[int(idx)] + me.ttfp.IdDelta[seg]) & 0xFFFF
}

func (me *SubsetFontObj) GlyphIndexToPdfWidth(glyphIndex uint64) uint64 {

	numberOfHMetrics := me.ttfp.NumberOfHMetrics()
	unitsPerEm := me.ttfp.UnitsPerEm()
	if glyphIndex >= numberOfHMetrics {
		glyphIndex = numberOfHMetrics - 1
	}

	width := me.ttfp.Widths()[glyphIndex]
	if unitsPerEm == 1000 {
		return width
	}
	return width * 1000 / unitsPerEm
}

func (me *SubsetFontObj) GetTTFParser() *core.TTFParser {
	return &me.ttfp
}

func (me *SubsetFontObj) GetUt() int64 {
	return me.ttfp.UnderlineThickness()
}
