package core

import (
	"bytes"
)

//ParseGDEF parse GDEF — Glyph Definition Table
//https://docs.microsoft.com/en-us/typography/opentype/spec/gdef
func (t *TTFParser) ParseGDEF(fd *bytes.Reader) error {

	header, err := t.parseGDEFHeader(fd)
	if err != nil {
		return err
	}

	glyphOfglyphClassDef := initParseClassDefinitionTableResult()
	if header.glyphClassDefOffset != 0 {
		result, err := t.parseClassDefinitionTable(fd, header.glyphClassDefOffset)
		if err != nil {
			return err
		}
		glyphOfglyphClassDef = result
	}

	glyphOfmarkAttachClassDef := initParseClassDefinitionTableResult()
	if header.markAttachClassDefOffset != 0 {
		result, err := t.parseClassDefinitionTable(fd, header.glyphClassDefOffset)
		if err != nil {
			return err
		}
		glyphOfmarkAttachClassDef = result
	}

	_ = glyphOfglyphClassDef
	_ = glyphOfmarkAttachClassDef

	return nil
}

func (t *TTFParser) parseGDEFHeader(fd *bytes.Reader) (gdefHeader, error) {
	err := t.Seek(fd, "GDEF")
	if err != nil {
		return gdefHeader{}, err
	}

	gdefOffset := int64(t.tables["GDEF"].Offset)
	var result gdefHeader

	majorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.majorVersion = majorVersion

	minorVersion, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.minorVersion = minorVersion

	glyphClassDefOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.glyphClassDefOffset = gdefOffset + int64(glyphClassDefOffset)

	attachListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.attachListOffset = gdefOffset + int64(attachListOffset)

	ligCaretListOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.ligCaretListOffset = gdefOffset + int64(ligCaretListOffset)

	markAttachClassDefOffset, err := t.ReadUShort(fd)
	if err != nil {
		return gdefHeader{}, err
	}
	result.markAttachClassDefOffset = gdefOffset + int64(markAttachClassDefOffset)

	if minorVersion == 2 || minorVersion == 3 {

		markGlyphSetsDefOffset, err := t.ReadUShort(fd)
		if err != nil {
			return gdefHeader{}, err
		}
		result.markGlyphSetsDefOffset = gdefOffset + int64(markGlyphSetsDefOffset)

		if minorVersion == 3 {
			itemVarStoreOffset, err := t.ReadUShort(fd)
			if err != nil {
				return gdefHeader{}, err
			}
			result.itemVarStoreOffset = gdefOffset + int64(itemVarStoreOffset)
		}
	}

	return result, nil
}
