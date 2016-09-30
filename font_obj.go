package gopdf

import (
	"bytes"
	"strconv"
)

//FontObj font obj
type FontObj struct {
	buffer bytes.Buffer
	Family string
	//Style string
	//Size int
	IsEmbedFont bool

	indexObjWidth          int
	indexObjFontDescriptor int
	indexObjEncoding       int

	Font        IFont
	CountOfFont int
}

func (f *FontObj) init(funcGetRoot func() *GoPdf) {
	f.IsEmbedFont = false
	//me.CountOfFont = -1
}

func (f *FontObj) build(objID int) error {

	baseFont := f.Family
	if f.Font != nil {
		baseFont = f.Font.GetName()
	}

	f.buffer.WriteString("<<\n")
	f.buffer.WriteString("  /Type /" + f.getType() + "\n")
	f.buffer.WriteString("  /Subtype /TrueType\n")
	f.buffer.WriteString("  /BaseFont /" + baseFont + "\n")
	if f.IsEmbedFont {
		f.buffer.WriteString("  /FirstChar 32 /LastChar 255\n")
		f.buffer.WriteString("  /Widths " + strconv.Itoa(f.indexObjWidth) + " 0 R\n")
		f.buffer.WriteString("  /FontDescriptor " + strconv.Itoa(f.indexObjFontDescriptor) + " 0 R\n")
		f.buffer.WriteString("  /Encoding " + strconv.Itoa(f.indexObjEncoding) + " 0 R\n")
	}
	f.buffer.WriteString(">>\n")
	return nil
}

func (f *FontObj) getType() string {
	return "Font"
}

func (f *FontObj) getObjBuff() *bytes.Buffer {
	return &(f.buffer)
}

func (f *FontObj) SetIndexObjWidth(index int) {
	f.indexObjWidth = index
}

func (f *FontObj) SetIndexObjFontDescriptor(index int) {
	f.indexObjFontDescriptor = index
}

func (f *FontObj) SetIndexObjEncoding(index int) {
	f.indexObjEncoding = index
}
