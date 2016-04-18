package core

import (
	"bytes"
	"fmt"
)

//KernTable https://www.microsoft.com/typography/otspec/kern.htm
type KernTable struct {
	Version uint //for debug
	NTables uint //for debug
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
type KernMap map[uint]KernValue

//KernValue kerning values  map[right]value
type KernValue map[uint]int16

//ValueByRight  get value by right
func (k KernValue) ValueByRight(right uint) (bool, int16) {
	if val, ok := k[uint(right)]; ok {
		return true, val
	}
	return false, 0
}
