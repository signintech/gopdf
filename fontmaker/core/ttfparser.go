package core

import (
	//"encoding/binary"
	//"encoding/hex"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ERROR_NO_UNICODE_ENCODING_FOUND = errors.New("No Unicode encoding found")
var ERROR_UNEXPECTED_SUBTABLE_FORMAT = errors.New("Unexpected subtable format")
var ERROR_INCORRECT_MAGIC_NUMBER = errors.New("Incorrect magic number")
var ERROR_POSTSCRIPT_NAME_NOT_FOUND = errors.New("PostScript name not found")

// TTFParser true type font parser
type TTFParser struct {
	tables map[string]TableDirectoryEntry
	//head
	unitsPerEm       uint
	xMin             int
	yMin             int
	xMax             int
	yMax             int
	indexToLocFormat int
	//Hhea
	numberOfHMetrics uint
	ascender         int
	descender        int
	//end Hhea

	numGlyphs      uint
	widths         []uint
	chars          map[int]uint
	postScriptName string

	//os2
	os2Version    uint
	Embeddable    bool
	Bold          bool
	typoAscender  int
	typoDescender int
	capHeight     int
	sxHeight      int

	//post
	italicAngle        int
	underlinePosition  int
	underlineThickness int
	isFixedPitch       bool
	sTypoLineGap       int
	usWinAscent        uint
	usWinDescent       uint

	//cmap
	IsShortIndex  bool
	LocaTable     []uint
	SegCount      uint
	StartCount    []uint
	EndCount      []uint
	IdRangeOffset []uint
	IdDelta       []uint
	GlyphIdArray  []uint
	symbol        bool

	//cmap format 12
	groupingTables []CmapFormat12GroupingTable

	//data of font
	cachedFontData []byte

	//kerning
	useKerning bool //user config for use or not use kerning
	kern       *KernTable
}

var Symbolic = 1 << 2
var Nonsymbolic = (1 << 5)

// Kern get KernTable
func (t *TTFParser) Kern() *KernTable {
	return t.kern
}

// UnderlinePosition position of underline
func (t *TTFParser) UnderlinePosition() int {
	return t.underlinePosition
}

// GroupingTables get cmap format12 grouping table
func (t *TTFParser) GroupingTables() []CmapFormat12GroupingTable {
	return t.groupingTables
}

// UnderlineThickness thickness of underline
func (t *TTFParser) UnderlineThickness() int {
	return t.underlineThickness
}

func (t *TTFParser) XHeight() int {
	if t.os2Version >= 2 && t.sxHeight != 0 {
		return t.sxHeight
	} else {
		return int((0.66) * float64(t.ascender))
	}
}

func (t *TTFParser) XMin() int {
	return t.xMin
}

func (t *TTFParser) YMin() int {
	return t.yMin
}

func (t *TTFParser) XMax() int {
	return t.xMax
}

func (t *TTFParser) YMax() int {
	return t.yMax
}

func (t *TTFParser) ItalicAngle() int {
	return t.italicAngle
}

func (t *TTFParser) Flag() int {
	flag := 0
	if t.symbol {
		flag |= Symbolic
	} else {
		flag |= Nonsymbolic
	}
	return flag
}

func (t *TTFParser) Ascender() int {
	if t.typoAscender == 0 {
		return t.ascender
	}
	return int(t.usWinAscent)
}

func (t *TTFParser) Descender() int {
	if t.typoDescender == 0 {
		return t.descender
	}
	descender := int(t.usWinDescent)
	if t.descender < 0 {
		descender = descender * (-1)
	}
	return descender
}

func (t *TTFParser) TypoAscender() int {
	return t.typoAscender
}

func (t *TTFParser) TypoDescender() int {
	return t.typoDescender
}

// CapHeight https://en.wikipedia.org/wiki/Cap_height
func (t *TTFParser) CapHeight() int {
	return t.capHeight
}

// NumGlyphs number of glyph
func (t *TTFParser) NumGlyphs() uint {
	return t.numGlyphs
}

func (t *TTFParser) UnitsPerEm() uint {
	return t.unitsPerEm
}

func (t *TTFParser) NumberOfHMetrics() uint {
	return t.numberOfHMetrics
}

func (t *TTFParser) Widths() []uint {
	return t.widths
}

func (t *TTFParser) Chars() map[int]uint {
	return t.chars
}

func (t *TTFParser) GetTables() map[string]TableDirectoryEntry {
	return t.tables
}

// SetUseKerning set useKerning must set before Parse
func (t *TTFParser) SetUseKerning(use bool) {
	t.useKerning = use
}

// Parse parse
func (t *TTFParser) Parse(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	return t.ParseFontData(data)
}

// ParseByReader parse by io.reader
func (t *TTFParser) ParseByReader(rd io.Reader) error {
	fontData, err := io.ReadAll(rd)
	if err != nil {
		return err
	}

	return t.ParseFontData(fontData)
}

// ParseFontData parses font data.
func (t *TTFParser) ParseFontData(fontData []byte) error {
	fd := bytes.NewReader(fontData)

	version, err := t.Read(fd, 4)
	if err != nil {
		return err
	}
	if !bytes.Equal(version, []byte{0x00, 0x01, 0x00, 0x00}) {
		return errors.New("Unrecognized file (font) format")
	}

	i := uint(0)
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	err = t.Skip(fd, 3*2) //searchRange, entrySelector, rangeShift
	if err != nil {
		return err
	}

	t.tables = make(map[string]TableDirectoryEntry)
	for i < numTables {

		tag, err := t.Read(fd, 4)
		if err != nil {
			return err
		}

		checksum, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		offset, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		length, err := t.ReadULong(fd)
		if err != nil {
			return err
		}
		var table TableDirectoryEntry
		table.Offset = uint(offset)
		table.CheckSum = checksum
		table.Length = length
		t.tables[t.BytesToString(tag)] = table
		i++
	}

	err = t.ParseHead(fd)
	if err != nil {
		return err
	}

	err = t.ParseHhea(fd)
	if err != nil {
		return err
	}

	err = t.ParseMaxp(fd)
	if err != nil {
		return err
	}
	err = t.ParseHmtx(fd)
	if err != nil {
		return err
	}
	err = t.ParseCmap(fd)
	if err != nil {
		return err
	}
	err = t.ParseName(fd)
	if err != nil {
		return err
	}
	err = t.ParseOS2(fd)
	if err != nil {
		return err
	}
	err = t.ParsePost(fd)
	if err != nil {
		return err
	}
	err = t.ParseLoca(fd)
	if err != nil {
		return err
	}

	if t.useKerning {
		err = t.Parsekern(fd)
		if err != nil {
			return err
		}
	}

	t.cachedFontData = fontData

	return nil
}

func (t *TTFParser) FontData() []byte {
	return t.cachedFontData
}

// ParseLoca parse loca table https://www.microsoft.com/typography/otspec/loca.htm
func (t *TTFParser) ParseLoca(fd *bytes.Reader) error {

	t.IsShortIndex = false
	if t.indexToLocFormat == 0 {
		t.IsShortIndex = true
	}

	err := t.Seek(fd, "loca")
	if err != nil {
		return err
	}
	var locaTable []uint
	table := t.tables["loca"]
	if t.IsShortIndex {
		//do ShortIndex
		entries := table.Length / 2
		i := uint(0)
		for i < entries {
			item, err := t.ReadUShort(fd)
			if err != nil {
				return err
			}
			locaTable = append(locaTable, item*2)
			i++
		}
	} else {
		entries := table.Length / 4
		i := uint(0)
		for i < entries {
			item, err := t.ReadULong(fd)
			if err != nil {
				return err
			}
			locaTable = append(locaTable, item)
			i++
		}
	}
	t.LocaTable = locaTable
	return nil
}

// ParsePost parse post table https://www.microsoft.com/typography/otspec/post.htm
func (t *TTFParser) ParsePost(fd *bytes.Reader) error {

	err := t.Seek(fd, "post")
	if err != nil {
		return err
	}

	err = t.Skip(fd, 4) // version
	if err != nil {
		return err
	}

	t.italicAngle, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 2) // Skip decimal part
	if err != nil {
		return err
	}

	t.underlinePosition, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.underlineThickness, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	isFixedPitch, err := t.ReadULong(fd)
	if err != nil {
		return err
	}
	t.isFixedPitch = (isFixedPitch != 0)

	return nil
}

