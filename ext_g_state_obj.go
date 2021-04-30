package gopdf

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// TODO: add all fields https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf 8.4.5 page 128
type ExtGState struct {
	Index      int
	ca         float64
	CA         float64
	BM         string
	SMaskIndex int
}

type ExtGStateOptions struct {
	StrokingCA    *float64
	NonStrokingCa *float64
	BlendMode     *string
	SMaskOptions  *SMaskOptions
}

func (extOpt ExtGStateOptions) GetId() string {
	id := ""
	if extOpt.StrokingCA != nil {
		id += fmt.Sprintf("CA_%.3f;", *extOpt.StrokingCA)
	}
	if extOpt.NonStrokingCa != nil {
		id += fmt.Sprintf("ca_%.3f;", *extOpt.NonStrokingCa)
	}
	if extOpt.BlendMode != nil {
		id += fmt.Sprintf("BM_%s;", *extOpt.BlendMode)
	}

	if extOpt.SMaskOptions != nil {
		id += "maskImgs_"

		maskImgs := make([]string, len(extOpt.SMaskOptions.Images))
		for ind, maskImg := range extOpt.SMaskOptions.Images {
			maskImgs[ind] = fmt.Sprintf("%d", maskImg.index)
		}

		id += strings.Join(maskImgs, ",") + ";"
	}

	//here

	return id
}

func NewExtGState(opts ExtGStateOptions, gp *GoPdf) (ExtGState, error) {
	state := ExtGState{}
	if opts.BlendMode != nil {
		state.BM = *opts.BlendMode
	}
	if opts.StrokingCA != nil {
		state.CA = *opts.StrokingCA
	}
	if opts.NonStrokingCa != nil {
		state.ca = *opts.NonStrokingCa
	}

	extGState, ok := gp.curr.extGStateMap.Find(opts)
	if !ok {
		if opts.SMaskOptions != nil {
			smask, err := NewSMask(*opts.SMaskOptions, gp)
			if err != nil {
				return ExtGState{}, err
			}

			state.SMaskIndex = smask.Index
		}

		state.Index = gp.addObj(state)

		pdfObj := gp.pdfObjs[gp.indexOfProcSet]
		procset, ok := pdfObj.(*ProcSetObj)
		if !ok {
			return ExtGState{}, errors.New("can't convert pdfobject to procsetobj")
		}
		procset.ExtGStates = append(procset.ExtGStates, ExtGS{Index: state.Index})

		gp.curr.extGStateMap.Save(opts.GetId(), state)

		extGState = state
	}

	return extGState, nil
}

func (egs ExtGState) init(func() *GoPdf) {}

func (egs ExtGState) getType() string {
	return "ExtGState"
}

func (egs ExtGState) write(w io.Writer, objID int) error {
	content := "<<\n"
	content += "\t/Type /ExtGState\n"
	content += fmt.Sprintf("\t/ca %.3F\n", egs.ca)
	content += fmt.Sprintf("\t/CA %.3F\n", egs.CA)
	content += fmt.Sprintf("\t/BM %s\n", egs.BM)

	if egs.SMaskIndex != 0 {
		content += fmt.Sprintf("\t/SMask %d 0 R\n", egs.SMaskIndex+1)
	}

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}

type ExtGStateMap struct {
	syncer sync.Mutex
	table  map[string]ExtGState
}

func NewExtGStateMap() ExtGStateMap {
	return ExtGStateMap{
		syncer: sync.Mutex{},
		table:  make(map[string]ExtGState),
	}
}

func (extm *ExtGStateMap) Find(extGState ExtGStateOptions) (ExtGState, bool) {
	key := extGState.GetId()

	extm.syncer.Lock()
	defer extm.syncer.Unlock()

	t, ok := extm.table[key]
	if !ok {
		return ExtGState{}, false
	}

	return t, ok

}

func (tm *ExtGStateMap) Save(id string, extGState ExtGState) ExtGState {
	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	tm.table[id] = extGState

	return extGState
}
