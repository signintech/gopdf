package malayalam

import (
	"testing"
)

func TestMalayalamReorder(t *testing.T) {
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
