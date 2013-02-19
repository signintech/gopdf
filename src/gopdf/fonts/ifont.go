package fonts

import (
	"fmt"
)

type IFont interface{
	Init()
	GetType() string
	GetName() string
	GetDesc() []FontDescItem
	GetUp() int
	GetUt()  int
	GetCw() map[string]int
	GetEnc() string
	GetDiff() string
	GetOriginalsize() int
	
	SetFamily(family string)
	GetFamily() string
}

type FontDescItem struct{
	Key string
	Val string
}

func Chr(n int) string{
	return fmt.Sprintf("%c",n)
}