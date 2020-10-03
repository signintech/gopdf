package core

//GSUBLookupSubTableTyper lookup sub table
type gsubLookupSubTableTyper interface {
	LookupType() uint //lookup type
	Format() uint     //Format identifier:
}
