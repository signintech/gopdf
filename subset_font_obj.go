package gopdf

import (
	"bytes"
	"log"

	"github.com/signintech/gopdf/fontmaker/core"
)

//PdfType0 Font
type SubsetFontObj struct {
	ttfp                  core.TTFParser
	Family                string
	CharacterToGlyphIndex map[rune]uint64
}

func (me *SubsetFontObj) Init(funcGetRoot func() *GoPdf) {
	me.CharacterToGlyphIndex = make(map[rune]uint64)
}

func (me *SubsetFontObj) Build() {
	//print PdfToUnicodeMap
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

	var buffer bytes.Buffer
	//buffer.Write(me.PdfToUnicodeMap().Bytes())
	//log.Printf("\n%s\n", buffer.String())
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
