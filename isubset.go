package gopdf

import (
	"errors"

	"github.com/signintech/gopdf/fontmaker/core"
)

//ErrCharNotFound char not found
var ErrCharNotFound = errors.New("char not found")

//ISubset ttf font
type ISubset interface {
	AddChars(txt string) error
	CharIndex(r rune) (uint64, error) //get char index
	CharWidth(r rune) (uint64, error) //find chear width
	GetUt() int64
	KernValueByLeft(r uint64) (bool, *core.KernValue)
}
