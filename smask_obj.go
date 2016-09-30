package gopdf

import (
	"bytes"
	"fmt"
)

//SMask smask
type SMask struct {
	buffer bytes.Buffer
	imgInfo
	data []byte
}

func (s *SMask) init(funcGetRoot func() *GoPdf) {

}

func (s *SMask) getType() string {
	return "smask"
}
func (s *SMask) getObjBuff() *bytes.Buffer {
	return &s.buffer
}

//สร้าง ข้อมูลใน pdf
func (s *SMask) build(objID int) error {

	buff, err := buildImgProp(s.imgInfo)
	if err != nil {
		return err
	}
	_, err = buff.WriteTo(&s.buffer)
	if err != nil {
		return err
	}

	s.buffer.WriteString(fmt.Sprintf("/Length %d\n>>\n", len(s.data))) // /Length 62303>>\n
	s.buffer.WriteString("stream\n")
	s.buffer.Write(s.data)
	s.buffer.WriteString("\nendstream\n")

	return nil
}
