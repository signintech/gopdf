package gopdf

import (
	"bytes"
	"strconv"
)

//DeviceRGBObj  DeviceRGB
type DeviceRGBObj struct {
	buffer bytes.Buffer
	data   []byte
}

func (d *DeviceRGBObj) init(funcGetRoot func() *GoPdf) {

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
	d.buffer.Write(d.data)
	d.buffer.WriteString("endstream\n")

	return nil
}
