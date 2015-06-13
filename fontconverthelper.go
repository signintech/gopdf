package gopdf

import (
	"strconv"
	//"fmt"
	"bytes"
)

func FontConvertHelper_Cw2Str(cw FontCw) string {
	buff := new(bytes.Buffer)
	buff.WriteString(" ")
	i := 32
	for i <= 255 {
		buff.WriteString(strconv.Itoa(cw[Chr(i)]) + " ")
		i++
	}
	return buff.String()
}
