package gopdf

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type SMaskSubtypes string

const (
	SMaskAlphaSubtype      = "Alpha"
	SMaskLuminositySubtype = "Luminosity"
)

//SMask smask
type SMask struct {
	imgInfo
	data []byte
	//getRoot func() *GoPdf
	pdfProtection *PDFProtection

	Index                         int
	TransparencyXObjectGroupIndex int
	S                             string
}

type SMaskOptions struct {
	TransparencyXObjectGroupIndex int
	Subtype                       SMaskSubtypes
}

func NewSMask(opts SMaskOptions, gp *GoPdf) (SMask, error) {
	smask := SMask{
		S: string(opts.Subtype),
	}

	smask.Index = gp.addObj(smask)

	pdfObj := gp.pdfObjs[gp.indexOfProcSet]
	procset, ok := pdfObj.(*ProcSetObj)
	if !ok {
		return SMask{}, errors.New("can't convert pdfobject to procsetobj")
	}
	procset.ExtGStates = append(procset.ExtGStates, ExtGS{Index: smask.Index})

	return smask, nil
}

func (s SMask) init(func() *GoPdf) {}

func (s *SMask) setProtection(p *PDFProtection) {
	s.pdfProtection = p
}

func (s SMask) protection() *PDFProtection {
	return s.pdfProtection
}

func (s SMask) getType() string {
	return "Mask"
}

func (s SMask) write(w io.Writer, objID int) error {

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
