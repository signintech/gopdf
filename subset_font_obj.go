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

func (s *SubsetFontObj) init(funcGetRoot func() *GoPdf) {
	s.CharacterToGlyphIndex = make(map[rune]uint64)
}

func (s *SubsetFontObj) build() error {
	//me.AddChars("จ")
	s.buffer.WriteString("<<\n")
	s.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(s.Family)))
	s.buffer.WriteString(fmt.Sprintf("/DescendantFonts [%d 0 R]\n", s.indexObjCIDFont+1)) //TODO fix
	s.buffer.WriteString("/Encoding /Identity-H\n")
	s.buffer.WriteString("/Subtype /Type0\n")
	s.buffer.WriteString(fmt.Sprintf("/ToUnicode %d 0 R\n", s.indexObjUnicodeMap+1)) //TODO fix
	s.buffer.WriteString("/Type /Font\n")
	s.buffer.WriteString(">>\n")
	return nil
}

func (s *SubsetFontObj) SetIndexObjCIDFont(index int) {
	s.indexObjCIDFont = index
}

func (s *SubsetFontObj) SetIndexObjUnicodeMap(index int) {
	s.indexObjUnicodeMap = index
}

func (s *SubsetFontObj) SetFamily(familyname string) {
	s.Family = familyname
}

func (s *SubsetFontObj) GetFamily() string {
	return s.Family
}

func (s *SubsetFontObj) SetTTFByPath(ttfpath string) error {
	err := s.ttfp.Parse(ttfpath)
	if err != nil {
		return err
	}
	return nil
}

func (s *SubsetFontObj) AddChars(txt string) {
	for _, runeValue := range txt {
		if _, ok := s.CharacterToGlyphIndex[runeValue]; ok {
			continue
		}
		glyphIndex := s.CharCodeToGlyphIndex(runeValue)
		s.CharacterToGlyphIndex[runeValue] = glyphIndex
	}
}

func (s *SubsetFontObj) CharIndex(r rune) (uint64, error) {
	if index, ok := s.CharacterToGlyphIndex[r]; ok {
		return index, nil
	}
	return 0, ErrCharNotFound
}

func (s *SubsetFontObj) CharWidth(r rune) (uint64, error) {
	glyphIndex := s.CharacterToGlyphIndex
	if index, ok := glyphIndex[r]; ok {
		return s.GlyphIndexToPdfWidth(index), nil
	}
	return 0, ErrCharNotFound
}

func (s *SubsetFontObj) getType() string {
	return "SubsetFont"
}

func (s *SubsetFontObj) getObjBuff() *bytes.Buffer {
	//fmt.Printf("%s\n", me.buffer.String())
	return &s.buffer
}

func (s *SubsetFontObj) CharCodeToGlyphIndex(r rune) uint64 {
	seg := uint64(0)
	value := uint64(r)
	segCount := s.ttfp.SegCount
	for seg < segCount {
		if value <= s.ttfp.EndCount[seg] {
			break
		}
		seg++
	}
	//fmt.Printf("\ncccc--->%#v\n", me.ttfp.Chars())
	if value < s.ttfp.StartCount[seg] {
		return 0
	}

	if s.ttfp.IdRangeOffset[seg] == 0 {
		return (value + s.ttfp.IdDelta[seg]) & 0xFFFF
	}
	//fmt.Printf("IdRangeOffset=%d\n", me.ttfp.IdRangeOffset[seg])
	idx := s.ttfp.IdRangeOffset[seg]/2 + (value - s.ttfp.StartCount[seg]) - (segCount - seg)

	if s.ttfp.GlyphIdArray[int(idx)] == uint64(0) {
		return 0
	}

	return (s.ttfp.GlyphIdArray[int(idx)] + s.ttfp.IdDelta[seg]) & 0xFFFF
}

func (s *SubsetFontObj) GlyphIndexToPdfWidth(glyphIndex uint64) uint64 {

	numberOfHMetrics := s.ttfp.NumberOfHMetrics()
	unitsPerEm := s.ttfp.UnitsPerEm()
	if glyphIndex >= numberOfHMetrics {
		glyphIndex = numberOfHMetrics - 1
	}

	width := s.ttfp.Widths()[glyphIndex]
	if unitsPerEm == 1000 {
		return width
	}
	return width * 1000 / unitsPerEm
}

func (s *SubsetFontObj) GetTTFParser() *core.TTFParser {
	return &s.ttfp
}

func (s *SubsetFontObj) GetUt() int64 {
	return s.ttfp.UnderlineThickness()
}
