package ucdn

import (
	"fmt"
	"testing"
)

func TestScript(t *testing.T) {
	var txt = "ന്മ"
	for _, r := range txt {
		sc := Script(r)
		fmt.Printf("%+v\n", sc)
	}
}
