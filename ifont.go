package gopdf

// IFont represents a font interface.
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

// FontCw maps characters to integers.
type FontCw map[byte]int

// FontDescItem is a (key, value) pair.
type FontDescItem struct {
	Key string
	Val string
}

// // Chr
// func Chr(n int) byte {
// 	return byte(n) //ToByte(fmt.Sprintf("%c", n ))
// }

// ToByte returns the first byte of a string.
func ToByte(chr string) byte {
	return []byte(chr)[0]
}
