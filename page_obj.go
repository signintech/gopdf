package gopdf

import (
	"fmt"
	"io"
)

//PageObj pdf page object
type PageObj struct { //impl IObj
	Contents        string
	ResourcesRelate string
	pageOption      PageOption
}

func (p *PageObj) init(funcGetRoot func() *GoPdf) {

}

func (p *PageObj) setOption(opt PageOption) {
	p.pageOption = opt
}

func (p *PageObj) write(w io.Writer, objID int) error {

	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", p.getType())
	io.WriteString(w, "  /Parent 2 0 R\n")
	fmt.Fprintf(w, "  /Resources %s\n", p.ResourcesRelate)
	/*me.buffer.WriteString("    /Font <<\n")
	i := 0
	max := len(me.Realtes)
	for i < max {
		realte := me.Realtes[i]
		me.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n",realte.CountOfFont + 1, realte.IndexOfObj + 1))
		i++
	}
	me.buffer.WriteString("    >>\n")*/
	//me.buffer.WriteString("  >>\n")
	fmt.Fprintf(w, "  /Contents %s\n", p.Contents) //sample  Contents 8 0 R
	if !p.pageOption.isEmpty() {
		fmt.Fprintf(w, " /MediaBox [ 0 0 %0.2f %0.2f ]\n", p.pageOption.PageSize.W, p.pageOption.PageSize.H)
	}
	io.WriteString(w, ">>\n")
	return nil
}

func (p *PageObj) getType() string {
	return "Page"
}
