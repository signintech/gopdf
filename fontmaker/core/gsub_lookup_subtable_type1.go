package core

import "bytes"

//GSUBLookupSubTableType1Format1 Single Substitution Format 1
type GSUBLookupSubTableType1Format1 struct {
	coverageOffset int64
	deltaGlyphID   int
}

//LookupType get lookup type
func (g *GSUBLookupSubTableType1Format1) LookupType() uint {
	return 1
}

//Format get format
func (g *GSUBLookupSubTableType1Format1) Format() uint {
	return 1
}

func (g *GSUBLookupSubTableType1Format1) processSubTable(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	gdefResult ParseGDEFResult,
) error {
	_, err := processGSUBLookupListTableSubTableLookupType1Format1(t, fd, table, *g, gdefResult)
	if err != nil {
		return err
	}
	return nil
}

func processGSUBLookupListTableSubTableLookupType1Format1(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType1Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableReplaceInfo, error) {

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableReplaceInfo{}, err
	}
	var result GSubLookupSubtableReplaceInfo
	glyphIDs := coverage.glyphIDs
	for _, glyphID := range glyphIDs {

		isIgnore, err := processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableReplaceInfo{}, err
		}
		if isIgnore {
			continue
		}
		var sub ReplaceRule
		sub.Substitute = []uint{glyphID + uint(subtable.deltaGlyphID)}
		sub.ReplaceGlyphIDs = append(sub.ReplaceGlyphIDs, glyphID)
		result.Rules = append(result.Rules, sub)
		//fmt.Printf("A ReplaceglyphIDs = %d Substitute =%d\n", glyphID, sub.Substitute)
	}
	return result, nil
}

//GSUBLookupSubTableType1Format2 Single Substitution Format 2
type GSUBLookupSubTableType1Format2 struct {
	coverageOffset     int64
	substituteGlyphIDs []uint
}

//LookupType get lookup type
func (g *GSUBLookupSubTableType1Format2) LookupType() uint {
	return 1
}

//Format get format
func (g *GSUBLookupSubTableType1Format2) Format() uint {
	return 2
}

func (g *GSUBLookupSubTableType1Format2) processSubTable(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	gdefResult ParseGDEFResult,
) error {
	_, err := processGSUBLookupListTableSubTableLookupType1Format2(t, fd, table, *g, gdefResult)
	if err != nil {
		return err
	}
	return nil
}

func processGSUBLookupListTableSubTableLookupType1Format2(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType1Format2,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableReplaceInfo, error) {

	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableReplaceInfo{}, err
	}
	var result GSubLookupSubtableReplaceInfo
	glyphIDs := coverage.glyphIDs
	for i, glyphID := range glyphIDs {

		isIgnore, err := processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
		if err != nil {
			return GSubLookupSubtableReplaceInfo{}, err
		}
		if isIgnore {
			continue
		}
		var sub ReplaceRule
		sub.Substitute = []uint{uint(subtable.substituteGlyphIDs[i])}
		sub.ReplaceGlyphIDs = append(sub.ReplaceGlyphIDs, glyphID)
		result.Rules = append(result.Rules, sub)
		//fmt.Printf("ReplaceglyphIDs = %d Substitute =%d\n", glyphID, sub.Substitute)
	}
	return result, nil
}
