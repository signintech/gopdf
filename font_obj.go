package gopdf

import (
	"fmt"
	"io"
)

//FontObj font obj
type FontObj struct {
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

func (f *FontObj) write(w io.Writer, objID int) error {
	baseFont := f.Family
	if f.Font != nil {
		baseFont = f.Font.GetName()
	}

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", f.getType())
	io.WriteString(w, "  /Subtype /TrueType\n")
	fmt.Fprintf(w, "  /BaseFont /%s\n", baseFont)
	if f.IsEmbedFont {
		io.WriteString(w, "  /FirstChar 32 /LastChar 255\n")
		fmt.Fprintf(w, "  /Widths %d 0 R\n", f.indexObjWidth)
		fmt.Fprintf(w, "  /FontDescriptor %d 0 R\n", f.indexObjFontDescriptor)
		fmt.Fprintf(w, "  /Encoding %d 0 R\n", f.indexObjEncoding)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (f *FontObj) getType() string {
	return "Font"
}

// SetIndexObjWidth sets the width of a font object.
func (f *FontObj) SetIndexObjWidth(index int) {
	f.indexObjWidth = index
}

// SetIndexObjFontDescriptor sets the font descriptor.
func (f *FontObj) SetIndexObjFontDescriptor(index int) {
	f.indexObjFontDescriptor = index
}

// SetIndexObjEncoding sets the encoding.
func (f *FontObj) SetIndexObjEncoding(index int) {
	f.indexObjEncoding = index
}
