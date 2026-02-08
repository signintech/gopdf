package gopdf

// Harf is the Arabic meaning of Letter, Harf holds the Arabic character with its different representation forms (glyphs).
type Harf struct {
	Unicode, Isolated, Beginning, Middle, Final rune
}

const (
	FATHA  rune = '\u064E' // َ (short a)
	DAMMA  rune = '\u064F' // ُ (short u)
	KASRA  rune = '\u0650' // ِ (short i)
	SUKUN  rune = '\u0652' // ْ (no vowel)
	SHADDA rune = '\u0651' // ّ (gemination/doubling)

	// Tanween (nutation)
	TANWEEN_FATH rune = '\u064B' // ً (an)
	TANWEEN_DAMM rune = '\u064C' // ٌ (un)
	TANWEEN_KASR rune = '\u064D' // ٍ (in)

	// Quranic / Extended marks
	SUPERSCRIPT_ALEF rune = '\u0670' // ٰ (dagger alef)
	MADDAH_ABOVE     rune = '\u0653' // ٓ (maddah)
	HAMZA_ABOVE      rune = '\u0654' // ٔ (hamza above)
	HAMZA_BELOW      rune = '\u0655' // ٕ (hamza below)
	SUBSCRIPT_ALEF   rune = '\u0656' // ٖ (subscript alef)
	INVERTED_DAMMA   rune = '\u0657' // ٗ (inverted damma)
	MARK_NOON_GHUNNA rune = '\u0658' // ٘ (noon ghunna)

	// Shadda + Vowel Ligatures (Arabic Presentation Forms-B)
	SHADDA_FATHA          rune = '\uFC60' // ﱠ
	SHADDA_DAMMA          rune = '\uFC61' // ﱡ
	SHADDA_KASRA          rune = '\uFC62' // ﱢ
	SHADDA_DAMMATAN       rune = '\uFC5E' // ﱞ (Shadda + Tanween Damm)
	SHADDA_KASRATAN       rune = '\uFC5F' // ﱟ (Shadda + Tanween Kasr)
	SHADDA_SUPERSCRIPT_ALEF rune = '\uFC63' // ﱣ
)

var tashkeelMarks = map[rune]bool{
	// Basic
	FATHA: true, DAMMA: true, KASRA: true,
	SHADDA: true, SUKUN: true,

	// Tanween
	TANWEEN_DAMM: true, TANWEEN_FATH: true, TANWEEN_KASR: true,

	// Quranic/Extended
	SUPERSCRIPT_ALEF: true, MADDAH_ABOVE: true,
	HAMZA_ABOVE: true, HAMZA_BELOW: true,
	SUBSCRIPT_ALEF: true, INVERTED_DAMMA: true,
	MARK_NOON_GHUNNA: true,
}

// shaddaLigatures maps vowels to their combined Shadda+Vowel ligature form
var shaddaLigatures = map[rune]rune{
	FATHA:            SHADDA_FATHA,
	DAMMA:            SHADDA_DAMMA,
	KASRA:            SHADDA_KASRA,
	TANWEEN_DAMM:     SHADDA_DAMMATAN,
	TANWEEN_KASR:     SHADDA_KASRATAN,
	SUPERSCRIPT_ALEF: SHADDA_SUPERSCRIPT_ALEF,
}

// GetShaddaLigature returns the combined Shadda+Vowel ligature for a given vowel.
// Returns 0 if no ligature exists for the vowel.
func GetShaddaLigature(vowel rune) rune {
	return shaddaLigatures[vowel]
}

