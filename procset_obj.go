package gopdf

import (
	"fmt"
	"io"
)

// ProcSetObj is a PDF procSet object.
type ProcSetObj struct {
	//Font
	Relates             RelateFonts
	RelateXobjs         RelateXobjects
	ExtGStates          []ExtGS
	ImportedTemplateIds map[string]int
	getRoot             func() *GoPdf
}

func (pr *ProcSetObj) init(funcGetRoot func() *GoPdf) {
	pr.getRoot = funcGetRoot
	pr.ImportedTemplateIds = make(map[string]int, 0)
	pr.ExtGStates = make([]ExtGS, 0)
}

func (pr *ProcSetObj) write(w io.Writer, objID int) error {

	io.WriteString(w, "<<\n")
	io.WriteString(w, "/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n")
	io.WriteString(w, "/Font <<\n")
	//me.buffer.WriteString("/F1 9 0 R
	//me.buffer.WriteString("/F2 12 0 R
	//me.buffer.WriteString("/F3 15 0 R
	i := 0
	max := len(pr.Relates)
	for i < max {
		realte := pr.Relates[i]
		fmt.Fprintf(w, "      /F%d %d 0 R\n", realte.CountOfFont+1, realte.IndexOfObj+1)
		i++
	}
	io.WriteString(w, ">>\n")
	io.WriteString(w, "/XObject <<\n")
	i = 0
	max = len(pr.RelateXobjs)
	for i < max {
		fmt.Fprintf(w, "/I%d %d 0 R\n", pr.getRoot().curr.CountOfL+1, pr.RelateXobjs[i].IndexOfObj+1)
		pr.getRoot().curr.CountOfL++
		i++
	}

	// Write imported template name and their ids
	for tplName, objID := range pr.ImportedTemplateIds {
		io.WriteString(w, fmt.Sprintf("%s %d 0 R\n", tplName, objID))
	}

	io.WriteString(w, ">>\n")

	io.WriteString(w, "/ExtGState <<\n")

	for _, extGState := range pr.ExtGStates {
		gsIndex := extGState.Index + 1
		fmt.Fprintf(w, "/GS%d %d 0 R\n", gsIndex, gsIndex)
		pr.getRoot().curr.CountOfL++
	}
	io.WriteString(w, ">>\n")

	io.WriteString(w, ">>\n")
	return nil
}

func (pr *ProcSetObj) getType() string {
	return "ProcSet"
}

// RelateFonts is a slice of RelateFont.
type RelateFonts []RelateFont

// IsContainsFamily checks if font family exists.
func (re *RelateFonts) IsContainsFamily(family string) bool {
	for _, rf := range *re {
		if rf.Family == family {
			return true
		}
	}
	return false
}

// IsContainsFamilyAndStyle checks if font with same name and style already exists .
func (re *RelateFonts) IsContainsFamilyAndStyle(family string, style int) bool {
	for _, rf := range *re {
		if rf.Family == family && rf.Style == style {
			return true
		}
	}
	return false
}

// RelateFont is a metadata index for fonts?
type RelateFont struct {
	Family string
	//etc /F1
	CountOfFont int
	//etc  5 0 R
	IndexOfObj int
	Style      int // Regular|Bold|Italic
}

// RelateXobjects is a slice of RelateXobject.
type RelateXobjects []RelateXobject

// RelateXobject is an index for ???
type RelateXobject struct {
	IndexOfObj int
}

// ExtGS is ???
type ExtGS struct {
	Index int
}
