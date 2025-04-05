package gopdf

import (
	"fmt"
	"io"
)

// ProcSetObj is a PDF procSet object.
type ProcSetObj struct {
	//Font
	Relates             RelateFonts
	RelateColorSpaces   RelateColorSpaces
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
	content := "<<\n"
	content += "\t/ProcSet [/PDF /Text /ImageB /ImageC /ImageI]\n"

	fonts := "\t/Font <<\n"
	for _, relate := range pr.Relates {
		fonts += fmt.Sprintf("\t\t/F%d %d 0 R\n", relate.CountOfFont+1, relate.IndexOfObj+1)
	}
	fonts += "\t>>\n"

	content += fonts

	colorSpaces := "\t/ColorSpace <<\n"
	for _, relate := range pr.RelateColorSpaces {
		colorSpaces += fmt.Sprintf("\t\t/CS%d %d 0 R\n", relate.CountOfColorSpace+1, relate.IndexOfObj+1)
	}
	colorSpaces += "\t>>\n"

	content += colorSpaces

	xobjects := "\t/XObject <<\n"
	for _, XObject := range pr.RelateXobjs {
		xobjects += fmt.Sprintf("\t\t/I%d %d 0 R\n", XObject.IndexOfObj+1, XObject.IndexOfObj+1)
	}
	// Write imported template name and their ids
	for tplName, objID := range pr.ImportedTemplateIds {
		xobjects += fmt.Sprintf("\t\t%s %d 0 R\n", tplName, objID)
	}
	xobjects += "\t>>\n"

	content += xobjects

	extGStates := "\t/ExtGState <<\n"
	for _, extGState := range pr.ExtGStates {
		extGStates += fmt.Sprintf("\t\t/GS%d %d 0 R\n", extGState.Index+1, extGState.Index+1)
	}
	extGStates += "\t>>\n"

	content += extGStates

	content += ">>\n"

	if _, err := io.WriteString(w, content); err != nil {
		return err
	}

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

type RelateColorSpaces []RelateColorSpace

type RelateColorSpace struct {
	Name string
	//etc /CS1
	CountOfColorSpace int
	//etc  5 0 R
	IndexOfObj int
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
