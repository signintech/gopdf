package gopdf

import (
	"fmt"
	"io"
)

type TransparencyXObjectGroup struct {
	Index            int
	BBox             [4]float64
	Matrix           [6]float64
	ExtGStateIndexes []int
	XObjects         []cacheContentImage

	getRoot       func() *GoPdf
	pdfProtection *PDFProtection
}

type TransparencyXObjectGroupOptions struct {
	Protection       *PDFProtection
	ExtGStateIndexes []int
	BBox             [4]float64
	XObjects         []cacheContentImage
}

func GetCachedTransparencyXObjectGroup(opts TransparencyXObjectGroupOptions, gp *GoPdf) (TransparencyXObjectGroup, error) {
	group := TransparencyXObjectGroup{
		BBox:             opts.BBox,
		XObjects:         opts.XObjects,
		pdfProtection:    opts.Protection,
		ExtGStateIndexes: opts.ExtGStateIndexes,
	}
	group.Index = gp.addObj(group)
	group.init(func() *GoPdf {
		return gp
	})

	return group, nil
}

func (s TransparencyXObjectGroup) init(funcGetRoot func() *GoPdf) {
	s.getRoot = funcGetRoot
}

func (s *TransparencyXObjectGroup) setProtection(p *PDFProtection) {
	s.pdfProtection = p
}

func (s TransparencyXObjectGroup) protection() *PDFProtection {
	return s.pdfProtection
}

func (s TransparencyXObjectGroup) getType() string {
	return "XObject"
}

func (s TransparencyXObjectGroup) write(w io.Writer, objId int) error {
	streamBuff := GetBuffer()
	defer PutBuffer(streamBuff)

	for _, XObject := range s.XObjects {
		if err := XObject.write(streamBuff, nil); err != nil {
			return err
		}
	}

	content := "<<\n"
	content += "\t/FormType 1\n"
	content += "\t/Subtype /Form\n"
	content += fmt.Sprintf("\t/Type /%s\n", s.getType())
	content += fmt.Sprintf("\t/Matrix [1 0 0 1 0 0]\n")
	content += fmt.Sprintf("\t/BBox [%.3F %.3F %.3F %.3F]\n", s.BBox[0], s.BBox[1], s.BBox[2], s.BBox[3])
	content += "\t/Group<</CS /DeviceGray /S /Transparency>>\n"

	content += fmt.Sprintf("\t/Length %d\n", len(streamBuff.Bytes()))
	content += ">>\n"
	content += "stream\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	if _, err := w.Write(streamBuff.Bytes()); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "endstream\n"); err != nil {
		return err
	}

	return nil
}
