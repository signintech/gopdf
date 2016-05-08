package gopdf

import (
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
		err = paesePng(&info, imgConfig)
		if err != nil {
			return info, err
		}
	}

	fmt.Printf("%#v\n", info)

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

func paesePng(info *imgInfo, imgConfig image.Config) error {
	return nil
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
