package gopdf

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"math/rand"
	"time"
)

const (
	//PermissionsPrint setProtection print
	PermissionsPrint = 4
	//PermissionsModify setProtection modify
	PermissionsModify = 8
	//PermissionsCopy setProtection copy
	PermissionsCopy = 16
	//PermissionsAnnotForms setProtection  annot-forms
	PermissionsAnnotForms = 32
)

var protectionPadding = []byte{
	0x28, 0xBF, 0x4E, 0x5E, 0x4E, 0x75, 0x8A, 0x41, 0x64, 0x00, 0x4E, 0x56, 0xFF, 0xFA, 0x01, 0x08,
	0x2E, 0x2E, 0x00, 0xB6, 0xD0, 0x68, 0x3E, 0x80, 0x2F, 0x0C, 0xA9, 0xFE, 0x64, 0x53, 0x69, 0x7A,
}

//PDFProtection protection in pdf
type PDFProtection struct {
	encrypted bool   //whether document is protected
	uValue    []byte //U entry in pdf document
	oValue    []byte //O entry in pdf document
	pValue    int    //P entry in pdf document
	//var $enc_obj_id;         //encryption object id

}

func (p *PDFProtection) setProtection(permissions int, userPass []byte, ownerPass []byte) error {
	protection := 192 | permissions
	if ownerPass == nil || len(ownerPass) == 0 {
		ownerPass = p.randomPass(24)
	}
	return p.generateencryptionkey(userPass, ownerPass, protection)
}

func (p *PDFProtection) generateencryptionkey(userPass []byte, ownerPass []byte, protection int) error {

	userPass = append(userPass, protectionPadding...)
	userPassWithPadding := userPass[0:32]
	ownerPass = append(ownerPass, protectionPadding...)
	ownerPassWithPadding := ownerPass[0:32]

	//oValue
	oValue, err := p.createOValue(userPassWithPadding, ownerPassWithPadding)
	if err != nil {
		return err
	}
	p.oValue = oValue

	uValue, err := p.createUValue(userPassWithPadding, oValue, protection)
	if err != nil {
		return err
	}
	p.uValue = uValue

	//fmt.Printf("%#v\n", uValue)

	return nil
}

func (p *PDFProtection) createOValue(userPassWithPadding []byte, ownerPassWithPadding []byte) ([]byte, error) {
	tmp := md5.Sum(ownerPassWithPadding)
	ownerRC4key := tmp[0:5]
	cip, err := rc4.NewCipher(ownerRC4key)
	if err != nil {
		return nil, err
	}
	dest := make([]byte, len(userPassWithPadding))
	cip.XORKeyStream(dest, userPassWithPadding)
	return dest, nil
}

func (p *PDFProtection) createUValue(userPassWithPadding []byte, oValue []byte, protection int) ([]byte, error) {

	var tmp bytes.Buffer
	tmp.Write(userPassWithPadding)
	tmp.Write(oValue)
	tmp.WriteByte(byte(protection))
	tmp.WriteByte(byte(0xff))
	tmp.WriteByte(byte(0xff))
	tmp.WriteByte(byte(0xff))

	encryptionKey := md5.Sum(tmp.Bytes())
	cip, err := rc4.NewCipher(encryptionKey[0:5])
	if err != nil {
		return nil, err
	}
	dest := make([]byte, len(protectionPadding))
	cip.XORKeyStream(dest, protectionPadding)
	return dest, nil
}

func (p *PDFProtection) randomPass(strlen int) []byte {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdef0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return result
}
