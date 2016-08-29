package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

const (
	protocolAddrTypeLen    = 1
	protocolDstAddrIPv4Len = net.IPv4len
	protocolDstAddrIPv6Len = net.IPv6len
	protocolDstAddrPortLen = 2
	protocolHostLen        = 1
	protocolHMACLen        = 10
)

const (
	addrTypeIPv4   = 1
	addrTypeDomain = 3
	addrTypeIPv6   = 4
)

const (
	AddrOneTimeAuthFlag = 0x10
	addrTypeFlag        = 0x0F
)

const (
	ParseStageIV = iota
	ParseStageAddrType
	ParseStageAddrAnddPort
	ParseStageUDPData
	ParseStageHMAC
	ParseStageDomainLen
)

//DstAddr description dest ip addr
type DstAddr struct {
	IP   net.IP
	Port int
}

//SSProtocol implement ss package
type SSProtocol struct {
	RespHeader  []byte
	IV          []byte
	AddrType    byte
	OneTimeAuth bool
	HostLen     int
	DstAddr     DstAddr
	Data        []byte
	HMAC        [10]byte
	ParseStage  int
}

func NewSSProcol(ivLen int) *SSProtocol {
	return &SSProtocol{
		RespHeader: make([]byte, 1, 32),
		IV:         make([]byte, ivLen),
		ParseStage: ParseStageIV,
	}
}

func (ssProtocal *SSProtocol) ParseReqIV(decBuffer []byte) int {
	copy(ssProtocal.IV, decBuffer[0:len(ssProtocal.IV)])

	return len(ssProtocal.IV)
}

func (ssProtocal *SSProtocol) ParseReqAddrType(decBuffer []byte) (int, error) {

	addr := decBuffer[0]
	ssProtocal.AddrType = addr & addrTypeFlag
	ssProtocal.OneTimeAuth = 0x10 == (addr & AddrOneTimeAuthFlag)

	ssProtocal.RespHeader = append(ssProtocal.RespHeader, decBuffer...)
	ssProtocal.RespHeader[0] = ssProtocal.AddrType

	var err error

	switch ssProtocal.AddrType {
	case addrTypeIPv4:
		ssProtocal.HostLen = protocolDstAddrIPv4Len
	case addrTypeIPv6:
		ssProtocal.HostLen = protocolDstAddrIPv6Len
	case addrTypeDomain:
		ssProtocal.HostLen = 0
	default:
		err = fmt.Errorf("[udp]addr type %v not supported\r\n", ssProtocal.AddrType)
	}

	return protocolAddrTypeLen, err
}

func (ssProtocal *SSProtocol) ParseReqDomainLen(decBuffer []byte) int {
	ssProtocal.RespHeader = append(ssProtocal.RespHeader, decBuffer...)

	hostlen := decBuffer[0]
	ssProtocal.HostLen, _ = strconv.Atoi(string(hostlen))

	return protocolHostLen
}

func (ssProtocal *SSProtocol) ParseReqAddrAndPort(decBuffer []byte) (int, error) {
	ssProtocal.RespHeader = append(ssProtocal.RespHeader, decBuffer...)

	switch ssProtocal.AddrType {
	case addrTypeIPv4, addrTypeIPv6:
		ssProtocal.DstAddr.IP = net.IP(decBuffer[:ssProtocal.HostLen])
	case addrTypeDomain:
		domainEndIdx := ssProtocal.HostLen
		domain := string(decBuffer[:domainEndIdx])

		dIP, err := net.ResolveIPAddr("ip", domain)
		if err != nil {
			return 0, fmt.Errorf("[udp]failed to resolve domain name: %s\n", domain)
		}
		ssProtocal.DstAddr.IP = dIP.IP

	default:
		return 0, fmt.Errorf("[udp]addr type %v not supported\r\n", ssProtocal.AddrType)
	}

	parseLen := ssProtocal.HostLen
	ssProtocal.DstAddr.Port = int(binary.BigEndian.Uint16(decBuffer[parseLen : parseLen+protocolDstAddrPortLen]))

	return parseLen, nil
}

func (ssProtocal *SSProtocol) ParseReqHMAC(decBuffer []byte) int {
	copy(ssProtocal.HMAC[:], decBuffer[:protocolHMACLen])

	return protocolHMACLen
}

func (ssProtocal *SSProtocol) ParseUDPReqData(decBuffer []byte) int {
	buffLen := len(decBuffer)

	ssProtocal.Data = make([]byte, buffLen)
	copy(ssProtocal.Data, decBuffer[:buffLen])

	return buffLen
}

