package gopdf

// Represents an RGB color with red, green, and blue components
type RGBColor struct {
	R uint8 // Red component (0-255)
	G uint8 // Green component (0-255)
	B uint8 // Blue component (0-255)
}

// Defines the border style for a cell or table
type BorderStyle struct {
	Top      bool     // Whether to draw the top border
	Left     bool     // Whether to draw the left border
	Right    bool     // Whether to draw the right border
	Bottom   bool     // Whether to draw the bottom border
	Width    float64  // Width of the border line
	RGBColor RGBColor // Color of the border
}

// Defines the style for a cell, including border, fill, text, and font properties
type CellStyle struct {
	BorderStyle BorderStyle // Border style for the cell
	FillColor   RGBColor    // Background color of the cell
	TextColor   RGBColor    // Color of the text in the cell
	Font        string      // Font name for the cell text
	FontSize    float64     // Font size for the cell text
}

type RowCell struct {
	content      string    // Content (display value) of the cell
	useCellStyle bool      // If true, use cellStyle instead of the style in the tableLayout
	cellStyle    CellStyle // Style of the cell
}

func NewRowCell(content string, cellStyle CellStyle) RowCell {
	return newStyledRowCell(content, true, cellStyle)

}

func newStyledRowCell(content string, useCellStyle bool, cellStyle CellStyle) RowCell {
	return RowCell{
		content:      content,
		useCellStyle: useCellStyle,
		cellStyle:    cellStyle,
	}
}

type TableLayout interface {
	AddColumn(header string, width float64, align string)
	AddRow(row []string)
	AddStyledRow(row []RowCell)
	SetTableStyle(style CellStyle)
	SetHeaderStyle(style CellStyle)
	SetCellStyle(style CellStyle)
	DrawTable() error
}

// Represents the layout of a table
type tableLayout struct {
	pdf       *GoPdf      // Reference to the GoPdf instance
	startX    float64     // Starting X coordinate of the table
	startY    float64     // Starting Y coordinate of the table
	rowHeight float64     // Height of each row in the table
	columns   []column    // Slice of column definitions
	rows      [][]RowCell // Slice of rows, each containing cell contents
	//styledRows   [][]RowCell // Slice of rows, each containing cell contents and styles.
	maxRows     int        // Maximum number of rows in the table
	padding     float64    // Padding inside each cell
	cellOption  CellOption // Options for cell content rendering
	tableStyle  CellStyle  // Style for the entire table
	headerStyle CellStyle  // Style for the header row
	cellStyle   CellStyle  // Style for regular cells
}

var _ TableLayout = (*tableLayout)(nil)

// Represents a column in the table
type column struct {
	header string  // Header text for the column
	width  float64 // Width of the column
	align  string  // Alignment of content within the column
}

// Creates a new table layout with the given parameters
func (gp *GoPdf) NewTableLayout(startX, startY, rowHeight float64, maxRows int) TableLayout {
	return &tableLayout{
		pdf:       gp,
		startX:    startX,
		startY:    startY,
		rowHeight: rowHeight,
		maxRows:   maxRows,
		padding:   2.0,
		cellOption: CellOption{
			BreakOption: &BreakOption{
				Mode:           BreakModeIndicatorSensitive,
				BreakIndicator: ' ',
			},
		},
		tableStyle: CellStyle{
			BorderStyle: BorderStyle{
				Top: true, Left: true, Right: true, Bottom: true,
				Width:    0.5,
				RGBColor: RGBColor{R: 0, G: 0, B: 0},
			},
		},
		headerStyle: CellStyle{
			BorderStyle: BorderStyle{
				Top: true, Left: true, Right: true, Bottom: true,
				Width:    0.5,
				RGBColor: RGBColor{R: 0, G: 0, B: 0},
			},
			FillColor: RGBColor{R: 240, G: 240, B: 240},
			TextColor: RGBColor{R: 0, G: 0, B: 0},
		},
		cellStyle: CellStyle{
			BorderStyle: BorderStyle{
				Top: true, Left: true, Right: true, Bottom: true,
				Width:    0.5,
				RGBColor: RGBColor{R: 0, G: 0, B: 0},
			},
			TextColor: RGBColor{R: 0, G: 0, B: 0},
		},
	}
}

// Adds a column to the table with the specified header, width, and alignment
func (t *tableLayout) AddColumn(header string, width float64, align string) {
	t.columns = append(t.columns, column{header, width, align})
}

// Adds a row of data to the table
func (t *tableLayout) AddRow(row []string) {
	rowCell := make([]RowCell, len(row))
	for i, cell := range row {
		rowCell[i] = newStyledRowCell(cell, false, CellStyle{})
	}
	t.rows = append(t.rows, rowCell)
}

// Adds a row of data to the table with individual styled cells
// Useful for styling individual cells in a row
func (t *tableLayout) AddStyledRow(row []RowCell) {
	t.rows = append(t.rows, row)
}

// Sets the style for the entire table
func (t *tableLayout) SetTableStyle(style CellStyle) {
	t.tableStyle = style
}

