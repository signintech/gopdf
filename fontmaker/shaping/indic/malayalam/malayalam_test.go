package malayalam

import (
	"testing"

	"github.com/signintech/gopdf/fontmaker/core"
)

func TestMalayalamReorder(t *testing.T) {

	parser := core.TTFParser{}
	parser.Parse("../test/res/NotoSansMalayalamUI-Regular.ttf")

	g, r := toFakeGlyphIndex("ബ്രഹ്മ")
	var m Malayalam
	_, _, err := m.Reorder(g, r)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func toFakeGlyphIndex(txt string) ([]uint, []rune) {
	var fakeGlyphindexs []uint
	var runes []rune
	for i, r := range txt {
		fakeGlyphindexs = append(fakeGlyphindexs, uint(i))
		runes = append(runes, r)
	}
	return fakeGlyphindexs, runes
}

/*
อันที่ถูก
pos= 4 cat=1
pos= 6 cat=4
pos= 4 cat=16
pos= 4 cat=1
pos= 6 cat=4
pos= 4 cat=1
*/
