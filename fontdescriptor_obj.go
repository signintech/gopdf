package gopdf

import (
	"bytes"
)

type FontDescriptorObj struct {
	buffer            bytes.Buffer
	font              IFont
	fontFileObjRelate string
}

func (f *FontDescriptorObj) init(funcGetRoot func() *GoPdf) {

}

func (f *FontDescriptorObj) build(objID int) error {

	f.buffer.WriteString("<</Type /FontDescriptor /FontName /" + f.font.GetName() + " ")
	descs := f.font.GetDesc()
	i := 0
	max := len(descs)
	for i < max {
		f.buffer.WriteString("/" + descs[i].Key + " " + descs[i].Val + " ")
		i++
	}

	if f.getType() == "Type1" {
		f.buffer.WriteString("/FontFile ")
	} else {
		f.buffer.WriteString("/FontFile2 ")
	}

	f.buffer.WriteString(f.fontFileObjRelate)
	f.buffer.WriteString(">>\n")

	return nil
}

func (f *FontDescriptorObj) getType() string {
	return "FontDescriptor"
}

func (f *FontDescriptorObj) getObjBuff() *bytes.Buffer {
	return &(f.buffer)
}

func (f *FontDescriptorObj) SetFont(font IFont) {
	f.font = font
}

func (f *FontDescriptorObj) GetFont() IFont {
	return f.font
}

func (f *FontDescriptorObj) SetFontFileObjRelate(relate string) {
	f.fontFileObjRelate = relate
}
