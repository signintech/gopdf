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
	GSUBLookupSubTableType1Format1,
	error,
) {

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType1Format1{}, err
	}

	deltaGlyphID, err := t.ReadShort(fd)
	if err != nil {
		return GSUBLookupSubTableType1Format1{}, err
	}

	return GSUBLookupSubTableType1Format1{
		CoverageOffset: coverageOffset,
		DeltaGlyphID:   deltaGlyphID,
	}, nil
}

//1.2 Single Substitution Format 2
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType1Format2(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	GSUBLookupSubTableType1Format2,
	error,
) {

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType1Format2{}, err
	}

	glyphCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType1Format2{}, err
	}

	var substituteGlyphIDs []uint
	for i := uint(0); i < glyphCount; i++ {
		glyphID, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType1Format2{}, err
		}
		substituteGlyphIDs = append(substituteGlyphIDs, glyphID)
	}

	return GSUBLookupSubTableType1Format2{
		CoverageOffset:     coverageOffset,
		SubstituteGlyphIDs: substituteGlyphIDs,
	}, nil
}
