package gopdf

import (
	"crypto/x509"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pkcs12"
)

//DigitalSignPkcs12 ข้อมูลการ sign ทั้งหมด
type DigitalSignPkcs12 struct {
	privateKey  interface{}
	certificate *x509.Certificate
}

//Cert get cert information
func (d *DigitalSignPkcs12) Cert() *x509.Certificate {
	return d.certificate
}

//Sign sign
func (d *DigitalSignPkcs12) Sign(data []byte) ([]byte, error) {
	s, err := newSignerFromKey(d.privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return s.Sign(data)
}

//SetCertFile set pkcs12 cert file (PKCS #12 files is ".p12" or ".pfx")
func (d *DigitalSignPkcs12) SetCertFile(file []byte, password string) error {
	privateKey, certificate, err := pkcs12.Decode(file, password)
	if err != nil {
		return errors.Wrapf(err, "pkcs12.Decode(...) fail")
	}
	d.privateKey = privateKey
	d.certificate = certificate
	return nil
}
