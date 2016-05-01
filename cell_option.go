package gopdf

//Left left
const Left = 8 //1000
//Top top
const Top = 4 //0100
//Right right
const Right = 2 //0010
//Bottom bottom
const Bottom = 1 //0001
//Center center
const Center = 0 //0000
//All Left | Top | Right | Bottom
const All = Left | Top | Right | Bottom //1111

//CellOption cell option
type CellOption struct {
	Align  int
	VAlign int
	Border int
}
