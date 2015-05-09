package gopdf

type SubsetFontObj struct {
	CharacterToGlyphIndex map[rune]int
}

func (me *SubsetFontObj) Build() {

}

func (me *SubsetFontObj) AddChars(txt string) {
	for _, runeValue := range txt {
		if _, ok := me.CharacterToGlyphIndex[runeValue]; ok {
			continue
		}
		me.CharCodeToGlyphIndex(runeValue)
	}
}

func (me *SubsetFontObj) CharCodeToGlyphIndex(char rune) {

}
