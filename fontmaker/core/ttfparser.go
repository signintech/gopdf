package core

import (
	//"encoding/binary"
	//"encoding/hex"
	"errors"
	"fmt"
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
	tables             map[string]uint64
	unitsPerEm         uint64
	xMin               int64
	yMin               int64
	xMax               int64
	yMax               int64
	numberOfHMetrics   uint64
	numGlyphs          uint64
	widths             []uint64
	chars              map[int]uint64
	postScriptName     string
	Embeddable         bool
	Bold               bool
	typoAscender       int64
	typoDescender      int64
	capHeight          int64
	italicAngle        int64
	underlinePosition  int64
	underlineThickness int64
	isFixedPitch       bool
}

func (me *TTFParser) Parse(fontpath string) error {
	//fmt.Printf("\nstart parse\n")
	fd, err := os.Open(fontpath)
	if err != nil {
		return err
	}
	defer fd.Close()
	version, err := me.Read(fd, 4)
	if err != nil {
		return err
	}

	if !me.CompareBytes(version, []byte{0x00, 0x01, 0x00, 0x00}) {
		return errors.New("Unrecognized file (font) format")
	}

	i := uint64(0)
	numTables, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	me.Skip(fd, 3*2) //searchRange, entrySelector, rangeShift
	me.tables = make(map[string]uint64)
	for i < numTables {

		tag, err := me.Read(fd, 4)
		if err != nil {
			return err
		}

		err = me.Skip(fd, 4)
		if err != nil {
			return err
		}

		offset, err := me.ReadULong(fd)
		if err != nil {
			return err
		}

		err = me.Skip(fd, 4)
		if err != nil {
			return err
		}
		//fmt.Printf("%s\n", me.BytesToString(tag))
		me.tables[me.BytesToString(tag)] = offset
		i++
	}

	//fmt.Printf("%+v\n", me.tables)

	err = me.ParseHead(fd)
	if err != nil {
		return err
	}

	err = me.ParseHhea(fd)
	if err != nil {
		return err
	}

	err = me.ParseMaxp(fd)
	if err != nil {
		return err
	}
	err = me.ParseHmtx(fd)
	if err != nil {
		return err
	}
	err = me.ParseCmap(fd)
	if err != nil {
		return err
	}
	err = me.ParseName(fd)
	if err != nil {
		return err
	}
	err = me.ParseOS2(fd)
	if err != nil {
		return err
	}
	err = me.ParsePost(fd)
	if err != nil {
		return err
	}
	//fmt.Printf("%#v\n", me.widths)
	return nil
}

func (me *TTFParser) ParsePost(fd *os.File) error {

	err := me.Seek(fd, "post")
	if err != nil {
		return err
	}

	err = me.Skip(fd, 4) // version
	if err != nil {
		return err
	}

	me.italicAngle, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	err = me.Skip(fd, 2) // Skip decimal part
	if err != nil {
		return err
	}

	me.underlinePosition, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	//fmt.Printf("start>>>>>>>\n")
	me.underlineThickness, err = me.ReadShort(fd)
	if err != nil {
		return err
	}
	//fmt.Printf("end>>>>>>>\n")
	//fmt.Printf(">>>>>>>%d\n", me.underlineThickness)

	isFixedPitch, err := me.ReadULong(fd)
	if err != nil {
		return err
	}
	me.isFixedPitch = (isFixedPitch != 0)

	return nil
}

func (me *TTFParser) ParseOS2(fd *os.File) error {
	err := me.Seek(fd, "OS/2")
	if err != nil {
		return err
	}
	version, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	err = me.Skip(fd, 3*2) // xAvgCharWidth, usWeightClass, usWidthClass
	if err != nil {
		return err
	}
	fsType, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	me.Embeddable = (fsType != 2) && ((fsType & 0x200) == 0)

	err = me.Skip(fd, (11*2)+10+(4*4)+4)
	if err != nil {
		return err
	}
	fsSelection, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	me.Bold = ((fsSelection & 32) != 0)
	err = me.Skip(fd, 2*2) // usFirstCharIndex, usLastCharIndex
	if err != nil {
		return err
	}
	me.typoAscender, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	me.typoDescender, err = me.ReadShort(fd)
	if err != nil {
		return err
	}
	if version >= 2 {
		err = me.Skip(fd, 3*2+2*4+2)
		if err != nil {
			return err
		}
		me.capHeight, err = me.ReadShort(fd)
		if err != nil {
			return err
		}
	} else {
		me.capHeight = 0
	}

	return nil
}

