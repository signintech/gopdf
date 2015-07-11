package gopdf

import (
	"bytes"
	"fmt"
	"log"

	"github.com/signintech/gopdf/fontmaker/core"
)

//PdfType0 Font
type SubsetFontObj struct {
	buffer                bytes.Buffer
	ttfp                  core.TTFParser
	Family                string
	CharacterToGlyphIndex map[rune]uint64
}

func (me *SubsetFontObj) Init(funcGetRoot func() *GoPdf) {
	me.CharacterToGlyphIndex = make(map[rune]uint64)
}

func (me *SubsetFontObj) Build() {
	//me.AddChars("จ")
	me.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(me.Family)))
	me.buffer.WriteString("/DescendantFonts [9 0 R]\n") //TODO fix
	me.buffer.WriteString("/Encoding /Identity-H\n")
	me.buffer.WriteString("/Subtype /Type0\n")
	me.buffer.WriteString("/ToUnicode 8 0 R\n") //TODO fix
	me.buffer.WriteString("/Type /Font\n")
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
