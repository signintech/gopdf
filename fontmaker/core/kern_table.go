package core

//KernTable https://www.microsoft.com/typography/otspec/kern.htm
type KernTable struct {
	Version uint //for debug
	NTables uint //for debug
	Kerning KernMap
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