func (me *TTFParser) ParseName(fd *os.File) error {

	//$this->Seek('name');
	err := me.Seek(fd, "name")
	if err != nil {
		return err
	}

	tableOffset, err := me.FTell(fd)
	if err != nil {
		return err
	}

	me.postScriptName = ""
	err = me.Skip(fd, 2) // format
	if err != nil {
		return err
	}

	count, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}

	stringOffset, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}

	for i := 0; i < int(count); i++ {
		err = me.Skip(fd, 3*2) // platformID, encodingID, languageID
		if err != nil {
			return err
		}
		nameID, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		length, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		offset, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		if nameID == 6 {
			// PostScript name
			_, err = fd.Seek(int64(tableOffset+stringOffset+offset), 0)
			if err != nil {
				return err
			}

			stmp, err := me.Read(fd, int(length))
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
			s, err = me.PregReplace("|[ \\[\\](){}<>/%]|", "", s)
			if err != nil {
				return err
			}
			me.postScriptName = s
			break
		}
	}

	if me.postScriptName == "" {
		return ERROR_POSTSCRIPT_NAME_NOT_FOUND
	}

	//fmt.Printf("%s\n", me.postScriptName)
	return nil
}

func (me *TTFParser) PregReplace(pattern string, replacement string, subject string) (string, error) {

	reg, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	str := reg.ReplaceAllString(subject, replacement)
	return str, nil
}

func (me *TTFParser) ParseCmap(fd *os.File) error {
	me.Seek(fd, "cmap")
	me.Skip(fd, 2) // version
	numTables, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}

	offset31 := uint64(0)
	for i := 0; i < int(numTables); i++ {
		platformID, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		encodingID, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		offset, err := me.ReadULong(fd)
		if err != nil {
			return err
		}
		if platformID == 3 && encodingID == 1 {
			offset31 = offset
		}
	} //end for

	if offset31 == 0 {
		//No Unicode encoding found
		return ERROR_NO_UNICODE_ENCODING_FOUND
	}

	var startCount, endCount, idDelta, idRangeOffset []uint64

	_, err = fd.Seek(int64(me.tables["cmap"]+offset31), 0)
	if err != nil {
		return err
	}

	format, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}

	if format != 4 {
		//Unexpected subtable format
		return ERROR_UNEXPECTED_SUBTABLE_FORMAT
	}

	err = me.Skip(fd, 2*2) // length, language
	if err != nil {
		return err
	}
	segCount, err := me.ReadUShort(fd)
	if err != nil {
		return err
	}
	segCount = segCount / 2
	err = me.Skip(fd, 3*2) // searchRange, entrySelector, rangeShift
	if err != nil {
		return err
	}

	for i := 0; i < int(segCount); i++ {
		tmp, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		endCount = append(endCount, tmp)
	}

	err = me.Skip(fd, 2) // reservedPad
	if err != nil {
		return err
	}

	for i := 0; i < int(segCount); i++ {
		tmp, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		startCount = append(startCount, tmp)
	}

	for i := 0; i < int(segCount); i++ {
		tmp, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		idDelta = append(idDelta, tmp)
	}

	offset, err := me.FTell(fd)
	if err != nil {
		return err
	}
	for i := 0; i < int(segCount); i++ {
		tmp, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		idRangeOffset = append(idRangeOffset, tmp)
	}
	//fmt.Printf("%d\n\n\n", offset)
	me.chars = make(map[int]uint64)
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
				gid, err = me.ReadUShort(fd)
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
				me.chars[int(c)] = gid
			}
		}
		//fmt.Printf("%#v\n\n", me.chars)
	}
	return nil
}

