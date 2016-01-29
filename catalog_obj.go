package gopdf

import (
	"bytes"
	//"fmt"
)

type CatalogObj struct { //impl IObj
	buffer bytes.Buffer
}

func (c *CatalogObj) Init(funcGetRoot func() *GoPdf) {

}

func (c *CatalogObj) Build() error {
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("  /Type /" + c.GetType() + "\n")
	c.buffer.WriteString("  /Pages 2 0 R\n")
	c.buffer.WriteString(">>\n")
	return nil
}

func (c *CatalogObj) GetType() string {
	return "Catalog"
}

func (c *CatalogObj) GetObjBuff() *bytes.Buffer {
	return &(c.buffer)
}
