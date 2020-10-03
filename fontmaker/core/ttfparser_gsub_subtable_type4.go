package core

import (
	"bytes"
	"fmt"
)

//Ligature (format 4.1)
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType4Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
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

func (t *TTFParser) processGSUBLookupListTableSubTableLookupType4Format1(fd *bytes.Reader,
	subtable GSUBLookupSubTableType4Format1,
) error {
	coverage, err := t.readCoverage(fd, subtable.coverageOffset)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", coverage.glyphIDs)

	return nil
}
