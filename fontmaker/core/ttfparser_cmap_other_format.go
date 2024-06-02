package core

import (
	"bytes"
	"errors"
)

// ParseCmapFormat12 parse cmap table format 12 https://www.microsoft.com/typography/otspec/cmap.htm
func (t *TTFParser) ParseCmapFormat12(fd *bytes.Reader) (bool, error) {

	err := t.Seek(fd, "cmap")
	if err != nil {
		return false, err
	}
	err = t.Skip(fd, 2) //skip version
	if err != nil {
		return false, err
	}
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return false, err
	}
	var cEncodingSubtables []cmapFormat12EncodingSubtable
	for i := 0; i < int(numTables); i++ {
		platformID, err := t.ReadUShort(fd)
		if err != nil {
			return false, err
		}
		encodingID, err := t.ReadUShort(fd)
		if err != nil {
			return false, err
		}
		offset, err := t.ReadULong(fd)
		if err != nil {
			return false, err
		}

		var ce cmapFormat12EncodingSubtable
		ce.platformID = platformID
		ce.encodingID = encodingID
		ce.offset = offset
		cEncodingSubtables = append(cEncodingSubtables, ce)
	}

	isFound := false
	offset := uint(0)
	for _, ce := range cEncodingSubtables {
		if ce.platformID == 3 && ce.encodingID == 10 {
			offset = ce.offset
			isFound = true
			break
		}
	}

	if !isFound {
		return false, nil
	}

	_, err = fd.Seek(int64(t.tables["cmap"].Offset+offset), 0)
	if err != nil {
		return false, err
	}

	format, err := t.ReadUShort(fd)
	if err != nil {
		return false, err
	}

	if format != 12 {
		return false, errors.New("format != 12")
	}

	reserved, err := t.ReadUShort(fd)
	if err != nil {
		return false, err
	}

	if reserved != 0 {
		return false, errors.New("reserved != 0")
	}

	err = t.Skip(fd, 4) //skip length
	if err != nil {
		return false, err
	}

	err = t.Skip(fd, 4) //skip language
	if err != nil {
		return false, err
	}

	nGroups, err := t.ReadULong(fd)
	if err != nil {
		return false, err
	}

	g := uint(0)
	for g < nGroups {
		startCharCode, err := t.ReadULong(fd)
		if err != nil {
			return false, err
		}

		endCharCode, err := t.ReadULong(fd)
		if err != nil {
			return false, err
		}

		glyphID, err := t.ReadULong(fd)
		if err != nil {
			return false, err
		}

		var gTb CmapFormat12GroupingTable
		gTb.StartCharCode = startCharCode
		gTb.EndCharCode = endCharCode
		gTb.GlyphID = glyphID
		t.groupingTables = append(t.groupingTables, gTb)
		g++
	}

	return true, nil
}

type cmapFormat12EncodingSubtable struct {
	platformID uint
	encodingID uint
	offset     uint
}

type CmapFormat12GroupingTable struct {
	StartCharCode, EndCharCode, GlyphID uint
}
