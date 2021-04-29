package gopdf

import (
	"bytes"
	"fmt"
	"io"
)

type TransparencyXObjectGroup struct {
	Index    int
	BBox     [4]float64
	Matrix   [6]float64
	XObjects []cacheContentImage

	pdfProtection *PDFProtection
}

type TransparencyXObjectGroupOptions struct {
	X          float64
	Y          float64
	BBox       *Rect
	Protection *PDFProtection
	XObjects   []cacheContentImage
}

func NewTransparencyXObjectGroup(opts TransparencyXObjectGroupOptions, gp *GoPdf) (TransparencyXObjectGroup, error) {
	group := TransparencyXObjectGroup{
		XObjects:      opts.XObjects,
		pdfProtection: opts.Protection,
		BBox:          [4]float64{0, 0, opts.X, opts.Y},
	}
	group.Index = gp.addObj(group)

	return group, nil
}

func (s TransparencyXObjectGroup) init(func() *GoPdf) {}

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
	streamBuff := new(bytes.Buffer)
	for _, XObject := range s.XObjects {
		if err := XObject.write(streamBuff, nil); err != nil {
			return err
		}
	}

	content := "<<\n"
	content += "\t/Group<</CS /DeviceGray /S /Transparency>>\n"
	content += fmt.Sprintf("\t/Type /%s\n", s.getType())
	content += "\t/Resources<<\n\t\t/XObject<<\n"

	for _, XObject := range s.XObjects {
		content += fmt.Sprintf("\t\t\t/I%d 0 R\n", XObject.index)
	}

	content += "\t\t>>\n\t>>\n"

	content += "\t/FormType 1\n"
	content += "\t/Subtype /Form\n"
	content += fmt.Sprintf("\t/Matrix [1 0 0 1 0 0]\n")
	content += fmt.Sprintf("\t/BBox [%.3F %.3F %.3F %.3F]\n", s.BBox[0], s.BBox[1], s.BBox[2], s.BBox[3])
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
