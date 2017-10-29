package gopdf

import (
	"fmt"
	"io"
)

//ProcSetObj pdf procSet object
type ProcSetObj struct {
	//Font
	Realtes     RelateFonts
	RealteXobjs RealteXobjects
	getRoot     func() *GoPdf
}

func (pr *ProcSetObj) init(funcGetRoot func() *GoPdf) {
	pr.getRoot = funcGetRoot
}

func (pr *ProcSetObj) write(w io.Writer, objID int) error {

	io.WriteString(w, "<<\n")
	io.WriteString(w, "/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n")
	io.WriteString(w, "/Font <<\n")
	//me.buffer.WriteString("/F1 9 0 R
	//me.buffer.WriteString("/F2 12 0 R
	//me.buffer.WriteString("/F3 15 0 R
	i := 0
	max := len(pr.Realtes)
	for i < max {
		realte := pr.Realtes[i]
		fmt.Fprintf(w, "      /F%d %d 0 R\n", realte.CountOfFont+1, realte.IndexOfObj+1)
		i++
	}
	io.WriteString(w, ">>\n")
	io.WriteString(w, "/XObject <<\n")
	i = 0
	max = len(pr.RealteXobjs)
	for i < max {
		fmt.Fprintf(w, "/I%d %d 0 R\n", pr.getRoot().curr.CountOfL+1, pr.RealteXobjs[i].IndexOfObj+1)
		pr.getRoot().curr.CountOfL++
		i++
	}
	io.WriteString(w, ">>\n")
	io.WriteString(w, ">>\n")
	return nil
}

func (pr *ProcSetObj) getType() string {
	return "ProcSet"
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

// IsContainsFamilyAndStyle - checks if already exists font with same name and style
func (re *RelateFonts) IsContainsFamilyAndStyle(family string, style int) bool {
	i := 0
	max := len(*re)
	for i < max {
		if (*re)[i].Family == family && (*re)[i].Style == style  {
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
	Style      int // Regular|Bold|Italic
}

type RealteXobjects []RealteXobject

type RealteXobject struct {
	IndexOfObj int
}
