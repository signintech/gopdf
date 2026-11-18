package gopdf

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"testing"
	"time"
	"unicode/utf16"
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

// TestEncodeUTF8NonBMP checks that characters outside the BMP (like emoji)
// are encoded as valid UTF-16 surrogate pairs, not broken 5-digit hex.
func TestEncodeUTF8NonBMP(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"emoji", "Report 😀", "005200650070006F007200740020D83DDE00"},
		{"cjk ext b", "𠀀", "D840DC00"},
		{"mixed bmp and emoji", "Aไทย😀", "00410E440E170E22D83DDE00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeUtf8(tt.input)
			if got != tt.want {
				t.Fatalf("encodeUtf8(%q) = %q, want %q", tt.input, got, tt.want)
			}
			raw, err := hex.DecodeString(got)
			if err != nil {
				t.Fatalf("hex.DecodeString(%q) failed: %v", got, err)
			}
			if len(raw)%2 != 0 {
				t.Fatalf("decoded UTF-16BE byte length %d is odd", len(raw))
			}
			u16 := make([]uint16, len(raw)/2)
			for i := range u16 {
				u16[i] = binary.BigEndian.Uint16(raw[i*2 : i*2+2])
			}
			if decoded := string(utf16.Decode(u16)); decoded != tt.input {
				t.Fatalf("round-trip decoded %q, want %q", decoded, tt.input)
			}
		})
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
