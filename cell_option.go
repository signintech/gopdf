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
	Align  int //Allows to align the text. Possible values are: Left,Center,Right
	VAlign int //Allows to  vertical align the text. Possible values are: Top,Center,Bottom
	Border int //Indicates if borders must be drawn around the cell. Possible values are: Left, Top, Right, Bottom, ALL
	Float  int //Indicates where the current position should go after the call. Possible values are: Right, Bottom
}
