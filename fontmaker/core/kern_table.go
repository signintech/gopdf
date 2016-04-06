package core

import (
	"bytes"
	"fmt"
)

//KernTable https://www.microsoft.com/typography/otspec/kern.htm
type KernTable struct {
	Version uint64 //for debug
	NTables uint64 //for debug
	Kerning KernMap
}

func (k KernTable) debug() string {
	var buff bytes.Buffer
	for left, kval := range k.Kerning {
		buff.WriteString(fmt.Sprintf("\nleft : %d\n", left))
		for right, val := range kval {
			buff.WriteString(fmt.Sprintf("\tright : %d value= %d\n", right, val))
		}
	}
	return buff.String()
}

/*
//KernSubTable kern sub table https://www.microsoft.com/typography/otspec/kern.htm
type KernSubTable struct {
	Version  uint64 //Kern subtable version number
	Length   uint64 //Length of the subtable, in bytes (including this header).
	Coverage uint64 //What type of information is contained in this table.

}*/

//KernMap kerning map   map[left]KernValue
type KernMap map[uint64]KernValue

//KernValue kerning values  map[right]value
type KernValue map[uint64]int64
