package core

import "bytes"

//2.1 Multiple Substitution Format 1
func (t *TTFParser) parseGSUBLookupListTableSubTableLookupType2Format1(
	fd *bytes.Reader,
	offset int64,
	substFormat uint,
	gdefResult ParseGDEFResult,
) (
	GSUBLookupSubTableType1Format1,
	error,
) {

	result := GSUBLookupSubTableType1Format1{}

	return result, nil
}
