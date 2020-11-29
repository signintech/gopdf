package core

import (
	"bytes"
)

//1.1 Single Substitution Format 1
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType1Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	*GSUBLookupSubTableType1Format1,
	error,
) {

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	deltaGlyphID, err := t.ReadShort(fd)
	if err != nil {
		return nil, err
	}

	return &GSUBLookupSubTableType1Format1{
		coverageOffset: int64(coverageOffset) + offset,
		deltaGlyphID:   deltaGlyphID,
	}, nil
}

//1.2 Single Substitution Format 2
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType1Format2(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	*GSUBLookupSubTableType1Format2,
	error,
) {

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	glyphCount, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	var substituteGlyphIDs []uint
	for i := uint(0); i < glyphCount; i++ {
		glyphID, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		substituteGlyphIDs = append(substituteGlyphIDs, glyphID)
		//fmt.Printf("glyphID %d\n", glyphID)
	}

	return &GSUBLookupSubTableType1Format2{
		coverageOffset:     int64(coverageOffset) + offset,
		substituteGlyphIDs: substituteGlyphIDs,
	}, nil
}