// Arabic Alphabet using the new Harf type.
var (
	ALEF_HAMZA_ABOVE = Harf{ // أ
		Unicode:   '\u0623',
		Isolated:  '\ufe83',
		Beginning: '\u0623',
		Middle:    '\ufe84',
		Final:     '\ufe84'}

	ALEF = Harf{ // ا
		Unicode:   '\u0627',
		Isolated:  '\ufe8d',
		Beginning: '\u0627',
		Middle:    '\ufe8e',
		Final:     '\ufe8e'}

	ALEF_MADDA_ABOVE = Harf{ // آ
		Unicode:   '\u0622',
		Isolated:  '\ufe81',
		Beginning: '\u0622',
		Middle:    '\ufe82',
		Final:     '\ufe82'}

	HAMZA = Harf{ // ء
		Unicode:   '\u0621',
		Isolated:  '\ufe80',
		Beginning: '\u0621',
		Middle:    '\u0621',
		Final:     '\u0621'}

	WAW_HAMZA_ABOVE = Harf{ // ؤ
		Unicode:   '\u0624',
		Isolated:  '\ufe85',
		Beginning: '\u0624',
		Middle:    '\ufe86',
		Final:     '\ufe86'}

	ALEF_HAMZA_BELOW = Harf{ // أ
		Unicode:   '\u0625',
		Isolated:  '\ufe87',
		Beginning: '\u0625',
		Middle:    '\ufe88',
		Final:     '\ufe88'}

	YEH_HAMZA_ABOVE = Harf{ // ئ
		Unicode:   '\u0626',
		Isolated:  '\ufe89',
		Beginning: '\ufe8b',
		Middle:    '\ufe8c',
		Final:     '\ufe8a'}

	BEH = Harf{ // ب
		Unicode:   '\u0628',
		Isolated:  '\ufe8f',
		Beginning: '\ufe91',
		Middle:    '\ufe92',
		Final:     '\ufe90'}

	PEH = Harf{ // پ
		Unicode:   '\u067e',
		Isolated:  '\ufb56',
		Beginning: '\ufb58',
		Middle:    '\ufb59',
		Final:     '\ufb57'}

	TEH = Harf{ // ت
		Unicode:   '\u062A',
		Isolated:  '\ufe95',
		Beginning: '\ufe97',
		Middle:    '\ufe98',
		Final:     '\ufe96'}

	TEH_MARBUTA = Harf{ // ة
		Unicode:   '\u0629',
		Isolated:  '\ufe93',
		Beginning: '\u0629',
		Middle:    '\u0629',
		Final:     '\ufe94'}

	THEH = Harf{ // ث
		Unicode:   '\u062b',
		Isolated:  '\ufe99',
		Beginning: '\ufe9b',
		Middle:    '\ufe9c',
		Final:     '\ufe9a'}

	JEEM = Harf{ // ج
		Unicode:   '\u062c', // ج
		Isolated:  '\ufe9d', // ج
		Beginning: '\ufe9f', // جـ
		Middle:    '\ufea0', // ـجـ
		Final:     '\ufe9e'} // ـج

	TCHEH = Harf{ // چ
		Unicode:   '\u0686',
		Isolated:  '\ufb7a',
		Beginning: '\ufb7c',
		Middle:    '\ufb7d',
		Final:     '\ufb7b'}

	HAH = Harf{ // ح
		Unicode:   '\u062d',
		Isolated:  '\ufea1',
		Beginning: '\ufea3',
		Middle:    '\ufea4',
		Final:     '\ufea2'}

	KHAH = Harf{ // خ
		Unicode:   '\u062e',
		Isolated:  '\ufea5',
		Beginning: '\ufea7',
		Middle:    '\ufea8',
		Final:     '\ufea6'}

	DAL = Harf{ // د
		Unicode:   '\u062f',
		Isolated:  '\ufea9',
		Beginning: '\u062f',
		Middle:    '\ufeaa',
		Final:     '\ufeaa'}

	THAL = Harf{ // ذ
		Unicode:   '\u0630',
		Isolated:  '\ufeab',
		Beginning: '\u0630',
		Middle:    '\ufeac',
		Final:     '\ufeac'}

	REH = Harf{ // ر
		Unicode:   '\u0631',
		Isolated:  '\ufead',
		Beginning: '\u0631',
		Middle:    '\ufeae',
		Final:     '\ufeae'}

	JEH = Harf{
		Unicode:   '\u0698',
		Isolated:  '\ufb8a',
		Beginning: '\u0698',
		Middle:    '\ufb8b',
		Final:     '\ufb8b',
	}

	ZAIN = Harf{ // ز
		Unicode:   '\u0632',
		Isolated:  '\ufeaf',
		Beginning: '\u0632',
		Middle:    '\ufeb0',
		Final:     '\ufeb0'}

	SEEN = Harf{ // س
		Unicode:   '\u0633',
		Isolated:  '\ufeb1',
		Beginning: '\ufeb3',
		Middle:    '\ufeb4',
		Final:     '\ufeb2'}

	SHEEN = Harf{ // ش
		Unicode:   '\u0634',
		Isolated:  '\ufeb5',
		Beginning: '\ufeb7',
		Middle:    '\ufeb8',
		Final:     '\ufeb6'}

	SAD = Harf{ // ص
		Unicode:   '\u0635',
		Isolated:  '\ufeb9',
		Beginning: '\ufebb',
		Middle:    '\ufebc',
		Final:     '\ufeba'}

	DAD = Harf{ // ض
		Unicode:   '\u0636',
		Isolated:  '\ufebd',
		Beginning: '\ufebf',
		Middle:    '\ufec0',
		Final:     '\ufebe'}

	TAH = Harf{ // ط
		Unicode:   '\u0637',
		Isolated:  '\ufec1',
		Beginning: '\ufec3',
		Middle:    '\ufec4',
		Final:     '\ufec2'}

	ZAH = Harf{ // ظ
		Unicode:   '\u0638',
		Isolated:  '\ufec5',
		Beginning: '\ufec7',
		Middle:    '\ufec8',
		Final:     '\ufec6'}

	AIN = Harf{ // ع
		Unicode:   '\u0639',
		Isolated:  '\ufec9',
		Beginning: '\ufecb',
		Middle:    '\ufecc',
		Final:     '\ufeca'}

	GHAIN = Harf{ // غ
		Unicode:   '\u063a',
		Isolated:  '\ufecd',
		Beginning: '\ufecf',
		Middle:    '\ufed0',
		Final:     '\ufece'}

	FEH = Harf{ // ف
		Unicode:   '\u0641',
		Isolated:  '\ufed1',
		Beginning: '\ufed3',
		Middle:    '\ufed4',
		Final:     '\ufed2'}

	QAF = Harf{ // ق
		Unicode:   '\u0642',
		Isolated:  '\ufed5',
		Beginning: '\ufed7',
		Middle:    '\ufed8',
		Final:     '\ufed6'}

	KAF = Harf{ // ك
		Unicode:   '\u0643',
		Isolated:  '\ufed9',
		Beginning: '\ufedb',
		Middle:    '\ufedc',
		Final:     '\ufeda'}

	KEHEH = Harf{ // ک
		Unicode:   '\u06a9',
		Isolated:  '\ufb8e',
		Beginning: '\ufb90',
		Middle:    '\ufb91',
		Final:     '\ufb8f',
	}

	GAF = Harf{ // گ
		Unicode:   '\u06af',
		Isolated:  '\ufb92',
		Beginning: '\ufb94',
		Middle:    '\ufb95',
		Final:     '\ufb93'}

	LAM = Harf{ // ل
		Unicode:   '\u0644',
		Isolated:  '\ufedd',
		Beginning: '\ufedf',
		Middle:    '\ufee0',
		Final:     '\ufede'}

	MEEM = Harf{ // م
		Unicode:   '\u0645',
		Isolated:  '\ufee1',
		Beginning: '\ufee3',
		Middle:    '\ufee4',
		Final:     '\ufee2'}

	NOON = Harf{ // ن
		Unicode:   '\u0646',
		Isolated:  '\ufee5',
		Beginning: '\ufee7',
		Middle:    '\ufee8',
		Final:     '\ufee6'}

	HEH = Harf{ // ه
		Unicode:   '\u0647',
		Isolated:  '\ufee9',
		Beginning: '\ufeeb',
		Middle:    '\ufeec',
		Final:     '\ufeea'}

	WAW = Harf{ // و
		Unicode:   '\u0648',
		Isolated:  '\ufeed',
		Beginning: '\u0648',
		Middle:    '\ufeee',
		Final:     '\ufeee'}

	YEH = Harf{ // ی
		Unicode:   '\u06cc',
		Isolated:  '\ufbfc',
		Beginning: '\ufbfe',
		Middle:    '\ufbff',
		Final:     '\ufbfd'}

	ARABICYEH = Harf{ // ي
		Unicode:   '\u064a',
		Isolated:  '\ufef1',
		Beginning: '\ufef3',
		Middle:    '\ufef4',
		Final:     '\ufef2'}

	ALEF_MAKSURA = Harf{ // ى
		Unicode:   '\u0649',
		Isolated:  '\ufeef',
		Beginning: '\u0649',
		Middle:    '\ufef0',
		Final:     '\ufef0'}

	TATWEEL = Harf{ // ـ
		Unicode:   '\u0640',
		Isolated:  '\u0640',
		Beginning: '\u0640',
		Middle:    '\u0640',
		Final:     '\u0640'}

	LAM_ALEF = Harf{ // لا
		Unicode:   '\ufefb',
		Isolated:  '\ufefb',
		Beginning: '\ufefb',
		Middle:    '\ufefc',
		Final:     '\ufefc'}

	LAM_ALEF_HAMZA_ABOVE = Harf{ // ﻷ
		Unicode:   '\ufef7',
		Isolated:  '\ufef7',
		Beginning: '\ufef7',
		Middle:    '\ufef8',
		Final:     '\ufef8'}
)

