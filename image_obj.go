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
	imginfo imgInfo
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

	i.buffer.WriteString("<</Type /XObject\n")
	i.buffer.WriteString("/Subtype /Image\n")
	i.buffer.WriteString(fmt.Sprintf("/Width %d\n", i.imginfo.w))  // /Width 675\n"
	i.buffer.WriteString(fmt.Sprintf("/Height %d\n", i.imginfo.h)) //  /Height 942\n"
	if i.imginfo.colspace == "Indexed" {
		//i.buffer.WriteString("/ColorSpace /DeviceRGB\n") //HARD CODE ไว้เป็น RGB
		//TODO fix this
		return errors.New("not suport Indexed yet")
	} else {
		i.buffer.WriteString(fmt.Sprintf("/ColorSpace /%s\n", i.imginfo.colspace))
		if i.imginfo.colspace == "DeviceCMYK" {
			i.buffer.WriteString("/Decode [1 0 1 0 1 0 1 0]\n")
		}
	}
	i.buffer.WriteString(fmt.Sprintf("/BitsPerComponent %s\n", i.imginfo.bitsPerComponent))
	if strings.TrimSpace(i.imginfo.filter) != "" {
		i.buffer.WriteString(fmt.Sprintf("/Filter /%s\n", i.imginfo.filter))
	}

	if strings.TrimSpace(i.imginfo.decodeParms) != "" {
		i.buffer.WriteString(fmt.Sprintf("/DecodeParms <<%s>>\n", i.imginfo.decodeParms))
	}

	if i.imginfo.trns != nil && len(i.imginfo.trns) > 0 {
		j := 0
		max := len(i.imginfo.trns)
		var trns bytes.Buffer
		for j < max {
			trns.WriteByte(i.imginfo.trns[j])
			trns.WriteString(" ")
			trns.WriteByte(i.imginfo.trns[j])
			trns.WriteString(" ")
			j++
		}
		i.buffer.WriteString(fmt.Sprintf("/Mask [%s]\n", trns.String()))
	}

	if i.haveSMask() {
		//TODO fix this
		return errors.New("not suport smask yet")
	}

	i.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(i.imgData))) // /Length 62303>>\n
	i.buffer.WriteString("stream\n")
	i.buffer.Write(i.imgData)
	i.buffer.WriteString("\nendstream\n")

	return nil
}

func (i *ImageObj) haveSMask() bool {
	if i.imginfo.smask != nil && len(i.imginfo.smask) > 0 {
		return true
	}
	return false
}

func (i ImageObj) createSMask() (*SMask, error) {
	var smk SMask
	smk.w = i.imginfo.w
	smk.h = i.imginfo.h
	smk.colspace = "DeviceGray"
	smk.bitsPerComponent = "8"
	smk.filter = i.imginfo.filter
	smk.data = i.imginfo.smask
	return &smk, nil
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

	imginfo, err := parseImg(i.imgData)
	if err != nil {
		return err
	}
	i.imginfo = imginfo

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
