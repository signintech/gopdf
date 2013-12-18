package gopdf


import (
	"bytes"
	"strconv"
)

type FontObj struct { //impl IObj
	buffer bytes.Buffer
	Family string
	//Style string
	//Size int
	IsEmbedFont bool
	
	indexObjWidth int
	indexObjFontDescriptor int
	indexObjEncoding int
	
	Font  IFont
	CountOfFont int
}

func (me *FontObj) Init(funcGetRoot func()(*GoPdf)) {
	me.IsEmbedFont = false
	//me.CountOfFont = -1
}

func (me *FontObj) Build() {
	
	baseFont := me.Family
	if me.Font != nil {
		baseFont =me.Font.GetName()
	}

	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Type /" + me.GetType() + "\n")
	me.buffer.WriteString("  /Subtype /TrueType\n")
	me.buffer.WriteString("  /BaseFont /"+baseFont+"\n")
	if me.IsEmbedFont {
		me.buffer.WriteString("  /FirstChar 32 /LastChar 255\n")
		me.buffer.WriteString("  /Widths "+  strconv.Itoa(me.indexObjWidth) +" 0 R\n")
		me.buffer.WriteString("  /FontDescriptor "+strconv.Itoa(me.indexObjFontDescriptor)+" 0 R\n")
		me.buffer.WriteString("  /Encoding "+strconv.Itoa(me.indexObjEncoding)+" 0 R\n")
	}
	me.buffer.WriteString(">>\n")
}

func (me *FontObj) GetType() string {
	return "Font"
}

func (me *FontObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *FontObj) SetIndexObjWidth( index int){
	me.indexObjWidth = index
}

func (me *FontObj) SetIndexObjFontDescriptor( index int){
	me.indexObjFontDescriptor = index
}

func (me *FontObj) SetIndexObjEncoding( index int){
	me.indexObjEncoding = index
}

