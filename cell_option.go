package gopdf

//AlignLeft align on left of cell
const AlignLeft = 0

//AlignCenter align on center of cell
const AlignCenter = 1

//AlignRight align on right of cell
const AlignRight = 2

//VAlignTop vertical align on top of cell
const VAlignTop = 0

//VAlignMiddle vertical align on middle of cell
const VAlignMiddle = 1

//VAlignBottom  vertical align on bottom of cell
const VAlignBottom = 2

type CellOption struct {
	Align  int
	VAlign int
}
