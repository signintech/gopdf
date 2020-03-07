package gopdf

import (
	"fmt"
	"io"
)

// ExtGStateObj is the graphics state parameter dictionary.
// TODO: add all fields https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf 8.4.5 page 128
type ExtGStateObj struct {
	ca float64
	CA float64
	BM string
}

func (egs *ExtGStateObj) init(func() *GoPdf) {}

func (egs *ExtGStateObj) getType() string {
	return "ExtGState"
}

func (egs *ExtGStateObj) write(w io.Writer, objID int) error {
	io.WriteString(w, "<<\n")
	io.WriteString(w, "\t/Type /ExtGState\n")
	//TODO make all fields nullable (reference)
	//if egs.ca != nil {
	fmt.Fprintf(w, "\t/ca %.3F\n", egs.ca)
	//}
	//if egs.CA != nil {
	fmt.Fprintf(w, "\t/CA %.3F\n", egs.CA)
	//}
	//if egs.BM != nil {
	fmt.Fprintf(w, "\t/BM %v\n", egs.BM)
	//}
	io.WriteString(w, ">>\n")
	return nil
}
