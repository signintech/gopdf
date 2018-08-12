package gopdf

import "io"

//acroForm now support only digital sign
type acroForm struct {
	signner DigitalSignner
}

func (a *acroForm) init(fn func() *GoPdf) {

}
func (a *acroForm) getType() string {
	return "AcroForm"
}
func (a *acroForm) write(w io.Writer, objID int) error {
	return nil
}

func (a *acroForm) setDigitalSign(signner DigitalSignner) {
	a.signner = signner
	
}

/*
func createSignAnnot(){
	annotType := initNodeKeyUseNameAndNodeContentUseString("Type", "/Annot")
	annotP := initNodeKeyUseNameAndNodeContentUseRefTo("P", refPage.content.refTo)
	annotSubtype := initNodeKeyUseNameAndNodeContentUseString("Subtype", "/Widget")
	annotV := initNodeKeyUseNameAndNodeContentUseRefTo("V", refSigID)
	annotT := initNodeKeyUseNameAndNodeContentUseString("T", "(Signature1)")
	annotFT := initNodeKeyUseNameAndNodeContentUseString("FT", "/Sig")
	annotF := initNodeKeyUseNameAndNodeContentUseString("F", "4")
	annotDA := initNodeKeyUseNameAndNodeContentUseString("DA", "(/Helvetica 0 Tf 0 g)")
}*/
