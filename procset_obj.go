package gopdf

import (
	"bytes"
	"fmt"
)

//ProcSetObj pdf procSet object
type ProcSetObj struct {
	buffer bytes.Buffer
	//Font
	Realtes     RelateFonts
	RealteXobjs RealteXobjects
	getRoot     func() *GoPdf
}

func (pr *ProcSetObj) init(funcGetRoot func() *GoPdf) {
	pr.getRoot = funcGetRoot
}

func (pr *ProcSetObj) build(objID int) error {

	pr.buffer.WriteString("<<\n")
	pr.buffer.WriteString("/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n")
	pr.buffer.WriteString("/Font <<\n")
	//me.buffer.WriteString("/F1 9 0 R
	//me.buffer.WriteString("/F2 12 0 R
	//me.buffer.WriteString("/F3 15 0 R
	i := 0
	max := len(pr.Realtes)
	for i < max {
		realte := pr.Realtes[i]
		pr.buffer.WriteString(fmt.Sprintf("      /F%d %d 0 R\n", realte.CountOfFont+1, realte.IndexOfObj+1))
		i++
	}
	pr.buffer.WriteString(">>\n")
	pr.buffer.WriteString("/XObject <<\n")
	i = 0
	max = len(pr.RealteXobjs)
	for i < max {
		pr.buffer.WriteString(fmt.Sprintf("/I%d %d 0 R\n", pr.getRoot().curr.CountOfL+1, pr.RealteXobjs[i].IndexOfObj+1))
		pr.getRoot().curr.CountOfL++
		i++
	}
	pr.buffer.WriteString(">>\n")
	pr.buffer.WriteString(">>\n")
	return nil
}

func (pr *ProcSetObj) getType() string {
	return "ProcSet"
}

func (pr *ProcSetObj) getObjBuff() *bytes.Buffer {
	return &(pr.buffer)
}

type RelateFonts []RelateFont

func (re *RelateFonts) IsContainsFamily(family string) bool {
	i := 0
	max := len(*re)
	for i < max {
		if (*re)[i].Family == family {
			return true
		}
		i++
	}
	return false
}

type RelateFont struct {
	Family string
	//etc /F1
	CountOfFont int
	//etc  5 0 R
	IndexOfObj int
}

type RealteXobjects []RealteXobject

type RealteXobject struct {
	IndexOfObj int
}
