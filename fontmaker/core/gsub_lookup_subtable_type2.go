package core

//GSUBLookupSubTableType2Format1 -> 2.1 Multiple Substitution Format 1
type GSUBLookupSubTableType2Format1 struct {
	substFormat     uint
	coverageOffset  int64   //Offset to Coverage table, from beginning of substitution subtable
	sequenceCount   uint    //Number of Sequence table offsets in the sequenceOffsets array
	sequenceOffsets []int64 //Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	sequenceTable   []GSUBLookupSequenceTable
}

//LookupType get lookup type
func (g GSUBLookupSubTableType2Format1) LookupType() uint {
	return 2
}

//Format get format
func (g GSUBLookupSubTableType2Format1) Format() uint {
	return 1
}

//GSUBLookupSequenceTable Sequence table
type GSUBLookupSequenceTable struct {
	glyphCount         uint   //Number of glyph IDs in the substituteGlyphIDs array. This must always be greater than 0.
	substituteGlyphIDs []uint //String of glyph IDs to substitute
}
