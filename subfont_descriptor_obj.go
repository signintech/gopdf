package gopdf

import "bytes"

type SubfontDescriptorObj struct {
	buffer bytes.Buffer
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

}
