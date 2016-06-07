package gopdf

import "bytes"

//SMask smask
type SMask struct {
	w, h             int
	colspace         string
	bitsPerComponent string
	filter           string
	decodeParms      string
	data             []byte
}

func (s *SMask) init(funcGetRoot func() *GoPdf) {

}

func (s *SMask) getType() string {
	return "smask"
}
func (s *SMask) getObjBuff() *bytes.Buffer {
	var buff bytes.Buffer
	return &buff
}

//สร้าง ข้อมูลใน pdf
func (s *SMask) build() error {
	return nil
}
