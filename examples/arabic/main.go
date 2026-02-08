package main

import (
	"log"

	"github.com/signintech/gopdf"
)

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	err := pdf.AddTTFFont("Amiri", "./examples/arabic/Amiri-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	err = pdf.SetFont("Amiri", "", 24)
	if err != nil {
		log.Fatal(err)
	}

	// Surah Al-Fatiha (The Opening)
	pdf.SetXY(50, 50)
	pdf.Cell(nil, gopdf.ToArabic("بِسْمِ اللهِ الرَّحْمَنِ الرَّحِيمِ"))

	pdf.SetXY(50, 90)
	pdf.Cell(nil, gopdf.ToArabic("الْحَمْدُ لله رَبِّ الْعَالَمِينَ"))

	pdf.SetXY(50, 130)
	pdf.Cell(nil, gopdf.ToArabic("الرَّحْمَنِ الرَّحِيمِ"))

	pdf.SetXY(50, 170)
	pdf.Cell(nil, gopdf.ToArabic("مَالِكِ يَوْمِ الدِّينِ"))

	pdf.SetXY(50, 210)
	pdf.Cell(nil, gopdf.ToArabic("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"))

	pdf.SetXY(50, 250)
	pdf.Cell(nil, gopdf.ToArabic("اهْدِنَا الصِّرَاطَ الْمُسْتَقِيمَ"))

	pdf.SetXY(50, 290)
	pdf.Cell(nil, gopdf.ToArabic("صِرَاطَ الَّذِينَ أَنْعَمْتَ عَلَيْهِمْ"))

	pdf.SetXY(50, 330)
	pdf.Cell(nil, gopdf.ToArabic("غَيْرِ الْمَغْضُوبِ عَلَيْهِمْ وَلَا الضَّالِّينَ"))

	err = pdf.WritePdf("arabic_example.pdf")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PDF created: arabic_example.pdf")
}
