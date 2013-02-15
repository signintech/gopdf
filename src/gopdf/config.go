package gopdf

import (

)


type Config struct{
	//pt , mm , cm , in
	Unit string 
	PageSize Rect
	K float64
}


/*
func  (me * Config) calK() float64 {
	k := 1.0 //pt
	if(me.Unit == "pt" ){
		k = 1.0
	}else if( me.Unit == "cm" ){
		k = 72.0/2.54;
	}else if( me.Unit =="in" ){
		k = 72.0/2.54
	}else if( me.Unit =="mm" ){
		k = 72.0/25.4
	}
	return k
}
*/

