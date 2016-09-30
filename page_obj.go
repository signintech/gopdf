package gopdf

import (
	"bytes"
	//"fmt"
)

//PageObj pdf page object
type PageObj struct { //impl IObj
	buffer          bytes.Buffer
	Contents        string
	ResourcesRelate string
}

func (p *PageObj) init(funcGetRoot func() *GoPdf) {

}

func (p *PageObj) build(objID int) error {

	p.buffer.WriteString("<<\n")
	p.buffer.WriteString("  /Type /" + p.getType() + "\n")
	p.buffer.WriteString("  /Parent 2 0 R\n")
	p.buffer.WriteString("  /Resources " + p.ResourcesRelate + "\n")
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
	p.buffer.WriteString("  /Contents " + p.Contents + "\n") //sample  Contents 8 0 R
	p.buffer.WriteString(">>\n")
	return nil
}

func (p *PageObj) getType() string {
	return "Page"
}

func (p *PageObj) getObjBuff() *bytes.Buffer {
	return &(p.buffer)
}
