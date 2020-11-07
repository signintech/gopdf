package core

import (
	"bytes"

	"github.com/signintech/gopdf/fontmaker/sliceutil"
)

//Ligature (format 4.1)
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType4Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	GSUBLookupSubTableType4Format1,
	error,
) {

	var result GSUBLookupSubTableType4Format1

	coverageOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType4Format1{}, err
	}
	result.coverageOffset = int64(coverageOffset) + offset //set result

	ligatureSetCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupSubTableType4Format1{}, err
	}
	result.ligatureSetCount = ligatureSetCount //set result

	var ligatureSetOffsets []int64
	for i := uint(0); i < ligatureSetCount; i++ {
		ligatureSetOffset, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType4Format1{}, err
		}
		ligatureSetOffsets = append(ligatureSetOffsets, int64(ligatureSetOffset)+offset)
	}
	result.ligatureSetOffsets = ligatureSetOffsets //set result

	for _, ligatureSetOffset := range ligatureSetOffsets {
		var ligatureSetTable LigatureSetTable
		_, err := fd.Seek(int64(ligatureSetOffset), 0)
		if err != nil {
			return GSUBLookupSubTableType4Format1{}, err
		}
		ligatureCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupSubTableType4Format1{}, err
		}
		ligatureSetTable.ligatureCount = ligatureCount //set result

		var ligatureOffsets []int64
		for j := uint(0); j < ligatureCount; j++ {
			ligatureOffset, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupSubTableType4Format1{}, err
			}
			ligatureOffsets = append(ligatureOffsets, ligatureSetOffset+int64(ligatureOffset))
		}
		ligatureSetTable.ligatureOffsets = ligatureOffsets //set result

		for j := uint(0); j < ligatureCount; j++ {
			var ligatureTables []LigatureTable
			for _, ligatureOffset := range ligatureOffsets {
				var ligatureTable LigatureTable
				_, err := fd.Seek(int64(ligatureOffset), 0)
				if err != nil {
					return GSUBLookupSubTableType4Format1{}, err
				}
				ligatureGlyph, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBLookupSubTableType4Format1{}, err
				}
				ligatureTable.ligatureGlyph = ligatureGlyph //set result

				componentCount, err := t.ReadUShort(fd)
				if err != nil {
					return GSUBLookupSubTableType4Format1{}, err
				}
				ligatureTable.componentCount = componentCount //set result

				var componentGlyphIDs []uint
				for s := uint(0); s < componentCount-1; s++ {
					componentGlyphID, err := t.ReadUShort(fd)
					if err != nil {
						return GSUBLookupSubTableType4Format1{}, err
					}
					componentGlyphIDs = append(componentGlyphIDs, componentGlyphID)
				}
				ligatureTable.componentGlyphIDs = componentGlyphIDs
				ligatureTables = append(ligatureTables, ligatureTable) //set result
			}
			ligatureSetTable.ligatureTables = ligatureTables
		}

		result.ligatureSetTables = append(result.ligatureSetTables, ligatureSetTable) //set result
	} //end for

	return result, nil
}

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType4Format1(
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType4Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableResult, error) {
	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableResult{}, err
	}
	glyphIDs := coverage.glyphIDs

	//fmt.Printf("%+v\n", coverage.glyphIDs)
	result := GSubLookupSubtableResult{}

	for x, ligatureSetTable := range subtable.ligatureSetTables {
		for _, ligatureTable := range ligatureSetTable.ligatureTables {
			var replaces []uint
			replaces = append(replaces, glyphIDs[x])
			isIgnore, err := t.processGSUBIsIgnore(table.lookupFlag, replaces[0], table.markFilteringSet, gdefResult)
			if err != nil {
				return GSubLookupSubtableResult{}, err
			}
			if isIgnore {
				continue
			}
			//dg1
			//fmt.Printf("---> %d %d\n", ligatureTable.componentCount, len(ligatureTable.componentGlyphIDs))
			for z := uint(1); z < ligatureTable.componentCount; z++ {
				glyphID := ligatureTable.componentGlyphIDs[z-1]
				isIgnore, err := t.processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
				if err != nil {
					return GSubLookupSubtableResult{}, err
				}
				if isIgnore {
					continue
				}
				replaces = append(replaces, glyphID)
			}

			sub := GSubLookupSubtableSub{
				ReplaceglyphIDs: replaces,
				Substitute:      []uint{ligatureTable.ligatureGlyph},
			}

			result.Subs = append(result.Subs, sub)
			/*
				if ligatureTable.ligatureGlyph == 187 {
					fmt.Printf("%d => %v\n", ligatureTable.ligatureGlyph, replaces)
				}
			*/

		}
	}

	return result, nil
}

/*
0x0001	rightToLeft					This bit relates only to the correct processing of the cursive attachment lookup type (GPOS lookup type 3). When this bit is set, the last glyph in a given sequence to which the cursive attachment lookup is applied, will be positioned on the baseline.
									Note: Setting of this bit is not intended to be used by operating systems or applications to determine text direction.
0x0002	ignoreBaseGlyphs			If set, skips over base glyphs
0x0004	ignoreLigatures				If set, skips over ligatures
0x0008	ignoreMarks					If set, skips over all combining marks
0x0010	useMarkFilteringSet			If set, indicates that the lookup table structure is followed by a MarkFilteringSet field. The layout engine skips over all mark glyphs not in the mark filtering set indicated.
0x00E0	reserved					For future use (Set to zero)
0xFF00	markAttachmentType			If not zero, skips over all marks of attachment type different from specified.
*/
//dj2
func (t *TTFParser) processGSUBIsIgnore(
	lookupFlag uint,
	glyphID uint,
	markFilteringSet uint,
	gdefResult ParseGDEFResult,
) (bool, error) {
	//TODO: ทำต่อ........................
	if lookupFlag&0x0002 == 0x0002 && sliceutil.IndexUint(gdefResult.glyphClassBases, glyphID) != -1 {
		return true, nil
	} else if lookupFlag&0x0004 == 0x0004 && sliceutil.IndexUint(gdefResult.glyphClassLigatures, glyphID) != -1 {
		return true, nil
	} else if (lookupFlag&0x0008 == 0x0008 && lookupFlag&0xFF00 == 0) && sliceutil.IndexUint(gdefResult.glyphClassMarks, glyphID) != -1 {
		return true, nil
	}

	return false, nil
}
