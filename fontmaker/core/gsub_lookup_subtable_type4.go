package core

//LookupType 4: Ligature Substitution Subtable

//GSUBLookupSubTableType4Format1 4.1 Ligature Substitution Format 1
type GSUBLookupSubTableType4Format1 struct {
	coverageOffset     int64   //Offset to Coverage table, from beginning of substitution subtable
	ligatureSetCount   uint    //Number of LigatureSet tables
	ligatureSetOffsets []int64 //Array of offsets to LigatureSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	//table
	ligatureSetTables []LigatureSetTable
}

//LookupType 4: Ligature Substitution Subtable
func (g GSUBLookupSubTableType4Format1) LookupType() uint {
	return 4
}

//Format identifier: format = 1
func (g GSUBLookupSubTableType4Format1) Format() uint {
	return 1
}

//LigatureSetTable LigatureSet table: All ligatures beginning with the same glyph
type LigatureSetTable struct {
	ligatureCount   uint    //Number of Ligature tables
	ligatureOffsets []int64 //Array of offsets to Ligature tables. Offsets are from beginning of LigatureSet table, ordered by preference.
	ligatureTables  []LigatureTable
}

//LigatureTable Ligature table: Glyph components for one ligature
type LigatureTable struct {
	ligatureGlyph     uint   //glyph ID of ligature to substitute
	componentCount    uint   //Number of components in the ligature
	componentGlyphIDs []uint //Array of component glyph IDs — start with the second component, ordered in writing direction
}
