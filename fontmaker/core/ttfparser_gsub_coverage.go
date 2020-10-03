package core

import (
	"bytes"
	"fmt"
)

//readCoverage read Coverage Table
//https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#coverage-table
func (t *TTFParser) readCoverage(fd *bytes.Reader, coverageOffset int64) (coverageResult, error) {

	var result coverageResult

	_, err := fd.Seek(coverageOffset, 0)
	if err != nil {
		return coverageResult{}, err
	}

	coverageFormat, err := t.ReadUShort(fd) //	Format identifier
	if err != nil {
		return coverageResult{}, err
	}

	if coverageFormat == 1 {
		glyphCount, err := t.ReadUShort(fd) //glyphCount Number of glyphs in the glyph array
		if err != nil {
			return coverageResult{}, err
		}
		var glyphArray []uint //Array of glyph IDs — in numerical order
		for i := uint(0); i < glyphCount; i++ {
			glyph, err := t.ReadUShort(fd)
			if err != nil {
				return coverageResult{}, err
			}
			glyphArray = append(glyphArray, glyph)
		}
		result.glyphIDs = glyphArray //set
	} else if coverageFormat == 2 {
		rangeCount, err := t.ReadUShort(fd) //Number of RangeRecords
		if err != nil {
			return coverageResult{}, err
		}
		//var rangeRecords []GSUBRangeCecord //Array of glyph ranges — ordered by startGlyphID.
		for i := uint(0); i < rangeCount; i++ {
			startGlyphID, err := t.ReadUShort(fd) //First glyph ID in the range
			if err != nil {
				return coverageResult{}, err
			}
			endGlyphID, err := t.ReadUShort(fd) //Last glyph ID in the range
			if err != nil {
				return coverageResult{}, err
			}
			_, err = t.ReadUShort(fd) //startCoverageIndex Coverage Index of first glyph ID in range (skip)
			if err != nil {
				return coverageResult{}, err
			}

			for j := startGlyphID; j <= endGlyphID; j++ {
				result.glyphIDs = append(result.glyphIDs, j) //set
			}
		}
	} else {
		return coverageResult{}, fmt.Errorf("Undefined Coverage Format")
	}

	return result, nil
}

//coverageResult read coverage result
type coverageResult struct {
	glyphIDs []uint //Array of glyph ID
}
