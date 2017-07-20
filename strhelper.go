package gopdf

import (
	"math/big"
	"strings"
)

//StrHelperGetStringWidth get string width
func StrHelperGetStringWidth(str string, fontSize int, ifont IFont) float64 {
	w := 0
	bs := []byte(str)
	i := 0
	max := len(bs)
	for i < max {
		w += ifont.GetCw()[bs[i]]
		i++
	}
	return float64(w) * (float64(fontSize) / 1000.0)
}

//CreateEmbeddedFontSubsetName create Embedded font (subset font) name
func CreateEmbeddedFontSubsetName(name string) string {
	name = strings.Replace(name, " ", "+", -1)
	name = strings.Replace(name, "/", "+", -1)
	return name
}

//ReadShortFromByte read short from byte array
func ReadShortFromByte(data []byte, offset int) (int64, int) {
	buff := data[offset : offset+2]
	num := big.NewInt(0)
	num.SetBytes(buff)
	u := num.Uint64()
	var v int64
	if u >= 0x8000 {
		v = int64(u) - 65536
	} else {
		v = int64(u)
	}
	return v, 2
}

//ReadUShortFromByte read ushort from byte array
func ReadUShortFromByte(data []byte, offset int) (uint64, int) {
	buff := data[offset : offset+2]
	num := big.NewInt(0)
	num.SetBytes(buff)
	return num.Uint64(), 2
}
