package gopdf

import (
	"fmt"
	"io"
)

//CatalogObj : catalog dictionary
type CatalogObj struct { //impl IObj
}

func (c *CatalogObj) init(funcGetRoot func() *GoPdf) {

}

func (c *CatalogObj) getType() string {
	return "Catalog"
}

func (c *CatalogObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", c.getType())
	io.WriteString(w, "  /Pages 2 0 R\n")
	io.WriteString(w, ">>\n")
	return nil
}
