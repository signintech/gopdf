package gopdf

import (
	"testing"
)

func TestToArabic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:  "Normal connected word (Salam)",
			input: "سلام", // Seen, Lam, Alef, Meem
			// Seen: Beg (since next is Lam) -> SEEN.Beginning
			// Lam+Alef: Ligature. Prev is Sen (connected). Next is Meem (disconnected from Alef/Ligature usually? No, LamAlef connects to right, but not left).
			// So LamAlef connects to Seen. Form: Final (since connected to right, not left). -> LAM_ALEF.Final
			// Meem: Isolated (since prev LamAlef doesn't connect left). -> MEEM.Isolated
			// Reverse order: Meem, LamAlef, Seen.
			expected: string([]rune{MEEM.Isolated, LAM_ALEF.Final, SEEN.Beginning}),
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
