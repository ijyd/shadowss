package permissions

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golib/pkg/util/crypto"
	"golib/pkg/util/dmidecode"
	"golib/pkg/util/license"
	"golib/pkg/util/network"
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

type Permissions struct {
	salt    string
	rsaUtil *crypto.RSAHelper
}

var PermissionsHandler = NewPermissions(salt, privateKeyData, publicKeyData)

func NewPermissions(salt, privateKeyData, publicKeyData string) *Permissions {
	return &Permissions{
		salt:    salt,
		rsaUtil: crypto.NewRSAHelper([]byte(privateKeyData), []byte(publicKeyData)),
	}
}

func (perm *Permissions) createLicenseHandler() (*license.License, error) {
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

	// glog.V(5).Infof("get system information sysSerialNum(%v) sysUUID(%v) sysVersion(%v)\r\n",
	// 	sysSerialNum, sysUUID, sysVersion)

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

	// glog.V(5).Infof("get board information baseBoardVer(%v) baseBoardSerialNum(%v) \r\n",
	// 	baseBoardVer, baseBoardSerialNum)

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
	// glog.V(5).Infof("get bios information biosVer(%v) biosReleaseDate(%v) biosVendor(%v)\r\n",
	// 	biosVer, biosReleaseDate, biosVendor)

	mac, err := network.ExternalMAC()

	licenseHelper := license.NewLicense(mac, perm.salt,
		sysUUID, sysSerialNum, sysVersion, biosVer, biosReleaseDate, biosVendor, baseBoardSerialNum, baseBoardVer)

	return licenseHelper, err
}

//PermissionsCheck check license when program check
func (perm *Permissions) PermissionsCheck(licenseStr string) bool {
	if len(licenseStr) == 0 {
		return false
	}

	hwcode, err := perm.GenerateHardWareCode()
	if err != nil {
		return false
	}

	licenseRsa, err := base64.StdEncoding.DecodeString(licenseStr)
	if err != nil {
		return false
	}

	/* Rsa to plaintext string */
	plainText, err := perm.rsaUtil.RsaDecrypt(licenseRsa)
	if err != nil {
		return false
	}

	if strings.Compare(hex.EncodeToString(plainText), hwcode) == 0 {
		return true
	}

	return false
}

//GenerateLicense generate license
func (perm *Permissions) GenerateHardWareCode() (string, error) {

	licenseHelper, err := perm.createLicenseHandler()
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
func (perm *Permissions) GenerateLicense(code []byte) (string, error) {

	license, err := perm.rsaUtil.RsaEncrypt([]byte(code))
	if err != nil {
		return string(""), err
	}

	licenseBase64 := base64.StdEncoding.EncodeToString(license)

	//PermissionsCheck(licenseBase64)

	return licenseBase64, nil
}
