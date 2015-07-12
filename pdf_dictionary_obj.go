package gopdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/signintech/gopdf/fontmaker/core"
)

var EntrySelectors = []int{0, 0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4}

type PdfDictionaryObj struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (me *PdfDictionaryObj) Init(funcGetRoot func() *GoPdf) {
}

func (me *PdfDictionaryObj) Build() {
	b, err := me.makeFont()
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	//zipvar buff bytes.Buffer
	var zbuff bytes.Buffer
	gzipwriter := zlib.NewWriter(&zbuff)
	_, err = gzipwriter.Write(b)
	if err != nil {
		log.Panicf("%s", err.Error())
		return
	}
	gzipwriter.Close()

	//fmt.Printf("\n%d\n", len(by))
	me.buffer.WriteString("<</Length " + strconv.Itoa(zbuff.Len()) + "\n")
	me.buffer.WriteString("/Filter /FlateDecode\n")
	me.buffer.WriteString("/Length1 " + strconv.Itoa(len(b)) + "\n")
	me.buffer.WriteString(">>\n")
	me.buffer.WriteString("stream\n")
	me.buffer.Write(zbuff.Bytes())
	me.buffer.WriteString("\nendstream\n")
}

func (me *PdfDictionaryObj) GetType() string {
	return "PdfDictionary"
}

func (me *PdfDictionaryObj) GetObjBuff() *bytes.Buffer {
	return &me.buffer
}

func (me *PdfDictionaryObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}

func (me *PdfDictionaryObj) makeGlyfAndLocaTable() ([]byte, []int, error) {
	ttfp := me.PtrToSubsetFontObj.GetTTFParser()
	var glyf core.TableDirectoryEntry //ttfp.GetTables()["glyf"]

	numGlyphs := int(ttfp.NumGlyphs())

	glyphs := me.completeGlyphClosure(me.PtrToSubsetFontObj.CharacterToGlyphIndex)
	glyphCount := len(glyphs)

	//copy
	var glyphArray []int
	for _, v := range me.PtrToSubsetFontObj.CharacterToGlyphIndex {
		glyphArray = append(glyphArray, int(v))
	}
	sort.Ints(glyphArray)

	size := 0
	for idx := 0; idx < glyphCount; idx++ {
		size += me.getGlyphSize(glyphArray[idx])
	}
	glyf.Length = uint64(size)
	//fmt.Printf("size---->%d\n", size)

	glyphTable := make([]byte, glyf.PaddedLength())
	locaTable := make([]int, numGlyphs+1)

	glyphOffset := 0
	glyphIndex := 0
	for idx := 0; idx < numGlyphs; idx++ {
		locaTable[idx] = glyphOffset
		if glyphIndex < glyphCount && glyphArray[glyphIndex] == idx {
			glyphIndex++
			bytes := me.getGlyphData(idx)
			length := len(bytes)
			if length > 0 {
				for i := 0; i < length; i++ {
					glyphTable[glyphOffset+i] = bytes[i]
				}
				glyphOffset += length
			}
		}
	} //end for
	locaTable[numGlyphs] = glyphOffset
	//fmt.Printf("---->%d\n", len(glyphTable))
	return glyphTable, locaTable, nil
}

func (me *PdfDictionaryObj) getGlyphSize(glyph int) int {
	ttfp := me.PtrToSubsetFontObj.GetTTFParser()
	glyf := ttfp.GetTables()["glyf"]
	start := int(glyf.Offset + ttfp.LocaTable[glyph])
	next := int(glyf.Offset + ttfp.LocaTable[glyph+1])
	return next - start
}

func (me *PdfDictionaryObj) getGlyphData(glyph int) []byte {
	ttfp := me.PtrToSubsetFontObj.GetTTFParser()
	glyf := ttfp.GetTables()["glyf"]
	start := int(glyf.Offset + ttfp.LocaTable[glyph])
	next := int(glyf.Offset + ttfp.LocaTable[glyph+1])
	count := next - start
	var data []byte
	i := 0
	for i < count {
		data = append(data, ttfp.FontData()[start+i])
		i++
	}
	return data
}

