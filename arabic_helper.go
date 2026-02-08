package gopdf

import "strings"

// ALLAH_LIGATURE is the Unicode character for the Allah ligature (U+FDF2 ﷲ)
const ALLAH_LIGATURE rune = 0xFDF2

// convertAllahToLigature replaces the word "الله" (Allah) with the Allah ligature U+FDF2 (ﷲ)
func convertAllahToLigature(text string) string {
	// الله without tashkeel: Alef + Lam + Lam + Heh
	allah := string([]rune{ALEF.Unicode, LAM.Unicode, LAM.Unicode, HEH.Unicode})
	// Replace with the Allah ligature character
	return strings.ReplaceAll(text, allah, string(ALLAH_LIGATURE))
}

// reverseWithTashkeel reverses Arabic text while keeping tashkeel attached to base characters
func reverseWithTashkeel(runes []rune) string {
	if len(runes) == 0 {
		return ""
	}

	// Group base characters with their following tashkeel
	type hrofGroup struct {
		base     rune
		tashkeel []rune
	}

	var groups []hrofGroup
	var currentGroup *hrofGroup

	for _, r := range runes {
		if IsTashkeel(r) {
			if currentGroup != nil {
				currentGroup.tashkeel = append(currentGroup.tashkeel, r)
			}
		} else {
			groups = append(groups, hrofGroup{base: r})
			currentGroup = &groups[len(groups)-1]
		}
	}

	// Reverse the groups and rebuild
	// Output tashkeel BEFORE base for proper RTL rendering in PDF
	result := make([]rune, 0, len(runes))
	for i := len(groups) - 1; i >= 0; i-- {
		result = append(result, groups[i].tashkeel...)
		result = append(result, groups[i].base)
	}
	return string(result)

}
func getHarf(char rune) Harf {
	for _, s := range arabicAlphabet {
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

	if _, ok := arabicAlphabet[previousChar]; ok {
		previousArabic = true
	}

	if _, ok := arabicAlphabet[nextChar]; ok {
		nextArabic = true
	}

	if _, ok := arabicAlphabet[currentChar]; !ok {
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

// findPreviousNonTashkeel finds the previous character that is not a tashkeel mark
func findPreviousNonTashkeelHarf(runes []rune, currentIndex int) rune {
	for i := currentIndex - 1; i >= 0; i-- {
		if !IsTashkeel(runes[i]) {
			return runes[i]
		}
	}
	return 0
}

// findNextNonTashkeel finds the next character that is not a tashkeel mark
func findNextNonTashkeelHarf(runes []rune, currentIndex int) rune {
	for i := currentIndex + 1; i < len(runes); i++ {
		if !IsTashkeel(runes[i]) {
			return runes[i]
		}
	}
	return 0
}

// IsTashkeel returns true if the rune is an Arabic diacritical mark
func IsTashkeel(r rune) bool {
	return tashkeelMarks[r]
}

func ToArabic(text string) string {
	// Preprocess: convert "الله" to the Allah ligature U+FDF2 (ﷲ)
	text = convertAllahToLigature(text)

	hrof := []rune(text)    // hrof is arabic letters
	hrofLength := len(hrof) // hrof length is the number of arabic letters

	arabicSentence := make([]rune, 0, hrofLength)
	for i := 0; i < hrofLength; i++ {
		currentHarf := hrof[i]

		// If current char is tashkeel
		if IsTashkeel(currentHarf) {
			// Check if vowel followed by SHADDA - output combined ligature
			if i+1 < hrofLength && hrof[i+1] == SHADDA && currentHarf != SHADDA {
				if ligature := GetShaddaLigature(currentHarf); ligature != 0 {
					arabicSentence = append(arabicSentence, ligature)
					i++ // skip the shadda we already added
					continue
				}
			}
			arabicSentence = append(arabicSentence, currentHarf)
			continue
		}
		// Find previous non-tashkeel character
		previousHarf := findPreviousNonTashkeelHarf(hrof, i)

		// Find next non-tashkeel character
		nextHarf := findNextNonTashkeelHarf(hrof, i)

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
				// Collect tashkeel between Lam and Alef
				var tashkeelBetween []rune
				for i++; i < hrofLength && hrof[i] != nextHarf; i++ {
					if IsTashkeel(hrof[i]) {
						tashkeelBetween = append(tashkeelBetween, hrof[i])
					}
				}
				nextHarf = findNextNonTashkeelHarf(hrof, i)

				// Append ligature shape first, then tashkeel (so tashkeel attaches to ligature after reversal)
				harfShape := getCharShape(previousHarf, currentHarf, nextHarf)
				arabicSentence = append(arabicSentence, harfShape)
				arabicSentence = append(arabicSentence, tashkeelBetween...)
				continue
			}
		}

		harfShape := getCharShape(previousHarf, currentHarf, nextHarf)
		arabicSentence = append(arabicSentence, harfShape)
	}
	arabicSentenceRTL := reverseWithTashkeel(arabicSentence)
	return arabicSentenceRTL
}
