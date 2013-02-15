package gopdf

import (

)

func StrHelper_GetStringWidth(str string,fontSize int) float64{
	cw := 600.0 //temp
	w := float64(len(str)) * cw 
	return  w*(float64(fontSize)/1000.0);
}