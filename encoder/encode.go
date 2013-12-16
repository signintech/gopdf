package encoder

import (
	"strings"
)


type EnCodeInfo struct{
	Utf8 int64
	Ucode int64
}


func TranslateFromUtf8(enCodeInfos  *[]EnCodeInfo,str string) string{

	chars := strings.Split(str,"")
	outputBuffer := make([]byte, len([]byte(str))*2)
	i := 0
	j := 0
	num := int64(0)
	multiply := int64(0)
	max := len(chars)
	clen := 0
	for i < max {
		b := []byte(chars[i])
		j = len(b)
		num = 0
		multiply = 0x1
		for j > 0 {
			num = num + ( int64(b[j - 1 ]) * multiply ) 
			multiply = multiply * 0x100
			j--
		}
		if  num >= 0xe0b880 && num <= 0xe0b99b	 {
			j = 0
			clen = len(*enCodeInfos)
			for j < clen {
				if( (*enCodeInfos)[j].Utf8 == num ){
					num = (*enCodeInfos)[j].Ucode
					break
				}
				j++
			}
		}
		if num <= 255 {
			outputBuffer[i] = byte(num)
		}
		i++
	}
	return  string(outputBuffer[0:i] )
}