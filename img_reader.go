package gopdf

import (
	"crypto/md5"
	"fmt"
)

type ImgReader interface {
	UniqueKey() (string, error)
	Bytes() []byte
}

//ImgBytes read image from byte
type ImgBytes struct {
	Data []byte
}

func (i *ImgBytes) UniqueKey() (string, error) {
	h := md5.New()
	_, err := h.Write(i.Data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (i *ImgBytes) Bytes() []byte {
	return i.Data
}
