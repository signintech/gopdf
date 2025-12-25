package gopdf

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func getHarf(char rune) Harf {
	for _, s := range arabic_alphabet {
		if s.equals(char) {
			return s
		}
	}
	return Harf{Unicode: char, Isolated: char, Midlle: char, Final: char}
}

// equals() return if true if the given Arabic char is alphabetically equal to
// the current Harf regardless its shape
func (c *Harf) equals(char rune) bool {
	switch char {
	case c.Unicode:
		return true
	case c.Beggining:
		return true
	case c.Isolated:
		return true
	case c.Midlle:
		return true
	case c.Final:
		return true
	default:
		return false
	}
}

// This checks the Harf place as it decide how will be the Harf Shape
func getCharShape(previousChar, currentChar, nextChar rune) rune {
	shape := currentChar
	nextArabic := false
	previousArabic := false

	for _, s := range arabic_alphabet {
		if s.equals(previousChar) { // prvious char is Arabic Harf
			previousArabic = true
		}
		if s.equals(nextChar) {
			nextArabic = true
		}
	}

	for _, s := range arabic_alphabet {
		if !s.equals(currentChar) {
			continue
		}

		if previousArabic && nextArabic { // in middle
			for s, _ := range rightJoiningOnlyLetters { // if its an Harf which is must be separated
				if s.equals(previousChar) {
					return getHarf(currentChar).Beggining // return the default shape
				}
			}
			return getHarf(currentChar).Midlle
		}

		if nextArabic { // first letter and in start of a word
			return getHarf(currentChar).Beggining
		}

		if previousArabic { // final letter as it's in the end of a word
			for s, _ := range rightJoiningOnlyLetters {
				if s.equals(previousChar) {
					return getHarf(currentChar).Isolated
				}
			}
			return getHarf(currentChar).Final
	}

}
}