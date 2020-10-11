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
