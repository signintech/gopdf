package gopdf

import (
	"bytes"
	"fmt"
	"strconv"
)

//PagesObj pdf pages object
type PagesObj struct { //impl IObj
	buffer    bytes.Buffer
	PageCount int
	Kids      string
	getRoot   func() *GoPdf
}

func (p *PagesObj) init(funcGetRoot func() *GoPdf) {
	p.PageCount = 0
	p.getRoot = funcGetRoot
}

func (p *PagesObj) build(objID int) error {

	height := fmt.Sprintf("%0.2f", p.getRoot().config.PageSize.H)
	width := fmt.Sprintf("%0.2f", p.getRoot().config.PageSize.W)
	p.buffer.WriteString("<<\n")
	p.buffer.WriteString("  /Type /" + p.getType() + "\n")
	p.buffer.WriteString("  /MediaBox [ 0 0 " + width + " " + height + " ]\n")
	p.buffer.WriteString("  /Count " + strconv.Itoa(p.PageCount) + "\n")
	p.buffer.WriteString("  /Kids [ " + p.Kids + " ]\n") //sample Kids [ 3 0 R ]
	p.buffer.WriteString(">>\n")
	return nil
}

func (p *PagesObj) getType() string {
	return "Pages"
}

func (p *PagesObj) getObjBuff() *bytes.Buffer {
	return &(p.buffer)
}

func (p *PagesObj) test() {
	fmt.Print(p.getType() + "\n")
}
