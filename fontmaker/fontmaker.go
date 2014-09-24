package fontmaker

import (
	"errors"
	"github.com/signintech/gopdf"
	"os"
	"path/filepath"
	"strings"
)

type FontMaker struct {
}

func NewFontMaker() *FontMaker {
	return new(FontMaker)
}

func (me *FontMaker) MakeFont(path string, encoding string) (gopdf.IFont, error) {

	//read font file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	fileext := filepath.Ext(path)
	if strings.ToLower(fileext) != "ttf" {
		//now support only ttf :-P
		return nil, errors.New("support only ttf")
	}

	return nil, nil
}

func (me *FontMaker) CompressFont(path string) ([]byte, error) {
	return nil, nil
}