// ParseOS2 parse OS2 table https://www.microsoft.com/typography/otspec/OS2.htm
func (t *TTFParser) ParseOS2(fd *bytes.Reader) error {
	err := t.Seek(fd, "OS/2")
	if err != nil {
		return err
	}
	version, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.os2Version = version

	err = t.Skip(fd, 3*2) // xAvgCharWidth, usWeightClass, usWidthClass
	if err != nil {
		return err
	}
	fsType, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.Embeddable = (fsType != 2) && ((fsType & 0x200) == 0)

	err = t.Skip(fd, (11*2)+10+(4*4)+4)
	if err != nil {
		return err
	}
	fsSelection, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.Bold = ((fsSelection & 32) != 0)
	err = t.Skip(fd, 2*2) // usFirstCharIndex, usLastCharIndex
	if err != nil {
		return err
	}
	t.typoAscender, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.typoDescender, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.sTypoLineGap, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.usWinAscent, err = t.ReadUShort(fd)
	if err != nil {
		return err
	}

	t.usWinDescent, err = t.ReadUShort(fd)
	if err != nil {
		return err
	}

	if version >= 2 {

		err = t.Skip(fd, 2*4)
		if err != nil {
			return err
		}

		t.sxHeight, err = t.ReadShort(fd)
		if err != nil {
			return err
		}

		t.capHeight, err = t.ReadShort(fd)
		if err != nil {
			return err
		}

	} else {
		t.capHeight = t.ascender
	}

	return nil
}

