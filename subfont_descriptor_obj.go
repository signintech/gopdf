package gopdf

import (
	"bytes"
	"fmt"
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
	me.buffer.WriteString(fmt.Sprintf(">>>%d    %d\n", ttfp.TypoAscender(), ttfp.TypoDescender()))
	me.buffer.WriteString(">>\n")
}

func (me *SubfontDescriptorObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}