// Sets the style for the header row
func (t *tableLayout) SetHeaderStyle(style CellStyle) {
	t.headerStyle = style
}

// Sets the style for regular cells
func (t *tableLayout) SetCellStyle(style CellStyle) {
	t.cellStyle = style
}

// DrawTable the entire table on the PDF
func (t *tableLayout) DrawTable() error {
	x := t.startX
	y := t.startY

	// Draw the header row
	for _, col := range t.columns {
		if err := t.drawCell(
			x,
			y,
			col.width,
			t.rowHeight,
			col.header,
			"center",
			true, /*isHeader*/
			t.headerStyle,
		); err != nil {
			return err
		}
		x += col.width
	}
	y += t.rowHeight

	// Draw the data rows
	for _, row := range t.rows {
		x = t.startX
		for i, cell := range row {
			cellStyle := t.cellStyle
			if cell.useCellStyle {
				cellStyle = cell.cellStyle
			}
			if err := t.drawCell(
				x,
				y,
				t.columns[i].width,
				t.rowHeight,
				cell.content,
				t.columns[i].align,
				false, /*isHeader*/
				cellStyle,
			); err != nil {
				return err
			}
			x += t.columns[i].width
		}
		y += t.rowHeight
	}

	// Fill any remaining rows with empty cells
	for i := len(t.rows); i < t.maxRows; i++ {
		x = t.startX
		for _, col := range t.columns {
			if err := t.drawCell(
				x,
				y,
				col.width,
				t.rowHeight,
				"",
				col.align,
				false, /*isHeader*/
				t.cellStyle,
			); err != nil {
				return err
			}
			x += col.width
		}
		y += t.rowHeight
	}

	// Draw the outer border of the table and header
	if err := t.drawTableAndHeaderBorder(); err != nil {
		return err
	}

	return nil
}

// Draws the outer border of the table and header
func (t *tableLayout) drawTableAndHeaderBorder() error {
	x1 := t.startX
	y1 := t.startY
	x2 := t.startX
	y2 := t.startY + float64(t.maxRows+1)*t.rowHeight

	for _, col := range t.columns {
		x2 += col.width
	}

	// Draw borders of the table
	err := t.drawBorder(x1, y1, x2, y2, t.tableStyle.BorderStyle)
	if err != nil {
		return err
	}

	// Draw borders of the header
	return t.drawBorder(x1, y1, x2, y1+t.rowHeight, t.headerStyle.BorderStyle)
}

// Draws a single cell of the table
func (t *tableLayout) drawCell(
	x float64,
	y float64,
	width float64,
	height float64,
	content string,
	align string,
	isHeader bool,
	style CellStyle,
) error {
	// Fill the cell background if a fill color is specified
	if style.FillColor != (RGBColor{}) {
		t.pdf.SetFillColor(style.FillColor.R, style.FillColor.G, style.FillColor.B)
		t.pdf.RectFromUpperLeftWithStyle(x, y, width, height, "F")
	}

	if !isHeader {
		// Draw the cell border
		if err := t.drawBorder(x, y, x+width, y+height, style.BorderStyle); err != nil {
			return err
		}
	}

	// Calculate the text area within the cell
	textX := x + t.padding
	textY := y + t.padding
	textWidth := width - (2 * t.padding)
	textHeight := height - (2 * t.padding)

	t.pdf.SetXY(textX, textY)

	// Set the text alignment
	var textOption = t.cellOption
	if align == "right" {
		textOption.Align = Right | Middle
	} else if align == "center" {
		textOption.Align = Center | Middle
	} else {
		textOption.Align = Left | Middle
	}

	// Set the text color and font
	t.pdf.SetTextColor(style.TextColor.R, style.TextColor.G, style.TextColor.B)
	if style.Font != "" {
		t.pdf.SetFont(style.Font, "", style.FontSize)
	}

	// Draw the cell content
	err := t.pdf.MultiCellWithOption(&Rect{W: textWidth, H: textHeight}, content, textOption)
	if err != nil && err.Error() == "empty string" {
		err = nil
	}

	return err
}

// Draws a border around a rectangular area
func (t *tableLayout) drawBorder(x1, y1, x2, y2 float64, borderStyle BorderStyle) error {
	if borderStyle.Width <= 0 {
		return nil
	}
	t.pdf.SetLineWidth(borderStyle.Width)
	t.pdf.SetStrokeColor(borderStyle.RGBColor.R, borderStyle.RGBColor.G, borderStyle.RGBColor.B)
	half := borderStyle.Width / 2.0

	// Draw each side of the border if specified
	if borderStyle.Top {
		t.pdf.Line(x1-half, y1, x2+half, y1)
	}
	if borderStyle.Bottom {
		t.pdf.Line(x1-half, y2, x2+half, y2)
	}
	if borderStyle.Left {
		t.pdf.Line(x1, y1-half, x1, y2+half)
	}
	if borderStyle.Right {
		t.pdf.Line(x2, y1-half, x2, y2+half)
	}

	return nil
}
