package gopdf

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

// pngChunk builds a PNG chunk with a valid CRC.
func pngChunk(typ string, data []byte) []byte {
	var b bytes.Buffer
	if err := binary.Write(&b, binary.BigEndian, uint32(len(data))); err != nil {
		panic(err)
	}
	b.WriteString(typ)
	b.Write(data)
	crc := crc32.NewIEEE()
	crc.Write([]byte(typ))
	crc.Write(data)
	if err := binary.Write(&b, binary.BigEndian, crc.Sum32()); err != nil {
		panic(err)
	}
	return b.Bytes()
}

// pngWithTrns encodes img as PNG and inserts a tRNS chunk right after IHDR.
func pngWithTrns(t *testing.T, img image.Image, trns []byte) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := png.Encode(&b, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	raw := b.Bytes()
	// A PNG signature is 8 bytes and the IHDR chunk is always first and 25 bytes long.
	const sigAndIHDR = 8 + 25
	var out bytes.Buffer
	out.Write(raw[:sigAndIHDR])
	out.Write(pngChunk("tRNS", trns))
	out.Write(raw[sigAndIHDR:])
	return out.Bytes()
}

func TestParsePngShortTrnsReturnsError(t *testing.T) {
	gray := image.NewGray(image.Rect(0, 0, 1, 1))
	gray.SetGray(0, 0, color.Gray{Y: 7})
	rgba := image.NewRGBA(image.Rect(0, 0, 1, 1))
	rgba.SetRGBA(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})

	cases := []struct {
		name string
		img  image.Image
		trns []byte
	}{
		{"grayscale truncated", gray, []byte{0x00}},
		{"grayscale empty", gray, nil},
		{"truecolor truncated", rgba, []byte{0x00, 0x01}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := pngWithTrns(t, tc.img, tc.trns)
			if _, err := parseImg(bytes.NewReader(raw)); err == nil {
				t.Fatal("expected an error for a truncated tRNS chunk, got nil")
			}
		})
	}
}

func TestParsePngTrnsColorKey(t *testing.T) {
	gray := image.NewGray(image.Rect(0, 0, 1, 1))
	gray.SetGray(0, 0, color.Gray{Y: 7})
	rgba := image.NewRGBA(image.Rect(0, 0, 1, 1))
	rgba.SetRGBA(0, 0, color.RGBA{R: 1, G: 2, B: 3, A: 255})

	cases := []struct {
		name string
		img  image.Image
		trns []byte
		want []byte
	}{
		{"grayscale", gray, []byte{0x00, 0x2a}, []byte{0x2a}},
		{"truecolor", rgba, []byte{0x00, 0x11, 0x00, 0x22, 0x00, 0x33}, []byte{0x11, 0x22, 0x33}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw := pngWithTrns(t, tc.img, tc.trns)
			info, err := parseImg(bytes.NewReader(raw))
			if err != nil {
				t.Fatalf("parseImg: %v", err)
			}
			if !bytes.Equal(info.trns, tc.want) {
				t.Fatalf("trns = %v, want %v", info.trns, tc.want)
			}
		})
	}
}

func TestImageTrnsFiles(t *testing.T) {
	badFilePaths := []string{
		"test/res/bad_trns_grayscale_truncated.png",
		"test/res/bad_trns_grayscale_empty.png",
		"test/res/bad_trns_truecolor_truncated.png",
	}

	goodFilePaths := []string{
		"test/res/good_trns_grayscale_valid.png",
	}
	for _, path := range badFilePaths {
		t.Run(path, func(t *testing.T) {
			var pdf GoPdf
			pdf.Start(Config{PageSize: *PageSizeA4})
			pdf.AddPage()
			err := pdf.Image(path, 0, 0, &Rect{W: 10, H: 10})
			if err == nil {
				t.Fatalf("expected error for %s, got nil", path)
			}
			if !strings.Contains(err.Error(), "invalid tRNS chunk length") {
				t.Fatalf("expected invalid tRNS chunk length error for %s, got: %v", path, err)
			}
		})
	}

	for _, path := range goodFilePaths {
		t.Run(path, func(t *testing.T) {
			var pdf GoPdf
			pdf.Start(Config{PageSize: *PageSizeA4})
			pdf.AddPage()
			if err := pdf.Image(path, 0, 0, &Rect{W: 10, H: 10}); err != nil {
				t.Fatalf("unexpected error for %s: %v", path, err)
			}
		})
	}
}
