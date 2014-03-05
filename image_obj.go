package gopdf

import (
	"bytes"
	"os"
	"image"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
)

type ImageObj struct{
	buffer  bytes.Buffer
	imagepath string
}

func (me *ImageObj) Init(funcGetRoot func()(*GoPdf)) {
	//me.getRoot = funcGetRoot
}

func (me *ImageObj) Build() {
	
	
	file, err := os.Open(me.imagepath)
	if err != nil {
		//fmt.Printf("0--%+v\n",err)
		return
	}
	
	m , _ ,  err := image.Decode(file)
	if err != nil {
		return 
	} 
	
	imageRect := m.Bounds()

	b, _  := ioutil.ReadFile(me.imagepath)
	
	me.buffer.WriteString("<</Type /XObject\n")
	me.buffer.WriteString("/Subtype /Image\n")
	me.buffer.WriteString(fmt.Sprintf("/Width %d\n",imageRect.Dx())) // /Width 675\n"
	me.buffer.WriteString(fmt.Sprintf("/Height %d\n",imageRect.Dy())) //  /Height 942\n"
	me.buffer.WriteString("/ColorSpace /DeviceRGB\n") //HARD CODE ไว้เป็น RGB
	me.buffer.WriteString("/BitsPerComponent 8\n") //HARD CODE ไว้เป็น 8 bit
	me.buffer.WriteString("/Filter /DCTDecode\n")
	//me.buffer.WriteString("/Filter /FlateDecode\n")
	//me.buffer.WriteString("/DecodeParms <</Predictor 15 /Colors 3 /BitsPerComponent 8 /Columns 675>>\n")
	me.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n",len(b)))// /Length 62303>>\n
	me.buffer.WriteString("stream\n")
	me.buffer.Write(b)
	me.buffer.WriteString("\nendstream\n")
}

func (me *ImageObj) GetType() string {
	return "Image"
}

func (me *ImageObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *ImageObj) SetImagePath(path string){
	me.imagepath = path
}

func (me *ImageObj) GetRect() *Rect {
	file, err := os.Open(me.imagepath)
	m , _ ,  err := image.Decode(file)
	if err != nil {
		return nil
	} 
	
	imageRect := m.Bounds()
	k := 1
	w := -128 //init
	h := -128 //init
	if(w < 0){
		w = -imageRect.Dx()*72/w/k
	}
	if(h < 0){
		h = -imageRect.Dy()*72/h/k
	}
	if(w == 0){
		w = h*imageRect.Dx()/imageRect.Dy()
	}
	if(h == 0){
		h = w*imageRect.Dy()/imageRect.Dx()
	}
	
	var rect = new(Rect)
	rect.H = float64(h)
	rect.W = float64(w)
	
	return rect
}

