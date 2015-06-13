package gopdf

import (
	"bytes"

	"github.com/signintech/gopdf/fontmaker/core"
)

//PdfType0 Font
type SubsetFontObj struct {
	ttfp                  core.TTFParser
	familyName            string
	CharacterToGlyphIndex map[rune]int
}

func (me *SubsetFontObj) Init(funcGetRoot func() *GoPdf) {

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
		me.CharCodeToGlyphIndex(runeValue)
	}
}

func (me *SubsetFontObj) GetType() string {
	return "SubsetFont"
}

func (me *SubsetFontObj) GetObjBuff() *bytes.Buffer {
	return nil
}

func (me *SubsetFontObj) CharCodeToGlyphIndex(r rune) {
	/*seg := uint64(0)
	value := uint64(r)
	segCount := ttfp.SegCount
	for seg < segCount {
		if value <= ttfp.EndCount[seg] {
			break
		}
		seg++
	}

	if value < ttfp.StartCount[seg] {
		return 0
	}

	if ttfp.IdRangeOffset[seg] == 0 {
		return (value + ttfp.IdDelta[seg]) & 0xFFFF
	}
	//idx := uint64(ttfp.IdRangeOffset[seg]/2 + (value - ttfp.StartCount[seg]) - (segCount - seg))
	log.Panic("unsupport yet!")
	return 0*/
}
