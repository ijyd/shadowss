package connection

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

const (
	dstAddrIPv4Len   = net.IPv4len
	dstAddrIPv6Len   = net.IPv6len
	dstAddrPortLen   = 2
	protolcolHostLen = 1
	protocalHMACLen  = 10
)

const (
	addrTypeIPv4   = 1
	addrTypeDomain = 3
	addrTypeIPv6   = 4
)

const (
	addrOneTimeAuthFlag = 0x10
	addrTypeFlag        = 0x0F
)

//DstAddr description dest ip addr
type DstAddr struct {
	IP   net.IP
	Port int
}

//SSProtocol implement ss package
type SSProtocol struct {
	IV          []byte
	AddrType    byte
	OneTimeAuth bool
	HostLen     int
	DstAddr     DstAddr
	RespHeader  []byte
	Data        []byte
	HMAC        [10]byte
}

//Decrypt decrypt data to plain text
func decrypt(encBuffer []byte, cipher *Cipher) ([]byte, error) {

	byteLen := len(encBuffer)
	ivLen := cipher.info.ivLen
	if byteLen < ivLen {
		return nil, fmt.Errorf("request body too short\r\n")
	}

	iv := make([]byte, ivLen)
	copy(iv, encBuffer[0:cipher.info.ivLen])

	glog.V(5).Infof("Got  decrypt cipher ivlen(%d) iv: \r\n%s", ivLen, util.DumpHex(iv[:]))
	glog.V(5).Infof("Got  decrypt datalen(%d) data:\r\n%s", byteLen, util.DumpHex(encBuffer[ivLen:byteLen]))

	if err := cipher.initDecrypt(iv); err != nil {
		glog.Errorf("init decrypt failure:%v\r\n", err)
		return nil, err
	}

	if len(cipher.iv) == 0 {
		cipher.iv = iv
	}

	decBuffer := make([]byte, byteLen-ivLen)

	cipher.decrypt(decBuffer[:], encBuffer[ivLen:byteLen])
	return decBuffer, nil
}

//Parse ss packet into SSProtocol
func Parse(input []byte, byteLen int, cipher *Cipher) (*SSProtocol, error) {

	decBuffer, err := decrypt(input[0:byteLen], cipher)
	if err != nil {
		return nil, err
	}

	ssProtocal := new(SSProtocol)
	ivLen := cipher.info.ivLen
	ssProtocal.IV = make([]byte, ivLen)
	copy(ssProtocal.IV, input[0:ivLen])

	parseLen := 0
	addr := decBuffer[parseLen]
	ssProtocal.AddrType = addr & addrTypeFlag
	ssProtocal.OneTimeAuth = 0x10 == (addr & addrOneTimeAuthFlag)
	parseLen++

	validBufferLen := byteLen - cipher.info.ivLen

	glog.V(5).Infof("Got decrypt plain text data buffer \r\n%s \r\n", util.DumpHex(decBuffer[0:validBufferLen]))

	glog.V(5).Infof("Got AddrType:%d  one time auth:%t\r\n", ssProtocal.AddrType, ssProtocal.OneTimeAuth)

	switch ssProtocal.AddrType {
	case addrTypeIPv4:
		ssProtocal.DstAddr.IP = net.IP(decBuffer[parseLen : parseLen+dstAddrIPv4Len])
		parseLen += dstAddrIPv4Len
	case addrTypeIPv6:
		ssProtocal.DstAddr.IP = net.IP(decBuffer[parseLen : parseLen+dstAddrIPv6Len])
		parseLen += dstAddrIPv6Len
	case addrTypeDomain:
		hostlen := decBuffer[parseLen]
		parseLen += protolcolHostLen

		ssProtocal.HostLen, _ = strconv.Atoi(string(hostlen))

		domainEndIdx := ssProtocal.HostLen + parseLen
		domain := string(decBuffer[parseLen:domainEndIdx])
		parseLen += ssProtocal.HostLen

		dIP, err := net.ResolveIPAddr("ip", domain)
		if err != nil {
			return nil, fmt.Errorf("[udp]failed to resolve domain name: %s\n", domain)
		}
		ssProtocal.DstAddr.IP = dIP.IP
	default:
		return nil, fmt.Errorf("[udp]addr type %v not supported\r\n", ssProtocal.AddrType)
	}

	ssProtocal.DstAddr.Port = int(binary.BigEndian.Uint16(decBuffer[parseLen : parseLen+dstAddrPortLen]))
	parseLen += dstAddrPortLen

	ssProtocal.RespHeader = make([]byte, parseLen)
	copy(ssProtocal.RespHeader, decBuffer[:parseLen])
	//need to fix resp header
	ssProtocal.RespHeader[0] = ssProtocal.AddrType

	dataBufferLen := 0
	if ssProtocal.OneTimeAuth {
		dataBufferLen = validBufferLen - parseLen - protocalHMACLen

		copy(ssProtocal.HMAC[:], decBuffer[parseLen+dataBufferLen:validBufferLen])

		glog.V(5).Infof("Got decrypt HMAC buffer \r\n%s \r\n", util.DumpHex(ssProtocal.HMAC[:]))
	} else {
		dataBufferLen = validBufferLen - parseLen
	}
	dataBufferEndIdex := dataBufferLen + parseLen
	ssProtocal.Data = make([]byte, dataBufferLen)
	copy(ssProtocal.Data, decBuffer[parseLen:dataBufferEndIdex])

	glog.V(5).Infof("Got decrypt data buffer \r\n%s \r\n", util.DumpHex(ssProtocal.Data[:]))
	return ssProtocal, nil
}

//encodeUDPResp encode buffer for resp.  n = iv + payload
func encodeUDPResp(b []byte, byteLen int, cipher *Cipher) ([]byte, error) {
	dataStart := 0

	iv, err := cipher.initEncryptFake()
	if err != nil {
		glog.Errorf("init encrypt failure %v\r\n", err)
		return nil, err
	}

	dataSize := byteLen + len(iv) // for addr type
	cipherData := make([]byte, dataSize)
	copy(cipherData[0:], iv)
	dataStart = len(iv)

	plainText := make([]byte, byteLen)
	copy(plainText[:], b[:])

	glog.V(5).Infof("encrypt cipher ivlen(%d) iv: \r\n%s \r\n", len(iv), util.DumpHex(iv))
	glog.V(5).Infof("encrypt plainText data : \r\n%s \r\n", util.DumpHex(plainText[:]))

	cipher.encrypt(cipherData[dataStart:], plainText)

	glog.V(5).Infof("encrypt data: \r\n%s \r\n", util.DumpHex(cipherData[dataStart:]))

	return cipherData, nil
}
