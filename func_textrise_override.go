package gopdf

//FuncTextriseOverride override text rise
type FuncTextriseOverride func(
	leftRune rune,
	rightRune rune,
	fontsize int,
	allText string,
	currTextIndex int,
) float32
