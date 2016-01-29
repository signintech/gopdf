package gopdf

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
)

type ImageObj struct {
	buffer    bytes.Buffer
	imagepath string
}

func (i *ImageObj) init(funcGetRoot func() *GoPdf) {
	//me.getRoot = funcGetRoot
}

func (i *ImageObj) build() error {

	file, err := os.Open(i.imagepath)
	if err != nil {
		//fmt.Printf("0--%+v\n",err)
		return err
	}
	defer file.Close()

	m, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	imageRect := m.Bounds()

	b, _ := ioutil.ReadFile(i.imagepath)

	i.buffer.WriteString("<</Type /XObject\n")
	i.buffer.WriteString("/Subtype /Image\n")
	i.buffer.WriteString(fmt.Sprintf("/Width %d\n", imageRect.Dx()))  // /Width 675\n"
	i.buffer.WriteString(fmt.Sprintf("/Height %d\n", imageRect.Dy())) //  /Height 942\n"
	i.buffer.WriteString("/ColorSpace /DeviceRGB\n")                  //HARD CODE ไว้เป็น RGB
	i.buffer.WriteString("/BitsPerComponent 8\n")                     //HARD CODE ไว้เป็น 8 bit
	i.buffer.WriteString("/Filter /DCTDecode\n")
	//me.buffer.WriteString("/Filter /FlateDecode\n")
	//me.buffer.WriteString("/DecodeParms <</Predictor 15 /Colors 3 /BitsPerComponent 8 /Columns 675>>\n")
	i.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(b))) // /Length 62303>>\n
	i.buffer.WriteString("stream\n")
	i.buffer.Write(b)
	i.buffer.WriteString("\nendstream\n")

	return nil
}

func (i *ImageObj) getType() string {
	return "Image"
}

func (i *ImageObj) getObjBuff() *bytes.Buffer {
	return &(i.buffer)
}

func (i *ImageObj) SetImagePath(path string) {
	i.imagepath = path
}

func (i *ImageObj) GetRect() *Rect {
	file, err := os.Open(i.imagepath)
	m, _, err := image.Decode(file)
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
