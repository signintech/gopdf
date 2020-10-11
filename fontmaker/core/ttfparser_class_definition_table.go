package core

import (
	"bytes"
	"fmt"
)

//parseClassDefinitionTable parse Glyph Class Definition Table
//https://docs.microsoft.com/en-us/typography/opentype/spec/chapter2#classDefTbl
func (t *TTFParser) parseClassDefinitionTable(fd *bytes.Reader, offset int64) (parseClassDefinitionTableResult, error) {
	_, err := fd.Seek(offset, 0)
	if err != nil {
		return parseClassDefinitionTableResult{}, err
	}

	result := initParseClassDefinitionTableResult()
	classFormat, err := t.ReadUShort(fd)
	if err != nil {
		return parseClassDefinitionTableResult{}, err
	}

	if classFormat == 1 {
		startGlyphID, err := t.ReadUShort(fd)
		if err != nil {
			return parseClassDefinitionTableResult{}, err
		}
		glyphCount, err := t.ReadUShort(fd)
		if err != nil {
			return parseClassDefinitionTableResult{}, err
		}

		for i := uint(0); i < glyphCount; i++ {
			classValue, err := t.ReadUShort(fd)
			if err != nil {
				return parseClassDefinitionTableResult{}, err
			}
			result.append(classValue, startGlyphID+i)
		}

	} else if classFormat == 2 {
		classRangeCount, err := t.ReadUShort(fd)
		if err != nil {
			return parseClassDefinitionTableResult{}, err
		}
		for i := uint(0); i < classRangeCount; i++ {
			startGlyphID, err := t.ReadUShort(fd)
			if err != nil {
				return parseClassDefinitionTableResult{}, err
			}
			endGlyphID, err := t.ReadUShort(fd)
			if err != nil {
				return parseClassDefinitionTableResult{}, err
			}
			classValue, err := t.ReadUShort(fd)
			if err != nil {
				return parseClassDefinitionTableResult{}, err
			}

			for k := startGlyphID; k <= endGlyphID; k++ {
				result.append(classValue, k)
			}
		}
	}

	return result, nil
}

//
type parseClassDefinitionTableResult struct {
	mapClassWithGlyphIDs map[uint]([]uint) // map[class] array of glyphID

}

func initParseClassDefinitionTableResult() parseClassDefinitionTableResult {
	p := parseClassDefinitionTableResult{}
	p.mapClassWithGlyphIDs = make(map[uint]([]uint))
	return p
}

func (p *parseClassDefinitionTableResult) append(class uint, glyphID uint) {
	p.mapClassWithGlyphIDs[class] = append(p.mapClassWithGlyphIDs[class], glyphID)
}

func (p *parseClassDefinitionTableResult) isContainClass(class uint) bool {
	if mc, ok := p.mapClassWithGlyphIDs[class]; ok {
		if len(mc) > 0 {
			return true
		}
	}
	return false
}

func (p *parseClassDefinitionTableResult) GlyphClassBases() []uint {
	if glyphIDs, ok := p.mapClassWithGlyphIDs[1]; ok {
		return glyphIDs
	}
	return []uint{}
}

func (p *parseClassDefinitionTableResult) GlyphClassLigatures() []uint {
	if glyphIDs, ok := p.mapClassWithGlyphIDs[2]; ok {
		return glyphIDs
	}
	return []uint{}
}

func (p *parseClassDefinitionTableResult) GlyphClassMarks() []uint {
	if glyphIDs, ok := p.mapClassWithGlyphIDs[3]; ok {
		return glyphIDs
	}
	return []uint{}
}

func (p *parseClassDefinitionTableResult) GlyphClassComponents() []uint {
	if glyphIDs, ok := p.mapClassWithGlyphIDs[4]; ok {
		return glyphIDs
	}
	return []uint{}
}

func (p *parseClassDefinitionTableResult) debug() string {
	var buff bytes.Buffer
	for cls, glyphIDs := range p.mapClassWithGlyphIDs {
		fmt.Fprintf(&buff, "{\nclass:%d,\n", cls)
		fmt.Fprintf(&buff, " glyphIDs:[")
		for i, glyphID := range glyphIDs {
			if i > 0 {
				fmt.Fprintf(&buff, ",")
			}
			fmt.Fprintf(&buff, "%d", glyphID)
		}
		fmt.Fprintf(&buff, "]")
		fmt.Fprintf(&buff, "\n}\n")
	}

	return buff.String()
}
