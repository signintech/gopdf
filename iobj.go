package gopdf

import (
	"bytes"
)

type IObj interface {
	Init(func()(*GoPdf))
	GetType() string
	GetObjBuff() *bytes.Buffer
	//สร้าง ข้อมูลใน pdf
	Build()
	
}
