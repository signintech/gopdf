package gopdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf/fontmaker/core"
)

//SubfontDescriptorObj pdf subfont descriptorObj object
type SubfontDescriptorObj struct {
	buffer                bytes.Buffer
	PtrToSubsetFontObj    *SubsetFontObj
	indexObjPdfDictionary int
}

func (s *SubfontDescriptorObj) init(func() *GoPdf) {}

func (s *SubfontDescriptorObj) getType() string {
	return "SubFontDescriptor"
}
func (s *SubfontDescriptorObj) getObjBuff() *bytes.Buffer {
	return &s.buffer
}

func (s *SubfontDescriptorObj) build(objID int) error {

	ttfp := s.PtrToSubsetFontObj.GetTTFParser()
	//fmt.Printf("-->%d\n", ttfp.UnitsPerEm())
	s.buffer.WriteString("<<\n")
	s.buffer.WriteString("/Type /FontDescriptor\n")
	s.buffer.WriteString(fmt.Sprintf("/Ascent %d\n", DesignUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/CapHeight %d\n", DesignUnitsToPdf(ttfp.CapHeight(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/Descent %d\n", DesignUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(fmt.Sprintf("/Flags %d\n", ttfp.Flag()))
	s.buffer.WriteString(fmt.Sprintf("/FontBBox [%d %d %d %d]\n",
		DesignUnitsToPdf(ttfp.XMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.XMax(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMax(), ttfp.UnitsPerEm()),
	))
	s.buffer.WriteString(fmt.Sprintf("/FontFile2 %d 0 R\n", s.indexObjPdfDictionary+1))
	s.buffer.WriteString(fmt.Sprintf("/FontName /%s\n", CreateEmbeddedFontSubsetName(s.PtrToSubsetFontObj.GetFamily())))
	s.buffer.WriteString(fmt.Sprintf("/ItalicAngle %d\n", ttfp.ItalicAngle()))
	s.buffer.WriteString("/StemV 0\n")
	s.buffer.WriteString(fmt.Sprintf("/XHeight %d\n", DesignUnitsToPdf(ttfp.XHeight(), ttfp.UnitsPerEm())))
	s.buffer.WriteString(">>\n")
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

//GetObjBuff get buffer
func (s *SubfontDescriptorObj) GetObjBuff() *bytes.Buffer {
	return s.getObjBuff()
}

//Build build buffer
func (s *SubfontDescriptorObj) Build(objID int) error {
	return s.build(objID)
}
