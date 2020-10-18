package sliceutil

import (
	"testing"
)

func TestContainUint(t *testing.T) {
	a := []uint{1, 2, 3, 2, 3, 4, 2, 3}
	aSub := []uint{2, 3}
	result := ContainSliceUint(a, aSub)
	if len(result) != 3 {
		t.Fatalf("len(result) != 3 A")
		return
	}

}
