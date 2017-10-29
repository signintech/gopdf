package gopdf

import (
	"fmt"
	"io"
)

//SMask smask
type SMask struct {
	imgInfo
	data []byte
	//getRoot func() *GoPdf
	pdfProtection *PDFProtection
}

func (s *SMask) init(funcGetRoot func() *GoPdf) {
	//s.getRoot = funcGetRoot
}

func (s *SMask) setProtection(p *PDFProtection) {
	s.pdfProtection = p
}

func (s *SMask) protection() *PDFProtection {
	return s.pdfProtection
}

func (s *SMask) getType() string {
	return "smask"
}

func (s *SMask) write(w io.Writer, objID int) error {

	err := writeImgProp(w, s.imgInfo)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "/Length %d\n>>\n", len(s.data)) // /Length 62303>>\n
	io.WriteString(w, "stream\n")
	if s.protection() != nil {
		tmp, err := rc4Cip(s.protection().objectkey(objID), s.data)
		if err != nil {
			return err
		}
		w.Write(tmp)
		io.WriteString(w, "\n")
	} else {
		w.Write(s.data)
	}
	io.WriteString(w, "\nendstream\n")

	return nil
}
