package gopdf

import (
	"bytes"
	//"fmt"
)

type CatalogObj struct { //impl IObj
	buffer bytes.Buffer
}

func (me *CatalogObj) Init(funcGetRoot func()(*GoPdf)) {
	
}

func (me *CatalogObj) Build() {
	me.buffer.WriteString("<<\n")
	me.buffer.WriteString("  /Type /" + me.GetType() + "\n")
	me.buffer.WriteString("  /Pages 2 0 R\n")
	me.buffer.WriteString(">>\n")
}

func (me *CatalogObj) GetType() string {
	return "Catalog"
}

func (me *CatalogObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

