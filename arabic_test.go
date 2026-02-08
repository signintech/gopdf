package gopdf

import (
	"testing"
)

func TestGetShaddaLigature(t *testing.T) {
	tests := []struct {
		name     string
		vowel    rune
		expected rune
	}{
		{"Fatha + Shadda", FATHA, SHADDA_FATHA},
		{"Damma + Shadda", DAMMA, SHADDA_DAMMA},
		{"Kasra + Shadda", KASRA, SHADDA_KASRA},
		{"Tanween Damm + Shadda", TANWEEN_DAMM, SHADDA_DAMMATAN},
		{"Tanween Kasr + Shadda", TANWEEN_KASR, SHADDA_KASRATAN},
		{"Superscript Alef + Shadda", SUPERSCRIPT_ALEF, SHADDA_SUPERSCRIPT_ALEF},
		{"Tanween Fath (no ligature)", TANWEEN_FATH, 0},
		{"Sukun (no ligature)", SUKUN, 0},
		{"Non-tashkeel (no ligature)", 'a', 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetShaddaLigature(tt.vowel)
			if got != tt.expected {
				t.Errorf("GetShaddaLigature(%U) = %U, want %U", tt.vowel, got, tt.expected)
			}
		})
	}
}

func TestIsTashkeel(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"Fatha is tashkeel", FATHA, true},
		{"Damma is tashkeel", DAMMA, true},
		{"Kasra is tashkeel", KASRA, true},
		{"Shadda is tashkeel", SHADDA, true},
		{"Sukun is tashkeel", SUKUN, true},
		{"Tanween Fath is tashkeel", TANWEEN_FATH, true},
		{"Tanween Damm is tashkeel", TANWEEN_DAMM, true},
		{"Tanween Kasr is tashkeel", TANWEEN_KASR, true},
		{"Superscript Alef is tashkeel", SUPERSCRIPT_ALEF, true},
		{"Beh is not tashkeel", BEH.Unicode, false},
		{"Alef is not tashkeel", ALEF.Unicode, false},
		{"Space is not tashkeel", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTashkeel(tt.input)
			if got != tt.expected {
				t.Errorf("IsTashkeel(%U) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToArabicWithTashkeel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Word with fatha - letters should join",
			input: "بَاب", // Beh + Fatha + Alef + Beh
			// Tashkeel now comes BEFORE its base letter in output
			// Reversed: Beh(iso), Alef(final), Fatha+Beh(beg)
			expected: string([]rune{BEH.Isolated, ALEF.Final, FATHA, BEH.Beginning}),
		},
		{
			name:  "Word with shadda",
			input: "شدّة", // Sheen + Dal + Shadda + Teh_Marbuta
			// Tashkeel now comes BEFORE its base letter in output
			// Reversed: Teh_Marbuta(iso), Shadda+Dal(mid), Sheen(beg)
			expected: string([]rune{TEH_MARBUTA.Isolated, SHADDA, DAL.Middle, SHEEN.Beginning}),
		},
		{
			name:  "Lam-Alef ligature with fatha between",
			input: "لَا", // Lam + Fatha + Alef
			// Tashkeel now comes BEFORE its base letter
			// Reversed: Fatha+LAM_ALEF(isolated)
			expected: string([]rune{FATHA, LAM_ALEF.Isolated}),
		},
		{
			name:  "Bismillah start",
			input: "بِسْمِ", // Beh + Kasra + Seen + Sukun + Meem + Kasra
			// Tashkeel now comes BEFORE its base letter
			// Reversed: Kasra+Meem(final), Sukun+Seen(mid), Kasra+Beh(beg)
			expected: string([]rune{KASRA, MEEM.Final, SUKUN, SEEN.Middle, KASRA, BEH.Beginning}),
		},
		{
			name:  "Multiple tashkeel on one letter",
			input: "بًّ", // Beh + Tanween_Fath + Shadda
			// No ligature for Tanween_Fath+Shadda, output separately
			expected: string([]rune{TANWEEN_FATH, SHADDA, BEH.Isolated}),
		},
		{
			name:  "Shadda + Fatha ligature",
			input: "رَّ", // Reh + Fatha + Shadda
			// Output: Combined Shadda+Fatha ligature
			expected: string([]rune{SHADDA_FATHA, REH.Isolated}),
		},
		{
			name:  "Shadda + Kasra ligature",
			input: "بِّ", // Beh + Kasra + Shadda
			// Output: Combined Shadda+Kasra ligature
			expected: string([]rune{SHADDA_KASRA, BEH.Isolated}),
		},
		{
			name:  "Allah ligature",
			input: "الله", // Alef + Lam + Lam + Heh (without tashkeel)
			// Now converted to Allah ligature U+FDF2
			expected: string([]rune{ALLAH_LIGATURE}),
		},
		{
			name:  "Allah in sentence - بسم الله",
			input: "بسم الله", // Beh + Seen + Meem + Space + Alef + Lam + Lam + Heh
			// Allah becomes ligature U+FDF2
			expected: string([]rune{
				ALLAH_LIGATURE,
				' ',
				MEEM.Final, SEEN.Middle, BEH.Beginning,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToArabic(tt.input)
			if got != tt.expected {
				t.Errorf("ToArabic(%q) = %x, want %x", tt.input, []rune(got), []rune(tt.expected))
			}
			t.Log(got)
		})
	}
}

func TestToArabic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Normal connected word (Salam)",
			input: "سلامً", // Seen, Lam, Alef, Meem, Tanween_Fath
			// Seen: Beg (since next is Lam) -> SEEN.Beginning
			// Lam+Alef: Ligature. Prev is Seen (connected). -> LAM_ALEF.Final
			// Meem: Isolated (since prev LamAlef doesn't connect left). -> MEEM.Isolated
			// Tanween_Fath: now comes BEFORE Meem in output
			// Reverse order: Tanween+Meem, LamAlef, Seen.
			expected: string([]rune{TANWEEN_FATH, MEEM.Isolated, LAM_ALEF.Final, SEEN.Beginning}),
		},
		{
			name:  "Isolated Ligature (La)",
			input: "لا", // Lam, Alef
			// Ligature Isolated.
			expected: string([]rune{LAM_ALEF.Isolated}),
		},
		{
			name:  "Word without ligature (Kitab)",
			input: "كتاب", // Kaf, Teh, Alef, Beh
			// Kaf: Beg -> KAF.Beginning
			// Teh: Mid -> TEH.Middle
			// Alef: Final (connects to Teh) -> ALEF.Final
			// Beh: Isolated (Alef doesn't connect left) -> BEH.Isolated
			// Reverse: Beh, Alef, Teh, Kaf
			expected: string([]rune{BEH.Isolated, ALEF.Final, TEH.Middle, KAF.Beginning}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToArabic(tt.input)
			if got != tt.expected {
				t.Errorf("ToArabic(%q) = %x, want %x", tt.input, []rune(got), []rune(tt.expected))
			}
		})
	}
}
