package gopdf

import (
	"io"
)

// EncodingObj is a font object.
type EncodingObj struct {
	font IFont
}

func (e *EncodingObj) init(funcGetRoot func() *GoPdf) {

}
func (e *EncodingObj) getType() string {
	return "Encoding"
}
func (e *EncodingObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<</Type /Encoding /BaseEncoding /WinAnsiEncoding /Differences [")
	io.WriteString(w, e.font.GetDiff())
	io.WriteString(w, "]>>\n")
	return nil
}

// SetFont sets the font of an encoding object.
func (e *EncodingObj) SetFont(font IFont) {
	e.font = font
}

// GetFont gets the font from an encoding object.
func (e *EncodingObj) GetFont() IFont {
	return e.font
}
