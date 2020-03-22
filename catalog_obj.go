package gopdf

import (
	"fmt"
	"io"
)

//CatalogObj : catalog dictionary
type CatalogObj struct { //impl IObj
	outlinesObjID int
}

func (c *CatalogObj) init(funcGetRoot func() *GoPdf) {
        c.outlinesObjID = -1

}

func (c *CatalogObj) getType() string {
	return "Catalog"
}

func (c *CatalogObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", c.getType())
	io.WriteString(w, "  /Pages 2 0 R\n")
	if c.outlinesObjID >= 0 {
		io.WriteString(w, "  /PageMode /UseOutlines\n")
		fmt.Fprintf(w, "  /Outlines %d 0 R\n", c.outlinesObjID)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (c *CatalogObj) SetIndexObjOutlines(index int) {
	c.outlinesObjID = index + 1
}
