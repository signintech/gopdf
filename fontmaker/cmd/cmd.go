package main

import (
	"fmt"
	"github.com/signintech/gopdf/fontmaker"
)

func main() {

	fontpath := "/data/CODES/GOPATH/src/github.com/oneplus1000/gopdfusecase/res/ttf/tahoma.ttf"
	encodingpath := "/var/www/html/fpdfGo/makefont/cp874.map"

	fmk := fontmaker.NewFontMaker()
	_, err := fmk.MakeFont(fontpath, encodingpath, "./tmp")
	if err != nil {
		fmt.Printf("Err: %s\n", err.Error())
	}
}
