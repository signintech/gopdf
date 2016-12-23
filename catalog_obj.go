package gopdf

import (
	"bytes"
	"fmt"
	//"fmt"
)

//CatalogObj : catalog dictionary
type CatalogObj struct { //impl IObj
	buffer  bytes.Buffer
	getRoot func() *GoPdf
}

func (c *CatalogObj) init(funcGetRoot func() *GoPdf) {
	c.getRoot = funcGetRoot
}

func (c *CatalogObj) build(objID int) error {
	c.buffer.WriteString("<<\n")
	c.buffer.WriteString("  /Type /" + c.getType() + "\n")
	c.buffer.WriteString("  /Pages 2 0 R\n")
	if c.getRoot != nil && c.getRoot().indexOfEmbedJs != 0 {
		c.buffer.WriteString(fmt.Sprintf(" /Names <</JavaScript %d 0 R>>\n", c.getRoot().indexOfEmbedJs+1))
	}
	c.buffer.WriteString(">>\n")
	return nil
}

func (c *CatalogObj) getType() string {
	return "Catalog"
}

func (c *CatalogObj) getObjBuff() *bytes.Buffer {
	return &(c.buffer)
}
