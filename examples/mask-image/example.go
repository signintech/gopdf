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

	//image bytes
	b, err := os.ReadFile(resourcesPath + "/gopher01.jpg")
	if err != nil {
		log.Panic(err.Error())
	}

	imgH1, err := gopdf.ImageHolderByBytes(b)
	if err != nil {
		log.Panic(err.Error())
	}
	if err := pdf.ImageByHolder(imgH1, 200, 250, nil); err != nil {
		log.Panic(err.Error())
	}

	//image io.Reader
	file, err := os.Open(resourcesPath + "/chilli.jpg")
	if err != nil {
		log.Panic(err.Error())
	}

	imgH2, err := gopdf.ImageHolderByReader(file)
	if err != nil {
		log.Panic(err.Error())
	}

	maskHolder, err := gopdf.ImageHolderByPath(resourcesPath + "/mask.png")
	if err != nil {
		log.Panic(err.Error())
	}

	maskOpts := gopdf.MaskOptions{
		Holder: maskHolder,
		ImageOptions: gopdf.ImageOptions{
			X: 0,
			Y: 0,
			Rect: &gopdf.Rect{
				W: 300,
				H: 300,
			},
		},
	}

	transparency, err := gopdf.NewTransparency(0.5, "")
	if err != nil {
		log.Panic(err.Error())
	}

	imOpts := gopdf.ImageOptions{
		X:            0,
		Y:            0,
		Mask:         &maskOpts,
		Transparency: &transparency,
		Rect:         &gopdf.Rect{W: 400, H: 400},
	}
	if err := pdf.ImageByHolderWithOptions(imgH2, imOpts); err != nil {
		log.Panic(err.Error())
	}

	pdf.SetCompressLevel(0)
	if err := pdf.WritePdf("image.pdf"); err != nil {
		log.Panic(err.Error())
	}
}
