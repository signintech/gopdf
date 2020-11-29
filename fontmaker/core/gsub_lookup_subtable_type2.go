package core

import "bytes"

//GSUBLookupSubTableType2Format1 -> 2.1 Multiple Substitution Format 1
type GSUBLookupSubTableType2Format1 struct {
	substFormat     uint
	coverageOffset  int64   //Offset to Coverage table, from beginning of substitution subtable
	sequenceCount   uint    //Number of Sequence table offsets in the sequenceOffsets array
	sequenceOffsets []int64 //Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	sequenceTable   []GSUBLookupSequenceTable
}

//LookupType get lookup type
func (g *GSUBLookupSubTableType2Format1) LookupType() uint {
	return 2
}

//Format get format
func (g *GSUBLookupSubTableType2Format1) Format() uint {
	return 1
}

func (g *GSUBLookupSubTableType2Format1) processSubTable(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	gdefResult ParseGDEFResult,
) error {
	_, err := processGSUBLookupListTableSubTableLookupType2Format1(t, fd, table, *g, gdefResult)
	if err != nil {
		return err
	}
	return nil
}

func processGSUBLookupListTableSubTableLookupType2Format1(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType2Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableReplaceInfo, error) {
	var result GSubLookupSubtableReplaceInfo

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableReplaceInfo{}, err
	}

	glyphIDs := coverage.glyphIDs
	for i, glyphID := range glyphIDs {
		isIgnore, err := processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableReplaceInfo{}, err
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
		//fmt.Printf("sss %d %+v\n", glyphID, substitutes)
		sub := ReplaceRule{
			ReplaceGlyphIDs: []uint{glyphID},
			Substitute:      substitutes,
		}

		result.Rules = append(result.Rules, sub)
	}

	return result, nil
}

//GSUBLookupSequenceTable Sequence table
type GSUBLookupSequenceTable struct {
	glyphCount         uint   //Number of glyph IDs in the substituteGlyphIDs array. This must always be greater than 0.
	substituteGlyphIDs []uint //String of glyph IDs to substitute
}
