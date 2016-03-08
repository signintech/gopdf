package core

import (
	"errors"
	"fmt"
	"os"
)

//ParseCmapFormat12 cmap format 12
func (t *TTFParser) ParseCmapFormat12(fd *os.File) error {

	t.Seek(fd, "cmap")
	t.Skip(fd, 2) // version
	numTables, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	var cEncodingSubtables []cmapFormat12EncodingSubtable
	for i := 0; i < int(numTables); i++ {
		platformID, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		encodingID, err := t.ReadUShort(fd)
		if err != nil {
			return err
		}
		offset, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		var ce cmapFormat12EncodingSubtable
		ce.platformID = platformID
		ce.encodingID = encodingID
		ce.offset = offset
		cEncodingSubtables = append(cEncodingSubtables, ce)
	}

	isFound := false
	offset := uint64(0)
	for _, ce := range cEncodingSubtables {
		if ce.platformID == 3 && ce.encodingID == 10 {
			offset = ce.offset
			isFound = true
			break
		}
	}

	if !isFound {
		return errors.New("not found Encoding Identifiers Unicode UCS-4")
	}

	_, err = fd.Seek(int64(t.tables["cmap"].Offset+offset), 0)
	if err != nil {
		return err
	}

	format, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	if format != 12 {
		return errors.New("format != 12")
	}

	reserved, err := t.ReadUShort(fd)
	if err != nil {
		return err
	}

	if reserved != 0 {
		return errors.New("reserved != 0")
	}

	length, err := t.ReadULong(fd)
	if err != nil {
		return err
	}

	err = t.Skip(fd, 4)
	if err != nil {
		return err
	}

	nGroups, err := t.ReadULong(fd)
	if err != nil {
		return err
	}

	fmt.Printf("length = %d , nGroups = %d\n", length, nGroups)

	g := uint64(0)
	for g < nGroups {
		startCharCode, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		endCharCode, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		glyphID, err := t.ReadULong(fd)
		if err != nil {
			return err
		}

		var gTb CmapFormat12GroupingTable
		gTb.StartCharCode = startCharCode
		gTb.EndCharCode = endCharCode
		gTb.GlyphID = glyphID
		t.groupingTables = append(t.groupingTables, gTb)
		g++
	}

	return nil
}

type cmapFormat12EncodingSubtable struct {
	platformID uint64
	encodingID uint64
	offset     uint64
}

type CmapFormat12GroupingTable struct {
	StartCharCode, EndCharCode, GlyphID uint64
}
