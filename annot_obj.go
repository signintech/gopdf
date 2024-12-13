package gopdf

import (
	"fmt"
	"io"
	"strings"
)

type annotObj struct {
	linkOption
	GetRoot func() *GoPdf
}

func (o annotObj) init(f func() *GoPdf) {
}

func (o annotObj) getType() string {
	return "Annot"
}

func (o annotObj) write(w io.Writer, objID int) error {
	if o.url != "" {
		return o.writeExternalLink(w, o.linkOption, objID)
	} else {
		return o.writeInternalLink(w, o.linkOption)
	}
}

func (o annotObj) writeExternalLink(w io.Writer, l linkOption, objID int) error {
	protection := o.GetRoot().protection()
	url := l.url
	if protection != nil {
		tmp, err := rc4Cip(protection.objectkey(objID), []byte(url))
		if err != nil {
			return err
		}
		url = string(tmp)
	}
	url = strings.Replace(url, "\\", "\\\\", -1)
	url = strings.Replace(url, "(", "\\(", -1)
	url = strings.Replace(url, ")", "\\)", -1)
	url = strings.Replace(url, "\r", "\\r", -1)

	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /A <</S /URI /URI (%s)>>>>\n",
		l.x, l.y, l.x+l.w, l.y-l.h, url)
	return err
}

func (o annotObj) writeInternalLink(w io.Writer, l linkOption) error {
	a, ok := o.GetRoot().anchors[l.anchor]
	if !ok {
		return nil
	}
	_, err := fmt.Fprintf(w, "<</Type /Annot /Subtype /Link /Rect [%.2f %.2f %.2f %.2f] /Border [0 0 0] /Dest [%d 0 R /XYZ 0 %.2f null]>>\n",
		l.x, l.y, l.x+l.w, l.y-l.h, a.page+1, a.y)
	return err
}
