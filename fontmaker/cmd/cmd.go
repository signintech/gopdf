package main

import (
	"fmt"
	"github.com/signintech/gopdf/fontmaker"
	"runtime/debug"
)

func main() {

	fontpath := "/data/CODES/GOPATH/src/github.com/oneplus1000/gopdfusecase/res/ttf/tahoma.ttf"
	mappath := "/var/www/html/fpdfGo/makefont"
	encoding := "cp874"

	fmk := fontmaker.NewFontMaker()
	_, err := fmk.MakeFont(fontpath, mappath, encoding, "./tmp")
	if err != nil {
		fmt.Printf("Err: %s\n %s\n", err.Error(), string(debug.Stack()))
	}
}
