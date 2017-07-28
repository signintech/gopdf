package gopdf

import (
	"bytes"
	"errors"
	"fmt"

	"io"

	"github.com/signintech/gopdf/fontmaker/core"
)

//ErrCharNotFound char not found
var ErrCharNotFound = errors.New("char not found")

//SubsetFontObj pdf subsetFont object
type SubsetFontObj struct {
	buffer                bytes.Buffer
	ttfp                  core.TTFParser
	Family                string
	CharacterToGlyphIndex *MapOfCharacterToGlyphIndex
	CountOfFont           int
	indexObjCIDFont       int
	indexObjUnicodeMap    int
	ttfFontOption         TtfOption
	funcKernOverride      FuncKernOverride
}

func (s *SubsetFontObj) init(funcGetRoot func() *GoPdf) {
	s.CharacterToGlyphIndex = NewMapOfCharacterToGlyphIndex() //make(map[rune]uint)
	s.funcKernOverride = nil
}

func (s *SubsetFontObj) build(objID int) error {
	//me.AddChars("จ")
	s.buffer.WriteString("<<\n")
	s.buffer.WriteString(fmt.Sprintf("/BaseFont /%s\n", CreateEmbeddedFontSubsetName(s.Family)))
	s.buffer.WriteString(fmt.Sprintf("/DescendantFonts [%d 0 R]\n", s.indexObjCIDFont+1))
	s.buffer.WriteString("/Encoding /Identity-H\n")
	s.buffer.WriteString("/Subtype /Type0\n")
	s.buffer.WriteString(fmt.Sprintf("/ToUnicode %d 0 R\n", s.indexObjUnicodeMap+1))
	s.buffer.WriteString("/Type /Font\n")
	s.buffer.WriteString(">>\n")
	return nil
}

//SetIndexObjCIDFont set IndexObjCIDFont
func (s *SubsetFontObj) SetIndexObjCIDFont(index int) {
	s.indexObjCIDFont = index
}

//SetIndexObjUnicodeMap set IndexObjUnicodeMap
func (s *SubsetFontObj) SetIndexObjUnicodeMap(index int) {
	s.indexObjUnicodeMap = index
}

//SetFamily set font family name
func (s *SubsetFontObj) SetFamily(familyname string) {
	s.Family = familyname
}

//GetFamily get font family name
func (s *SubsetFontObj) GetFamily() string {
	return s.Family
}

//SetTtfFontOption set TtfOption must set before SetTTFByPath
func (s *SubsetFontObj) SetTtfFontOption(option TtfOption) {
	s.ttfFontOption = option
}

//KernValueByLeft find kern value from kern table by left
func (s *SubsetFontObj) KernValueByLeft(left uint) (bool, *core.KernValue) {

	if !s.ttfFontOption.UseKerning {
		return false, nil
	}

	k := s.ttfp.Kern()
	if k == nil {
		return false, nil
	}

	if kval, ok := k.Kerning[left]; ok {
		return true, &kval
	}

	return false, nil
}

//SetTTFByPath set ttf
func (s *SubsetFontObj) SetTTFByPath(ttfpath string) error {
	useKerning := s.ttfFontOption.UseKerning
	s.ttfp.SetUseKerning(useKerning)
	err := s.ttfp.Parse(ttfpath)
	if err != nil {
		return err
	}
	return nil
}

//SetTTFByReader set ttf
func (s *SubsetFontObj) SetTTFByReader(rd io.Reader) error {
	useKerning := s.ttfFontOption.UseKerning
	s.ttfp.SetUseKerning(useKerning)
	err := s.ttfp.ParseByReader(rd)
	if err != nil {
		return err
	}
	return nil
}

//AddChars add char to map CharacterToGlyphIndex
func (s *SubsetFontObj) AddChars(txt string) error {
	for _, runeValue := range txt {
		if s.CharacterToGlyphIndex.KeyExists(runeValue) {
			continue
		}
		glyphIndex, err := s.CharCodeToGlyphIndex(runeValue)
		if err != nil {
			return err
		}
		s.CharacterToGlyphIndex.Set(runeValue, glyphIndex) // [runeValue] = glyphIndex
	}
	return nil
}

