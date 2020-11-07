package core

import (
	"bytes"
	"fmt"
)

//2.1 Multiple Substitution Format 1
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType2Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	GSUBLookupSubTableType2Format1,
	error,
) {
	result := GSUBLookupSubTableType2Format1{}

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType2Format1{}, err
	}
	result.coverageOffset = int64(coverageOffset) + offset //set result

	sequenceCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType2Format1{}, err
	}
	result.sequenceCount = sequenceCount //set result

	var sequenceOffsets []int64
	for x := uint(0); x < sequenceCount; x++ {
		sequenceOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType2Format1{}, err
		}
		sequenceOffsets = append(sequenceOffsets, int64(sequenceOffset)+offset)
	}
	result.sequenceOffsets = sequenceOffsets //set result

	var sequenceTables []GSUBLookupSequenceTable
	for _, offset := range sequenceOffsets {

		_, err := fd.Seek(offset, 0)
		if err != nil {
			return GSUBLookupSubTableType2Format1{}, err
		}
		glyphCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType2Format1{}, err
		}
		var substituteGlyphIDs []uint
		for y := uint(0); y < glyphCount; y++ {
			substituteGlyphID, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupSubTableType2Format1{}, err
			}
			substituteGlyphIDs = append(substituteGlyphIDs, substituteGlyphID)
		}

		sequenceTables = append(sequenceTables, GSUBLookupSequenceTable{
			glyphCount:         glyphCount,
			substituteGlyphIDs: substituteGlyphIDs,
		})
	}

	result.sequenceTable = sequenceTables

	return result, nil
}

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType2Format1(
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType2Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {
	var result GSubLookupSubtableResult

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}

	glyphIDs := coverage.glyphIDs
	for i, glyphID := range glyphIDs {
		isIgnore, err := t.processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableResult{}, err
		}
		if isIgnore {
			continue
		}

		if len(subtable.sequenceTable[i].substituteGlyphIDs) <= 0 {
			continue
		}
		var substitutes []uint
		for _, substituteGlyphID := range subtable.sequenceTable[i].substituteGlyphIDs {
			substitutes = append(substitutes, substituteGlyphID)
		}
		fmt.Printf("sss %d %+v\n", glyphID, substitutes)
		sub := GSubLookupSubtableSub{
			ReplaceglyphIDs: []uint{glyphID},
			Substitute:      substitutes,
		}

		result.Subs = append(result.Subs, sub)
	}

	return result, nil
}
