package malayalam

import (
	"bytes"
	"regexp"
	"unicode"

	"github.com/signintech/gopdf/fontmaker/shaping/indic"
)

// syllable_type
const ConsonantSyllable = 0
const VowelSyllable = 1
const StandaloneCluster = 2
const BrokenCluster = 3
const NonIndicCluster = 4

//Malayalam https://docs.microsoft.com/en-us/typography/script-development/malayalam
type Malayalam struct {
}

//Reorder order
func (m Malayalam) Reorder(glyphindexs []uint, runes []rune) ([]uint, []rune, error) {
	otlinfos, err := m.findOTLInfo(glyphindexs, runes)
	if err != nil {
		return nil, nil, err
	}

	otlinfos, _, err = m.syllables(otlinfos)
	if err != nil {
		return nil, nil, err
	}

	err = m.initReordering(otlinfos)
	if err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}

func (m Malayalam) initReordering(otlInfos []otlInfo) error {
	var count = len(otlInfos)
	for i := 0; i < count; i++ {
		if otlInfos[i].pos == indic.PosBaseC {
			//c := otlInfos[i].char
		}
	}
	return nil
}

func (m Malayalam) updateConsonantPositions(otlInfos []otlInfo) error {

	return nil
}

func (m Malayalam) findOTLInfo(glyphindexs []uint, runes []rune) ([]otlInfo, error) {
	var posAndCats []otlInfo
	for i, r := range runes {
		g := glyphindexs[i]
		indicType, err := m.indicGetCategories(g, r)
		if err != nil {
			return nil, err
		}

		indicCat := indicType & 0x7F
		indicPos := indicType >> 8

		if r == 0x17D1 {
			indicCat = indic.OtX
		}

		if indicCat == indic.OtX && m.inRange(uint(r), 0x17CB, 0x17D3) { /* Khmer Various signs */
			/* These are like Top Matras. */
			indicCat = indic.OtM
			indicPos = indic.PosPostC
		}

		if r == 0x17C6 {
			indicCat = indic.OtN
		} /* Khmer Bindu doesn't like to be repositioned. */

		if r == 0x17D2 {
			indicCat = indic.OtCOENG
		} /* Khmer coeng */

		/* The spec says U+0952 is OT_A.	However, testing shows that Uniscribe
		 * treats U+0951..U+0952 all as OT_VD.
		 * TESTS:
		 * U+092E,U+0947,U+0952
		 * U+092E,U+0952,U+0947
		 * U+092E,U+0947,U+0951
		 * U+092E,U+0951,U+0947
		 * */
		if m.inRange(uint(r), 0x0951, 0x0954) {
			indicCat = indic.OtVD
		}

		if r == 0x200C {
			indicCat = indic.OtZWNJ
		} else if r == 0x200D {
			indicCat = indic.OtZWJ
		} else if r == 0x25CC {
			indicCat = indic.OtDOTTEDCIRCLE
		} else if r == 0x0A71 {
			indicCat = indic.OtSM
		}

		/* GURMUKHI ADDAK.	More like consonant medial. like 0A75. */
		if indicCat == indic.OtREPHA {
			if unicode.In(r, unicode.Mn) {
				indicCat = indic.OtN
			}
		}

		// Re-assign position.
		if (m.flag(int(indicCat)) & (m.flag(indic.OtC) | m.flag(indic.OtCM) | m.flag(indic.OtRA) | m.flag(indic.OtV) | m.flag(indic.OtNBSP) | m.flag(indic.OtDOTTEDCIRCLE))) > 0 { // = CONSONANT_FLAGS like is_consonant
			//if ($scriptblock == Ucdn::SCRIPT_KHMER) {
			if unicode.In(r, unicode.Khmer) {
				indicPos = indic.PosBelowC
			} else { /* Khmer differs from Indic here. */
				indicPos = indic.PosBaseC
			} /* Will recategorize later based on font lookups. */
			if indic.IsRuneRA(r) {
				indicCat = indic.OtRA
			}
		} else if indicCat == indic.OtM {
			indicPos = indic.MatraPosition(r, indicPos)
		} else if indicCat == indic.OtSM || indicCat == indic.OtVD {
			indicPos = indic.PosSMVD
		}

		if r == 0x0B01 {
			indicPos = indic.PosBeforeSUB
		} /* Oriya Bindu is BeforeSub in the spec. */

		posAndCats = append(posAndCats, otlInfo{
			glyphindex: g,
			char:       r,
			pos:        indicPos,
			cat:        indicCat,
		})
	}
	return posAndCats, nil
}

