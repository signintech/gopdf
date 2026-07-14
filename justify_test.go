package gopdf

import "testing"

func TestJustifyAdjustment(t *testing.T) {
	tests := []struct {
		name     string
		slack    float64
		gapCount int
		fontSize float64
		want     int
	}{
		{"zero gaps", 30, 0, 10, 0},
		{"zero slack", 0, 3, 10, 0},
		{"negative slack", -5, 3, 10, 0},
		{"zero font size", 30, 3, 0, 0},
		{"even distribution", 30, 3, 10, -1000},   // 10pt/gap -> -(10*1000/10)
		{"rounded distribution", 10, 3, 12, -278}, // 3.3333pt/gap -> -(277.78) rounded
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := justifyAdjustment(tt.slack, tt.gapCount, tt.fontSize)
			if got != tt.want {
				t.Fatalf("justifyAdjustment(%v,%d,%v) = %d, want %d",
					tt.slack, tt.gapCount, tt.fontSize, got, tt.want)
			}
		})
	}
}

func TestNonSpaceBounds(t *testing.T) {
	first, last := nonSpaceBounds("  hi  ")
	if first != 2 || last != 3 {
		t.Fatalf("nonSpaceBounds = (%d,%d), want (2,3)", first, last)
	}
	if f, l := nonSpaceBounds("   "); f != -1 || l != -1 {
		t.Fatalf("all-spaces bounds = (%d,%d), want (-1,-1)", f, l)
	}
	if f, l := nonSpaceBounds(""); f != -1 || l != -1 {
		t.Fatalf("empty bounds = (%d,%d), want (-1,-1)", f, l)
	}
}

func TestInteriorSpaceCount(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"hello world", 1},
		{"a b c", 2},
		{"  hello   world  ", 3}, // 3 interior spaces; leading/trailing excluded
		{"word", 0},
		{"", 0},
		{"   ", 0},
	}
	for _, tt := range tests {
		if got := interiorSpaceCount(tt.text); got != tt.want {
			t.Fatalf("interiorSpaceCount(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}
