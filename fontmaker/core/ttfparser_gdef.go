package core

import (
	"bytes"

	"github.com/signintech/gopdf/fontmaker/sliceutil"
)

//ParseGDEF parse GDEF — Glyph Definition Table
//https://docs.microsoft.com/en-us/typography/opentype/spec/gdef
func (t *TTFParser) ParseGDEF(fd *bytes.Reader) (ParseGDEFResult, error) {

	result := InitParseGDEFResult()

	header, err := t.parseGDEFHeader(fd)
	if err != nil {
		return ParseGDEFResult{}, err
	}

	if header.glyphClassDefOffset != 0 {
		r, err := t.parseClassDefinitionTable(fd, header.glyphClassDefOffset)
		if err != nil {
			return ParseGDEFResult{}, err
		}
		result.glyphClassBases = r.GlyphClassBases()           //set
		result.glyphClassComponents = r.GlyphClassComponents() //set
		result.glyphClassLigatures = r.GlyphClassLigatures()   //set
		result.glyphClassMarks = r.GlyphClassMarks()           //set
	}

	markAttachmentType := make(map[uint]([]uint)) // map[class]( array of glyphId)
	if header.markAttachClassDefOffset != 0 {
		r, err := t.parseClassDefinitionTable(fd, header.markAttachClassDefOffset)
		if err != nil {
			return ParseGDEFResult{}, err
		}
		marks := r.GlyphClassMarks()
		for cls, glyphIDs := range r.mapClassWithGlyphIDs {
			if len(marks) > 0 {
				diff := sliceutil.DiffUInt(marks, glyphIDs)
				markAttachmentType[cls] = diff
			}
		}
		result.markAttachmentType = markAttachmentType //set
	}

	return result, nil
}

func (t *TTFParser) parseGDEFHeader(fd *bytes.Reader) (gdefHeader, error) {
	err := t.Seek(fd, "GDEF")
	if err != nil {
		return gdefHeader{}, err
	}

	gdefOffset := int64(t.tables["GDEF"].Offset)
	var result gdefHeader

	majorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.majorVersion = majorVersion

	minorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.minorVersion = minorVersion

	glyphClassDefOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.glyphClassDefOffset = gdefOffset + int64(glyphClassDefOffset)

	attachListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.attachListOffset = gdefOffset + int64(attachListOffset)

	ligCaretListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.ligCaretListOffset = gdefOffset + int64(ligCaretListOffset)

	markAttachClassDefOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.markAttachClassDefOffset = gdefOffset + int64(markAttachClassDefOffset)

	if minorVersion == 2 || minorVersion == 3 {

		markGlyphSetsDefOffset, err := t.ReadUShort(fd)
		if err != nil {
			return gdefHeader{}, err
		}
		result.markGlyphSetsDefOffset = gdefOffset + int64(markGlyphSetsDefOffset)

		if minorVersion == 3 {
			itemVarStoreOffset, err := t.ReadUShort(fd)
			if err != nil {
				return gdefHeader{}, err
			}
			result.itemVarStoreOffset = gdefOffset + int64(itemVarStoreOffset)
		}
	}

	return result, nil
}

//ParseGDEFResult result form ParseGDEF
type ParseGDEFResult struct {
	glyphClassBases      []uint //class 1
	glyphClassLigatures  []uint //class 2
	glyphClassMarks      []uint //class 3
	glyphClassComponents []uint //class 4
	markAttachmentType   map[uint]([]uint)
}

//InitParseGDEFResult init ParseGDEFResult
func InitParseGDEFResult() ParseGDEFResult {
	var r ParseGDEFResult
	r.markAttachmentType = make(map[uint]([]uint))
	return r
}
