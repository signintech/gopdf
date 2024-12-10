package gopdf

import (
	"testing"
)

func TestImagePares(t *testing.T) {

	var err error
	_, err = parseImgByPath("test/res/gopher01.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImgByPath("test/res/gopher01_g_mode.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImgByPath("test/res/gopher01_i_mode.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	//Channel_digital_image_CMYK_color.jpg
	_, err = parseImgByPath("test/res/Channel_digital_image_CMYK_color.jpg")
	if err != nil {
		t.Error(err)
		//return
	}

	_, err = parseImgByPath("test/res/gopher02.png")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = parseImgByPath("test/res/gopher02.png")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = parseImgByPath("test/res/gopher03.gif")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = parseImgByPath("test/res/gopher03_color.gif")
	if err != nil {
		t.Error(err)
		return
	}

}

func TestImage02Pares(t *testing.T) {

	var err error
	//_, err = parseImgByPath("test/res/OpenOffice.org_1.1_official_main_logo_2col_trans.png")

	_, err = parseImgByPath("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		//return
	}
}

func TestImage03Pares(t *testing.T) {

	var err error
	_, err = parseImgByPath("test/res/OpenOffice.org_1.1_official_main_logo_2col_trans.png")

	//_, err = parseImgByPath("./test/res/PNG_transparency_demonstration_1.png")
	if err != nil {
		t.Error(err)
		//return
	}
}
