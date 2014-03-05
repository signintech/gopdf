package gopdf

import (
	
)

type Current struct{
	X float64
	Y float64
	
	//font
	IndexOfFontObj int
	CountOfFont int
	CountOfL int
	
	Font_Size int
	Font_Style string
	Font_IFont IFont
	Font_FontCount int
	
	//page
	IndexOfPageObj int

	//img
	CountOfImg int
	//paths ของรูปที่แสดงใน pdf ใบนั้นๆ เพราะเราจะไม่ยัดรูปซ้ำลงใน pdf เพื่อขนาดของ pdf จะเล็กลง
	ImgCaches []ImageCache
}

type ImageCache struct{
	Path string
	Index int
}