var arabicAlphabet = map[rune]Harf{}

var arabicAlphabetCollection = []Harf{
	ALEF_HAMZA_ABOVE,
	ALEF,
	ALEF_MADDA_ABOVE,
	HAMZA,
	WAW_HAMZA_ABOVE,
	ALEF_HAMZA_BELOW,
	YEH_HAMZA_ABOVE,
	BEH,
	PEH,
	TEH,
	TEH_MARBUTA,
	THEH,
	JEEM,
	TCHEH,
	HAH,
	KHAH,
	DAL,
	THAL,
	REH,
	JEH,
	ZAIN,
	SEEN,
	SHEEN,
	SAD,
	DAD,
	TAH,
	ZAH,
	AIN,
	GHAIN,
	FEH,
	QAF,
	KAF,
	KEHEH,
	GAF,
	LAM,
	MEEM,
	NOON,
	HEH,
	WAW,
	YEH,
	ARABICYEH,
	ALEF_MAKSURA,
	TATWEEL,
	LAM_ALEF,
	LAM_ALEF_HAMZA_ABOVE,
}

func init() {
	for _, harf := range arabicAlphabetCollection {
		// Map all forms to the Harf struct
		arabicAlphabet[harf.Unicode] = harf
		arabicAlphabet[harf.Isolated] = harf
		arabicAlphabet[harf.Beginning] = harf
		arabicAlphabet[harf.Middle] = harf
		arabicAlphabet[harf.Final] = harf
	}
}

// use map for faster lookups.
var rightJoiningOnlyLetters = map[Harf]bool{
	ALEF_HAMZA_ABOVE: true,
	ALEF_MADDA_ABOVE: true,
	ALEF:             true,
	HAMZA:            true,
	WAW_HAMZA_ABOVE:  true,
	ALEF_HAMZA_BELOW: true,
	TEH_MARBUTA:      true,
	DAL:              true,
	THAL:             true,
	REH:              true,
	ZAIN:             true,
	WAW:              true,
	ALEF_MAKSURA:     true}
