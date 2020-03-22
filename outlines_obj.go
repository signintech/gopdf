package gopdf

import (
	"fmt"
	"io"
)

//OutlinesObj : outlines dictionary
type OutlinesObj struct { //impl IObj
	getRoot func() *GoPdf

	index   int
	first   int
	last    int
	count   int
	lastObj *OutlineObj
}

func (o *OutlinesObj) init(funcGetRoot func() *GoPdf) {
	o.getRoot = funcGetRoot
	o.first = -1
	o.last = -1
}

func (o *OutlinesObj) getType() string {
	return "Outlines"
}

func (o *OutlinesObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Type /%s\n", o.getType())
	fmt.Fprintf(w, "  /Count %d\n", o.count)
	fmt.Fprintf(w, "  /First %d 0 R\n", o.first)
	fmt.Fprintf(w, "  /Last %d 0 R\n", o.last)
	io.WriteString(w, ">>\n")
	return nil
}

func (o *OutlinesObj) SetIndexObjOutlines(index int) {
	o.index = index
}

func (o *OutlinesObj) AddOutline(dest int, title string) {
	oo := &OutlineObj{title: title, dest: dest, parent: o.index, prev: o.last, next: -1}
	o.last = o.getRoot().addObj(oo) + 1
	if o.first <= 0 {
		o.first = o.last
	}
	if o.lastObj != nil {
		o.lastObj.next = o.last
	}
	o.lastObj = oo
	o.count++
}

func (o *OutlinesObj) Count() int {
	return o.count
}

type OutlineObj struct { //impl IObj
	title string

	dest   int
	parent int
	prev   int
	next   int
}

func (o *OutlineObj) init(funcGetRoot func() *GoPdf) {
}

func (o *OutlineObj) getType() string {
	return "Outline"
}

func (o *OutlineObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	fmt.Fprintf(w, "  /Parent %d 0 R\n", o.parent + 1)
	if o.prev >= 0 {
		fmt.Fprintf(w, "  /Prev %d 0 R\n", o.prev)
	}
	if o.next >= 0 {
		fmt.Fprintf(w, "  /Next %d 0 R\n", o.next)
	}
	fmt.Fprintf(w, "  /Dest [ %d 0 R /XYZ null null null ]\n", o.dest)
	fmt.Fprintf(w, "  /Title <FEFF%s>\n", encodeUtf8(o.title))
	io.WriteString(w, ">>\n")
	return nil
}
