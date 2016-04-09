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
		buff.WriteString(fmt.Sprintf("\nleft : %c\n", left))
		for right, val := range kval {
			buff.WriteString(fmt.Sprintf("\tright : %c value= %d\n", right, val))
		}
	}
	return buff.String()
}

//KernMap kerning map   map[left]KernValue
type KernMap map[uint64]KernValue

//KernValue kerning values  map[right]value
type KernValue map[uint64]int64

/*
func (k KernValue) Debug() string {
	var buff bytes.Buffer
	for right, val := range k {
		buff.WriteString(fmt.Sprintf("\tright : %d value= %d\n", right, val))
	}
	return buff.String()
}
*/

//ValueByRight  get value by right
func (k KernValue) ValueByRight(right uint64) (bool, int64) {
	if val, ok := k[uint64(right)]; ok {
		return true, val
	}
	return false, 0
}
