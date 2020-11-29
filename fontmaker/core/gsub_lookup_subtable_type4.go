package core

import (
	"bytes"

	"github.com/signintech/gopdf/fontmaker/sliceutil"
)

//LookupType 4: Ligature Substitution Subtable

//GSUBLookupSubTableType4Format1 4.1 Ligature Substitution Format 1
type GSUBLookupSubTableType4Format1 struct {
	coverageOffset     int64   //Offset to Coverage table, from beginning of substitution subtable
	ligatureSetCount   uint    //Number of LigatureSet tables
	ligatureSetOffsets []int64 //Array of offsets to LigatureSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	//table
	ligatureSetTables []LigatureSetTable
	replaseInfo       GSubLookupSubtableReplaceInfo
}

//LookupType 4: Ligature Substitution Subtable
func (g *GSUBLookupSubTableType4Format1) LookupType() uint {
	return 4
}

//Format identifier: format = 1
func (g *GSUBLookupSubTableType4Format1) Format() uint {
	return 1
}

func (g *GSUBLookupSubTableType4Format1) processSubTable(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	gdefResult ParseGDEFResult,
) error {
	replaseInfo, err := processGSUBLookupListTableSubTableLookupType4Format1(t, fd, table, *g, gdefResult)
	if err != nil {
		return err
	}
	g.replaseInfo = replaseInfo
	return nil
}

func (g *GSUBLookupSubTableType4Format1) process(glyphindexs []uint) ([]uint, error) {

	for _, rule := range g.replaseInfo.Rules {
		matchs := sliceutil.ContainSliceUint(glyphindexs, rule.ReplaceGlyphIDs)
		if len(matchs) > 0 {
			diffLen := 0
			for _, match := range matchs {
				start := glyphindexs[0 : match.FirstIndex+diffLen]
				end := glyphindexs[match.FirstIndex+match.Length+diffLen:]
				temp := start
				temp = append(temp, rule.Substitute...)
				temp = append(temp, end...)
				diffLen = len(temp) - len(glyphindexs)
				glyphindexs = temp
			}
			break
		}
	}

	return glyphindexs, nil
}

func processGSUBLookupListTableSubTableLookupType4Format1(
	t *TTFParser,
	fd *bytes.Reader,
	table GSUBLookupTable,
	subtable GSUBLookupSubTableType4Format1,
	gdefResult ParseGDEFResult,
) (GSubLookupSubtableReplaceInfo, error) {
	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return GSubLookupSubtableReplaceInfo{}, err
	}
	glyphIDs := coverage.glyphIDs

	//fmt.Printf("%+v\n", coverage.glyphIDs)
	result := GSubLookupSubtableReplaceInfo{}

	for x, ligatureSetTable := range subtable.ligatureSetTables {
		for _, ligatureTable := range ligatureSetTable.ligatureTables {
			var replaces []uint
			replaces = append(replaces, glyphIDs[x])
			isIgnore, err := processGSUBIsIgnore(table.lookupFlag, replaces[0], table.markFilteringSet, gdefResult)
			if err != nil {
				return GSubLookupSubtableReplaceInfo{}, err
			}
			if isIgnore {
				continue
			}
			//dg1
			//fmt.Printf("---> %d %d\n", ligatureTable.componentCount, len(ligatureTable.componentGlyphIDs))
			for z := uint(1); z < ligatureTable.componentCount; z++ {
				glyphID := ligatureTable.componentGlyphIDs[z-1]
				isIgnore, err := processGSUBIsIgnore(table.lookupFlag, glyphID, table.markFilteringSet, gdefResult)
				if err != nil {
					return GSubLookupSubtableReplaceInfo{}, err
				}
				if isIgnore {
					continue
				}
				replaces = append(replaces, glyphID)
			}

			sub := ReplaceRule{
				ReplaceGlyphIDs: replaces,
				Substitute:      []uint{ligatureTable.ligatureGlyph},
			}

			result.Rules = append(result.Rules, sub)

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
func processGSUBIsIgnore(
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

//LigatureSetTable LigatureSet table: All ligatures beginning with the same glyph
type LigatureSetTable struct {
	ligatureCount   uint    //Number of Ligature tables
	ligatureOffsets []int64 //Array of offsets to Ligature tables. Offsets are from beginning of LigatureSet table, ordered by preference.
	ligatureTables  []LigatureTable
}

//LigatureTable Ligature table: Glyph components for one ligature
type LigatureTable struct {
	ligatureGlyph     uint   //glyph ID of ligature to substitute
	componentCount    uint   //Number of components in the ligature
	componentGlyphIDs []uint //Array of component glyph IDs — start with the second component, ordered in writing direction
}
