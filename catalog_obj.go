package gopdf

import (
	"bytes"
	//"fmt"
)

//CatalogObj : catalog dictionary
type CatalogObj struct { //impl IObj
	buffer bytes.Buffer
}

func (c *CatalogObj) init(funcGetRoot func() *GoPdf) {

}

func (c *CatalogObj) build(objID int) error {
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("  /Type /" + c.getType() + "\n")
	c.buffer.WriteString("  /Pages 2 0 R\n")
	c.buffer.WriteString(">>\n")
	return nil
}

func (c *CatalogObj) getType() string {
	return "Catalog"
}

func (c *CatalogObj) getObjBuff() *bytes.Buffer {
	return &(c.buffer)
}
