package core

import "bytes"

//GSUBLookupSubTableType3Format1 3.1 Alternate Substitution Format 1
type GSUBLookupSubTableType3Format1 struct {
	coverageOffset      int64   //Offset to Coverage table, from beginning of substitution subtable
	alternateSetCount   uint    //	Number of AlternateSet tables
	alternateSetOffsets []int64 //Array of offsets to AlternateSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	alternateSetTables  []AlternateSetTable
}

//LookupType get lookup type
func (g *GSUBLookupSubTableType3Format1) LookupType() uint {
	return 3
}

//Format get format
func (g *GSUBLookupSubTableType3Format1) Format() uint {
	return 1
}

func (g *GSUBLookupSubTableType3Format1) processSubTable(t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	gdefResult ParseGDEFResult) error {
	_, err := processGSUBLookupListTableSubTableLookupType3Format1(t, fd, table, *g, gdefResult)
	if err != nil {
		return err
	}
	return nil
}

func processGSUBLookupListTableSubTableLookupType3Format1(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType3Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {
	var result GSubLookupSubtableResult
	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}
	glyphIDs := coverage.glyphIDs
	for i, glyphID := range glyphIDs {
		isIgnore, err := processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableResult{}, err
		}
		if isIgnore {
			continue
		}
		if len(subtable.alternateSetTables[i].alternateGlyphIDs) > 0 {
			continue
		}

		var sub GSubLookupSubtableSub
		sub.Substitute = []uint{
			subtable.alternateSetTables[i].alternateGlyphIDs[0],
		}
		sub.ReplaceglyphIDs = []uint{
			glyphID,
		}
		result.Subs = append(result.Subs, sub)
		//fmt.Printf(">>> %+v\n", result)
	}
	return result, nil
}

//AlternateSetTable AlternateSet table
type AlternateSetTable struct {
	glyphCount        uint   //	Number of glyph IDs in the alternateGlyphIDs array
	alternateGlyphIDs []uint //Array of alternate glyph IDs, in arbitrary order
}