// ParseName parse name table https://www.microsoft.com/typography/otspec/name.htm
func (t *TTFParser) ParseName(fd *bytes.Reader) error {

	//$this->Seek('name');
	err := t.Seek(fd, "name")
	if err != nil {
		return err
	}

	tableOffset, err := t.FTell(fd)
	if err != nil {
		return err
	}

	t.postScriptName = ""
	err = t.Skip(fd, 2) // format
	if err != nil {
		return err
	}

	count, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	stringOffset, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	for i := 0; i < int(count); i++ {
		err = t.Skip(fd, 3*2) // platformID, encodingID, languageID
		if err != nil {
			return err
		}
		nameID, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		length, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		offset, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		if nameID == 6 {
			// PostScript name
			_, err = fd.Seek(int64(tableOffset+stringOffset+offset), 0)
			if err != nil {
				return err
			}

			stmp, err := t.Read(fd, int(length))
			if err != nil {
				return err
			}

			var tmpStmp []byte
			for _, v := range stmp {
				if v != 0 {
					tmpStmp = append(tmpStmp, v)
				}
			}
			//s := fmt.Sprintf("%s", string(tmpStmp)) //strings(stmp)
			s := string(tmpStmp)
			s = strings.Replace(s, strconv.Itoa(0), "", -1)
			s, err = t.PregReplace("|[ \\[\\](){}<>/%]|", "", s)
			if err != nil {
				return err
			}
			t.postScriptName = s
			break
		}
	}

	if t.postScriptName == "" {
		return ERROR_POSTSCRIPT_NAME_NOT_FOUND
	}

	return nil
}

func (t *TTFParser) PregReplace(pattern string, replacement string, subject string) (string, error) {

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	str := reg.ReplaceAllString(subject, replacement)
	return str, nil
}

