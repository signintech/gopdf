package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
)

var resourcesPath string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	resourcesPath = filepath.Join(cwd, "test/res")
}

func main() {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	if err := pdf.AddTTFFont("loma", resourcesPath+"/LiberationSerif-Regular.ttf"); err != nil {
		log.Panic(err.Error())
	}

	if err := pdf.SetFont("loma", "", 14); err != nil {
		log.Panic(err.Error())
	}

	file, err := os.Open(resourcesPath + "/chilli.jpg")
	if err != nil {
		log.Panic(err.Error())
	}

	chili, err := gopdf.ImageHolderByReader(file)
	if err != nil {
		log.Panic(err.Error())
	}

	if err := pdf.ImageByHolderWithOptions(chili, gopdf.ImageOptions{
		X:    100,
		Y:    100,
		Rect: &gopdf.Rect{W: 200, H: 200},
	}); err != nil {
		log.Panic(err.Error())
	}

	//When the image is rotated 90 degrees, the ratio of the image changes.
	if err := pdf.ImageByHolderWithOptions(chili, gopdf.ImageOptions{
		DegreeAngle: 90,
		X:           200,
		Y:           400,
		Rect:        &gopdf.Rect{W: 200, H: 200},
	}); err != nil {
		log.Panic(err.Error())
	}

	pdf.SetCompressLevel(0)
	if err := pdf.WritePdf("image.pdf"); err != nil {
		log.Panic(err.Error())
	}

}
