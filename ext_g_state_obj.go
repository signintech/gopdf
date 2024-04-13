package gopdf

import (
	"fmt"
	"io"
	"sync"

	"errors"
)

// TODO: add all fields https://www.adobe.com/content/dam/acom/en/devnet/acrobat/pdfs/PDF32000_2008.pdf 8.4.5 page 128
type ExtGState struct {
	Index      int
	ca         *float64
	CA         *float64
	BM         *BlendModeType
	SMaskIndex *int
}

type ExtGStateOptions struct {
	StrokingCA    *float64
	NonStrokingCa *float64
	BlendMode     *BlendModeType
	SMaskIndex    *int
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
	if extOpt.SMaskIndex != nil {
		id += fmt.Sprintf("SMask_%d_0_R;", *extOpt.SMaskIndex)
	}

	return id
}

func GetCachedExtGState(opts ExtGStateOptions, gp *GoPdf) (ExtGState, error) {
	extGState, ok := gp.curr.extGStatesMap.Find(opts)
	if !ok {
		extGState = ExtGState{
			BM:         opts.BlendMode,
			CA:         opts.StrokingCA,
			ca:         opts.NonStrokingCa,
			SMaskIndex: opts.SMaskIndex,
		}

		extGState.Index = gp.addObj(extGState)

		pdfObj := gp.pdfObjs[gp.indexOfProcSet]
		procset, ok := pdfObj.(*ProcSetObj)
		if !ok {
			return ExtGState{}, errors.New("can't convert pdfobject to procsetobj")
		}

		procset.ExtGStates = append(procset.ExtGStates, ExtGS{Index: extGState.Index})

		gp.curr.extGStatesMap.Save(opts.GetId(), extGState)

		//extGState = extGState
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

	if egs.ca != nil {
		content += fmt.Sprintf("\t/ca %.3F\n", *egs.ca)
	}
	if egs.CA != nil {
		content += fmt.Sprintf("\t/CA %.3F\n", *egs.CA)
	}
	if egs.BM != nil {
		content += fmt.Sprintf("\t/BM %s\n", *egs.BM)
	}

	if egs.SMaskIndex != nil {
		content += fmt.Sprintf("\t/SMask %d 0 R\n", *egs.SMaskIndex+1)
	}

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

	return nil
}

type ExtGStatesMap struct {
	syncer sync.Mutex
	table  map[string]ExtGState
}

func NewExtGStatesMap() ExtGStatesMap {
	return ExtGStatesMap{
		syncer: sync.Mutex{},
		table:  make(map[string]ExtGState),
	}
}

func (extm *ExtGStatesMap) Find(extGState ExtGStateOptions) (ExtGState, bool) {
	key := extGState.GetId()

	extm.syncer.Lock()
	defer extm.syncer.Unlock()

	t, ok := extm.table[key]
	if !ok {
		return ExtGState{}, false
	}

	return t, ok

}

func (tm *ExtGStatesMap) Save(id string, extGState ExtGState) ExtGState {
	tm.syncer.Lock()
	defer tm.syncer.Unlock()

	tm.table[id] = extGState

	return extGState
}
