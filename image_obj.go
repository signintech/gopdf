package gopdf

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//ImageObj image object
type ImageObj struct {
	buffer bytes.Buffer
	//imagepath string

	imgData []byte
}

func (i *ImageObj) init(funcGetRoot func() *GoPdf) {
	//me.getRoot = funcGetRoot
}

func (i *ImageObj) build() error {

	/*m, _, err := image.Decode(bytes.NewBuffer(i.imgData))
	if err != nil {
		return err
	}

	imageRect := m.Bounds()*/
	imgInfo, err := parseImg(i.imgData)
	if err != nil {
		return err
	}

	i.buffer.WriteString("<</Type /XObject\n")
	i.buffer.WriteString("/Subtype /Image\n")
	i.buffer.WriteString(fmt.Sprintf("/Width %d\n", imgInfo.w))  // /Width 675\n"
	i.buffer.WriteString(fmt.Sprintf("/Height %d\n", imgInfo.h)) //  /Height 942\n"
	if imgInfo.colspace == "Indexed" {
		//i.buffer.WriteString("/ColorSpace /DeviceRGB\n") //HARD CODE ไว้เป็น RGB
		return errors.New("not suport Indexed yet")
	} else {
		i.buffer.WriteString(fmt.Sprintf("/ColorSpace /%s\n", imgInfo.colspace))
		if imgInfo.colspace == "DeviceCMYK" {
			i.buffer.WriteString("/Decode [1 0 1 0 1 0 1 0]\n")
		}
	}
	i.buffer.WriteString(fmt.Sprintf("/BitsPerComponent %s\n", imgInfo.bitsPerComponent))
	if strings.TrimSpace(imgInfo.filter) != "" {
		i.buffer.WriteString(fmt.Sprintf("/Filter /%s\n", imgInfo.filter))
	}

	if strings.TrimSpace(imgInfo.decodeParms) != "" {
		i.buffer.WriteString(fmt.Sprintf("/DecodeParms <<%s>>\n", imgInfo.decodeParms))
	}

	if imgInfo.trns != nil && len(imgInfo.trns) > 0 {
		//TODO ต่อ
	}

	i.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(i.imgData))) // /Length 62303>>\n
	i.buffer.WriteString("stream\n")
	i.buffer.Write(i.imgData)
	i.buffer.WriteString("\nendstream\n")

	return nil
}

func (i *ImageObj) getType() string {
	return "Image"
}

func (i *ImageObj) getObjBuff() *bytes.Buffer {
	return &(i.buffer)
}

//SetImagePath set image path
func (i *ImageObj) SetImagePath(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = i.SetImage(file)
	if err != nil {
		return err
	}
	return nil
}

//SetImage set image
func (i *ImageObj) SetImage(r io.Reader) error {
	var err error
	i.imgData, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return nil
}

//GetRect get rect of img
func (i *ImageObj) GetRect() *Rect {

	m, _, err := image.Decode(bytes.NewBuffer(i.imgData))
	if err != nil {
		return nil
	}

	imageRect := m.Bounds()
	k := 1
	w := -128 //init
	h := -128 //init
	if w < 0 {
		w = -imageRect.Dx() * 72 / w / k
	}
	if h < 0 {
		h = -imageRect.Dy() * 72 / h / k
	}
	if w == 0 {
		w = h * imageRect.Dx() / imageRect.Dy()
	}
	if h == 0 {
		h = w * imageRect.Dy() / imageRect.Dx()
	}

	var rect = new(Rect)
	rect.H = float64(h)
	rect.W = float64(w)

	return rect
}

func (i *ImageObj) parse() error {
	return nil
}

//GetObjBuff get buffer
func (i *ImageObj) GetObjBuff() *bytes.Buffer {
	return i.getObjBuff()
}

//Build build buffer
func (i *ImageObj) Build() error {
	return i.build()
}
