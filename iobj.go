package gopdf

import (
	"bytes"
)

type IObj interface {
	init(func() *GoPdf)
	getType() string
	getObjBuff() *bytes.Buffer
	//สร้าง ข้อมูลใน pdf
	build() error
}
