package malayalam

import (
	"fmt"
	"unicode"

	"github.com/signintech/gopdf/fontmaker/shaping/indic"
)

//Malayalam https://docs.microsoft.com/en-us/typography/script-development/malayalam
type Malayalam struct {
}

//Reorder order
func (m Malayalam) Reorder(glyphindexs []uint, runes []rune) ([]uint, []rune, error) {

	for i, r := range runes {
		g := glyphindexs[i]
		indicType, err := m.indicGetCategories(g, r)
		if err != nil {
			return nil, nil, err
		}

		indicCat := indicType & 0x7F
		indicPos := indicType >> 8

		if r == 0x17D1 {
			indicCat = indic.OtX
		}

		xxx := m.inRange(0x17CE, 0x17CB, 0x17D3)
		fmt.Printf("%v\n", xxx)

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

		 if ((self::FLAG($cat) & (self::FLAG(self::OT_C) | self::FLAG(self::OT_CM) | self::FLAG(self::OT_RA) | self::FLAG(self::OT_V) | self::FLAG(self::OT_NBSP) | self::FLAG(self::OT_DOTTEDCIRCLE)))) { // = CONSONANT_FLAGS like is_consonant
		 	if ($scriptblock == Ucdn::SCRIPT_KHMER) {
		 		$pos = self::POS_BELOW_C;
		 	} /* Khmer differs from Indic here. */
		 	else {
		 		$pos = self::POS_BASE_C;
		 	} /* Will recategorize later based on font lookups. *
		 	if (self::is_ra($u)) {
		 		$cat = self::OT_RA;
		 	}
		 } elseif ($cat == self::OT_M) {
		 	$pos = self::matra_position($u, $pos);
		 } elseif ($cat == self::OT_SM || $cat == self::OT_VD) {
		 	$pos = self::POS_SMVD;
		 }

		// if ($u == 0x0B01) {
		// 	$pos = self::POS_BEFORE_SUB;
		// } /* Oriya Bindu is BeforeSub in the spec. */

		// $info['indic_category'] = $cat;
		// $info['indic_position'] = $pos;
	}

	return nil, nil, nil
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
