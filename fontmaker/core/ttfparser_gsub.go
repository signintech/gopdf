package core

import (
	"bytes"
	"fmt"
)

//ParseGSUB paese GSUB table https://docs.microsoft.com/en-us/typography/opentype/spec/gsub
func (t *TTFParser) ParseGSUB(fd *bytes.Reader, gdefResult ParseGDEFResult) error {
	err := t.Seek(fd, "GSUB")
	if err == ErrTableNotFound {
		return nil
	} else if err != nil {
		return err
	}
	gsubOffset := t.tables["GSUB"].Offset

	header, err := t.parseGSBHeader(fd, int64(gsubOffset))
	if err != nil {
		return err
	}

	parseScriptListResult, err := t.parseScriptList(fd, header.scriptListOffset)
	if err != nil {
		return err
	}

	parseFeatureListResult, err := t.parseFeatureList(fd, header.featureListOffset, parseScriptListResult)
	if err != nil {
		return err
	}

	_ = parseFeatureListResult
	//TEST

	//END TEST

	lookupTables, err := t.parseGSUBLookupListTable(fd, header.lookupListOffset, gdefResult)
	if err != nil {
		return err
	}

	resultGsubLk, err := t.processGSUBLookupListTable(fd, lookupTables, gdefResult)
	if err != nil {
		return err
	}
	t.GSubLookupSubtable = resultGsubLk //set global value

	return nil
}

func (t *TTFParser) parseGSBHeader(fd *bytes.Reader, offset int64) (GSUBHeader, error) {

	_, err := fd.Seek(offset, 0)
	if err != nil {
		return GSUBHeader{}, err
	}

	majorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBHeader{}, err
	}

	minorVersion, err := t.ReadShort(fd)
	if err != nil {
		return GSUBHeader{}, err
	}

	scriptListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBHeader{}, err
	}

	featureListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBHeader{}, err
	}

	lookupListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBHeader{}, err
	}

	return GSUBHeader{
		gsubOffset:        offset,
		majorVersion:      majorVersion,
		minorVersion:      minorVersion,
		scriptListOffset:  int64(scriptListOffset) + offset,
		featureListOffset: int64(featureListOffset) + offset,
		lookupListOffset:  int64(lookupListOffset) + offset,
	}, nil
}

func (t *TTFParser) processGSUBLookupListTable(fd *bytes.Reader, lookupTables []GSUBLookupTable, gdefResult ParseGDEFResult) (GSubLookupSubtableResult, error) {
	var result GSubLookupSubtableResult
	for _, lookupTable := range lookupTables {
		for _, subtable := range lookupTable.gsubLookupSubTables {
			//_ = subtable
			if subtable == nil {
				continue
			}
			//TODO: add other type
			if subtable.LookupType() == 1 && subtable.Format() == 1 {
				if subtable1F1, ok := subtable.(GSUBLookupSubTableType1Format1); ok {
					resultType1F1, err := t.processGSUBLookupListTableSubTableLookupType1Format1(fd, lookupTable, subtable1F1, gdefResult)
					if err != nil {
						return GSubLookupSubtableResult{}, err
					}
					result.merge(resultType1F1)
				}
			} else if subtable.LookupType() == 1 && subtable.Format() == 2 {
				if subtable1F2, ok := subtable.(GSUBLookupSubTableType1Format2); ok {
					resultType1F2, err := t.processGSUBLookupListTableSubTableLookupType1Format2(fd, lookupTable, subtable1F2, gdefResult)
					if err != nil {
						return GSubLookupSubtableResult{}, err
					}
					result.merge(resultType1F2)
				}
			} else if subtable.LookupType() == 2 && subtable.Format() == 1 {
				if subtable2F1, ok := subtable.(GSUBLookupSubTableType2Format1); ok {
					resultType2F1, err := t.processGSUBLookupListTableSubTableLookupType2Format1(fd, lookupTable, subtable2F1, gdefResult)
					if err != nil {
						return GSubLookupSubtableResult{}, err
					}
					result.merge(resultType2F1)
				}
			} else if subtable.LookupType() == 3 && subtable.Format() == 1 {
				if subtable3F1, ok := subtable.(GSUBLookupSubTableType3Format1); ok {
					resultType3F1, err := t.processGSUBLookupListTableSubTableLookupType3Format1(fd, lookupTable, subtable3F1, gdefResult)
					if err != nil {
						return GSubLookupSubtableResult{}, err
					}
					result.merge(resultType3F1)
				}
			} else if subtable.LookupType() == 4 && subtable.Format() == 1 {
				if subtable4F1, ok := subtable.(GSUBLookupSubTableType4Format1); ok {
					resultType4F1, err := t.processGSUBLookupListTableSubTableLookupType4Format1(fd, lookupTable, subtable4F1, gdefResult)
					if err != nil {
						return GSubLookupSubtableResult{}, err
					}
					result.merge(resultType4F1)
				} else {
					return GSubLookupSubtableResult{}, fmt.Errorf("subtable not GSUBLookupSubTableType4Format1")
				}
			}

		}
	}

	return result, nil
}

