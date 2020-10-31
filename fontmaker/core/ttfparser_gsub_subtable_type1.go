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
		//fmt.Printf("glyphID %d\n", glyphID)
	}

	return GSUBLookupSubTableType1Format2{
		coverageOffset:     int64(coverageOffset) + offset,
		substituteGlyphIDs: substituteGlyphIDs,
	}, nil
}

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType1Format1(
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType1Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}
	var result GSubLookupSubtableResult
	glyphIDs := coverage.glyphIDs
	for _, glyphID := range glyphIDs {

		isIgnore, err := t.processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableResult{}, err
		}
		if isIgnore {
			continue
		}
		var sub GSubLookupSubtableSub
		sub.Substitute = uint(subtable.deltaGlyphID)
		sub.ReplaceglyphIDs = append(sub.ReplaceglyphIDs, glyphID)
		result.Subs = append(result.Subs, sub)
		//fmt.Printf("A ReplaceglyphIDs = %d Substitute =%d\n", glyphID, sub.Substitute)
	}
	return result, nil
}

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType1Format2(
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType1Format2,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}
	var result GSubLookupSubtableResult
	glyphIDs := coverage.glyphIDs
	for i, glyphID := range glyphIDs {

		isIgnore, err := t.processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableResult{}, err
		}
		if isIgnore {
			continue
		}
		var sub GSubLookupSubtableSub
		sub.Substitute = uint(subtable.substituteGlyphIDs[i])
		sub.ReplaceglyphIDs = append(sub.ReplaceglyphIDs, glyphID)
		result.Subs = append(result.Subs, sub)
		//fmt.Printf("ReplaceglyphIDs = %d Substitute =%d\n", glyphID, sub.Substitute)
	}
	return result, nil
}
