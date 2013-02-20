package gopdf

import (
	"gopdf/fonts"
	"fmt"
	//iconv "github.com/djimenez/iconv-go"
	//"strings"
)

func StrHelper_GetStringWidth(str string,fontSize int,ifont fonts.IFont) float64{
	/*cw := 600.0 //temp
	ifont.GetCw()[
	w := float64(len(str)) * cw 
	return  w*(float64(fontSize)/1000.0);*/
	/*tmp, _ := iconv.ConvertString(str,"cp874","utf-8")
	fmt.Printf("----%s\n",tmp[0:3])
	i := 0
	w := 0.0
	max := len(tmp)
	for i < max {
		//w += float64(ifont.GetCw()[str[i:i+1]])
		//fmt.Printf("%c\n",3585)
		x := float64(ifont.GetCw()[tmp[i:i+1]])
		w += x
		//fmt.Printf("    %f\n",x)
		i++
	}
	fmt.Printf("%d\n\n\n",max)
	//fmt.Printf("------------------------%d--%f\n",fontSize,w*(float64(fontSize)/1000.0))*/
	//str,_ = iconv.ConvertString(str,"cp874","utf-8")
	//arrays := strings.Split(str, "")
	bs := []byte(str)
	i := 0
	max := len(bs)
	for i < max {
		//chr := fmt.Sprintf("%s",arrays[i:i+1])
		//fmt.Printf("%s %s %d\n", fonts.Chr(161) , chr, ifont.GetCw()[chr])
		//if arrays[i] == fonts.Chr(161 ) {
		//	fmt.Printf("%s  %s -xxxx\n",arrays , str)
		//}
		fmt.Printf("%d\n",bs[i])
		i++
	}
	fmt.Printf("%d\n\n\n",max)
	return  10//w*(float64(fontSize)/1000.0);
}