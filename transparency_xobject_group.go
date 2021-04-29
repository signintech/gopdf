package gopdf

import (
	"fmt"
	"io"
	"math"

	"github.com/pkg/errors"
)

type TransparencyXObjectGroup struct {
	Index    int
	BBox     Rect
	XObjects []cacheContentImage

	pdfProtection *PDFProtection
}

type TransparencyXObjectGroupOptions struct {
	BBox       Rect
	Protection *PDFProtection
	XObjects   []cacheContentImage
}

func NewTransparencyXObjectGroup(opts TransparencyXObjectGroupOptions, gp *GoPdf) (TransparencyXObjectGroup, error) {
	group := TransparencyXObjectGroup{
		BBox:     opts.BBox,
		XObjects: opts.XObjects,
		pdfProtection: opts.Protection,
	}
	group.Index = gp.addObj(group)

	pdfObj := gp.pdfObjs[gp.indexOfProcSet]
	procset, ok := pdfObj.(*ProcSetObj)
	if !ok {
		return TransparencyXObjectGroup{}, errors.New("can't convert pdfobject to procsetobj")
	}
	procset.ExtGStates = append(procset.ExtGStates, ExtGS{Index: group.Index})

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
	stream := "stream\n"
	for _, XObject := range s.XObjects {
		if err := XObject.write(w, nil); err != nil {
			return err
		}
	}

	stream += "endstream\n"

	content := fmt.Sprintf("%d 0 obj<<\n", objId)
	content += "/Group<</CS /DeviceGray /S /Transparency>>\n"
	content += "/Type /" + s.getType() + "\n"
	content += "/Resources<</XObject<<\n"

	for _, XObject := range s.XObjects {
		content += fmt.Sprintf("/I%d 0 R ", XObject.index)
	}

	content += "/Subtype /Form\n"
	content += "/FormType 1\n"
	content += fmt.Sprintf("/BBox [0 0 %d %d]\n", int(math.Round(s.BBox.W)), int(math.Round(s.BBox.H)))
	content += fmt.Sprintf("/Matrix [1 0 0 1 0 0]\n")
	content += fmt.Sprintf("/Length %d\n", len(stream))
	content += ">>\n"
	content += stream

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}