//CharIndex index of char in glyph table
func (s *SubsetFontObj) CharIndex(r rune) (uint, error) {
	/*
		if index, ok := s.CharacterToGlyphIndex[r]; ok {
			return index, nil
		}
		return 0, ErrCharNotFound
	*/
	glyIndex, ok := s.CharacterToGlyphIndex.Val(r)
	if ok {
		return glyIndex, nil
	}
	return 0, ErrCharNotFound
}

//CharWidth with of char
func (s *SubsetFontObj) CharWidth(r rune) (uint, error) {
	/*glyphIndex := s.CharacterToGlyphIndex
	if index, ok := glyphIndex[r]; ok {
		return s.GlyphIndexToPdfWidth(index), nil
	}
	return 0, ErrCharNotFound*/
	glyIndex, ok := s.CharacterToGlyphIndex.Val(r)
	if ok {
		return s.GlyphIndexToPdfWidth(glyIndex), nil
	}
	return 0, ErrCharNotFound
}

func (s *SubsetFontObj) getType() string {
	return "SubsetFont"
}

func (s *SubsetFontObj) getObjBuff() *bytes.Buffer {
	return &s.buffer
}

func (s *SubsetFontObj) charCodeToGlyphIndexFormat12(r rune) (uint, error) {

	value := uint(r)
	gTbs := s.ttfp.GroupingTables()
	for _, gTb := range gTbs {
		if value >= gTb.StartCharCode && value < gTb.EndCharCode {
			gIndex := (value - gTb.StartCharCode) + gTb.GlyphID
			return gIndex, nil
		}
	}

	return uint(0), errors.New("not found glyph")
}

func (s *SubsetFontObj) charCodeToGlyphIndexFormat4(r rune) (uint, error) {
	value := uint(r)
	seg := uint(0)
	segCount := s.ttfp.SegCount
	for seg < segCount {
		if value <= s.ttfp.EndCount[seg] {
			break
		}
		seg++
	}
	//fmt.Printf("\ncccc--->%#v\n", me.ttfp.Chars())
	if value < s.ttfp.StartCount[seg] {
		return 0, nil
	}

	if s.ttfp.IdRangeOffset[seg] == 0 {

		return (value + s.ttfp.IdDelta[seg]) & 0xFFFF, nil
	}
	//fmt.Printf("IdRangeOffset=%d\n", me.ttfp.IdRangeOffset[seg])
	idx := s.ttfp.IdRangeOffset[seg]/2 + (value - s.ttfp.StartCount[seg]) - (segCount - seg)

	if s.ttfp.GlyphIdArray[int(idx)] == uint(0) {
		return 0, nil
	}

	return (s.ttfp.GlyphIdArray[int(idx)] + s.ttfp.IdDelta[seg]) & 0xFFFF, nil
}

//CharCodeToGlyphIndex get glyph index from char code
func (s *SubsetFontObj) CharCodeToGlyphIndex(r rune) (uint, error) {

	value := uint64(r)
	if value <= 0xFFFF {
		gIndex, err := s.charCodeToGlyphIndexFormat4(r)
		if err != nil {
			return 0, err
		}
		return gIndex, nil
	} else {
		gIndex, err := s.charCodeToGlyphIndexFormat12(r)
		if err != nil {
			return 0, err
		}
		return gIndex, nil
	}

}

//GlyphIndexToPdfWidth get with from glyphIndex
func (s *SubsetFontObj) GlyphIndexToPdfWidth(glyphIndex uint) uint {

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

//GetTTFParser get TTFParser
func (s *SubsetFontObj) GetTTFParser() *core.TTFParser {
	return &s.ttfp
}

//GetUt underlineThickness
func (s *SubsetFontObj) GetUt() int {
	return s.ttfp.UnderlineThickness()
}

//GetUp underline postion
func (s *SubsetFontObj) GetUp() int {
	return s.ttfp.UnderlinePosition()
}

//GetObjBuff get buffer
func (s *SubsetFontObj) GetObjBuff() *bytes.Buffer {
	return s.getObjBuff()
}

//Build build buffer
func (s *SubsetFontObj) Build(objID int) error {
	return s.build(objID)
}
