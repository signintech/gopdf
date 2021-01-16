package indic

var raChars = []uint{
	0x0930, /* Devanagari */
	0x09B0, /* Bengali */
	0x09F0, /* Bengali (Assamese) */
	0x0A30, /* Gurmukhi */ /* No Reph */
	0x0AB0, /* Gujarati */
	0x0B30, /* Oriya */
	0x0BB0, /* Tamil */  /* No Reph */
	0x0C30, /* Telugu */ /* Reph formed only with ZWJ */
	0x0CB0, /* Kannada */
	0x0D30, /* Malayalam */ /* No Reph, Logical Repha */
	0x0DBB, /* Sinhala */   /* Reph formed only with ZWJ */
	0x179A, /* Khmer */     /* No Reph, Visual Repha */
}

func IsRuneRA(r rune) bool {
	return isRuneRA(r)
}

func isRuneRA(r rune) bool {
	for _, ra := range raChars {
		if ra == uint(r) {
			return true
		}
	}
	return false
}
