package core

type GSUBLookupTable struct {
	lookupType          uint    //Different enumerations for GSUB and GPOS
	lookupFlag          uint    //Lookup qualifiers
	subTableCount       uint    //Number of subtables for this lookup
	subtableOffsets     []int64 //Array of offsets to lookup subtables, from beginning of Lookup table
	markFilteringSet    uint    //Index (base 0) into GDEF mark glyph sets structure. This field is only present if bit useMarkFilteringSet of lookup flags is set.
	gsubLookupSubTables []gsubLookupSubTableTyper
}