// ParseCmap parse cmap table format 4 https://www.microsoft.com/typography/otspec/cmap.htm
func (t *TTFParser) ParseCmap(fd *bytes.Reader) error {
	err := t.Seek(fd, "cmap")
	if err != nil {
		return err
	}
	err = t.Skip(fd, 2) // version
	if err != nil {
		return err
	}
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	offset31 := uint(0)
	for i := 0; i < int(numTables); i++ {
		platformID, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		encodingID, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		offset, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		t.symbol = false //init
		if platformID == 3 && encodingID == 1 {
			if encodingID == 0 {
				t.symbol = true
			}
			offset31 = offset
		}
	} //end for

	if offset31 == 0 {
		//No Unicode encoding found
		return ERROR_NO_UNICODE_ENCODING_FOUND
	}

	var startCount, endCount, idDelta, idRangeOffset, glyphIDArray []uint

	_, err = fd.Seek(int64(t.tables["cmap"].Offset+offset31), 0)
	if err != nil {
		return err
	}

	format, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	if format != 4 {
		//Unexpected subtable format
		return ERROR_UNEXPECTED_SUBTABLE_FORMAT
	}

	length, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 2) // language
	if err != nil {
		return err
	}
	segCount, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	segCount = segCount / 2
	t.SegCount = segCount
	err = t.Skip(fd, 3*2) // searchRange, entrySelector, rangeShift
	if err != nil {
		return err
	}

	glyphCount := (length - (16 + 8*segCount)) / 2

	for i := 0; i < int(segCount); i++ {
		tmp, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		endCount = append(endCount, tmp)
	}
	t.EndCount = endCount

	err = t.Skip(fd, 2) // reservedPad
	if err != nil {
		return err
	}

	for i := 0; i < int(segCount); i++ {
		tmp, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		startCount = append(startCount, tmp)
	}
	t.StartCount = startCount

	for i := 0; i < int(segCount); i++ {
		tmp, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		idDelta = append(idDelta, tmp)
	}
	t.IdDelta = idDelta

	offset, err := t.FTell(fd)
	if err != nil {
		return err
	}
	for i := 0; i < int(segCount); i++ {
		tmp, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		idRangeOffset = append(idRangeOffset, tmp)
	}
	t.IdRangeOffset = idRangeOffset
	//_ = glyphIdArray
	for i := 0; i < int(glyphCount); i++ {
		tmp, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		glyphIDArray = append(glyphIDArray, tmp)
	}
	t.GlyphIdArray = glyphIDArray

	t.chars = make(map[int]uint)
	for i := 0; i < int(segCount); i++ {
		c1 := startCount[i]
		c2 := endCount[i]
		d := idDelta[i]
		ro := idRangeOffset[i]
		if ro > 0 {
			_, err = fd.Seek(int64(offset+uint(2*i)+ro), 0)
			if err != nil {
				return err
			}
		}

		for c := c1; c <= c2; c++ {
			var gid uint
			if c == 0xFFFF {
				break
			}
			if ro > 0 {
				gid, err = t.ReadUShort(fd)
				if err != nil {
					return err
				}
				if gid > 0 {
					gid += d
				}
			} else {
				gid = c + d
			}

			if gid >= 65536 {
				gid -= 65536
			}
			if gid > 0 {
				t.chars[int(c)] = gid
			}
		}

	}

	_, err = t.ParseCmapFormat12(fd)
	if err != nil {
		return err
	}

	return nil
}

func (t *TTFParser) FTell(fd *bytes.Reader) (uint, error) {
	offset, err := fd.Seek(0, io.SeekCurrent)
	return uint(offset), err
}

