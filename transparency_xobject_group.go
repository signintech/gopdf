package gopdf

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type TransparencyXObjectGroup struct {
	Index            int
	BBox             [4]float64
	Matrix           [6]float64
	ExtGStateIndexes []int
	XObjects         []cacheContentImage

	pdfProtection *PDFProtection
}

type TransparencyXObjectGroupOptions struct {
	Protection       *PDFProtection
	ExtGStateIndexes []int
	BBox             [4]float64
	XObjects         []cacheContentImage
}

func (groupOpts TransparencyXObjectGroupOptions) GetId() string {
	extGStateId := "ExtGState"
	for _, extGStateInd := range groupOpts.ExtGStateIndexes {
		extGStateId += fmt.Sprintf("_%d", extGStateInd)
	}

	xObjectId := "XObject"
	for _, xObject := range groupOpts.XObjects {
		xObjectId += fmt.Sprintf("_%d", xObject.index)
	}

	id := fmt.Sprintf("%s;%s", extGStateId, xObjectId)

	return id
}

func GetCachedTransparencyXObjectGroup(opts TransparencyXObjectGroupOptions, gp *GoPdf) (TransparencyXObjectGroup, error) {
	group, ok := gp.curr.transparencyXObjectGroupsMap.Find(opts)
	if !ok {
		group = TransparencyXObjectGroup{
			BBox:             opts.BBox,
			XObjects:         opts.XObjects,
			pdfProtection:    opts.Protection,
			ExtGStateIndexes: opts.ExtGStateIndexes,
		}
		group.Index = gp.addObj(group)

		gp.curr.transparencyXObjectGroupsMap.Save(opts.GetId(), group)
	}

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
	content += "\t/FormType 1\n"
	content += "\t/Subtype /Form\n"
	content += fmt.Sprintf("\t/Type /%s\n", s.getType())
	content += fmt.Sprintf("\t/Matrix [1 0 0 1 0 0]\n")
	content += fmt.Sprintf("\t/BBox [%.3F %.3F %.3F %.3F]\n", s.BBox[0], s.BBox[1], s.BBox[2], s.BBox[3])
	content += "\t/Group<</CS /DeviceGray /S /Transparency>>\n"
	content += "\t/Resources<<\n"

	xObjects := "\t\t/XObject<<\n"
	for _, XObject := range s.XObjects {
		xObjects += fmt.Sprintf("\t\t\t/I%d %d 0 R\n", XObject.index+1, XObject.index+1)
	}
	xObjects += "\t\t>>\n"

	extGStates := "\t\t/ExtGState<<\n"
	for _, extGStateIndex := range s.ExtGStateIndexes {
		extGStates += fmt.Sprintf("\t\t\t/GS%d %d 0 R\n", extGStateIndex+1, extGStateIndex)
	}
	extGStates += "\t\t>>\n"

	content += xObjects
	content += extGStates

	content += "\t>>\n"

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

type TransparencyXObjectGroupsMap struct {
	syncer sync.Mutex
	table  map[string]TransparencyXObjectGroup
}

func NewTransparencyXObjectGroupsMap() TransparencyXObjectGroupsMap {
	return TransparencyXObjectGroupsMap{
		syncer: sync.Mutex{},
		table:  make(map[string]TransparencyXObjectGroup),
	}
}

func (tgm *TransparencyXObjectGroupsMap) Find(tgmOpts TransparencyXObjectGroupOptions) (TransparencyXObjectGroup, bool) {
	key := tgmOpts.GetId()

	tgm.syncer.Lock()
	defer tgm.syncer.Unlock()

	t, ok := tgm.table[key]
	if !ok {
		return TransparencyXObjectGroup{}, false
	}

	return t, ok

}

func (tgm *TransparencyXObjectGroupsMap) Save(id string, group TransparencyXObjectGroup) TransparencyXObjectGroup {
	tgm.syncer.Lock()
	defer tgm.syncer.Unlock()

	tgm.table[id] = group

	return group
}
