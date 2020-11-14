package core

//GSUBLookupSubTableType5Format1 5.1 Context Substitution Format 1: Simple Glyph Contexts
type GSUBLookupSubTableType5Format1 struct {
	coverageOffset    int64   //Offset to Coverage table, from beginning of substitution subtable
	subRuleSetCount   uint    //Number of SubRuleSet tables — must equal glyphCount in Coverage table
	subRuleSetOffsets []int64 //Array of offsets to SubRuleSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	subRuleSets       []SubRuleSetTable
}

//LookupType 5
func (g GSUBLookupSubTableType5Format1) LookupType() uint {
	return 5
}

//Format identifier: format = 1
func (g GSUBLookupSubTableType5Format1) Format() uint {
	return 1
}

//SubRuleSetTable table: All contexts beginning with the same glyph
type SubRuleSetTable struct {
	subRuleCount   uint    //Number of SubRule tables
	subRuleOffsets []int64 //Array of offsets to SubRule tables. Offsets are from beginning of SubRuleSet table, ordered by preference
	subRules       []SubRuleTable
}

//SubRuleTable One simple context definition
type SubRuleTable struct {
	glyphCount         uint                //Total number of glyphs in input glyph sequence — includes the first glyph.
	substitutionCount  uint                //Number of SubstLookupRecords
	inputSequence      []uint              //Array of input glyph IDs — start with second glyph
	substLookupRecords []SubstLookupRecord //Array of SubstLookupRecords, in design order
}

//SubstLookupRecord Substitution Lookup Record
type SubstLookupRecord struct {
	glyphSequenceIndex uint
	lookupListIndex    uint
}
