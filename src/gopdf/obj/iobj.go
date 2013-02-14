package obj

import (
	"bytes"
)

type IObj interface {
	Init()
	GetType() string
	GetObjBuff() *bytes.Buffer
	Build()
}