// //Parse ss packet into SSProtocol
// func Parse(decBuffer []byte, byteLen, ivLen int) (*SSProtocol, error) {
//
// 	ssProtocal := new(SSProtocol)
//
// 	ssProtocal.ParseReqIV(decBuffer, ivLen)
// 	for index := byteLen; index > 0; index++ {
//
// 	}
//
// 	parseLen := 0 + ivLen
// 	ssProtocal.ParseReqIV(decBuffer, byteLen, ivLen)
//
// 	parseLen++
//
// 	glog.V(5).Infof("Got decrypt plain text data buffer \r\n%s \r\n", util.DumpHex(decBuffer[ivLen:]))
//
// 	glog.V(5).Infof("Got AddrType:%d  one time auth:%t\r\n", ssProtocal.AddrType, ssProtocal.OneTimeAuth)
//
// 	switch ssProtocal.AddrType {
// 	case addrTypeIPv4:
// 		ssProtocal.DstAddr.IP = net.IP(decBuffer[parseLen : parseLen+protocaldstAddrIPv4Len])
// 		parseLen += dstAddrIPv4Len
// 	case addrTypeIPv6:
// 		ssProtocal.DstAddr.IP = net.IP(decBuffer[parseLen : parseLen+protocaldstAddrIPv6Len])
// 		parseLen += dstAddrIPv6Len
// 	case addrTypeDomain:
// 		hostlen := decBuffer[parseLen]
// 		parseLen += protolcolHostLen
//
// 		ssProtocal.HostLen, _ = strconv.Atoi(string(hostlen))
//
// 		domainEndIdx := ssProtocal.HostLen + parseLen
// 		domain := string(decBuffer[parseLen:domainEndIdx])
// 		parseLen += ssProtocal.HostLen
//
// 		dIP, err := net.ResolveIPAddr("ip", domain)
// 		if err != nil {
// 			return nil, fmt.Errorf("[udp]failed to resolve domain name: %s\n", domain)
// 		}
// 		ssProtocal.DstAddr.IP = dIP.IP
// 	default:
// 		return nil, fmt.Errorf("[udp]addr type %v not supported\r\n", ssProtocal.AddrType)
// 	}
//
// 	ssProtocal.DstAddr.Port = int(binary.BigEndian.Uint16(decBuffer[parseLen : parseLen+protocaldstAddrPortLen]))
// 	parseLen += dstAddrPortLen
//
// 	ssProtocal.RespHeader = make([]byte, parseLen-ivLen)
// 	copy(ssProtocal.RespHeader, decBuffer[ivLen:parseLen])
// 	//need to fix resp header
// 	ssProtocal.RespHeader[0] = ssProtocal.AddrType
//
// 	dataBufferLen := 0
// 	if ssProtocal.OneTimeAuth {
// 		dataBufferLen = byteLen - parseLen - protocalHMACLen
//
// 		copy(ssProtocal.HMAC[:], decBuffer[parseLen+dataBufferLen:])
//
// 		glog.V(5).Infof("Got decrypt HMAC buffer \r\n%s \r\n", util.DumpHex(ssProtocal.HMAC[:]))
// 	} else {
// 		dataBufferLen = byteLen - parseLen
// 	}
// 	dataBufferEndIdex := dataBufferLen + parseLen
// 	ssProtocal.Data = make([]byte, dataBufferLen)
// 	copy(ssProtocal.Data, decBuffer[parseLen:dataBufferEndIdex])
//
// 	glog.V(5).Infof("Got decrypt data buffer \r\n%s \r\n", util.DumpHex(ssProtocal.Data[:]))
// 	return ssProtocal, nil
// }

func (ssProtocol *SSProtocol) CheckHMAC(key []byte) bool {
	if ssProtocol.OneTimeAuth {
		authKey := append(ssProtocol.IV, key...)
		reqHeader := make([]byte, len(ssProtocol.RespHeader))
		copy(reqHeader, ssProtocol.RespHeader)
		reqHeader[0] = ssProtocol.AddrType | (AddrOneTimeAuthFlag)

		authData := reqHeader
		glog.V(5).Infof("request auth data: \r\n %s \r\n  authKey:\r\n %s \r\n", util.DumpHex(authData), util.DumpHex(authKey))

		hmac := util.HmacSha1(authKey, authData)
		if !bytes.Equal(ssProtocol.HMAC[:], hmac) {
			glog.Errorf("Unauthorized request\r\n")
			return false
		}
	} else {
		glog.Warningf("invalid request with auth \r\n")
		return false
	}

	return true
}
