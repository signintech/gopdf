package gopdf

import (
	"fmt"
	"io"
)

type SMaskSubtypes string

const (
	SMaskAlphaSubtype      = "/Alpha"
	SMaskLuminositySubtype = "/Luminosity"
)

//SMask smask
type SMask struct {
	imgInfo
	data []byte
	//getRoot func() *GoPdf
	pdfProtection                 *PDFProtection
	Index                         int
	TransparencyXObjectGroupIndex int
	S                             string
}

type SMaskOptions struct {
	X                float64
	Y                float64
	Subtype          SMaskSubtypes
	ExtGStateIndexes []int
	Images           []cacheContentImage
}

func NewSMask(opts SMaskOptions, gp *GoPdf) (SMask, error) {
	groupOpts := TransparencyXObjectGroupOptions{
		X:                opts.X,
		Y:                opts.Y,
		XObjects:         opts.Images,
		ExtGStateIndexes: opts.ExtGStateIndexes,
	}
	transparencyXObjectGroup, err := NewTransparencyXObjectGroup(groupOpts, gp)
	if err != nil {
		return SMask{}, err
	}

	smask := SMask{
		S:                             string(opts.Subtype),
		TransparencyXObjectGroupIndex: transparencyXObjectGroup.Index,
	}

	smask.Index = gp.addObj(smask)

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
	if s.TransparencyXObjectGroupIndex != 0 {
		content := "<<\n"
		content += "\t/Type /Mask\n"
		content += fmt.Sprintf("\t/S %s\n", s.S)
		content += fmt.Sprintf("\t/G %d 0 R\n", s.TransparencyXObjectGroupIndex+1)
		content += ">>\n"

		if _, err := io.WriteString(w, content); err != nil {
			return err
		}

	} else {
		err := writeImgProps(w, s.imgInfo)
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
	}

	return nil
}
