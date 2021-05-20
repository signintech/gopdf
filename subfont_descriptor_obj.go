package gopdf

import (
	"fmt"
	"io"

	"github.com/signintech/gopdf/fontmaker/core"
)

//SubfontDescriptorObj pdf subfont descriptorObj object
type SubfontDescriptorObj struct {
	PtrToSubsetFontObj    *SubsetFontObj
	indexObjPdfDictionary int
}

func (s *SubfontDescriptorObj) init(func() *GoPdf) {}

func (s *SubfontDescriptorObj) getType() string {
	return "SubFontDescriptor"
}

func (s *SubfontDescriptorObj) write(w io.Writer, objID int) error {
	ttfp := s.PtrToSubsetFontObj.GetTTFParser()
	//fmt.Printf("-->%d\n", ttfp.UnitsPerEm())
	io.WriteString(w, "<<\n")
	io.WriteString(w, "/Type /FontDescriptor\n")
	fmt.Fprintf(w, "/Ascent %d\n", DesignUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/CapHeight %d\n", DesignUnitsToPdf(ttfp.CapHeight(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/Descent %d\n", DesignUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm()))
	fmt.Fprintf(w, "/Flags %d\n", ttfp.Flag())
	fmt.Fprintf(w, "/FontBBox [%d %d %d %d]\n",
		DesignUnitsToPdf(ttfp.XMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.XMax(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMax(), ttfp.UnitsPerEm()),
	)
	fmt.Fprintf(w, "/FontFile2 %d 0 R\n", s.indexObjPdfDictionary+1)
	fmt.Fprintf(w, "/FontName /%s\n", CreateEmbeddedFontSubsetName(s.PtrToSubsetFontObj.GetFamily()))
	fmt.Fprintf(w, "/ItalicAngle %d\n", ttfp.ItalicAngle())
	io.WriteString(w, "/StemV 0\n")
	fmt.Fprintf(w, "/XHeight %d\n", DesignUnitsToPdf(ttfp.XHeight(), ttfp.UnitsPerEm()))
	io.WriteString(w, ">>\n")
	return nil
}

//SetIndexObjPdfDictionary set PdfDictionary pointer
func (s *SubfontDescriptorObj) SetIndexObjPdfDictionary(index int) {
	s.indexObjPdfDictionary = index
}

//SetPtrToSubsetFontObj set SubsetFont pointer
func (s *SubfontDescriptorObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	s.PtrToSubsetFontObj = ptr
}

//DesignUnitsToPdf convert unit
func DesignUnitsToPdf(val int, unitsPerEm uint) int {
	return core.Round(float64(float64(val) * 1000.00 / float64(unitsPerEm)))
}
