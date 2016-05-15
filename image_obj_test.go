package gopdf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"testing"
)

func TestImagePares(t *testing.T) {
	var err error

	_, err = parseImg("test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImg("test/res/gopher01_g_mode.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImg("test/res/gopher01_i_mode.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	//Channel_digital_image_CMYK_color.jpg
	_, err = parseImg("test/res/Channel_digital_image_CMYK_color.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImg("test/res/gopher02.png")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = parseImg("test/res/gopher02_g_mode.png")
	if err != nil {
		t.Error(err)
		return
	}

}

type imgInfo struct {
	src              string
	formatName       string
	colspace         string
	bitsPerComponent string
	filter           string
}

func parseImg(src string) (imgInfo, error) {
	var info imgInfo
	info.src = src
	file, err := os.Open(src)
	if err != nil {
		return info, err
	}
	defer file.Close()

	imgConfig, formatname, err := image.DecodeConfig(file)
	if err != nil {
		return info, err
	}
	info.formatName = formatname
	if formatname == "jpeg" {
		err = parseImgJpg(&info, imgConfig)
		if err != nil {
			return info, err
		}
	} else if formatname == "png" {

		err = paesePng(file, &info, imgConfig)
		if err != nil {
			return info, err
		}
	}

	//fmt.Printf("%#v\n", info)

	return info, nil
}

func parseImgJpg(info *imgInfo, imgConfig image.Config) error {
	if imgConfig.ColorModel == color.YCbCrModel {
		info.colspace = "DeviceRGB"
	} else if imgConfig.ColorModel == color.GrayModel {
		info.colspace = "DeviceGray"
	} else if imgConfig.ColorModel == color.CMYKModel {
		info.colspace = "DeviceCMYK"
	} else {
		return errors.New("color model not support")
	}
	info.bitsPerComponent = "8"
	info.filter = "DCTDecode"
	return nil
}

var pngMagicNumber = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
var pngIHDR = []byte{0x49, 0x48, 0x44, 0x52}

func paesePng(f *os.File, info *imgInfo, imgConfig image.Config) error {

	f.Seek(0, 0)
	b, err := readBytes(f, 8)
	if err != nil {
		return err
	}
	if !compareBytes(b, pngMagicNumber) {
		return errors.New("Not a PNG file")
	}

	f.Seek(4, 1) //skip header chunk
	b, err = readBytes(f, 4)
	if err != nil {
		return err
	}
	if !compareBytes(b, pngIHDR) {
		return errors.New("Incorrect PNG file")
	}

	w, err := readInt(f)
	if err != nil {
		return err
	}
	h, err := readInt(f)
	if err != nil {
		return err
	}
	fmt.Printf("w=%d h=%d\n", w, h)

	bpc, err := readBytes(f, 1)
	if err != nil {
		return err
	}

	if bpc[0] > 8 {
		return errors.New("16-bit depth not supported")
	}

	ct, err := readBytes(f, 1)
	if err != nil {
		return err
	}

	if ct[0] == 0 || ct[0] == 4 {
		info.colspace = "DeviceGray"
	} else if ct[0] == 2 || ct[0] == 6 {
		info.colspace = "DeviceRGB"
	} else if ct[0] == 3 {
		info.colspace = "Indexed"
	} else {
		return errors.New("Unknown color type")
	}

	return nil
}

func readUInt(f *os.File) (uint, error) {
	buff, err := readBytes(f, 4)
	fmt.Printf("%#v\n\n", buff)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint32(buff)
	return uint(n), nil
}

func readInt(f *os.File) (int, error) {

	u, err := readUInt(f)
	if err != nil {
		return 0, err
	}
	var v int
	if u >= 0x8000 {
		v = int(u) - 65536
	} else {
		v = int(u)
	}
	return v, nil
}

func readBytes(f *os.File, len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := f.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func compareBytes(a []byte, b []byte) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	i := 0
	max := len(a)
	for i < max {
		if a[i] != b[i] {
			return false
		}
		i++
	}
	return true
}

func isDeviceRGB(formatname string, img *image.Image) bool {
	if _, ok := (*img).(*image.YCbCr); ok {
		return true
	} else if _, ok := (*img).(*image.NRGBA); ok {
		return true
	}
	return false
}

func ImgReactagleToWH(imageRect image.Rectangle) (float64, float64) {
	k := 1
	w := -128 //init
	h := -128 //init
	if w < 0 {
		w = -imageRect.Dx() * 72 / w / k
	}
	if h < 0 {
		h = -imageRect.Dy() * 72 / h / k
	}
	if w == 0 {
		w = h * imageRect.Dx() / imageRect.Dy()
	}
	if h == 0 {
		h = w * imageRect.Dy() / imageRect.Dx()
	}
	return float64(w), float64(h)
}
