package gopdf

//FuncTextriseOverride override text rise
type FuncTextriseOverride func(
	leftRune rune,
	rightRune rune,
	leftPair uint,
	rightPair uint,
	fontsize int,
) float32
