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
	return Harf{Unicode: char, Isolated: char, Middle: char, Final: char}
}

// equals() return if true if the given Arabic char is alphabetically equal to
// the current Harf regardless its shape
func (c *Harf) equals(char rune) bool {
	switch char {
	case c.Unicode:
		return true
	case c.Beginning:
		return true
	case c.Isolated:
		return true
	case c.Middle:
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

	if _, ok := arabic_alphabet[previousChar]; ok {
		previousArabic = true
	}

	if _, ok := arabic_alphabet[nextChar]; ok {
		nextArabic = true
	}

	if _, ok := arabic_alphabet[currentChar]; !ok {
		return shape
	}

	if previousArabic && nextArabic { // in middle
		for s := range rightJoiningOnlyLetters { // if its an Harf which is must be separated
			if s.equals(previousChar) {
				return getHarf(currentChar).Beginning // return the default shape
			}
		}
		return getHarf(currentChar).Middle
	}

	if nextArabic { // first letter and in start of a word
		return getHarf(currentChar).Beginning
	}

	if previousArabic { // final letter as it's in the end of a word
		for s := range rightJoiningOnlyLetters {
			if s.equals(previousChar) {
				return getHarf(currentChar).Isolated
			}
		}
		return getHarf(currentChar).Final
	}

	if !previousArabic && !nextArabic { // single isolated letter
		return getHarf(currentChar).Isolated
	}

	return shape
}

func ToArabic(text string) string {
	var nextHarf, previousHarf rune

	hrof := []rune(text)    // hrof is arabic letters
	hrofLength := len(hrof) // hrof length is the number of arabic letters

	arabicSentence := make([]rune, 0, hrofLength)
	for i := 0; i < hrofLength; i++ {
		currentHarf := hrof[i]

		if i == 0 {
			previousHarf = 0
		} else {
			previousHarf = hrof[i-1]
		}

		if i == hrofLength-1 {
			nextHarf = 0
		} else {
			nextHarf = hrof[i+1]
		}

		// Lam-Alef Ligature Check
		if currentHarf == LAM.Unicode && nextHarf != 0 {
			var ligatureHarf rune
			foundLigature := false
			switch nextHarf {
			case ALEF.Unicode:
				ligatureHarf = LAM_ALEF.Unicode
				foundLigature = true
			case ALEF_HAMZA_ABOVE.Unicode:
				ligatureHarf = LAM_ALEF_HAMZA_ABOVE.Unicode
				foundLigature = true
			}
			if foundLigature {
				currentHarf = ligatureHarf
				i++
				// We need to update nextHarf to the one *after* the Alef for correct shaping of the ligature itself
				if i == hrofLength-1 {
					nextHarf = 0
				} else {
					nextHarf = hrof[i+1]
				}
			}
		}

		harfShape := getCharShape(previousHarf, currentHarf, nextHarf)
		arabicSentence = append(arabicSentence, harfShape)
	}
	arabicSentenceRTL := Reverse(string(arabicSentence))
	return arabicSentenceRTL
}
