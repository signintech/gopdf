package gopdf

import (
	"bytes"
	"compress/zlib"
	"errors"
	"sort"
	"strconv"

	"github.com/signintech/gopdf/fontmaker/core"
)

var EntrySelectors = []int{
	0, 0, 1, 1, 2, 2,
	2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 4, 4,
	4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4}

var ErrNotSupportShortIndexYet = errors.New("not suport none short index yet!")

type PdfDictionaryObj struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (p *PdfDictionaryObj) Init(funcGetRoot func() *GoPdf) {
}

func (p *PdfDictionaryObj) Build() error {
	b, err := p.makeFont()
	if err != nil {
		//log.Panicf("%s", err.Error())
		return err
	}

	//zipvar buff bytes.Buffer
	var zbuff bytes.Buffer
	gzipwriter := zlib.NewWriter(&zbuff)
	_, err = gzipwriter.Write(b)
	if err != nil {
		return err
	}
	gzipwriter.Close()

	p.buffer.WriteString("<</Length " + strconv.Itoa(zbuff.Len()) + "\n")
	p.buffer.WriteString("/Filter /FlateDecode\n")
	p.buffer.WriteString("/Length1 " + strconv.Itoa(len(b)) + "\n")
	p.buffer.WriteString(">>\n")
	p.buffer.WriteString("stream\n")
	p.buffer.Write(zbuff.Bytes())
	p.buffer.WriteString("\nendstream\n")
	return nil
}

func (p *PdfDictionaryObj) GetType() string {
	return "PdfDictionary"
}

func (p *PdfDictionaryObj) GetObjBuff() *bytes.Buffer {
	return &p.buffer
}

func (p *PdfDictionaryObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	p.PtrToSubsetFontObj = ptr
}

func (p *PdfDictionaryObj) makeGlyfAndLocaTable() ([]byte, []int, error) {
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	var glyf core.TableDirectoryEntry

	numGlyphs := int(ttfp.NumGlyphs())

	_, glyphArray := p.completeGlyphClosure(p.PtrToSubsetFontObj.CharacterToGlyphIndex)
	glyphCount := len(glyphArray)
	/*glyphCount := len(glyphs)
	//copy
	var glyphArray []int
	for _, v := range p.PtrToSubsetFontObj.CharacterToGlyphIndex {
		glyphArray = append(glyphArray, int(v))
	}*/
	sort.Ints(glyphArray)

	size := 0
	for idx := 0; idx < glyphCount; idx++ {
		size += p.getGlyphSize(glyphArray[idx])
	}
	glyf.Length = uint64(size)

	glyphTable := make([]byte, glyf.PaddedLength())
	locaTable := make([]int, numGlyphs+1)

	glyphOffset := 0
	glyphIndex := 0
	for idx := 0; idx < numGlyphs; idx++ {
		locaTable[idx] = glyphOffset
		if glyphIndex < glyphCount && glyphArray[glyphIndex] == idx {
			glyphIndex++
			bytes := p.getGlyphData(idx)
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
	return glyphTable, locaTable, nil
}

func (p *PdfDictionaryObj) getGlyphSize(glyph int) int {

	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	glyf := ttfp.GetTables()["glyf"]
	start := int(glyf.Offset + ttfp.LocaTable[glyph])
	next := int(glyf.Offset + ttfp.LocaTable[glyph+1])
	return next - start
}

func (p *PdfDictionaryObj) getGlyphData(glyph int) []byte {
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
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

func (p *PdfDictionaryObj) makeFont() ([]byte, error) {
	var buff Buff
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
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

	glyphTable, locaTable, err := p.makeGlyfAndLocaTable()
	if err != nil {
		return nil, err
	}

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
		offset := uint64(tablePosition)
		buff.SetPosition(tablePosition)
		if tags[idx] == "glyf" {
			entry.Length = uint64(len(glyphTable))
			entry.CheckSum = CheckSum(glyphTable)
			WriteBytes(&buff, glyphTable, 0, entry.PaddedLength())
		} else if tags[idx] == "loca" {
			if ttfp.IsShortIndex {
				entry.Length = uint64(len(locaTable) * 2)
			} else {
				entry.Length = uint64(len(locaTable) * 4)
			}

			data := make([]byte, entry.PaddedLength())
			length := len(locaTable)
			byteIdx := 0
			if ttfp.IsShortIndex {
				for idx := 0; idx < length; idx++ {
					val := locaTable[idx] / 2
					data[byteIdx] = byte(val >> 8)
					byteIdx++
					data[byteIdx] = byte(val)
					byteIdx++
				}
			} else {
				for idx := 0; idx < length; idx++ {
					val := locaTable[idx]
					data[byteIdx] = byte(val >> 24)
					byteIdx++
					data[byteIdx] = byte(val >> 16)
					byteIdx++
					data[byteIdx] = byte(val >> 8)
					byteIdx++
					data[byteIdx] = byte(val)
					byteIdx++
				}
			}
			entry.CheckSum = CheckSum(data)
			WriteBytes(&buff, data, 0, len(data))
		} else {
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

func (p *PdfDictionaryObj) completeGlyphClosure(glyphs map[rune]uint64) (map[rune]uint64, []int) {
	//count := len(glyphs)
	//fmt.Printf(">>>>>>%#v\n", glyphs)
	//runtime.Breakpoint()
	/*var glyphArray []int
	isContainZero := false
	for _, v := range glyphArray {
		fmt.Printf(">>>>%d\n", v)
		glyphArray = append(glyphArray, v)
		if v == 0 {
			isContainZero = true
		}
	}

	if !isContainZero {
		glyphs[0] = 0 //ผิด
		//glyphs = append(glyphs,
	}*/
	/*TODO ทำต่อ*/
	/*for idx := 0; idx < count; idx++ {
		me.AddCompositeGlyphs(glyphs, glyphArray[idx])
	}*/
	var glyphArray []int
	//copy
	isContainZero := false
	for _, v := range glyphs {
		glyphArray = append(glyphArray, int(v))
		if v == 0 {
			isContainZero = true
		}
	}
	if !isContainZero {
		glyphArray = append(glyphArray, 0)
	}

	i := 0
	count := len(glyphs)
	for i < count {
		p.AddCompositeGlyphs(glyphArray, glyphArray[i])
		i++
	}

	//return glyphs, []int{131, 0, 36, 118}
	return glyphs, glyphArray
}

func (p *PdfDictionaryObj) AddCompositeGlyphs(glyphArray []int, glyph int) []int {
	start := p.GetOffset(int(glyph))
	if start == p.GetOffset(int(glyph)+1) {
		return glyphArray
	}
	//ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	//fontData := ttfp.FontData()
	//fmt.Printf("--->%d\n", len(fontData))
	return glyphArray
}

func (p *PdfDictionaryObj) GetOffset(glyph int) int {
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	glyf := ttfp.GetTables()["glyf"]
	offset := int(glyf.Offset + ttfp.LocaTable[glyph])
	return offset
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
	return uint64(result)
}
