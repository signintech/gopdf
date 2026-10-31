package gopdf

import (
	"bytes"
	"strings"
	"testing"
)

func TestUnicodeMapWritesSurrogatePairs(t *testing.T) {
	subset := &SubsetFontObj{CharacterToGlyphIndex: NewMapOfCharacterToGlyphIndex()}
	subset.CharacterToGlyphIndex.Set('A', 36)
	subset.CharacterToGlyphIndex.Set(0x10300, 500) // OLD ITALIC LETTER A
	subset.CharacterToGlyphIndex.Set(0x1F600, 501) // GRINNING FACE
	u := UnicodeMap{PtrToSubsetFontObj: subset}

	var buf bytes.Buffer
	if err := u.write(&buf, 1); err != nil {
		t.Fatal(err)
	}
	cmap := buf.String()
	for _, want := range []string{"<0024><0024><0041>", "<01F4><01F4><D800DF00>", "<01F5><01F5><D83DDE00>"} {
		if !strings.Contains(cmap, want) {
			t.Errorf("ToUnicode CMap lacks %s:\n%s", want, cmap)
		}
	}
}
