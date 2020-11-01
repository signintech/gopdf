package core

//GSUBLookupSubTableType2Format1 -> 2.1 Multiple Substitution Format 1
type GSUBLookupSubTableType2Format1 struct {
	substFormat     uint
	coverageOffset  int64                     //Offset to Coverage table, from beginning of substitution subtable
	sequenceCount   uint                      //Number of Sequence table offsets in the sequenceOffsets array
	sequenceOffsets []GSUBLookupSequenceTable //Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
}

//GSUBLookupSequenceTable Sequence table
type GSUBLookupSequenceTable struct {
	glyphCount         uint   //Number of glyph IDs in the substituteGlyphIDs array. This must always be greater than 0.
	substituteGlyphIDs []uint //String of glyph IDs to substitute
}