func (me *TTFParser) FTell(fd *os.File) (uint64, error) {
	offset, err := fd.Seek(0, os.SEEK_CUR)
	return uint64(offset), err
}

func (me *TTFParser) ParseHmtx(fd *os.File) error {

	me.Seek(fd, "hmtx")
	i := uint64(0)
	for i < me.numberOfHMetrics {
		advanceWidth, err := me.ReadUShort(fd)
		if err != nil {
			return err
		}
		err = me.Skip(fd, 2)
		if err != nil {
			return err
		}
		me.widths = append(me.widths, advanceWidth)
		i++
	}
	if me.numberOfHMetrics < me.numGlyphs {
		var err error
		lastWidth := me.widths[me.numberOfHMetrics-1]
		me.widths, err = me.ArrayPadUint(me.widths, me.numGlyphs, lastWidth)
		if err != nil {
			return err
		}
	}

	return nil
}

func (me *TTFParser) ArrayPadUint(arr []uint64, size uint64, val uint64) ([]uint64, error) {
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

func (me *TTFParser) ParseHead(fd *os.File) error {

	//fmt.Printf("\nParseHead\n")
	err := me.Seek(fd, "head")
	if err != nil {
		return err
	}

	err = me.Skip(fd, 3*4) // version, fontRevision, checkSumAdjustment
	if err != nil {
		return err
	}
	magicNumber, err := me.ReadULong(fd)
	if err != nil {
		return err
	}

	//fmt.Printf("\nmagicNumber = %d\n", magicNumber)
	if magicNumber != 0x5F0F3CF5 {
		return ERROR_INCORRECT_MAGIC_NUMBER
	}

	err = me.Skip(fd, 2)
	if err != nil {
		return err
	}

	me.unitsPerEm, err = me.ReadUShort(fd)
	if err != nil {
		return err
	}

	err = me.Skip(fd, 2*8) // created, modified
	if err != nil {
		return err
	}

	me.xMin, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	me.yMin, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	me.xMax, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	me.yMax, err = me.ReadShort(fd)
	if err != nil {
		return err
	}

	return nil
}

func (me *TTFParser) ParseHhea(fd *os.File) error {

	err := me.Seek(fd, "hhea")
	if err != nil {
		return err
	}

	err = me.Skip(fd, 4+15*2)
	if err != nil {
		return err
	}

	me.numberOfHMetrics, err = me.ReadUShort(fd)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) ParseMaxp(fd *os.File) error {
	err := me.Seek(fd, "maxp")
	if err != nil {
		return err
	}
	err = me.Skip(fd, 4)
	if err != nil {
		return err
	}
	me.numGlyphs, err = me.ReadUShort(fd)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) Seek(fd *os.File, tag string) error {
	val, ok := me.tables[tag]
	if !ok {
		return errors.New("me.tables not contain key=" + tag)
	}
	_, err := fd.Seek(int64(val), 0)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) BytesToString(b []byte) string {
	return string(b)
}

func (me *TTFParser) ReadUShort(fd *os.File) (uint64, error) {
	buff, err := me.Read(fd, 2)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (me *TTFParser) ReadShort(fd *os.File) (int64, error) {
	buff, err := me.Read(fd, 2)
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

func (me *TTFParser) ReadULong(fd *os.File) (uint64, error) {
	buff, err := me.Read(fd, 4)
	if err != nil {
		return 0, err
	}
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), nil
}

func (me *TTFParser) Skip(fd *os.File, length int64) error {
	_, err := fd.Seek(int64(length), 1)
	if err != nil {
		return err
	}
	return nil
}

func (me *TTFParser) Read(fd *os.File, length int) ([]byte, error) {
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

func (me *TTFParser) CompareBytes(a []byte, b []byte) bool {

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
