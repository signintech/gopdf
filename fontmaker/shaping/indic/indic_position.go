package indic

// indic_position
const PosSTART = 0

const PosRaTOBecomeReph = 1
const PosPreM = 2
const PosPreC = 3

const PosBaseC = 4
const PosAfterMain = 5

const PosAboveC = 6

const PosBeforeSUB = 7
const PosBelowC = 8
const PosAfterSUB = 9

const PosBeforePOST = 10
const PosPostC = 11
const PosAfterPOST = 12

const PosFinalC = 13
const PosSMVD = 14

const PosEND = 15

func matraPosLeft(r rune) uint {
	return PosPreM
}

func MatraPosition(r rune, pos uint) uint {
	return matraPosition(r, pos)
}

func matraPosition(r rune, pos uint) uint {
	if pos == PosPreC {
		return matraPosLeft(r)
	} else if pos == PosPostC {
		return matraPosRight(r)
	} else if pos == PosAboveC {
		return matraPosTop(r)
	} else if pos == PosBelowC {
		return matraPosTopBottom(r)
	}
	return pos
}

func inHalfBlock(r rune, base uint) bool {
	return ((uint(r) & ^uint(0x7F)) == base)
}

func isDeva(r rune) bool {
	return inHalfBlock(r, 0x0900)
}

func isBeng(r rune) bool {
	return inHalfBlock(r, 0x0980)
}

func isGuru(r rune) bool {
	return inHalfBlock(r, 0x0A00)
}

func isGujr(r rune) bool {
	return inHalfBlock(r, 0x0A80)
}

func isOrya(r rune) bool {
	return inHalfBlock(r, 0x0B00)
}

func isTaml(r rune) bool {
	return inHalfBlock(r, 0x0B80)
}

func isTelu(r rune) bool {
	return inHalfBlock(r, 0x0C00)
}

func isKnda(r rune) bool {
	return inHalfBlock(r, 0x0C80)
}

func isMlym(r rune) bool {
	return inHalfBlock(r, 0x0D00)
}

func isSinh(r rune) bool {
	return inHalfBlock(r, 0x0D80)
}

func isKhmr(r rune) bool {
	return inHalfBlock(r, 0x1780)
}

func matraPosRight(r rune) uint {
	if isDeva(r) {
		return PosAfterSUB
	}
	if isBeng(r) {
		return PosAfterPOST
	}
	if isGuru(r) {
		return PosAfterPOST
	}
	if isGujr(r) {
		return PosAfterPOST
	}
	if isOrya(r) {
		return PosAfterPOST
	}
	if isTaml(r) {
		return PosAfterPOST
	}
	if isTelu(r) {
		if r <= 0x0C42 {
			return PosBeforeSUB
		}
		return PosAfterSUB
	}
	if isKnda(r) {
		if r < 0x0CC3 || r > 0xCD6 {
			return PosBeforeSUB
		}
		return PosAfterSUB
	}
	if isMlym(r) {
		return PosAfterPOST
	}
	if isSinh(r) {
		return PosAfterSUB
	}
	if isKhmr(r) {
		return PosAfterPOST
	}
	return PosAfterSUB
}

func matraPosTop(r rune) uint {
	if isDeva(r) {
		return PosAfterSUB
	}
	if isGuru(r) {
		return PosAfterPOST
	}
	if isGujr(r) {
		return PosAfterSUB
	}
	if isOrya(r) {
		return PosAfterMain
	}
	if isTaml(r) {
		return PosAfterSUB
	}
	if isTelu(r) {
		return PosBeforeSUB
	}
	if isKnda(r) {
		return PosBeforeSUB
	}
	if isSinh(r) {
		return PosAfterSUB
	}
	if isKhmr(r) {
		return PosAfterPOST
	}
	return PosAfterSUB

}

func matraPosTopBottom(r rune) uint {

	if isDeva(r) {
		return PosAfterSUB
	}
	if isBeng(r) {
		return PosAfterSUB
	}
	if isGuru(r) {
		return PosAfterPOST
	}
	if isGujr(r) {
		return PosAfterPOST
	}
	if isOrya(r) {
		return PosAfterSUB
	}
	if isTaml(r) {
		return PosAfterPOST
	}
	if isTelu(r) {
		return PosBeforeSUB
	}
	if isKnda(r) {
		return PosBeforeSUB
	}
	if isMlym(r) {
		return PosAfterPOST
	}
	if isSinh(r) {
		return PosAfterSUB
	}
	if isKhmr(r) {
		return PosAfterPOST
	}
	return PosAfterSUB
}
