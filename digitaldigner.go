package gopdf

import "crypto/x509"

//DigitalSignner pdf digital signage
type DigitalSignner interface {
	Cert() *x509.Certificate          //return cert
	Sign(data []byte) ([]byte, error) //sign
}
