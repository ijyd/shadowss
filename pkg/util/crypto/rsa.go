package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"github.com/golang/glog"
)

func RsaEncrypt(origData []byte, cafile string) ([]byte, error) {
	data, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	glog.V(8).Infoln("...input origData ", string(origData))

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("public key pem Decode error!")
	}

	pubInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	pub := pubInterface.PublicKey.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext []byte, prikeyfile string) ([]byte, error) {
	data, err := ioutil.ReadFile(prikeyfile)
	if err != nil {
		return nil, err
	}
	glog.V(8).Infoln("...input ciphertext ", ciphertext)

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func DecodePassword(passwordB64 string, prikeyfile string) ([]byte, error) {
	if len(passwordB64) == 0 {
		return []byte(""), nil
	}
	/* base64 to RSA */
	pwdRsa, err := base64.StdEncoding.DecodeString(passwordB64)
	if err != nil {
		return nil, err
	}

	/* Rsa to plaintext string */
	pwdPlaintext, err := RsaDecrypt(pwdRsa, prikeyfile)
	if err != nil {
		return nil, err
	}

	return pwdPlaintext, err
}

func EncodePassword(password []byte, cafile string) (string, error) {
	if len(string(password)) == 0 {
		return string(""), nil
	}
	/* plaintext to RSA */
	pwdRsa, err := RsaEncrypt(password, cafile)
	if err != nil {
		return "", err
	}

	/* Rsa to base64 string */
	pwdB64 := base64.StdEncoding.EncodeToString(pwdRsa)
	return pwdB64, err
}
