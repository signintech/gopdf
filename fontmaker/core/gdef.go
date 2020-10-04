package core

type gdefHeader struct {
	majorVersion             uint  // Major version of the GDEF table, = 1
	minorVersion             uint  //Minor version of the GDEF table, = 3
	glyphClassDefOffset      int64 //Offset to class definition table for glyph type, from beginning of GDEF header (may be NULL)
	attachListOffset         int64 //Offset to attachment point list table, from beginning of GDEF header (may be NULL)
	ligCaretListOffset       int64 //Offset to ligature caret list table, from beginning of GDEF header (may be NULL)
	markAttachClassDefOffset int64 //Offset to class definition table for mark attachment type, from beginning of GDEF header (may be NULL)
	//for minorVersion 2,3
	markGlyphSetsDefOffset int64 //Offset to the table of mark glyph set definitions, from beginning of GDEF header (may be NULL)
	//for minorVersion 3
	itemVarStoreOffset int64 //	Offset to the Item Variation Store table, from beginning of GDEF header (may be NULL)
}
