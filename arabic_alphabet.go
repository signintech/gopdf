package gopdf

// Harf is the Arabic meaning of Letter, Harf holds the Arabic character with its different representation forms (glyphs).
type Harf struct {
	Unicode, Isolated, Beggining, Midlle, Final rune
}

// Arabic Alphabet using the new Harf type.
var (
	ALEF_HAMZA_ABOVE = Harf{ // أ
		Unicode:   '\u0623',
		Isolated:  '\ufe83',
		Beggining: '\u0623',
		Midlle:    '\ufe84',
		Final:     '\ufe84'}

	ALEF = Harf{ // ا
		Unicode:   '\u0627',
		Isolated:  '\ufe8d',
		Beggining: '\u0627',
		Midlle:    '\ufe8e',
		Final:     '\ufe8e'}

	ALEF_MADDA_ABOVE = Harf{ // آ
		Unicode:   '\u0622',
		Isolated:  '\ufe81',
		Beggining: '\u0622',
		Midlle:    '\ufe82',
		Final:     '\ufe82'}

	HAMZA = Harf{ // ء
		Unicode:   '\u0621',
		Isolated:  '\ufe80',
		Beggining: '\u0621',
		Midlle:    '\u0621',
		Final:     '\u0621'}

	WAW_HAMZA_ABOVE = Harf{ // ؤ
		Unicode:   '\u0624',
		Isolated:  '\ufe85',
		Beggining: '\u0624',
		Midlle:    '\ufe86',
		Final:     '\ufe86'}

	ALEF_HAMZA_BELOW = Harf{ // أ
		Unicode:   '\u0625',
		Isolated:  '\ufe87',
		Beggining: '\u0625',
		Midlle:    '\ufe88',
		Final:     '\ufe88'}

	YEH_HAMZA_ABOVE = Harf{ // ئ
		Unicode:   '\u0626',
		Isolated:  '\ufe89',
		Beggining: '\ufe8b',
		Midlle:    '\ufe8c',
		Final:     '\ufe8a'}

	BEH = Harf{ // ب
		Unicode:   '\u0628',
		Isolated:  '\ufe8f',
		Beggining: '\ufe91',
		Midlle:    '\ufe92',
		Final:     '\ufe90'}

	PEH = Harf{ // پ
		Unicode:   '\u067e',
		Isolated:  '\ufb56',
		Beggining: '\ufb58',
		Midlle:    '\ufb59',
		Final:     '\ufb57'}

	TEH = Harf{ // ت
		Unicode:   '\u062A',
		Isolated:  '\ufe95',
		Beggining: '\ufe97',
		Midlle:    '\ufe98',
		Final:     '\ufe96'}

	TEH_MARBUTA = Harf{ // ة
		Unicode:   '\u0629',
		Isolated:  '\ufe93',
		Beggining: '\u0629',
		Midlle:    '\u0629',
		Final:     '\ufe94'}

	THEH = Harf{ // ث
		Unicode:   '\u062b',
		Isolated:  '\ufe99',
		Beggining: '\ufe9b',
		Midlle:    '\ufe9c',
		Final:     '\ufe9a'}

	JEEM = Harf{ // ج
		Unicode:   '\u062c',
		Isolated:  '\ufe9d',
		Beggining: '\ufe9f',
		Midlle:    '\ufea0',
		Final:     '\ufe9e'}

	TCHEH = Harf{ // چ
		Unicode:   '\u0686',
		Isolated:  '\ufb7a',
		Beggining: '\ufb7c',
		Midlle:    '\ufb7d',
		Final:     '\ufb7b'}

	HAH = Harf{ // ح
		Unicode:   '\u062d',
		Isolated:  '\ufea1',
		Beggining: '\ufea3',
		Midlle:    '\ufea4',
		Final:     '\ufea2'}

	KHAH = Harf{ // خ
		Unicode:   '\u062e',
		Isolated:  '\ufea5',
		Beggining: '\ufea7',
		Midlle:    '\ufea8',
		Final:     '\ufea6'}

	DAL = Harf{ // د
		Unicode:   '\u062f',
		Isolated:  '\ufea9',
		Beggining: '\u062f',
		Midlle:    '\ufeaa',
		Final:     '\ufeaa'}

	THAL = Harf{ // ذ
		Unicode:   '\u0630',
		Isolated:  '\ufeab',
		Beggining: '\u0630',
		Midlle:    '\ufeac',
		Final:     '\ufeac'}

	REH = Harf{ // ر
		Unicode:   '\u0631',
		Isolated:  '\ufead',
		Beggining: '\u0631',
		Midlle:    '\ufeae',
		Final:     '\ufeae'}

	JEH = Harf{
		Unicode:   '\u0698',
		Isolated:  '\ufb8a',
		Beggining: '\u0698',
		Midlle:    '\ufb8b',
		Final:     '\ufb8b',
	}

	ZAIN = Harf{ // ز
		Unicode:   '\u0632',
		Isolated:  '\ufeaf',
		Beggining: '\u0632',
		Midlle:    '\ufeb0',
		Final:     '\ufeb0'}

	SEEN = Harf{ // س
		Unicode:   '\u0633',
		Isolated:  '\ufeb1',
		Beggining: '\ufeb3',
		Midlle:    '\ufeb4',
		Final:     '\ufeb2'}

	SHEEN = Harf{ // ش
		Unicode:   '\u0634',
		Isolated:  '\ufeb5',
		Beggining: '\ufeb7',
		Midlle:    '\ufeb8',
		Final:     '\ufeb6'}

	SAD = Harf{ // ص
		Unicode:   '\u0635',
		Isolated:  '\ufeb9',
		Beggining: '\ufebb',
		Midlle:    '\ufebc',
		Final:     '\ufeba'}

	DAD = Harf{ // ض
		Unicode:   '\u0636',
		Isolated:  '\ufebd',
		Beggining: '\ufebf',
		Midlle:    '\ufec0',
		Final:     '\ufebe'}

	TAH = Harf{ // ط
		Unicode:   '\u0637',
		Isolated:  '\ufec1',
		Beggining: '\ufec3',
		Midlle:    '\ufec4',
		Final:     '\ufec2'}

	ZAH = Harf{ // ظ
		Unicode:   '\u0638',
		Isolated:  '\ufec5',
		Beggining: '\ufec7',
		Midlle:    '\ufec8',
		Final:     '\ufec6'}

	AIN = Harf{ // ع
		Unicode:   '\u0639',
		Isolated:  '\ufec9',
		Beggining: '\ufecb',
		Midlle:    '\ufecc',
		Final:     '\ufeca'}

	GHAIN = Harf{ // غ
		Unicode:   '\u063a',
		Isolated:  '\ufecd',
		Beggining: '\ufecf',
		Midlle:    '\ufed0',
		Final:     '\ufece'}

	FEH = Harf{ // ف
		Unicode:   '\u0641',
		Isolated:  '\ufed1',
		Beggining: '\ufed3',
		Midlle:    '\ufed4',
		Final:     '\ufed2'}

	QAF = Harf{ // ق
		Unicode:   '\u0642',
		Isolated:  '\ufed5',
		Beggining: '\ufed7',
		Midlle:    '\ufed8',
		Final:     '\ufed6'}

	KAF = Harf{ // ك
		Unicode:   '\u0643',
		Isolated:  '\ufed9',
		Beggining: '\ufedb',
		Midlle:    '\ufedc',
		Final:     '\ufeda'}

	KEHEH = Harf{ // ک
		Unicode:   '\u06a9',
		Isolated:  '\ufb8e',
		Beggining: '\ufb90',
		Midlle:    '\ufb91',
		Final:     '\ufb8f',
	}

	GAF = Harf{ // گ
		Unicode:   '\u06af',
		Isolated:  '\ufb92',
		Beggining: '\ufb94',
		Midlle:    '\ufb95',
		Final:     '\ufb93'}

	LAM = Harf{ // ل
		Unicode:   '\u0644',
		Isolated:  '\ufedd',
		Beggining: '\ufedf',
		Midlle:    '\ufee0',
		Final:     '\ufede'}

	MEEM = Harf{ // م
		Unicode:   '\u0645',
		Isolated:  '\ufee1',
		Beggining: '\ufee3',
		Midlle:    '\ufee4',
		Final:     '\ufee2'}

	NOON = Harf{ // ن
		Unicode:   '\u0646',
		Isolated:  '\ufee5',
		Beggining: '\ufee7',
		Midlle:    '\ufee8',
		Final:     '\ufee6'}

	HEH = Harf{ // ه
		Unicode:   '\u0647',
		Isolated:  '\ufee9',
		Beggining: '\ufeeb',
		Midlle:    '\ufeec',
		Final:     '\ufeea'}

	WAW = Harf{ // و
		Unicode:   '\u0648',
		Isolated:  '\ufeed',
		Beggining: '\u0648',
		Midlle:    '\ufeee',
		Final:     '\ufeee'}

	YEH = Harf{ // ی
		Unicode:   '\u06cc',
		Isolated:  '\ufbfc',
		Beggining: '\ufbfe',
		Midlle:    '\ufbff',
		Final:     '\ufbfd'}

	ARABICYEH = Harf{ // ي
		Unicode:   '\u064a',
		Isolated:  '\ufef1',
		Beggining: '\ufef3',
		Midlle:    '\ufef4',
		Final:     '\ufef2'}

	ALEF_MAKSURA = Harf{ // ى
		Unicode:   '\u0649',
		Isolated:  '\ufeef',
		Beggining: '\u0649',
		Midlle:    '\ufef0',
		Final:     '\ufef0'}

	TATWEEL = Harf{ // ـ
		Unicode:   '\u0640',
		Isolated:  '\u0640',
		Beggining: '\u0640',
		Midlle:    '\u0640',
		Final:     '\u0640'}

	LAM_ALEF = Harf{ // لا
		Unicode:   '\ufefb',
		Isolated:  '\ufefb',
		Beggining: '\ufefb',
		Midlle:    '\ufefc',
		Final:     '\ufefc'}

	LAM_ALEF_HAMZA_ABOVE = Harf{ // ﻷ
		Unicode:   '\ufef7',
		Isolated:  '\ufef7',
		Beggining: '\ufef7',
		Midlle:    '\ufef8',
		Final:     '\ufef8'}
)

var arabic_alphabet = []Harf{
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
