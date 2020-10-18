package sliceutil

//DiffUInt diff slice
func DiffUInt(aa []uint, bb []uint) []uint {
	var result []uint
	for _, a := range aa {
		found := false
		for _, b := range bb {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			result = append(result, a)
		}
	}
	return result
}

//IndexUint find index b inside a ( return -1 if not found)
func IndexUint(aa []uint, b uint) int {
	for i, a := range aa {
		if a == b {
			return i
		}
	}

	return -1
}

//ContainSliceUint find slice inside slice
func ContainSliceUint(slice []uint, subSlice []uint) []ContainSliceUintResult {

	sliceLen := len(slice)
	subSliceLen := len(subSlice)
	if subSliceLen <= 0 {
		return []ContainSliceUintResult{} //empty
	}

	var result []ContainSliceUintResult
	for i, s := range slice {
		sliceLenLeft := sliceLen - i
		//fmt.Printf("%d %d\n", sliceLenLeft, subSliceLen)
		if s == subSlice[0] && sliceLenLeft >= subSliceLen {
			match := true
			for j := 1; j < subSliceLen; j++ {
				if subSlice[j] != slice[i+j] {
					match = false
					break
				}
			}
			if match {
				result = append(result, ContainSliceUintResult{
					FirstIndex: i,
					Length:     subSliceLen,
				})
			}
		}
	}
	return result
}

type ContainSliceUintResult struct {
	FirstIndex int
	Length     int
}
