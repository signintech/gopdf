package gopdf

// FuncKernOverride  return your custome pair value
type FuncKernOverride func(
	leftRune rune,
	rightRune rune,
	leftPair uint,
	rightPair uint,
	pairVal int16,
) int16
