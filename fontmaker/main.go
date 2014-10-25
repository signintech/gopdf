package main

import (
	"bytes"
	"fmt"
	"github.com/signintech/gopdf/fontmaker/core"
	"os"
	//"runtime/debug"
)

func main() {

	lenarg := len(os.Args)
	if lenarg < 5 {
		echoUsage()
		return
	}
	/*
		fontpath := "/var/www/html/fpdfGo/ttf/tahoma.ttf"
		mappath := "/var/www/html/fpdfGo/makefont"
		encoding := "cp874"
	*/
	i := 1
	encoding := os.Args[i+0]
	mappath := os.Args[i+1]
	fontpath := os.Args[i+2]
	outputpath := os.Args[i+3]

	fmk := core.NewFontMaker()
	err := fmk.MakeFont(fontpath, mappath, encoding, outputpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nERROR: %s\n\n", err.Error())
		echoUsage()
		return
	}
	fmt.Printf("Finish.\n")
}

func echoUsage() {
	var buff bytes.Buffer
	//buff.WriteString("\n")
	buff.WriteString("fontmaker is tool for making font file to use with gopdf.\n")
	buff.WriteString("\nUsage:\n")
	buff.WriteString("\tfontmaker encoding map_folder font_file output_folder\n")
	buff.WriteString("\nExample:\n")
	buff.WriteString("\tfontmaker cp874 ../map  ../ttf/Loma.ttf ./tmp\n")
	buff.WriteString("\n")
	fmt.Print(buff.String())
}
