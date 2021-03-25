package gopdf

import (
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"sort"

	"github.com/signintech/gopdf/fontmaker/core"
)

//EntrySelectors entry selectors
var EntrySelectors = []int{
	0, 0, 1, 1, 2, 2,
	2, 2, 3, 3, 3, 3,
	3, 3, 3, 3, 4, 4,
	4, 4, 4, 4, 4, 4,
	4, 4, 4, 4, 4, 4, 4,
}

//ErrNotSupportShortIndexYet not suport none short index yet
var ErrNotSupportShortIndexYet = errors.New("not suport none short index yet")

//PdfDictionaryObj pdf dictionary object
type PdfDictionaryObj struct {
	PtrToSubsetFontObj *SubsetFontObj
	//getRoot            func() *GoPdf
	pdfProtection *PDFProtection
}

func (p *PdfDictionaryObj) init(funcGetRoot func() *GoPdf) {
	//p.getRoot = funcGetRoot
}

func (p *PdfDictionaryObj) setProtection(pr *PDFProtection) {
	p.pdfProtection = pr
}

func (p *PdfDictionaryObj) protection() *PDFProtection {
	return p.pdfProtection
}

func (p *PdfDictionaryObj) write(w io.Writer, objID int) error {
	b, err := p.makeFont()
	if err != nil {
		//log.Panicf("%s", err.Error())
		return err
	}

	//zipvar buff bytes.Buffer
	zbuff := GetBuffer()
	defer PutBuffer(zbuff)

	gzipwriter := zlib.NewWriter(zbuff)
	_, err = gzipwriter.Write(b)
	if err != nil {
		return err
	}
	gzipwriter.Close()

	fmt.Fprintf(w, "<</Length %d\n", zbuff.Len())
	io.WriteString(w, "/Filter /FlateDecode\n")
	fmt.Fprintf(w, "/Length1 %d\n", len(b))
	io.WriteString(w, ">>\n")
	io.WriteString(w, "stream\n")
	if p.protection() != nil {
		tmp, err := rc4Cip(p.protection().objectkey(objID), zbuff.Bytes())
		if err != nil {
			return err
		}
		w.Write(tmp)
		//p.buffer.WriteString("\n")
	} else {
		w.Write(zbuff.Bytes())
	}
	io.WriteString(w, "\nendstream\n")

	return nil
}

func (p *PdfDictionaryObj) getType() string {
	return "PdfDictionary"
}

//SetPtrToSubsetFontObj set subsetFontObj pointer
func (p *PdfDictionaryObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	p.PtrToSubsetFontObj = ptr
}

func (p *PdfDictionaryObj) makeGlyfAndLocaTable() ([]byte, []int, error) {
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	var glyf core.TableDirectoryEntry

	numGlyphs := int(ttfp.NumGlyphs())

	glyphArray := p.completeGlyphClosure(p.PtrToSubsetFontObj.CharacterToGlyphIndex)
	glyphCount := len(glyphArray)
	sort.Ints(glyphArray)

	size := 0
	for idx := 0; idx < glyphCount; idx++ {
		size += p.getGlyphSize(glyphArray[idx])
	}
	glyf.Length = uint(size)

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
	for tag := range tables {
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
			entry.Length = uint(len(glyphTable))
			entry.CheckSum = CheckSum(glyphTable)
			WriteBytes(&buff, glyphTable, 0, entry.PaddedLength())
		} else if tags[idx] == "loca" {
			if ttfp.IsShortIndex {
				entry.Length = uint(len(locaTable) * 2)
			} else {
				entry.Length = uint(len(locaTable) * 4)
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

func (p *PdfDictionaryObj) completeGlyphClosure(mapOfglyphs *MapOfCharacterToGlyphIndex) []int {
	var glyphArray []int
	//copy
	isContainZero := false
	glyphs := mapOfglyphs.AllVals()
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
		p.AddCompositeGlyphs(&glyphArray, glyphArray[i])
		i++
	}
	return glyphArray
}

//AddCompositeGlyphs add composite glyph
//composite glyph is a Unicode entity that can be defined as a sequence of one or more other characters.
func (p *PdfDictionaryObj) AddCompositeGlyphs(glyphArray *[]int, glyph int) {
	start := p.GetOffset(int(glyph))
	if start == p.GetOffset(int(glyph)+1) {
		return
	}

	offset := start
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	fontData := ttfp.FontData()
	numContours, step := ReadShortFromByte(fontData, offset)
	offset += step
	if numContours >= 0 {
		return
	}

	offset += 8
	for {
		flags, step1 := ReadUShortFromByte(fontData, offset)
		offset += step1
		cGlyph, step2 := ReadUShortFromByte(fontData, offset)
		offset += step2
		//check cGlyph is contain in glyphArray?
		glyphContainsKey := false
		for _, g := range *glyphArray {
			if g == int(cGlyph) {
				glyphContainsKey = true
				break
			}
		}
		if !glyphContainsKey {
			*glyphArray = append(*glyphArray, int(cGlyph))
		}

		if (flags & moreComponents) == 0 {
			return
		}
		offsetAppend := 4
		if (flags & arg1and2areWords) == 0 {
			offsetAppend = 2
		}
		if (flags & hasScale) != 0 {
			offsetAppend += 2
		} else if (flags & xAndYScale) != 0 {
			offsetAppend += 4
		}
		if (flags & twoByTwo) != 0 {
			offsetAppend += 8
		}
		offset += offsetAppend
	}
}

const hasScale = 8
const moreComponents = 32
const arg1and2areWords = 1
const xAndYScale = 64
const twoByTwo = 128

//GetOffset get offset from glyf table
func (p *PdfDictionaryObj) GetOffset(glyph int) int {
	ttfp := p.PtrToSubsetFontObj.GetTTFParser()
	glyf := ttfp.GetTables()["glyf"]
	offset := int(glyf.Offset + ttfp.LocaTable[glyph])
	return offset
}

//CheckSum check sum
func CheckSum(data []byte) uint {

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
	return uint(result)
}
