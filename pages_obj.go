package gopdf

import (
	"fmt"
	"io"
)

// PagesObj pdf pages object
type PagesObj struct { //impl IObj
	PageCount int
	Kids      string
	getRoot   func() *GoPdf
}

func (p *PagesObj) init(funcGetRoot func() *GoPdf) {
	p.PageCount = 0
	p.getRoot = funcGetRoot
}

func (p *PagesObj) write(w io.Writer, objID int) error {

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.getType())

	rootConfig := p.getRoot().config
	fmt.Fprintf(w, "  /MediaBox [ 0 0 %0.2f %0.2f ]\n", rootConfig.PageSize.W, rootConfig.PageSize.H)
	fmt.Fprintf(w, "  /Count %d\n", p.PageCount)
	fmt.Fprintf(w, "  /Kids [ %s ]\n", p.Kids) //sample Kids [ 3 0 R ]
	io.WriteString(w, ">>\n")
	return nil
}

func (p *PagesObj) getType() string {
	return "Pages"
}

func (p *PagesObj) test() {
	fmt.Print(p.getType() + "\n")
}
