package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/signintech/gopdf"
)

const resourcesPath = "../../test/res"

func main() {
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}) //595.28, 841.89 = A4
	pdf.AddPage()

	if err := pdf.AddTTFFont("loma", resourcesPath+"/LiberationSerif-Regular.ttf"); err != nil {
		log.Panic(err.Error())
	}

	if err := pdf.SetFont("loma", "", 14); err != nil {
		log.Panic(err.Error())
	}

	//image bytes
	b, err := ioutil.ReadFile(resourcesPath + "/gopher01.jpg")
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

	transparency, err := gopdf.NewTransparency(0.5, "")
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
			X: 100,
			Y: 450,
			Rect: &gopdf.Rect{
				W: 300,
				H: 300,
			},
			Transparency: &transparency,
		},
	}

	imOpts := gopdf.ImageOptions{
		X:    100,
		Y:    450,
		Mask: &maskOpts,
		Rect: &gopdf.Rect{W: 400, H: 400},
	}
	if err := pdf.ImageByHolderWithOptions(imgH2, imOpts); err != nil {
		log.Panic(err.Error())
	}

	pdf.SetCompressLevel(0)
	if err := pdf.WritePdf("mask-image.pdf"); err != nil {
		log.Panic(err.Error())
	}
}
