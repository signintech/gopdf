package gopdf

import "bytes"

type PdfDictionaryObj struct {
	buffer             bytes.Buffer
	PtrToSubsetFontObj *SubsetFontObj
}

func (me *PdfDictionaryObj) Init(funcGetRoot func() *GoPdf) {
}

func (me *PdfDictionaryObj) Build() {

}
func (me *PdfDictionaryObj) GetType() string {
	return "PdfDictionary"
}

func (me *PdfDictionaryObj) GetObjBuff() *bytes.Buffer {
	return &me.buffer
}

func (me *PdfDictionaryObj) SetPtrToSubsetFontObj(ptr *SubsetFontObj) {
	me.PtrToSubsetFontObj = ptr
}
