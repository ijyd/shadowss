package permissions

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"shadowss/pkg/util/crypto"
	"shadowss/pkg/util/dmidecode"
	"shadowss/pkg/util/license"
	"shadowss/pkg/util/network"

	"github.com/golang/glog"
)

const (
	SysSerialNumber       = "Serial Number"
	SysUUID               = "UUID"
	SysVersion            = "Version"
	BIOSVersion           = "Version"
	BIOSReleaseDate       = "Release Date"
	BIOSVendor            = "Vendor"
	BaseBoardVersion      = "Version"
	BaseBoardSerialNumber = "Serial Number"
)

const (
	salt = "b34953e8b3428438fed88ba38b02a456d19e1b5100e266f70eae8083dff0348e"
)

func CreateLicenseHandler() (*license.License, error) {
	dmi := dmidecode.NewDMI()
	if err := dmi.Run(); err != nil {
		return nil, err
	}

	notFound := fmt.Errorf("Not found dmi code")
	systemInfo, err := dmi.SearchByName("System Information")
	if err != nil {
		return nil, err
	}
	sysSerialNum, found := systemInfo[SysSerialNumber]
	if !found {
		return nil, notFound
	}

	sysUUID, found := systemInfo[SysUUID]
	if !found {
		return nil, notFound
	}
	sysVersion, found := systemInfo[SysVersion]
	if !found {
		return nil, notFound
	}

	glog.V(5).Infof("get system information sysSerialNum(%v) sysUUID(%v) sysVersion(%v)\r\n",
		sysSerialNum, sysUUID, sysVersion)

	baseBoardInfo, err := dmi.SearchByName("Base Board Information")
	if err != nil {
		return nil, err
	}
	baseBoardVer, found := baseBoardInfo[BaseBoardVersion]
	if !found {
		return nil, notFound
	}

	baseBoardSerialNum, found := baseBoardInfo[BaseBoardSerialNumber]
	if !found {
		return nil, notFound
	}

	glog.V(5).Infof("get board information baseBoardVer(%v) baseBoardSerialNum(%v) \r\n",
		baseBoardVer, baseBoardSerialNum)

	biosInfo, err := dmi.SearchByName("BIOS Information")
	if err != nil {
		return nil, err
	}
	biosVer, found := biosInfo[BIOSVersion]
	if !found {
		return nil, notFound
	}

	biosReleaseDate, found := biosInfo[BIOSReleaseDate]
	if !found {
		return nil, notFound
	}

	biosVendor, found := biosInfo[BIOSVendor]
	if !found {
		return nil, notFound
	}
	glog.V(5).Infof("get bios information biosVer(%v) biosReleaseDate(%v) biosVendor(%v)\r\n",
		biosVer, biosReleaseDate, biosVendor)

	mac, err := network.ExternalMAC()
	glog.V(5).Infof("get mac information %+v\r\n", mac)

	licenseHelper := license.NewLicense(mac, salt,
		sysUUID, sysSerialNum, sysVersion, biosVer, biosReleaseDate, biosVendor, baseBoardSerialNum, baseBoardVer)

	return licenseHelper, err
}

//PermissionsCheck check license when program check
func PermissionsCheck(licenseStr string) bool {
	if len(licenseStr) == 0 {
		return false
	}

	hwcode, err := GenerateHardWareCode()
	if err != nil {
		return false
	}

	licenseRsa, err := base64.StdEncoding.DecodeString(licenseStr)
	if err != nil {
		return false
	}

	/* Rsa to plaintext string */
	rsaHelper := crypto.NewRSAHelper([]byte(privateKeyData), []byte(publicKeyData))
	plainText, err := rsaHelper.RsaDecrypt(licenseRsa)
	if err != nil {
		return false
	}

	if strings.Compare(hex.EncodeToString(plainText), hwcode) == 0 {
		glog.V(5).Infof("get  license %v true\r\n", hex.EncodeToString(plainText))
		return true
	}
	glog.V(5).Infof("get  license %v false\r\n", hex.EncodeToString(plainText))

	return false
}

//GenerateLicense generate license
func GenerateHardWareCode() (string, error) {

	licenseHelper, err := CreateLicenseHandler()
	if err != nil {
		return "", err
	}

	hwcode, err := licenseHelper.CreateSignature()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hwcode), err
}

//GenerateLicense generate license
func GenerateLicense(code []byte) (string, error) {

	rsaHelper := crypto.NewRSAHelper([]byte(privateKeyData), []byte(publicKeyData))

	license, err := rsaHelper.RsaEncrypt([]byte(code))
	if err != nil {
		return string(""), err
	}

	licenseBase64 := base64.StdEncoding.EncodeToString(license)

	PermissionsCheck(licenseBase64)

	return licenseBase64, nil
}