// ParseHmtx parse hmtx table  https://www.microsoft.com/typography/otspec/hmtx.htm
func (t *TTFParser) ParseHmtx(fd *bytes.Reader) error {
	err := t.Seek(fd, "hmtx")
	if err != nil {
		return err
	}

	i := uint(0)
	for i < t.numberOfHMetrics {
		advanceWidth, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		err = t.Skip(fd, 2)
		if err != nil {
			return err
		}
		t.widths = append(t.widths, advanceWidth)
		i++
	}
	if t.numberOfHMetrics < t.numGlyphs {
		lastWidth := t.widths[t.numberOfHMetrics-1]
		t.widths, err = t.ArrayPadUint(t.widths, t.numGlyphs, lastWidth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TTFParser) ArrayPadUint(arr []uint, size uint, val uint) ([]uint, error) {
	var result []uint
	i := uint(0)
	for i < size {
		if int(i) < len(arr) {
			result = append(result, arr[i])
		} else {
			result = append(result, val)
		}
		i++
	}

	return result, nil
}

// ParseHead parse head table  https://www.microsoft.com/typography/otspec/Head.htm
func (t *TTFParser) ParseHead(fd *bytes.Reader) error {

	err := t.Seek(fd, "head")
	if err != nil {
		return err
	}

	err = t.Skip(fd, 3*4) // version, fontRevision, checkSumAdjustment
	if err != nil {
		return err
	}
	magicNumber, err := t.ReadULong(fd)
	if err != nil {
		return err
	}

	if magicNumber != 0x5F0F3CF5 {
		return ERROR_INCORRECT_MAGIC_NUMBER
	}

	err = t.Skip(fd, 2)
	if err != nil {
		return err
	}

	t.unitsPerEm, err = t.ReadUShort(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 2*8) // created, modified
	if err != nil {
		return err
	}

	t.xMin, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.yMin, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.xMax, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.yMax, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 2*3) //skip macStyle,lowestRecPPEM,fontDirectionHint
	if err != nil {
		return err
	}

	t.indexToLocFormat, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	return nil
}

// ParseHhea parse hhea table  https://www.microsoft.com/typography/otspec/hhea.htm
func (t *TTFParser) ParseHhea(fd *bytes.Reader) error {

	err := t.Seek(fd, "hhea")
	if err != nil {
		return err
	}

	err = t.Skip(fd, 4) //skip version
	if err != nil {
		return err
	}

	t.ascender, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	t.descender, err = t.ReadShort(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 13*2)
	if err != nil {
		return err
	}

	t.numberOfHMetrics, err = t.ReadUShort(fd)
	if err != nil {
		return err
	}

	return nil
}

// ParseMaxp parse maxp table  https://www.microsoft.com/typography/otspec/Maxp.htm
func (t *TTFParser) ParseMaxp(fd *bytes.Reader) error {
	err := t.Seek(fd, "maxp")
	if err != nil {
		return err
	}
	err = t.Skip(fd, 4)
	if err != nil {
		return err
	}
	t.numGlyphs, err = t.ReadUShort(fd)
	if err != nil {
		return err
	}
	return nil
}

// ErrTableNotFound error table not found
var ErrTableNotFound = errors.New("table not found")

// Seek seek by tag
func (t *TTFParser) Seek(fd *bytes.Reader, tag string) error {
	table, ok := t.tables[tag]
	if !ok {
		return ErrTableNotFound
	}
	val := table.Offset
	_, err := fd.Seek(int64(val), 0)
	if err != nil {
		return err
	}
	return nil
}

// BytesToString convert bytes to string
func (t *TTFParser) BytesToString(b []byte) string {
	return string(b) //strings.TrimSpace(string(b))
}

// ReadUShort read ushort
func (t *TTFParser) ReadUShort(fd *bytes.Reader) (uint, error) {
	buff, err := t.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint16(buff)
	return uint(n), nil
}

// ReadShort read short
func (t *TTFParser) ReadShort(fd *bytes.Reader) (int, error) {
	u, err := t.ReadUShort(fd)
	if err != nil {
		return 0, err
	}

	var v int
	if u >= 0x8000 {
		v = int(u) - 65536
	} else {
		v = int(u)
	}
	return v, nil
}

// ReadShortInt16 read short return int16
func (t *TTFParser) ReadShortInt16(fd *bytes.Reader) (int16, error) {
	n, err := t.ReadShort(fd)
	if err != nil {
		return 0, err
	}
	return int16(n), nil
}

// ReadULong read ulong
func (t *TTFParser) ReadULong(fd *bytes.Reader) (uint, error) {
	buff, err := t.Read(fd, 4)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint32(buff)
	return uint(n), nil
}

// Skip skip
func (t *TTFParser) Skip(fd *bytes.Reader, length int) error {
	_, err := fd.Seek(int64(length), 1)
	if err != nil {
		return err
	}
	return nil
}

// Read read
func (t *TTFParser) Read(fd *bytes.Reader, length int) ([]byte, error) {
	buff := make([]byte, length)
	readlength, err := fd.Read(buff)
	if err != nil {
		return nil, err
	}
	if readlength != length {
		return nil, errors.New("file out of length")
	}
	return buff, nil
}
