package gopdf

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
)

//ImageObj image object
type ImageObj struct {
	buffer bytes.Buffer
	//imagepath string

	raw           []byte
	imginfo       imgInfo
	pdfProtection *PDFProtection
	//getRoot func() *GoPdf
}

func (i *ImageObj) init(funcGetRoot func() *GoPdf) {

}

func (i *ImageObj) setProtection(p *PDFProtection) {
	i.pdfProtection = p
}

func (i *ImageObj) protection() *PDFProtection {
	return i.pdfProtection
}

func (i *ImageObj) build(objID int) error {

	buff, err := buildImgProp(i.imginfo)
	if err != nil {
		return err
	}
	_, err = buff.WriteTo(&i.buffer)
	if err != nil {
		return err
	}

	i.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(i.imginfo.data))) // /Length 62303>>\n
	i.buffer.WriteString("stream\n")
	if i.protection() != nil {
		tmp, err := rc4Cip(i.protection().objectkey(objID), i.imginfo.data)
		if err != nil {
			return err
		}
		i.buffer.Write(tmp)
		i.buffer.WriteString("\n")
	} else {
		i.buffer.Write(i.imginfo.data)
	}
	i.buffer.WriteString("\nendstream\n")

	return nil
}

func (i *ImageObj) isColspaceIndexed() bool {
	return isColspaceIndexed(i.imginfo)
}

func (i *ImageObj) haveSMask() bool {
	return haveSMask(i.imginfo)
}

func (i *ImageObj) createSMask() (*SMask, error) {
	var smk SMask
	smk.setProtection(i.protection())
	smk.w = i.imginfo.w
	smk.h = i.imginfo.h
	smk.colspace = "DeviceGray"
	smk.bitsPerComponent = "8"
	smk.filter = i.imginfo.filter
	smk.data = i.imginfo.smask
	smk.decodeParms = fmt.Sprintf("/Predictor 15 /Colors 1 /BitsPerComponent 8 /Columns %d", i.imginfo.w)
	return &smk, nil
}

func (i *ImageObj) createDeviceRGB() (*DeviceRGBObj, error) {
	var dRGB DeviceRGBObj
	dRGB.data = i.imginfo.pal
	return &dRGB, nil
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
	i.raw, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	return nil
}

//GetRect get rect of img
func (i *ImageObj) GetRect() *Rect {

	m, _, err := image.Decode(bytes.NewBuffer(i.raw))
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

	imginfo, err := parseImg(i.raw)
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

//Parse parse img
func (i *ImageObj) Parse() error {
	return i.parse()
}

//Build build buffer
func (i *ImageObj) Build(objID int) error {
	return i.build(objID)
}
