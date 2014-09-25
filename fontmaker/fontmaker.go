package fontmaker

import (
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/signintech/gopdf"
	"io/ioutil"
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
	if strings.ToLower(fileext) != ".ttf" {
		//now support only ttf :-P
		return nil, errors.New("support only ttf ")
	}

	var parser TTFParser
	err := parser.Parse(path)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (me *FontMaker) CompressFont(path string) (*bytes.Buffer, error) {
	rawbytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	gw := gzip.NewWriter(&buff)
	_, err = gw.Write(rawbytes)
	if err != nil {
		return nil, err
	}
	gw.Close()
	return &buff, nil
}
