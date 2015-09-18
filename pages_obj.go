package gopdf

import (
	"bytes"
	"fmt"
	"strconv"
)

type PagesObj struct { //impl IObj
	buffer    bytes.Buffer
	PageCount int
	Kids      string
	getRoot   func() *GoPdf
}

func (p *PagesObj) Init(funcGetRoot func() *GoPdf) {
	p.PageCount = 0
	p.getRoot = funcGetRoot
}

func (p *PagesObj) Build() error {

	height := fmt.Sprintf("%0.2f", p.getRoot().config.PageSize.H)
	width := fmt.Sprintf("%0.2f", p.getRoot().config.PageSize.W)
	p.buffer.WriteString("<<\n")
	p.buffer.WriteString("  /Type /" + p.GetType() + "\n")
	p.buffer.WriteString("  /MediaBox [ 0 0 " + width + " " + height + " ]\n")
	p.buffer.WriteString("  /Count " + strconv.Itoa(p.PageCount) + "\n")
	p.buffer.WriteString("  /Kids [ " + p.Kids + " ]\n") //sample Kids [ 3 0 R ]
	p.buffer.WriteString(">>\n")
	return nil
}

func (p *PagesObj) GetType() string {
	return "Pages"
}

func (p *PagesObj) GetObjBuff() *bytes.Buffer {
	return &(p.buffer)
}

func (p *PagesObj) Test() {
	fmt.Print(p.GetType() + "\n")
}
