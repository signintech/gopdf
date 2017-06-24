package gopdf

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
)

//ImageHolder hold image data
type ImageHolder interface {
	ID() string
	io.Reader
}

//ImageHolderByBytes create ImageHolder by []byte
func ImageHolderByBytes(b []byte) (ImageHolder, error) {
	return newImageBuff(b)
}

//ImageHolderByPath create ImageHolder by image path
func ImageHolderByPath(path string) (ImageHolder, error) {
	return newImageBuffByPath(path)
}

//ImageHolderByReader create ImageHolder by io.Reader
func ImageHolderByReader(r io.Reader) (ImageHolder, error) {
	return newImageBuffByReader(r)
}

//imageBuff image holder (impl ImageHolder)
type imageBuff struct {
	id string
	bytes.Buffer
}

func newImageBuff(b []byte) (*imageBuff, error) {
	h := md5.New()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	var i imageBuff
	i.id = fmt.Sprintf("%x", h.Sum(nil))
	i.Write(b)
	return &i, nil
}

func newImageBuffByPath(path string) (*imageBuff, error) {
	var i imageBuff
	i.id = path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i.Write(b)
	return &i, nil
}

func newImageBuffByReader(r io.Reader) (*imageBuff, error) {

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	_, err = h.Write(b)
	if err != nil {
		return nil, err
	}
	var i imageBuff
	i.id = fmt.Sprintf("%x", h.Sum(nil))
	i.Write(b)
	return &i, nil
}

func (i *imageBuff) ID() string {
	return i.id
}
