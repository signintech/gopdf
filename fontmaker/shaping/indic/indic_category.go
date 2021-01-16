package indic

// indic_category
const OtX = 0
const OtC = 1
const OtV = 2
const OtN = 3
const OtH = 4
const OtZWNJ = 5
const OtZWJ = 6
const OtM = 7 /* Matra or Dependent Vowel */
const OtSM = 8
const OtVD = 9
const OtA = 10
const OtNBSP = 11
const OtDOTTEDCIRCLE = 12 /* Not in the spec, but special in Uniscribe. /Very very/ special! */
const OtRS = 13           /* Register Shifter, used in Khmer OT spec */
const OtCOENG = 14
const OtREPHA = 15

const OtRA = 16 /* Not explicitly listed in the OT spec, but used in the grammar. */
const OtCM = 17

var CategoryChar = []string{
	"x",
	"C",
	"V",
	"N",
	"H",
	"Z",
	"J",
	"M",
	"S",
	"v",
	"A", /* Spec gives Andutta U+0952 as OT_A. However, testing shows that Uniscribe
	 * treats U+0951..U+0952 all as OT_VD - see set_indic_properties */
	"s",
	"D",
	"F", /* Register shift Khmer only */
	"G", /* Khmer only */
	"r", /* 0D4E (dot reph) only one in Malayalam */
	"R",
	"m", /* Consonant medial only used in Indic 0A75 in Gurmukhi  (0A00..0A7F)  : also in Lao, Myanmar, Tai Tham, Javanese & Cham  */
}
