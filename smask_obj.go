package gopdf

import (
	"fmt"
	"io"
	"sync"
)

type SMaskSubtypes string

const (
	SMaskAlphaSubtype      = "/Alpha"
	SMaskLuminositySubtype = "/Luminosity"
)

// SMask smask
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
	TransparencyXObjectGroupIndex int
	Subtype                       SMaskSubtypes
}

func (smask SMaskOptions) GetId() string {
	id := fmt.Sprintf("S_%s;G_%d_0_R", smask.Subtype, smask.TransparencyXObjectGroupIndex)

	return id
}

func GetCachedMask(opts SMaskOptions, gp *GoPdf) SMask {
	smask, ok := gp.curr.sMasksMap.Find(opts)
	if !ok {
		smask = SMask{
			S:                             string(opts.Subtype),
			TransparencyXObjectGroupIndex: opts.TransparencyXObjectGroupIndex,
		}
		smask.Index = gp.addObj(smask)

		gp.curr.sMasksMap.Save(opts.GetId(), smask)
	}

	return smask
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
		err := writeImgProps(w, s.imgInfo, false)
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

type SMaskMap struct {
	syncer sync.Mutex
	table  map[string]SMask
}

func NewSMaskMap() SMaskMap {
	return SMaskMap{
		syncer: sync.Mutex{},
		table:  make(map[string]SMask),
	}
}

func (smask *SMaskMap) Find(sMask SMaskOptions) (SMask, bool) {
	key := sMask.GetId()

	smask.syncer.Lock()
	defer smask.syncer.Unlock()

	t, ok := smask.table[key]
	if !ok {
		return SMask{}, false
	}

	return t, ok

}

func (smask *SMaskMap) Save(id string, sMask SMask) SMask {
	smask.syncer.Lock()
	defer smask.syncer.Unlock()

	smask.table[id] = sMask

	return sMask
}
