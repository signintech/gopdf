package gopdf

import (
	"bytes"
)

type FontDescriptorObj struct{
	buffer bytes.Buffer
	font IFont
	fontFileObjRelate string
}

func (me *FontDescriptorObj) Init(funcGetRoot func()(*GoPdf)) {
	
}

func (me *FontDescriptorObj) Build() {
	
	me.buffer.WriteString("<</Type /FontDescriptor /FontName /"+me.font.GetName() +" ")
	descs := me.font.GetDesc()
	i := 0
	max := len(descs)
	for i < max {
		me.buffer.WriteString("/"+ descs[i].Key +" " + descs[i].Val +" ")
		i++
	}
	
	if me.GetType() == "Type1" {
		me.buffer.WriteString("/FontFile ")
	}else {
		me.buffer.WriteString("/FontFile2 ")
	}
	
	me.buffer.WriteString(me.fontFileObjRelate)
	me.buffer.WriteString(">>\n")
}

func (me *FontDescriptorObj) GetType() string {
	return "FontDescriptor"
}

func (me *FontDescriptorObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me * FontDescriptorObj) SetFont(font IFont) {
	me.font = font
}

func (me * FontDescriptorObj) GetFont() IFont{
	return me.font
}

func (me * FontDescriptorObj) SetFontFileObjRelate(relate string){
	me.fontFileObjRelate = relate
}

