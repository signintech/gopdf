package gopdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf/fontmaker/core"
)

type SubfontDescriptorObj struct {
	buffer                bytes.Buffer
	PtrToSubsetFontObj    *SubsetFontObj
	indexObjPdfDictionary int
}

func (me *SubfontDescriptorObj) Init(func() *GoPdf) {

}
func (me *SubfontDescriptorObj) GetType() string {
	return "SubFontDescriptor"
}
func (me *SubfontDescriptorObj) GetObjBuff() *bytes.Buffer {
	return &me.buffer
}

func (me *SubfontDescriptorObj) Build() {
	ttfp := me.PtrToSubsetFontObj.GetTTFParser()
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Type /FontDescriptor\n")
	me.buffer.WriteString(fmt.Sprintf("/Ascent %d\n", DesignUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm())))
	me.buffer.WriteString(fmt.Sprintf("/CapHeight %d\n", DesignUnitsToPdf(ttfp.CapHeight(), ttfp.UnitsPerEm())))
	me.buffer.WriteString(fmt.Sprintf("/Descent %d\n", DesignUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm())))
	me.buffer.WriteString(fmt.Sprintf("/Flags %d\n", ttfp.Flag()))
	me.buffer.WriteString(fmt.Sprintf("/FontBBox [%d %d %d %d]\n",
		DesignUnitsToPdf(ttfp.XMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMin(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.XMax(), ttfp.UnitsPerEm()),
		DesignUnitsToPdf(ttfp.YMax(), ttfp.UnitsPerEm()),
	))
	me.buffer.WriteString(fmt.Sprintf("/FontFile2 %d 0 R\n", me.indexObjPdfDictionary+1))
	me.buffer.WriteString(fmt.Sprintf("/FontName %s\n", CreateEmbeddedFontSubsetName(me.PtrToSubsetFontObj.GetFamily())))
	me.buffer.WriteString(fmt.Sprintf("/ItalicAngle %d\n", ttfp.ItalicAngle()))
	me.buffer.WriteString("/StemV 0\n")
	me.buffer.WriteString(fmt.Sprintf("/XHeight %d\n", DesignUnitsToPdf(ttfp.XHeight(), ttfp.UnitsPerEm())))
	me.buffer.WriteString(">>\n")
}

func (me *SubfontDescriptorObj) SetIndexObjPdfDictionary(index int) {
	me.indexObjPdfDictionary = index
}

func (me *SubfontDescriptorObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}

func DesignUnitsToPdf(val int64, unitsPerEm uint64) int64 {
	return core.Round(float64(float64(val) * 1000.00 / float64(unitsPerEm)))
}
