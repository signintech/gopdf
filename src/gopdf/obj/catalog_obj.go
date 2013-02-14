package obj

import (
	"bytes"
	//"fmt"
)

type CatalogObj struct { //impl IObj
	buffer bytes.Buffer
}

func (me *CatalogObj) Init() {
}

func (me *CatalogObj) Build() {
	me.buffer.WriteString("\t/Type /" + me.GetType() + "\n")
	me.buffer.WriteString("\t/Pages 2 0 R\n")
}

func (me *CatalogObj) GetType() string {
	return "Catalog"
}

func (me *CatalogObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}
