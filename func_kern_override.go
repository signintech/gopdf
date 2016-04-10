package gopdf

//FuncKernOverride  return your custome pair value
type FuncKernOverride func(
	leftRune rune,
	rightRune rune,
	leftPair uint64,
	rightPair uint64,
	pairVal int64,
) int64
