package core

import "bytes"

//GSUBLookupSubTableTyper lookup sub table
type gsubLookupSubTabler interface {
	LookupType() uint //lookup type
	Format() uint     //Format identifier:
	processSubTable(t *TTFParser,
		fd *bytes.Reader,
		table GSUBLookupTable,
		gdefResult ParseGDEFResult) error
}
