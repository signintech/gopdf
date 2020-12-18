package core

import (
	"bytes"
	"fmt"
)

//https://docs.microsoft.com/en-us/typography/script-development/malayalam
//ParseGSUB paese GSUB table https://docs.microsoft.com/en-us/typography/opentype/spec/gsub
func (t *TTFParser) ParseGSUB(fd *bytes.Reader, gdef ParseGDEFResult) error {
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
	t.gsubScriptList = parseScriptListResult //set result

	parseFeatureListResult, err := t.parseFeatureList(fd, header.featureListOffset)
	if err != nil {
		return err
	}
	t.gsubFeatureList = parseFeatureListResult //set result

	lookupTables, err := t.parseGSUBLookupListTable(fd, header.lookupListOffset, gdef)
	if err != nil {
		return err
	}
	err = t.processGSUBLookupListTable(fd, lookupTables, gdef)
	if err != nil {
		return err
	}

	t.gsubLookups = lookupTables //set result

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

func (t *TTFParser) processGSUBLookupListTable(fd *bytes.Reader, lookupTables GSUBLookupTables, gdef ParseGDEFResult) error {
	//var result GSubLookupSubtableResult
	for _, lookupTable := range lookupTables.lookups {
		for _, subtable := range lookupTable.subTables {
			//_ = subtable
			if subtable == nil {
				continue
			}
			err := subtable.processSubTable(t, fd, lookupTable, gdef)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (t *TTFParser) parseGSUBLookupListTable(fd *bytes.Reader, offset int64, gdefResult ParseGDEFResult) (GSUBLookupTables, error) {

	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return GSUBLookupTables{}, err
	}

	//Number of lookups in this table
	lookupCount, err := t.ReadUShort(fd)
	if err != nil {
		return GSUBLookupTables{}, err
	}

	//Array of offsets to Lookup tables,
	//from beginning of LookupList — zero based (first lookup is Lookup index = 0)
	var lookups []uint
	for i := uint(0); i < lookupCount; i++ {
		l, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupTables{}, err
		}
		lookups = append(lookups, l)
	}

	var lookupTables []GSUBLookupTable
	for _, l := range lookups {

		var lookupTable GSUBLookupTable

		offsetLookup := int64(offset) + int64(l)

		_, err := fd.Seek(offsetLookup, 0)
		if err != nil {
			return GSUBLookupTables{}, err
		}

		//lookupType: Different enumerations for GSUB and GPOS
		lookupType, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupTables{}, err
		}
		lookupTable.lookupType = lookupType //set

		//lookupFlag: Lookup qualifiers
		lookupFlag, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupTables{}, err
		}
		lookupTable.lookupFlag = lookupFlag //set

		subTableCount, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupTables{}, err
		}
		lookupTable.subTableCount = subTableCount //set

		var subtableOffsets []int64
		for s := uint(0); s < subTableCount; s++ {
			subtableOffset, err := t.ReadUShort(fd)
			if err != nil {
				return GSUBLookupTables{}, err
			}
			subtableOffsets = append(subtableOffsets, offsetLookup+int64(subtableOffset))
		}
		lookupTable.subtableOffsets = subtableOffsets //set

		markFilteringSet, err := t.ReadUShort(fd)
		if err != nil {
			return GSUBLookupTables{}, err
		}
		lookupTable.markFilteringSet = markFilteringSet //set

		//fmt.Printf("lookupIndex =%d lookupType %d\n", lookupIndex, lookupType)
		var subtables []gsubLookupSubTabler
		for _, subtableOffset := range subtableOffsets {
			subtable, err := t.parseGSUBLookupListTableSubTable(fd, subtableOffset, lookupType, gdefResult)
			if err != nil {
				return GSUBLookupTables{}, err
			}
			subtables = append(subtables, subtable)
		}
		lookupTable.subTables = subtables

		lookupTables = append(lookupTables, lookupTable) //set
	}

	return GSUBLookupTables{lookups: lookupTables}, nil
}

func (t *TTFParser) parseGSUBLookupListTableSubTable(
	fd *bytes.Reader,
	offset int64,
	lookupType uint,
	gdefResult ParseGDEFResult,

) (gsubLookupSubTabler, error) {
	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return nil, err
	}

	substFormat, err := t.ReadUShort(fd)
	if err != nil {
		return nil, err
	}

	var subtable gsubLookupSubTabler
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

type GSubLookupSubtableReplaceInfo struct {
	Rules []ReplaceRule
}

func (g *GSubLookupSubtableReplaceInfo) merge(r GSubLookupSubtableReplaceInfo) {
	g.Rules = append(g.Rules, r.Rules...)
}

type ReplaceRule struct {
	ReplaceGlyphIDs []uint
	Substitute      []uint
}
