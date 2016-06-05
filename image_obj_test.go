package gopdf

import (
	"io/ioutil"
	"testing"
)

func TestImagePares(t *testing.T) {
	var err error
	/*
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
	*/
	//data, err := ioutil.ReadFile("test/res/OpenOffice.org_1.1_official_main_logo_2col_trans.png")
	data, err := ioutil.ReadFile("test/res/gopher02.png")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = parseImg(data)
	if err != nil {
		t.Error(err)
		return
	}

}
