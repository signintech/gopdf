package core

//GSUBLookupSubTableType3Format1 3.1 Alternate Substitution Format 1
type GSUBLookupSubTableType3Format1 struct {
	coverageOffset      int64   //Offset to Coverage table, from beginning of substitution subtable
	alternateSetCount   uint    //	Number of AlternateSet tables
	alternateSetOffsets []int64 //Array of offsets to AlternateSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	alternateSetTables  []AlternateSetTable
}

//LookupType get lookup type
func (g GSUBLookupSubTableType3Format1) LookupType() uint {
	return 3
}

//Format get format
func (g GSUBLookupSubTableType3Format1) Format() uint {
	return 1
}

//AlternateSetTable AlternateSet table
type AlternateSetTable struct {
	glyphCount        uint   //	Number of glyph IDs in the alternateGlyphIDs array
	alternateGlyphIDs []uint //Array of alternate glyph IDs, in arbitrary order
}
