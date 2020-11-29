package core

import (
	"bytes"
)

//3.1 Alternate Substitution Format 1
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType3Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	*GSUBLookupSubTableType3Format1,
	error,
) {

	var result GSUBLookupSubTableType3Format1

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}
	result.coverageOffset = int64(coverageOffset) + offset //set result

	alternateSetCount, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}
	result.alternateSetCount = alternateSetCount //set result

	var alternateSetOffsets []int64
	for i := uint(0); i < alternateSetCount; i++ {
		alternateSetOffset, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		alternateSetOffsets = append(alternateSetOffsets, offset+int64(alternateSetOffset))
	}
	result.alternateSetOffsets = alternateSetOffsets //set result

	var alternateSetTables []AlternateSetTable
	for _, alternateSetOffset := range result.alternateSetOffsets {
		_, err := fd.Seek(alternateSetOffset, 0)
		if err != nil {
			return nil, err
		}

		glyphCount, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}

		var alternateGlyphIDs []uint
		for i := uint(0); i < glyphCount; i++ {
			alternateGlyphID, err := t.ReadUShort(fd)
			if err != nil {
				return nil, err
			}
			alternateGlyphIDs = append(alternateGlyphIDs, alternateGlyphID)
		}
		alternateSetTables = append(alternateSetTables, AlternateSetTable{
			glyphCount:        glyphCount,
			alternateGlyphIDs: alternateGlyphIDs,
		})
	}
	result.alternateSetTables = alternateSetTables //set result

	return &result, nil
}
