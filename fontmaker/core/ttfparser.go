package core

import (
	//"encoding/binary"
	//"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var ERROR_NO_UNICODE_ENCODING_FOUND = errors.New("No Unicode encoding found")
var ERROR_UNEXPECTED_SUBTABLE_FORMAT = errors.New("Unexpected subtable format")
var ERROR_INCORRECT_MAGIC_NUMBER = errors.New("Incorrect magic number")
var ERROR_POSTSCRIPT_NAME_NOT_FOUND = errors.New("PostScript name not found")

type TTFParser struct {
	tables map[string]TableDirectoryEntry
	//head
	unitsPerEm       uint64
	xMin             int64
	yMin             int64
	xMax             int64
	yMax             int64
	indexToLocFormat int64
	//Hhea
	numberOfHMetrics uint64
	ascender         int64
	descender        int64
	//end Hhea

	numGlyphs      uint64
	widths         []uint64
	chars          map[int]uint64
	postScriptName string

	//os2
	os2Version    uint64
	Embeddable    bool
	Bold          bool
	typoAscender  int64
	typoDescender int64
	capHeight     int64
	sxHeight      int64

	//post
	italicAngle        int64
	underlinePosition  int64
	underlineThickness int64
	isFixedPitch       bool
	sTypoLineGap       int64
	usWinAscent        uint64
	usWinDescent       uint64

	//cmap
	IsShortIndex  bool
	LocaTable     []uint64
	SegCount      uint64
	StartCount    []uint64
	EndCount      []uint64
	IdRangeOffset []uint64
	IdDelta       []uint64
	GlyphIdArray  []uint64
	symbol        bool
	//data of font
	cahceFontData []byte
}

var Symbolic = 1 << 2
var Nonsymbolic = (1 << 5)

func (t *TTFParser) UnderlinePosition() int64 {
	return t.underlinePosition
}

func (t *TTFParser) UnderlineThickness() int64 {
	return t.underlineThickness
}

func (t *TTFParser) XHeight() int64 {
	if t.os2Version >= 2 && t.sxHeight != 0 {
		return t.sxHeight
	} else {
		return int64((0.66) * float64(t.ascender))
	}
}

func (t *TTFParser) XMin() int64 {
	return t.xMin
}

func (t *TTFParser) YMin() int64 {
	return t.yMin
}

func (t *TTFParser) XMax() int64 {
	return t.xMax
}

func (t *TTFParser) YMax() int64 {
	return t.yMax
}

func (t *TTFParser) ItalicAngle() int64 {
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

func (t *TTFParser) Ascender() int64 {
	if t.typoAscender == 0 {
		return t.ascender
	}
	return int64(t.usWinAscent)
}

func (t *TTFParser) Descender() int64 {
	if t.typoDescender == 0 {
		return t.descender
	}
	descender := int64(t.usWinDescent)
	if t.descender < 0 {
		descender = descender * (-1)
	}
	return descender
}

func (t *TTFParser) TypoAscender() int64 {
	return t.typoAscender
}

func (t *TTFParser) TypoDescender() int64 {
	return t.typoDescender
}

func (t *TTFParser) CapHeight() int64 {
	//fmt.Printf("\n\n>>>>>%d\n\n\n", me.capHeight)
	return t.capHeight
}

func (t *TTFParser) NumGlyphs() uint64 {
	return t.numGlyphs
}

func (t *TTFParser) UnitsPerEm() uint64 {
	return t.unitsPerEm
}

func (t *TTFParser) NumberOfHMetrics() uint64 {
	return t.numberOfHMetrics
}

func (t *TTFParser) Widths() []uint64 {
	return t.widths
}

func (t *TTFParser) Chars() map[int]uint64 {
	return t.chars
}

func (t *TTFParser) GetTables() map[string]TableDirectoryEntry {
	return t.tables
}

func (t *TTFParser) Parse(fontpath string) error {
	//fmt.Printf("\nstart parse\n")
	fd, err := os.Open(fontpath)
	if err != nil {
		return err
	}
	defer fd.Close()
	version, err := t.Read(fd, 4)
	if err != nil {
		return err
	}
	if !t.CompareBytes(version, []byte{0x00, 0x01, 0x00, 0x00}) {
		return errors.New("Unrecognized file (font) format")
	}

	i := uint64(0)
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}
	t.Skip(fd, 3*2) //searchRange, entrySelector, rangeShift
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

		//fmt.Printf("offset\n")
		offset, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		length, err := t.ReadULong(fd)
		if err != nil {
			return err
		}
		//fmt.Printf("\n\ntag=%s  \nOffset = %d\n", tag, offset)
		var table TableDirectoryEntry
		table.Offset = uint64(offset)
		table.CheckSum = checksum
		table.Length = length
		//fmt.Printf("\n\ntag=%s  \nOffset = %d\nPaddedLength =%d\n\n ", tag, table.Offset, table.PaddedLength())
		t.tables[t.BytesToString(tag)] = table
		i++
	}

	//fmt.Printf("%+v\n", me.tables)

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
	//fmt.Printf("%#v\n", me.widths)
	t.cahceFontData, err = t.readFontData(fontpath)
	if err != nil {
		return err
	}

	return nil
}