func (t *TTFParser) parseGSUBLookupListTable(fd *bytes.Reader, offset int64, gdefResult ParseGDEFResult) ([]GSUBLookupTable, error) {

	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}

	//Number of lookups in this table
	lookupCount, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	//Array of offsets to Lookup tables,
	//from beginning of LookupList — zero based (first lookup is Lookup index = 0)
	var lookups []uint
	for i := uint(0); i < lookupCount; i++ {
		l, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		lookups = append(lookups, l)
	}

	var lookupTables []GSUBLookupTable
	for _, l := range lookups {

		var lookupTable GSUBLookupTable

		offsetLookup := int64(offset) + int64(l)

		_, err := fd.Seek(offsetLookup, 0)
		if err != nil {
			return nil, err
		}

		//lookupType: Different enumerations for GSUB and GPOS
		lookupType, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		lookupTable.lookupType = lookupType //set

		//lookupFlag: Lookup qualifiers
		lookupFlag, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		lookupTable.lookupFlag = lookupFlag //set

		subTableCount, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		lookupTable.subTableCount = subTableCount //set

		var subtableOffsets []int64
		for s := uint(0); s < subTableCount; s++ {
			subtableOffset, err := t.ReadUShort(fd)
			if err != nil {
				return nil, err
			}
			subtableOffsets = append(subtableOffsets, offsetLookup+int64(subtableOffset))
		}
		lookupTable.subtableOffsets = subtableOffsets //set

		markFilteringSet, err := t.ReadUShort(fd)
		if err != nil {
			return nil, err
		}
		lookupTable.markFilteringSet = markFilteringSet //set

		//fmt.Printf("lookupIndex =%d lookupType %d\n", lookupIndex, lookupType)
		var subtables []gsubLookupSubTableTyper
		for _, subtableOffset := range subtableOffsets {
			subtable, err := t.parseGSUBLookupListTableSubTable(fd, subtableOffset, lookupType, gdefResult)
			if err != nil {
				return nil, err
			}
			subtables = append(subtables, subtable)
		}
		lookupTable.gsubLookupSubTables = subtables

		lookupTables = append(lookupTables, lookupTable) //set
	}

	return lookupTables, nil
}

func (t *TTFParser) parseGSUBLookupListTableSubTable(
	fd *bytes.Reader,
	offset int64,
	lookupType uint,
	gdefResult ParseGDEFResult,

) (gsubLookupSubTableTyper, error) {
	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}

	substFormat, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	var subtable gsubLookupSubTableTyper
	//TODO: add other type
	if lookupType == 1 { //LookupType 1: Single Substitution Subtable

		if substFormat == 1 {
			subtable, err = t.parseGSUBLookupListTableSubTableLookupType1Format1(fd, offset, substFormat, gdefResult)
			if err != nil {
				return nil, err
			}
		} else if substFormat == 2 {
			subtable, err = t.parseGSUBLookupListTableSubTableLookupType1Format2(fd, offset, substFormat, gdefResult)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsuport lookup type %d format %d", lookupType, substFormat)
		}

	} else if lookupType == 2 {
		if substFormat == 1 {
			subtable, err = t.parseGSUBLookupListTableSubTableLookupType2Format1(fd, offset, substFormat, gdefResult)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsuport lookup type %d format %d", lookupType, substFormat)
		}
	} else if lookupType == 3 {
		if substFormat == 1 {
			subtable, err = t.parseGSUBLookupListTableSubTableLookupType3Format1(fd, offset, substFormat, gdefResult)
			if err != nil {
				return nil, err
			}
		}
	} else if lookupType == 4 {
		//LookupType 4: Ligature Substitution Subtable
		//4.1 Ligature Substitution Format 1
		if substFormat == 1 {
			subtable, err = t.parseGSUBLookupListTableSubTableLookupType4Format1(fd, offset, substFormat, gdefResult)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unsuport lookup type %d format %d", lookupType, substFormat)
		}
	}

	return subtable, nil
}

type GSubLookupSubtableResult struct {
	Subs []GSubLookupSubtableSub
}

func (g *GSubLookupSubtableResult) merge(r GSubLookupSubtableResult) {
	g.Subs = append(g.Subs, r.Subs...)
}

type GSubLookupSubtableSub struct {
	ReplaceglyphIDs []uint
	Substitute      []uint
}
