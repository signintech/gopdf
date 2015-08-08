package gopdf

type IFont interface {
	Init()
	GetType() string
	GetName() string
	GetDesc() []FontDescItem
	GetUp() int
	GetUt() int
	GetCw() FontCw
	GetEnc() string
	GetDiff() string
	GetOriginalsize() int

	SetFamily(family string)
	GetFamily() string
}

type FontCw map[byte]int

type FontDescItem struct {
	Key string
	Val string
}

func Chr(n int) byte {
	return byte(n) //ToByte(fmt.Sprintf("%c", n ))
}

func ToByte(chr string) byte {
	return []byte(chr)[0]
}
