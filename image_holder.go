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

//ImageHolderByBytes create ImageHolder by bytes
func ImageHolderByBytes(b []byte) (ImageHolder, error) {
	return newImageBuff(b)
}

//ImageHolderByPath create ImageHolder by bytes
func ImageHolderByPath(path string) (ImageHolder, error) {
	return newImageFile(path)
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

func (i *imageBuff) ID() string {
	return i.id
}

//imageFile image holder
type imageFile struct {
	id string
	bytes.Buffer
}

func newImageFile(path string) (*imageFile, error) {
	var i imageFile
	i.id = path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	i.Write(b)
	return &i, nil
}

func (i *imageFile) ID() string {
	return i.id
}
