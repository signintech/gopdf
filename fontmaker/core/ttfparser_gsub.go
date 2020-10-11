package core

import (
	"bytes"
)

//ParseGSUB paese GSUB table https://docs.microsoft.com/en-us/typography/opentype/spec/gsub
func (t *TTFParser) ParseGSUB(fd *bytes.Reader, gdefResult ParseGDEFResult) error {
	err := t.Seek(fd, "GSUB")
	if err == ErrTableNotFound {
		return nil
	} else if err != nil {
		return err
	}

	majorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	minorVersion, err := t.ReadShort(fd)
	if err != nil {
		return err
	}

	scriptListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	featureListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	lookupListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	featureVariationsOffset := uint(0)
	if minorVersion == 1 {
		featureVariationsOffset, err = t.ReadULong(fd)
		if err != nil {
			return err
		}
	}
	_ = featureVariationsOffset
	_ = majorVersion
	_ = minorVersion
	_ = scriptListOffset
	_ = featureListOffset
	_ = lookupListOffset
	//fmt.Printf("majorVersion=%d minorVersion=%d  scriptListOffset=%d featureListOffset=%d lookupListOffset=%d\n", majorVersion, minorVersion, scriptListOffset, featureListOffset, lookupListOffset)

	lookupTables, err := t.parseGSUBLookupListTable(fd, int64(t.tables["GSUB"].Offset+lookupListOffset), gdefResult)
	if err != nil {
		return err
	}

	err = t.processGSUBLookupListTable(fd, lookupTables, gdefResult)
	if err != nil {
		return err
	}

	return nil
}

func (t *TTFParser) processGSUBLookupListTable(fd *bytes.Reader, lookupTables []GSUBLookupTable, gdefResult ParseGDEFResult) error {

	for _, lookupTable := range lookupTables {
		for _, subtable := range lookupTable.gsubLookupSubTables {
			//_ = subtable
			if subtable == nil {
				continue
			}
			//TODO: add other type
			if subtable.LookupType() == 4 && subtable.Format() == 1 {
				if subtable41, ok := subtable.(GSUBLookupSubTableType4Format1); ok {
					err := t.processGSUBLookupListTableSubTableLookupType4Format1(fd, lookupTable, subtable41, gdefResult)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
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
	//LookupType 4: Ligature Substitution Subtable
	if lookupType == 4 && substFormat == 1 {
		//4.1 Ligature Substitution Format 1
		subtable, err = t.parseGSUBLookupListTableSubTableLookupType4Format1(fd, offset, substFormat, gdefResult)
		if err != nil {
			return nil, err
		}
	}

	return subtable, nil
}
