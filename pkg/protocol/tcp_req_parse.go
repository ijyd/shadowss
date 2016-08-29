package protocol

import (
	"io"
	"net"

	"shadowsocks-go/pkg/crypto"

	"github.com/golang/glog"
)

func parseTCPBuffer(conn net.Conn, cryp *crypto.Crypto, ssProtocol *SSProtocol, stage, expectedLen int) error {

	readLen := expectedLen
	readBuf := make([]byte, readLen)
	decryptBuf := make([]byte, readLen)

	if _, err := io.ReadFull(conn, readBuf[0:readLen]); err != nil {
		glog.Errorln("read address type : ", err)
		return err
	}

	if stage == ParseStageIV {
		ssProtocol.ParseReqIV(readBuf[0:readLen])
		return nil
	} else {
		iv := ssProtocol.IV
		cryp.Decrypt(iv, decryptBuf[0:readLen], readBuf[:readLen])
	}

	switch stage {
	case ParseStageAddrType:
		_, _ = ssProtocol.ParseReqAddrType(decryptBuf[0:readLen])
	case ParseStageDomainLen:
		_ = ssProtocol.ParseReqDomainLen(decryptBuf[0:readLen])
	case ParseStageAddrAnddPort:
		_, _ = ssProtocol.ParseReqAddrAndPort(decryptBuf[0:readLen])
	case ParseStageHMAC:
		ssProtocol.ParseReqHMAC(decryptBuf[0:readLen])
	}

	return nil
}

func ParseTcpReq(conn net.Conn, cryp *crypto.Crypto) (*SSProtocol, error) {

	ivLen := cryp.GetIVLen()
	ssProtocol := NewSSProcol(ivLen)

	parseTCPBuffer(conn, cryp, ssProtocol, ParseStageIV, ivLen)

	parseTCPBuffer(conn, cryp, ssProtocol, ParseStageAddrType, protocolAddrTypeLen)

	if ssProtocol.AddrType == addrTypeDomain {
		parseTCPBuffer(conn, cryp, ssProtocol, ParseStageDomainLen, protocolHostLen)
	}

	parseTCPBuffer(conn, cryp, ssProtocol, ParseStageAddrAnddPort, ssProtocol.HostLen+protocolDstAddrPortLen)

	if ssProtocol.OneTimeAuth {
		parseTCPBuffer(conn, cryp, ssProtocol, ParseStageHMAC, protocolHMACLen)
	}

	return ssProtocol, nil
}
