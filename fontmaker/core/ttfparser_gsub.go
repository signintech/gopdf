package core

import (
	"bytes"
	"fmt"
)

//ParseGSUB paese GSUB table https://docs.microsoft.com/en-us/typography/opentype/spec/gsub
func (t *TTFParser) ParseGSUB(fd *bytes.Reader) error {
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
	fmt.Printf("majorVersion=%d minorVersion=%d  scriptListOffset=%d featureListOffset=%d lookupListOffset=%d\n", majorVersion, minorVersion, scriptListOffset, featureListOffset, lookupListOffset)

	t.parseGSUBLookupListTable(fd, int64(t.tables["GSUB"].Offset+lookupListOffset))

	return nil
}

func (t *TTFParser) parseGSUBLookupListTable(fd *bytes.Reader, offset int64) error {

	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	//Number of lookups in this table
	lookupCount, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	//Array of offsets to Lookup tables,
	//from beginning of LookupList — zero based (first lookup is Lookup index = 0)
	var lookups []uint
	for i := uint(0); i < lookupCount; i++ {
		l, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		lookups = append(lookups, l)
	}

	for lookupIndex, l := range lookups {
		offsetLookup := int64(offset) + int64(l)

		_, err := fd.Seek(offsetLookup, 0)
		if err != nil {
			return err
		}

		//lookupType: Different enumerations for GSUB and GPOS
		lookupType, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		//lookupFlag: Lookup qualifiers
		lookupFlag, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		subTableCount, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}

		var subtableOffsets []uint
		for s := uint(0); s < subTableCount; s++ {
			subtableOffset, err := t.ReadUShort(fd)
			if err != nil {
				return err
			}
			subtableOffsets = append(subtableOffsets, subtableOffset)
		}

		markFilteringSet, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		_ = markFilteringSet
		_ = lookupType
		_ = lookupFlag

		fmt.Printf("lookupIndex =%d lookupType %d\n", lookupIndex, lookupType)
		for _, subtableOffset := range subtableOffsets {
			t.parseGSUBLookupListTableSubTable(fd, offsetLookup+int64(subtableOffset))
		}

		//fmt.Printf("\t\toffsetLookup =%d  lookupType=%d lookupFlag=%d  markFilteringSet=%d\n", offsetLookup, lookupType, lookupFlag, markFilteringSet)
	}

	return nil
}

func (t *TTFParser) parseGSUBLookupListTableSubTable(fd *bytes.Reader, offset int64) error {
	_, err := fd.Seek(int64(offset), 0)
	if err != nil {
		return err
	}

	coverageFormat, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	fmt.Printf("\tcoverageFormat =%d\n", coverageFormat)
	_ = coverageFormat
	return nil
}
