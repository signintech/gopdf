package gopdf

import (
	"bytes"
	"fmt"
	"image"

	// Packages image/jpeg and image/png are not used explicitly in the code below,
	// but are imported for their initialization side-effect, which allows
	// image.Decode to understand JPEG formatted images.
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
)

// ImageObj image object
type ImageObj struct {
	//imagepath string
	IsMask        bool
	SplittedMask  bool
	rawImgReader  *bytes.Reader
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

func (i *ImageObj) write(w io.Writer, objID int) error {
	data := i.imginfo.data

	if i.IsMask {
		data = i.imginfo.smask
		if err := writeMaskImgProps(w, i.imginfo); err != nil {
			return err
		}
	} else {
		if err := writeImgProps(w, i.imginfo, i.SplittedMask); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "\t/Length %d\n>>\n", len(data)); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "stream\n"); err != nil {
		return err
	}

	if i.protection() != nil {
		tmp, err := rc4Cip(i.protection().objectkey(objID), data)
		if err != nil {
			return err
		}

		if _, err := w.Write(tmp); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
	} else {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "\nendstream\n"); err != nil {
		return err
	}

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

// SetImagePath set image path
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

// SetImage set image
func (i *ImageObj) SetImage(r io.Reader) error {

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	i.rawImgReader = bytes.NewReader(data)

	return nil
}

// GetRect get rect of img
func (i *ImageObj) GetRect() *Rect {

	rect, err := i.getRect()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return rect
}

// GetRect get rect of img
func (i *ImageObj) getRect() (*Rect, error) {

	i.rawImgReader.Seek(0, 0)
	m, _, err := image.Decode(i.rawImgReader)
	if err != nil {
		return nil, err
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

	return rect, nil
}

func (i *ImageObj) parse() error {

	i.rawImgReader.Seek(0, 0)
	imginfo, err := parseImg(i.rawImgReader)
	if err != nil {
		return err
	}
	i.imginfo = imginfo

	return nil
}

// Parse parse img
func (i *ImageObj) Parse() error {
	return i.parse()
}
