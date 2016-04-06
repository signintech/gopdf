package core

//KernTable https://www.microsoft.com/typography/otspec/kern.htm
type KernTable struct {
	Version   uint64
	NTables   uint64
	Subtables []KernSubTable
}

//KernSubTable kern sub table https://www.microsoft.com/typography/otspec/kern.htm
type KernSubTable struct {
	Version  uint64 //Kern subtable version number
	Length   uint64 //Length of the subtable, in bytes (including this header).
	Coverage uint64 //What type of information is contained in this table.
}
