package gopdf

import (
	"bytes"
	"fmt"

	"github.com/signintech/gopdf/fontmaker/core"
)

type SubfontDescriptorObj struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
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
	_ = ttfp
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("/Type /FontDescriptor\n")
	//fake
	me.buffer.WriteString(fmt.Sprintf(">>>%d    %d\n", DesignUnitsToPdf(ttfp.Ascender(), ttfp.UnitsPerEm()), DesignUnitsToPdf(ttfp.Descender(), ttfp.UnitsPerEm())))
	me.buffer.WriteString(">>\n")
}

func (me *SubfontDescriptorObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}

func DesignUnitsToPdf(val int64, unitsPerEm uint64) int64 {
	return core.Round(float64(float64(val) * 1000.00 / float64(unitsPerEm)))
}
