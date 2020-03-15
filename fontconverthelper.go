package gopdf

import (
	"strconv"
	//"fmt"
	"bytes"
)

// FontConvertHelperCw2Str converts main ASCII characters of a FontCW to a string.
func FontConvertHelperCw2Str(cw FontCw) string {
	buff := new(bytes.Buffer)
	buff.WriteString(" ")
	i := 32
	for i <= 255 {
		buff.WriteString(strconv.Itoa(cw[byte(i)]) + " ")
		i++
	}
	return buff.String()
}

// FontConvertHelper_Cw2Str converts main ASCII characters of a FontCW to a string. (for backward compatibility)
// Deprecated: Use FontConvertHelperCw2Str(cw FontCw) instead
func FontConvertHelper_Cw2Str(cw FontCw) string {
	return FontConvertHelperCw2Str(cw)
}