func (me *PdfDictionaryObj) makeFont() ([]byte, error) {
	var buff Buff
	ttfp := me.PtrToSubsetFontObj.GetTTFParser()
	tables := make(map[string]core.TableDirectoryEntry)
	tables["cvt "] = ttfp.GetTables()["cvt "] //มีช่องว่างด้วยนะ
	tables["fpgm"] = ttfp.GetTables()["fpgm"]
	tables["glyf"] = ttfp.GetTables()["glyf"]
	tables["head"] = ttfp.GetTables()["head"]
	tables["hhea"] = ttfp.GetTables()["hhea"]
	tables["hmtx"] = ttfp.GetTables()["hmtx"]
	tables["loca"] = ttfp.GetTables()["loca"]
	tables["maxp"] = ttfp.GetTables()["maxp"]
	tables["prep"] = ttfp.GetTables()["prep"]
	tableCount := len(tables)
	selector := EntrySelectors[tableCount]

	glyphTable, locaTable, err := me.makeGlyfAndLocaTable()
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%#v", glyphTable)

	WriteUInt32(&buff, 0x00010000)
	WriteUInt16(&buff, uint(tableCount))
	WriteUInt16(&buff, ((1 << uint(selector)) * 16))
	WriteUInt16(&buff, uint(selector))
	WriteUInt16(&buff, (uint(tableCount)-(1<<uint(selector)))*16)

	var tags []string
	for tag, _ := range tables {
		tags = append(tags, tag) //copy all tag
	}
	sort.Strings(tags) //order
	idx := 0
	tablePosition := int(12 + 16*tableCount)
	for idx < tableCount {
		entry := tables[tags[idx]]
		//write data
		//entry.Offset = uint64(tablePosition)
		offset := uint64(tablePosition)
		buff.SetPosition(tablePosition)
		if tags[idx] == "glyf" {
			entry.Length = uint64(len(glyphTable))
			entry.CheckSum = CheckSum(glyphTable)
			WriteBytes(&buff, glyphTable, 0, entry.PaddedLength())
		} else if tags[idx] == "loca" {
			if !ttfp.IsShortIndex {
				log.Fatalf("not suport none short index yet!")
				return nil, nil
			}
			//entry.Offset = 0
			entry.Length = uint64(len(locaTable) * 2)
			data := make([]byte, entry.PaddedLength())
			length := len(locaTable)
			byteIdx := 0
			for idx := 0; idx < length; idx++ {
				val := locaTable[idx] / 2
				data[byteIdx] = byte(val >> 8)
				byteIdx++
				data[byteIdx] = byte(val)
				byteIdx++
			}
			entry.CheckSum = CheckSum(data)
			WriteBytes(&buff, data, 0, len(data))
			//fmt.Printf(">>>>%#v\n%#v\n\n %d \n %d\n", entry, data, len(data), entry.CheckSum)
		} else {
			fmt.Printf("tag=%s offset=%d\n ", tags[idx], int(entry.Offset))
			WriteBytes(&buff, ttfp.FontData(), int(entry.Offset), entry.PaddedLength())
		}
		endPosition := buff.Position()
		tablePosition = endPosition

		//write table
		buff.SetPosition(idx*16 + 12)
		WriteTag(&buff, tags[idx])
		WriteUInt32(&buff, uint(entry.CheckSum))
		WriteUInt32(&buff, uint(offset)) //offset
		WriteUInt32(&buff, uint(entry.Length))

		tablePosition = endPosition
		idx++
	}
	//DebugSubType(buff.Bytes())
	//me.buffer.Write(buff.Bytes())
	return buff.Bytes(), nil
}

func (me *PdfDictionaryObj) completeGlyphClosure(glyphs map[rune]uint64) map[rune]uint64 {
	//count := len(glyphs)
	var glyphArray []int
	isContainZero := false
	for _, v := range glyphArray {
		glyphArray = append(glyphArray, v)
		if v == 0 {
			isContainZero = true
		}
	}

	if !isContainZero {
		glyphs[0] = 0
	}
	/*TODO ทำต่อ
		for (int idx = 0; idx < count; idx++)
	        AddCompositeGlyphs(glyphs, glyphArray[idx]);
	*/
	return glyphs
}

func CheckSum(data []byte) uint64 {

	var byte3, byte2, byte1, byte0 uint64
	byte3 = 0
	byte2 = 0
	byte1 = 0
	byte0 = 0
	length := len(data)
	i := 0
	for i < length {
		byte3 += uint64(data[i])
		i++
		byte2 += uint64(data[i])
		i++
		byte1 += uint64(data[i])
		i++
		byte0 += uint64(data[i])
		i++
	}
	//var result uint32
	result := uint32(byte3<<24) + uint32(byte2<<16) + uint32(byte1<<8) + uint32(byte0)
	//fmt.Printf(">>>> %d ,%d,%d ,%d %d   ----%d\n", byte3, byte2, byte1, byte0, result, uint32(byte3<<24))
	return uint64(result)
}
