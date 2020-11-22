package ucdn

import (
	"unicode"
)

func Script(r rune) []string {
	var names []string
	for name, table := range unicode.Scripts {
		if unicode.Is(table, r) {
			names = append(names, name)
		}
	}
	return names
}

func Category(r rune) []string {
	var names []string
	for name, table := range unicode.Categories {
		if unicode.Is(table, r) {
			names = append(names, name)
		}
	}
	return names
}