func (t *TTFParser) FontData() []byte {
	return t.cahceFontData
}

func (t *TTFParser) readFontData(fontpath string) ([]byte, error) {
	b, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (t *TTFParser) ParseLoca(fd *os.File) error {

	t.IsShortIndex = false
	if t.indexToLocFormat == 0 {
		t.IsShortIndex = true
	}

	//fmt.Printf("indexToLocFormat = %d\n", me.indexToLocFormat)
	err := t.Seek(fd, "loca")
	if err != nil {
		return err
	}
	var locaTable []uint64
	table := t.tables["loca"]
	if t.IsShortIndex {
		//do ShortIndex
		entries := table.Length / 2
		i := uint64(0)
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
		i := uint64(0)
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

func (t *TTFParser) ParsePost(fd *os.File) error {

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

	//fmt.Printf("start>>>>>>>\n")
	t.underlineThickness, err = t.ReadShort(fd)
	if err != nil {
		return err
	}
	//fmt.Printf("end>>>>>>>\n")
	//fmt.Printf(">>>>>>>%d\n", me.underlineThickness)

	isFixedPitch, err := t.ReadULong(fd)
	if err != nil {
		return err
	}
	t.isFixedPitch = (isFixedPitch != 0)

	return nil
}

func (t *TTFParser) ParseOS2(fd *os.File) error {
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
	//fmt.Printf("\n\nme.capHeight=%d , me.usWinAscent=%d,me.usWinDescent=%d\n\n", me.capHeight, me.usWinAscent, me.usWinDescent)

	return nil
}

func (t *TTFParser) ParseName(fd *os.File) error {

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
			s := fmt.Sprintf("%s", string(tmpStmp)) //strings(stmp)
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

	//fmt.Printf("%s\n", me.postScriptName)
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

func (t *TTFParser) ParseCmap(fd *os.File) error {
	t.Seek(fd, "cmap")
	t.Skip(fd, 2) // version
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	offset31 := uint64(0)
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
		//fmt.Printf("me.symbol=%d\n", me.symbol)
	} //end for

	if offset31 == 0 {
		//No Unicode encoding found
		return ERROR_NO_UNICODE_ENCODING_FOUND
	}

	var startCount, endCount, idDelta, idRangeOffset, glyphIdArray []uint64

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
	//fmt.Printf("\nlength=%d\n", length)

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
	//fmt.Printf("\nglyphCount=%d\n", glyphCount)

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
		glyphIdArray = append(glyphIdArray, tmp)
	}
	t.GlyphIdArray = glyphIdArray

	t.chars = make(map[int]uint64)
	for i := 0; i < int(segCount); i++ {
		c1 := startCount[i]
		c2 := endCount[i]
		d := idDelta[i]
		ro := idRangeOffset[i]
		if ro > 0 {
			_, err = fd.Seek(int64(offset+uint64(2*i)+ro), 0)
			if err != nil {
				return err
			}
		}

		for c := c1; c <= c2; c++ {
			var gid uint64
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
				//fmt.Printf("%d gid = %d, ", int(c), gid)
				t.chars[int(c)] = gid
			}
		}

	}
	//fmt.Printf("len() = %d , me.chars[10] = %d , me.chars[56]  = %d \n", len(me.chars), me.chars[10], me.chars[56])
	//fmt.Printf("len() = %d , me.chars[99] = %d , me.chars[107]  = %d \n\n", len(me.chars), me.chars[99], me.chars[107])
	return nil
}

func (t *TTFParser) FTell(fd *os.File) (uint64, error) {
	offset, err := fd.Seek(0, os.SEEK_CUR)
	return uint64(offset), err
}

func (t *TTFParser) ParseHmtx(fd *os.File) error {

	t.Seek(fd, "hmtx")
	i := uint64(0)
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
		var err error
		lastWidth := t.widths[t.numberOfHMetrics-1]
		t.widths, err = t.ArrayPadUint(t.widths, t.numGlyphs, lastWidth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TTFParser) ArrayPadUint(arr []uint64, size uint64, val uint64) ([]uint64, error) {
	var result []uint64
	i := uint64(0)
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

func (t *TTFParser) ParseHead(fd *os.File) error {

	//fmt.Printf("\nParseHead\n")
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

	//fmt.Printf("\nmagicNumber = %d\n", magicNumber)
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

func (t *TTFParser) ParseHhea(fd *os.File) error {

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

	//fmt.Printf("---------me.numberOfHMetrics=%d,me.ascender=%d,me.descender = %d\n\n", me.numberOfHMetrics, me.ascender, me.descender)
	return nil
}

func (t *TTFParser) ParseMaxp(fd *os.File) error {
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

func (t *TTFParser) Seek(fd *os.File, tag string) error {
	table, ok := t.tables[tag]
	if !ok {
		return errors.New("me.tables not contain key=" + tag)
	}
	val := table.Offset
	_, err := fd.Seek(int64(val), 0)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTFParser) BytesToString(b []byte) string {
	return string(b) //strings.TrimSpace(string(b))
}

func (t *TTFParser) ReadUShort(fd *os.File) (uint64, error) {
	buff, err := t.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (t *TTFParser) ReadShort(fd *os.File) (int64, error) {
	buff, err := t.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	u := num.Uint64()

	//fmt.Printf("%#v\n", buff)
	var v int64
	if u >= 0x8000 {
		v = int64(u) - 65536
	} else {
		v = int64(u)
	}
	return v, nil
}

func (t *TTFParser) ReadULong(fd *os.File) (uint64, error) {
	buff, err := t.Read(fd, 4)
	//fmt.Printf("%#v\n", buff)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (t *TTFParser) Skip(fd *os.File, length int64) error {
	_, err := fd.Seek(int64(length), 1)
	if err != nil {
		return err
	}
	return nil
}

func (t *TTFParser) Read(fd *os.File, length int) ([]byte, error) {
	buff := make([]byte, length)
	readlength, err := fd.Read(buff)
	if err != nil {
		return nil, err
	}
	if readlength != length {
		return nil, errors.New("file out of length")
	}
	//fmt.Printf("%d,%s\n", readlength, string(buff))
	return buff, nil
}

func (t *TTFParser) CompareBytes(a []byte, b []byte) bool {

	if a == nil && b == nil {
		return true
	} else if a == nil && b != nil {
		return false
	} else if a != nil && b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	i := 0
	length := len(a)
	for i < length {
		if a[i] != b[i] {
			return false
		}
		i++
	}
	return true
}
