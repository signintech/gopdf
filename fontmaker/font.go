package fontmaker

import (
	"github.com/signintech/gopdf"
)

type Font struct {
	family   string
	fonttype string
	name     string
	desc     []gopdf.FontDescItem
	up       int
	ut       int
	cw       gopdf.FontCw
	enc      string
	diff     string
}

func (me *Font) Init() {

}

func (me *Font) GetType() string {
	return me.fonttype
}
func (me *Font) GetName() string {
	return me.name
}
func (me *Font) GetDesc() []gopdf.FontDescItem {
	return me.desc
}
func (me *Font) GetUp() int {
	return me.up
}
func (me *Font) GetUt() int {
	return me.ut
}
func (me *Font) GetCw() gopdf.FontCw {
	return me.cw
}
func (me *Font) GetEnc() string {
	return me.enc
}
func (me *Font) GetDiff() string {
	return me.diff
}
func (me *Font) GetOriginalsize() int {
	return 98764
}
func (me *Font) SetFamily(family string) {
	me.family = family
}
func (me *Font) GetFamily() string {
	return me.family
}
