package gopdf

import (
)

func StrHelper_GetStringWidth(str string,fontSize int,ifont IFont) float64{

	w := 0
	bs := []byte(str)
	i := 0
	max := len(bs)
	for i < max {
		w += ifont.GetCw()[bs[i]]
		i++
	}
	return  float64(w)*(float64(fontSize)/1000.0)
}
