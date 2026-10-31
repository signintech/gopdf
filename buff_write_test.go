package gopdf

import (
	"bytes"
	"testing"
	"time"
)

func TestWriteUInt32(t *testing.T) {
	var buff bytes.Buffer
	err := WriteUInt32(&buff, 65536)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	b := buff.Bytes()
	if b[0] != 0 || b[1] != 1 || b[2] != 0 || b[3] != 0 {
		t.Errorf("WriteUInt64 fail")
	}
}

/*
func TestFloat(t *testing.T){
	a := float64(1170.08)
    b := float64(1013.08)
    c := float64(54673.00)
    d := float64(131588.00)
    e := float64(54236.52)
	sum := a + b + c + d + e
	if sum != 242680.68 {
		t.Errorf("!242680.68")
	}
}*/

func TestEncodeUTF8(t *testing.T) {
	str := "Boonchai Manasirisuk"
	buff := encodeUtf8(str)
	if buff != "0042006F006F006E00630068006100690020004D0061006E0061007300690072006900730075006B" {
		t.Error("not match")
	}
}

func TestInfoDate(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		want   string
	}{
		{"UTC", 0, "20261029120000+00'00'"},
		{"whole hours east", 3 * 3600, "20261029120000+03'00'"},
		{"whole hours west", -5 * 3600, "20261029120000-05'00'"},
		{"half hour east", 5*3600 + 30*60, "20261029120000+05'30'"},
		{"quarter hour east", 5*3600 + 45*60, "20261029120000+05'45'"},
		{"half hour west", -(3*3600 + 30*60), "20261029120000-03'30'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			date := time.Date(2026, 10, 29, 12, 0, 0, 0, time.FixedZone(tt.name, tt.offset))
			if got := infodate(date); got != tt.want {
				t.Errorf("infodate() = %q, want %q", got, tt.want)
			}
		})
	}
}
