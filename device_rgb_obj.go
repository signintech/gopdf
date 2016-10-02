package gopdf

import (
	"bytes"
	"strconv"
)

//DeviceRGBObj  DeviceRGB
type DeviceRGBObj struct {
	buffer  bytes.Buffer
	data    []byte
	getRoot func() *GoPdf
}

func (d *DeviceRGBObj) init(funcGetRoot func() *GoPdf) {
	d.getRoot = funcGetRoot
}

func (d *DeviceRGBObj) protection() *PDFProtection {
	return d.getRoot().protection()
}

func (d *DeviceRGBObj) getType() string {
	return "devicergb"
}
func (d *DeviceRGBObj) getObjBuff() *bytes.Buffer {
	return &d.buffer
}

//สร้าง ข้อมูลใน pdf
func (d *DeviceRGBObj) build(objID int) error {

	d.buffer.WriteString("<<\n")
	d.buffer.WriteString("/Length " + strconv.Itoa(len(d.data)) + "\n")
	d.buffer.WriteString(">>\n")
	d.buffer.WriteString("stream\n")
	if d.protection() != nil {
		tmp, err := rc4Cip(d.protection().objectkey(objID), d.data)
		if err != nil {
			return err
		}
		d.buffer.Write(tmp)
		d.buffer.WriteString("\n")
	} else {
		d.buffer.Write(d.data)
	}
	d.buffer.WriteString("endstream\n")

	return nil
}
