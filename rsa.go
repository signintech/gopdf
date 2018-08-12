package gopdf

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

//signer return sign data
type signer interface {
	// Sign returns raw signature for the given data. This method
	// will apply the hash specified for the keytype to the data.
	Sign(data []byte) ([]byte, error)
}

func newSignerFromKey(k interface{}) (signer, error) {
	var sshKey signer
	switch t := k.(type) {
	case *rsa.PrivateKey:
		sshKey = &rsaPrivateKey{t}
	default:
		return nil, fmt.Errorf("ssh: unsupported key type %T", k)
	}
	return sshKey, nil
}

type rsaPrivateKey struct {
	*rsa.PrivateKey
}

// Sign signs data with rsa-sha256
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)

	//d2 := sha256.Sum224(data)
	//fmt.Printf("%x<<<<\n%xxxxx\n", d, d2)
	return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
}
