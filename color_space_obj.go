package gopdf

import (
	"fmt"
	"io"
)

type ColorSpaceObj struct {
	CountOfSpaceColor int
	Name              string
	spaceName         string
	colorString0      string
	colorString1      string
	space             string
}

func (cs *ColorSpaceObj) init(func() *GoPdf) {}

func (cs *ColorSpaceObj) getType() string {
	return "ColorSpace"
}

func (cs *ColorSpaceObj) write(w io.Writer, objID int) error {

	fmt.Fprintf(w, "[ /Separation %s %s\n", cs.spaceName, cs.space)
	io.WriteString(w, "		<< \n")
	io.WriteString(w, "			/FunctionType 2\n")
	io.WriteString(w, "			/Domain [0 1]\n")
	fmt.Fprintf(w, "			/C0 [%s]\n", cs.colorString0)
	fmt.Fprintf(w, "			/C1 [%s]\n", cs.colorString1)
	io.WriteString(w, "			/N 1\n")
	io.WriteString(w, "		>>\n")
	io.WriteString(w, "]\n")

	return nil
}

func (cs *ColorSpaceObj) SetColorRBG(r, g, b uint8) {
	cs.colorString0 = "0.0 0.0 0.0"
	cs.colorString1 = fmt.Sprintf("%.3f %.3f %.3f", float64(r)/255.0, float64(g)/255.0, float64(b)/255.0)
	cs.space = "/DeviceRGB"
	cs.spaceName = fmt.Sprintf("/%s", cs.Name)
}

func (cs *ColorSpaceObj) SetColorCMYK(c, m, y, k uint8) {
	cs.colorString0 = "0.0 0.0 0.0 0.0"
	cs.colorString1 = fmt.Sprintf("%.3f %.3f %.3f %.3f", float64(c)/100.0, float64(m)/100.0, float64(y)/100.0, float64(k)/100.0)
	cs.space = "/DeviceCMYK"
	cs.spaceName = fmt.Sprintf("/%s", cs.Name)
}
