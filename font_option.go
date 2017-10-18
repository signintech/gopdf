package gopdf

import (
	"strings"
)

//Regular - font style regular
const Regular = 0 //000000
//Italic - font style italic
const Italic = 1 //000001
//Bold - font style bold
const Bold = 2 //000010
//Underline - font style underline
const Underline = 4 //000100

func getConvertedStyle(fontStyle string) (style int) {
	fontStyle  = strings.ToUpper(fontStyle)
	if strings.Contains(fontStyle, "B") {
		style=style|Bold
	}
	if strings.Contains(fontStyle, "I") {
		style=style|Italic
	}
	if strings.Contains(fontStyle, "U"){
		style=style|Underline
	}
	return
}