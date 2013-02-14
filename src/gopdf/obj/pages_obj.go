package obj

import (
	"bytes"
	"fmt"
	"strconv"
)

type PagesObj struct { //impl IObj
	buffer    bytes.Buffer
	PageCount int
	Kids string
}

func (me *PagesObj) Init() {
	me.PageCount = 0
}

func (me *PagesObj) Build() {
	me.buffer.WriteString("\t/Type /" + me.GetType() + "\n")
	me.buffer.WriteString("\t/MediaBox [ 0 0 200 200 ]\n")
	me.buffer.WriteString("\t/Count "+strconv.Itoa(me.PageCount)+"\n")
	me.buffer.WriteString("\t/Kids [ "+me.Kids+" ]\n") //sample Kids [ 3 0 R ]
}

func (me *PagesObj) GetType() string {
	return "Pages"
}

func (me *PagesObj) GetObjBuff() *bytes.Buffer {
	return &(me.buffer)
}

func (me *PagesObj) Test() {
	fmt.Print(me.GetType() + "\n")
}

