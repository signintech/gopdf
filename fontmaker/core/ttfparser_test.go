package core

import "testing"

func TestParseGSUB(t *testing.T) {
	var tp TTFParser
	err := tp.Parse("../../test/res/NotoSansMalayalamUI-Regular.ttf")
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
