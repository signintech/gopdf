package gopdf

import (
	"bytes"
)

//IObj inteface for all pdf object
type IObj interface {
	init(func() *GoPdf)
	getType() string
	getObjBuff() *bytes.Buffer
	//สร้าง ข้อมูลใน pdf
	build(objID int) error
}
