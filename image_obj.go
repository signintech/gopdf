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

	imgData []byte
}

func (i *ImageObj) init(funcGetRoot func() *GoPdf) {
	//me.getRoot = funcGetRoot
}

func (i *ImageObj) build() error {

	m, _, err := image.Decode(bytes.NewBuffer(i.imgData))
	if err != nil {
		return err
	}

	imageRect := m.Bounds()

	i.buffer.WriteString("<</Type /XObject\n")
	i.buffer.WriteString("/Subtype /Image\n")
	i.buffer.WriteString(fmt.Sprintf("/Width %d\n", imageRect.Dx()))  // /Width 675\n"
	i.buffer.WriteString(fmt.Sprintf("/Height %d\n", imageRect.Dy())) //  /Height 942\n"
	i.buffer.WriteString("/ColorSpace /DeviceRGB\n")                  //HARD CODE ไว้เป็น RGB
	i.buffer.WriteString("/BitsPerComponent 8\n")                     //HARD CODE ไว้เป็น 8 bit
	i.buffer.WriteString("/Filter /DCTDecode\n")
	//me.buffer.WriteString("/Filter /FlateDecode\n")
	//me.buffer.WriteString("/DecodeParms <</Predictor 15 /Colors 3 /BitsPerComponent 8 /Columns 675>>\n")
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
