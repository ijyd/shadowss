package license

import (
	"bytes"
	"crypto/md5"
)

//License manager a license
type License struct {
	systemUUID         string
	systemSerialNum    string
	systemVersion      string
	biosVersion        string
	biosRelaseDate     string
	biosVendor         string
	baseBoardSerialNum string
	baseBoardVersion   string
	mac                string
	salt               string
}

//NewLicense create a licenses
func NewLicense(mac, salt, sysUUID, sysSerNum, sysVer, biosVer, biosReleaseDate, biosVendor, baseBoradSerNum, baseBoradVer string) *License {
	return &License{
		mac:                mac,
		salt:               salt,
		systemUUID:         sysUUID,
		systemSerialNum:    sysSerNum,
		systemVersion:      sysVer,
		biosVersion:        biosVer,
		biosRelaseDate:     biosReleaseDate,
		biosVendor:         biosVendor,
		baseBoardSerialNum: baseBoradSerNum,
		baseBoardVersion:   baseBoradVer,
	}
}

//CreateSignature create a licese by dmi information
func (lic *License) CreateSignature() ([]byte, error) {

	buffer := bytes.Buffer{}
	buffer.WriteString(lic.mac)
	buffer.WriteString(lic.systemUUID)
	buffer.WriteString(lic.systemSerialNum)
	buffer.WriteString(lic.systemVersion)
	buffer.WriteString(lic.biosVersion)
	buffer.WriteString(lic.biosRelaseDate)
	buffer.WriteString(lic.biosVendor)
	buffer.WriteString(lic.baseBoardSerialNum)
	buffer.WriteString(lic.baseBoardVersion)
	buffer.WriteString(lic.salt)

	h := md5.New()
	h.Write(buffer.Bytes())
	return h.Sum(nil), nil
}

func (lic *License) CheckLicense(licStr []byte) bool {
	sign, err := lic.CreateSignature()
	if err != nil {
		return false
	}

	if !bytes.Equal(sign, licStr) {
		return false
	}

	return true
}
