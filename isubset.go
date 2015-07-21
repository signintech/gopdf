package gopdf

import "errors"

var ErrCharNotFound = errors.New("char not found")

type ISubset interface {
	AddChars(txt string)
	CharIndex(r rune) (uint64, error)
}
