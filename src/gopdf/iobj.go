package gopdf

import (
	"bytes"
)

type IObj interface {
	Init(func()(*GoPdf))
	GetType() string
	GetObjBuff() *bytes.Buffer
	Build()
	
}