func (m Malayalam) syllables(info []otlInfo) ([]otlInfo, bool, error) {

	var buff bytes.Buffer
	for _, p := range info {
		buff.WriteString(indic.CategoryChar[p.cat])
	}

	var result = info
	brokenSyllables := false

	s := 0
	str := buff.String()
	size := len(str)
	syllableSerial := 1
	// CONSONANT_SYLLABLE Consonant syllable
	conRegex, err := regexp.Compile("^([CR]m*[N]?(H[ZJ]?|[ZJ]H))*[CR]m*[N]?[A]?(H[ZJ]?|[M]*[N]?[H]?)?[S]?[v]{0,2}")
	if err != nil {
		return nil, brokenSyllables, err
	}
	// VOWEL_SYLLABLE Vowel-based syllable
	vowelRegex, err := regexp.Compile("/^(RH|r)?V[N]?([ZJ]?H[CR]m*|J[CR]m*)?([M]*[N]?[H]?)?[S]?[v]{0,2}/")
	if err != nil {
		return nil, brokenSyllables, err
	}

	standaloneRegex, err := regexp.Compile("^(RH|r)?[sD][N]?([ZJ]?H[CR]m*)?([M]*[N]?[H]?)?[S]?[v]{0,2}")
	if err != nil {
		return nil, brokenSyllables, err
	}

	brokenRegex, err := regexp.Compile("^(RH|r)?[N]?([ZJ]?H[CR])?([M]*[N]?[H]?)?[S]?[v]{0,2}")
	if err != nil {
		return nil, brokenSyllables, err
	}
	//brokenSyllables := false
	for s < size {
		syllableLen := 1
		syllableType := NonIndicCluster //self::NON_INDIC_CLUSTER;
		sub := str[s:]
		if conRegex.MatchString(sub) {
			match := conRegex.FindAllString(sub, -1)
			syllableLen = len(match[0])
			syllableType = ConsonantSyllable
		} else if vowelRegex.MatchString(sub) {
			match := vowelRegex.FindAllString(sub, -1)
			syllableLen = len(match[0])
			syllableType = VowelSyllable
		} else if (s == 0 || (!unicode.Is(unicode.Categories["L"], info[s-1].char) && !unicode.Is(unicode.Categories["M"], info[s-1].char))) &&
			standaloneRegex.MatchString(sub) {
			match := standaloneRegex.FindAllString(sub, -1)
			syllableLen = len(match[0])
			syllableType = StandaloneCluster
		} else if brokenRegex.MatchString(sub) {
			match := brokenRegex.FindAllString(sub, -1)
			if len(match[0]) > 0 { // May match blank
				syllableLen = len(match[0])
				syllableType = BrokenCluster
				brokenSyllables = true
			}
		}

		for i := s; i < s+syllableLen; i++ {
			result[i].syllable = (syllableSerial << 4) | syllableType
		}

		s += syllableLen
		syllableSerial++
		if syllableSerial == 16 {
			syllableSerial = 1
		}
	}

	return result, brokenSyllables, nil
}

func (m Malayalam) indicGetCategories(glyph uint, r rune) (uint, error) {
	if 0x0900 <= r && r <= 0x0DFF {
		return indic.IndicTable[r-0x0900+0], nil // offset 0 for Most "indic"
	}
	if 0x1CD0 <= r && r <= 0x1D00 {
		return indic.IndicTable[r-0x1CD0+1152], nil // offset for Vedic extensions
	}
	if 0x1780 <= r && r <= 0x17FF {
		return indic.KhmerTable[r-0x1780], nil // Khmer
	}
	if r == 0x00A0 {
		return 3851, nil // (ISC_CP | (IMC_x << 8))
	}
	if r == 0x25CC {
		return 3851, nil // (ISC_CP | (IMC_x << 8))
	}
	return 3840, nil // (ISC_x | (IMC_x << 8))

}

func (m Malayalam) inRange(u uint, lo uint, hi uint) bool {
	if ((lo^hi)&lo) == 0 && ((lo^hi)&hi) == (lo^hi) && ((lo^hi)&((lo^hi)+1)) == 0 {
		return ((u & ^(lo ^ hi)) == lo)
	}
	return lo <= u && u <= hi
}

func (m Malayalam) flag(n int) int {
	return 1 << n
}

type otlInfo struct {
	glyphindex uint
	char       rune //rune
	pos        uint
	cat        uint
	syllable   int
}
