package connection

import (
	"encoding/binary"
	"net"
	"strconv"

	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

// const (
// 	idType  = 0 // address type index
// 	idIP0   = 1 // ip addres start index
// 	idDmLen = 1 // domain address length index
// 	idDm0   = 2 // domain address start index
//
// 	typeIPv4 = 1 // type is ipv4 address
// 	typeDm   = 3 // type is domain address
// 	typeIPv6 = 4 // type is ipv6 address
//
// 	lenIPv4         = 1 + net.IPv4len + 2 // 1addrType + ipv4 + 2port
// 	lenIPv6         = 1 + net.IPv6len + 2 // 1addrType + ipv6 + 2port
// 	lenDmBase       = 1 + 1 + 2           // 1addrType + 1addrLen + 2port, plus addrLen
// 	udpLeakyBufSize = 4108                // data.len(2) + hmacsha1(10) + data(4096)
// 	udpMaxNBuf      = 2048
// )

const (
	dstAddrIPv4Len   = net.IPv4len
	dstAddrIPv6Len   = net.IPv6len
	dstAddrPortLen   = 2
	protolcolHostLen = 1
)

const (
	addrTypeIPv4   = 1
	addrTypeDomain = 3
	addrTypeIPv6   = 4
)

//DstAddr description dest ip addr
type DstAddr struct {
	IP   net.IP
	Port int
}

//SSProtocol implement ss package
type SSProtocol struct {
	AddrType   byte
	OAuth      bool
	HostLen    int
	DstAddr    DstAddr
	RespHeader []byte
	Data       []byte
	HMAC       [10]byte
}

//Decrypt decrypt data to plain text
func decrypt(encBuffer []byte, cipher *Cipher) ([]byte, error) {

	byteLen := len(encBuffer)
	ivLen := cipher.info.ivLen
	iv := make([]byte, ivLen)
	copy(iv, encBuffer[0:cipher.info.ivLen])

	glog.V(5).Infof("Got  decrypt cipher ivlen(%d) iv:%s", ivLen, util.DumpHex(iv[:]))
	glog.V(5).Infof("Got  decrypt datalen(%d) data:%s", byteLen, util.DumpHex(encBuffer[ivLen:]))

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
func Parse(input []byte, byteLen int, cipher *Cipher) *SSProtocol {

	decBuffer, err := decrypt(input[0:byteLen], cipher)
	if err != nil {
		glog.Errorf("Decrypt failure %v\r\n", err)
	}
	validBufferLen := byteLen - cipher.info.ivLen

	glog.V(5).Infof("Got decrypt plain text buffer(%s) \r\n", util.DumpHex(decBuffer[0:validBufferLen]))

	ssProtocal := new(SSProtocol)
	parseLen := 0
	ssProtocal.AddrType = decBuffer[parseLen]
	parseLen++

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
			glog.Fatalf("[udp]failed to resolve domain name: %s\n", domain)
			return nil
		}
		ssProtocal.DstAddr.IP = dIP.IP
	default:
		glog.Fatalf("[udp]addr type %v not supported\r\n", ssProtocal.AddrType)
		return nil
	}

	ssProtocal.DstAddr.Port = int(binary.BigEndian.Uint16(decBuffer[parseLen : parseLen+dstAddrPortLen]))
	parseLen += dstAddrPortLen

	ssProtocal.RespHeader = make([]byte, parseLen)
	copy(ssProtocal.RespHeader, decBuffer[:parseLen])

	ssProtocal.Data = make([]byte, validBufferLen-parseLen)
	copy(ssProtocal.Data, decBuffer[parseLen:validBufferLen])

	return ssProtocal
}

// //parseUDPRequest parse udp package, return request buffer and remote ip in request
// func parseUDPRequest(input []byte, byteLen int, cipher *Cipher, output []byte) *net.UDPAddr {
// 	iv := make([]byte, cipher.info.ivLen)
//
// 	copy(iv, input[0:cipher.info.ivLen])
// 	glog.Infof("Got  decrypt cipher ivlen(%d) iv:%s", cipher.info.ivLen, util.DumpHex(iv[:]))
// 	glog.Infof("Got  decrypt datalen(%d) data:%s", byteLen, util.DumpHex(input[cipher.info.ivLen:byteLen]))
// 	if err := cipher.initDecrypt(iv); err != nil {
// 		glog.Errorf("init decrypt failure:%v\r\n", err)
// 		return nil
// 	}
// 	if len(cipher.iv) == 0 {
// 		cipher.iv = iv
// 	}
//
// 	var decBuffer [2048]byte
//
// 	cipher.decrypt(decBuffer[0:byteLen-cipher.info.ivLen], input[cipher.info.ivLen:byteLen])
// 	n := byteLen - cipher.info.ivLen
//
// 	glog.Infof("Got decrypt plain text buffer(%s) \r\n", util.DumpHex(decBuffer[0:byteLen-cipher.info.ivLen]))
//
// 	var dstIP net.IP
// 	var reqLen int
//
// 	switch decBuffer[idType] {
// 	case typeIPv4:
// 		reqLen = lenIPv4
// 		dstIP = net.IP(decBuffer[idIP0 : idIP0+net.IPv4len])
// 	case typeIPv6:
// 		reqLen = lenIPv6
// 		dstIP = net.IP(decBuffer[idIP0 : idIP0+net.IPv6len])
// 	case typeDm:
// 		reqLen = int(decBuffer[idDmLen]) + lenDmBase
// 		dIP, err := net.ResolveIPAddr("ip", string(decBuffer[idDm0:idDm0+decBuffer[idDmLen]]))
// 		if err != nil {
// 			glog.Fatalf("[udp]failed to resolve domain name: %s\n", string(decBuffer[idDm0:idDm0+decBuffer[idDmLen]]))
// 			return nil
// 		}
// 		dstIP = dIP.IP
// 	default:
// 		glog.Fatalf("[udp]addr type %v not supported\r\n", decBuffer[idType])
// 		return nil
// 	}
//
// 	remoteServer := &net.UDPAddr{
// 		IP:   dstIP,
// 		Port: int(binary.BigEndian.Uint16(decBuffer[reqLen-2 : reqLen])),
// 	}
//
// 	copy(output, decBuffer[reqLen:n])
//
// 	return remoteServer
// }

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

	glog.V(5).Infof("encrypt cipher ivlen(%d) iv:%s \r\n", len(iv), util.DumpHex(iv))
	glog.V(5).Infof("encrypt plainText data :%s \r\n", util.DumpHex(plainText[:]))

	cipher.encrypt(cipherData[dataStart:], plainText)

	glog.V(5).Infof("encrypt data %s \r\n", util.DumpHex(cipherData[dataStart:]))

	return cipherData, nil
}
