package gopdf

import (
	"fmt"
	"io"
)

// CIDFontObj is a CID-keyed font.
// cf. https://www.adobe.com/content/dam/acom/en/devnet/font/pdfs/5014.CIDFont_Spec.pdf
type CIDFontObj struct {
	PtrToSubsetFontObj        *SubsetFontObj
	indexObjSubfontDescriptor int
}

func (ci *CIDFontObj) init(funcGetRoot func() *GoPdf) {
}

//SetIndexObjSubfontDescriptor set  indexObjSubfontDescriptor
func (ci *CIDFontObj) SetIndexObjSubfontDescriptor(index int) {
	ci.indexObjSubfontDescriptor = index
}

func (ci *CIDFontObj) getType() string {
	return "CIDFont"
}

func (ci *CIDFontObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "/BaseFont /%s\n", CreateEmbeddedFontSubsetName(ci.PtrToSubsetFontObj.GetFamily()))
	io.WriteString(w, "/CIDSystemInfo\n")
	io.WriteString(w, "<<\n")
	io.WriteString(w, "  /Ordering (Identity)\n")
	io.WriteString(w, "  /Registry (Adobe)\n")
	io.WriteString(w, "  /Supplement 0\n")
	io.WriteString(w, ">>\n")
	fmt.Fprintf(w, "/FontDescriptor %d 0 R\n", ci.indexObjSubfontDescriptor+1) //TODO fix
	io.WriteString(w, "/Subtype /CIDFontType2\n")
	io.WriteString(w, "/Type /Font\n")
	glyphIndexs := ci.PtrToSubsetFontObj.CharacterToGlyphIndex.AllVals()
	io.WriteString(w, "/W [")
	for _, v := range glyphIndexs {
		width := ci.PtrToSubsetFontObj.GlyphIndexToPdfWidth(v)
		fmt.Fprintf(w, "%d[%d]", v, width)
	}
	io.WriteString(w, "]\n")
	io.WriteString(w, ">>\n")
	return nil
}

//SetPtrToSubsetFontObj set PtrToSubsetFontObj
func (ci *CIDFontObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	ci.PtrToSubsetFontObj = ptr
}
