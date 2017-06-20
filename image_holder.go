package gopdf

import (
	"crypto/md5"
	"fmt"
)

//ImageHolder hold image data
type ImageHolder interface {
	ID() string
	Bytes() []byte
}

//ImageHolderByBytes create ImageHolder by bytes
func ImageHolderByBytes(b []byte) (ImageHolder, error) {
	return newImageHolderByByte(b)
}

//imageHolderByByte read image from byte
type imageHolderByByte struct {
	id   string
	data []byte
}

func newImageHolderByByte(b []byte) (*imageHolderByByte, error) {
	h := md5.New()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	var imgb imageHolderByByte
	imgb.data = b
	imgb.id = hash
	return &imgb, nil
}

func (i *imageHolderByByte) ID() string {
	return i.id
}

func (i *imageHolderByByte) Bytes() []byte {
	return i.data
}
