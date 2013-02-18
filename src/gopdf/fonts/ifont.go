package fonts

import (
	"fmt"
)

type IFont interface{
	Init()
	GetType() string
	GetName() string
	GetDesc() map[string]string
	GetUp() int
	GetUt()  int
	GetCw() map[string]int
	GetEnc() string
	GetDiff() string
	GetOriginalsize() int
}


func Chr(n int) string{
	return fmt.Sprintf("%c",n)
}