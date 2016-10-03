package gopdf

import (
	"bytes"
	"fmt"
)

//SMask smask
type SMask struct {
	buffer bytes.Buffer
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
func (s *SMask) getObjBuff() *bytes.Buffer {
	return &s.buffer
}

func (s *SMask) build(objID int) error {

	buff, err := buildImgProp(s.imgInfo)
	if err != nil {
		return err
	}
	_, err = buff.WriteTo(&s.buffer)
	if err != nil {
		return err
	}

	s.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(s.data))) // /Length 62303>>\n
	s.buffer.WriteString("stream\n")
	if s.protection() != nil {
		tmp, err := rc4Cip(s.protection().objectkey(objID), s.data)
		if err != nil {
			return err
		}
		s.buffer.Write(tmp)
		s.buffer.WriteString("\n")
	} else {
		s.buffer.Write(s.data)
	}
	s.buffer.WriteString("\nendstream\n")

	return nil
}
