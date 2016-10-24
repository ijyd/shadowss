package permissions

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"golib/pkg/util/crypto"
)

type Permissions struct {
	featureCode string
}

func NewPermissions(featureCode, salt string) *Permissions {
	buffer := bytes.Buffer{}
	buffer.WriteString(featureCode)
	buffer.WriteString(salt)
	h := md5.New()
	h.Write(buffer.Bytes())

	return &Permissions{
		featureCode: hex.EncodeToString(h.Sum(nil)),
	}
}

func NewPermissionsWithCode(featureCode string) *Permissions {
	return &Permissions{
		featureCode: featureCode,
	}
}

func (perm *Permissions) GetFeatureCode() string {
	return perm.featureCode
}

//PermissionsCheck check license when program check
func (perm *Permissions) PermissionsCheck(licenseStr string, privateKeyData string) bool {
	if len(licenseStr) == 0 {
		return false
	}

	licenseRsa, err := base64.StdEncoding.DecodeString(licenseStr)
	if err != nil {
		return false
	}

	/* Rsa to plaintext string */
	rsaUtil := crypto.NewRSAHelper([]byte(privateKeyData), nil)
	plainText, err := rsaUtil.RsaDecrypt(licenseRsa)
	if err != nil {
		return false
	}

	if strings.Compare(hex.EncodeToString(plainText), perm.featureCode) == 0 {
		return true
	}

	return false
}

//GenerateLicense generate license
func (perm *Permissions) GenerateLicense(publicKeyData string) (string, error) {

	rsaUtil := crypto.NewRSAHelper(nil, []byte(publicKeyData))
	license, err := rsaUtil.RsaEncrypt([]byte(perm.featureCode))
	if err != nil {
		return string(""), err
	}

	licenseBase64 := base64.StdEncoding.EncodeToString(license)

	return licenseBase64, nil
}
