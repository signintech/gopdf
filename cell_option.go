package gopdf

//Left left
const Left = 8 //001000
//Top top
const Top = 4 //000100
//Right right
const Right = 2 //000010
//Bottom bottom
const Bottom = 1 //000001
//Center center
const Center = 16 //010000
//Middle middle
const Middle = 32 //100000

//CellOption cell option
type CellOption struct {
	Align  int //Allows to align the text. Possible values are: Left,Center,Right,Top,Bottom,Middle
	Border int //Indicates if borders must be drawn around the cell. Possible values are: Left, Top, Right, Bottom, ALL
	Float  int //Indicates where the current position should go after the call. Possible values are: Right, Bottom
}
