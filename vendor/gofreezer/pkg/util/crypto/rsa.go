package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

type RSAHelper struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewRSAHelper(privateKey, publicKey []byte) *RSAHelper {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil
	}
	pubInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil
	}
	pub := pubInterface.PublicKey.(*rsa.PublicKey)

	block, _ = pem.Decode(privateKey)
	if block == nil {
		return nil
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}

	return &RSAHelper{
		privateKey: priv,
		publicKey:  pub,
	}
}

func (rsaHelper *RSAHelper) RsaEncrypt(origData []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, rsaHelper.publicKey, origData)
}

func (rsaHelper *RSAHelper) RsaDecrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, rsaHelper.privateKey, ciphertext)
}
