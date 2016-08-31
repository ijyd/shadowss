package protocol

import (
	"io"
	"net"

	"shadowsocks-go/pkg/crypto"
	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

func parseTCPBuffer(conn net.Conn, cryp *crypto.Crypto, ssProtocol *SSProtocol, stage, expectedLen int) error {

	readLen := expectedLen
	readBuf := make([]byte, readLen)
	decryptBuf := make([]byte, readLen)

	if _, err := io.ReadFull(conn, readBuf[0:readLen]); err != nil {
		glog.Errorln("read err : ", err)
		return err
	}

	glog.V(5).Infof("read data\r\n%s\r\n",
		util.DumpHex(readBuf[0:readLen]))

	if stage == ParseStageIV {
		ssProtocol.ParseReqIV(readBuf[0:readLen])
		return nil
		// } else if stage == ParseStageAddrType {
		// 	iv := ssProtocol.IV
		// 	cryp.Decrypt(decryptBuf[0:readLen], readBuf[0:readLen], iv)
		//
	} else {
		cryp.Decrypt(decryptBuf[0:readLen], readBuf[0:readLen])
	}
	glog.V(5).Infof("plainText data\r\n%s\r\n",
		util.DumpHex(decryptBuf[0:readLen]))

	var err error
	switch stage {
	case ParseStageAddrType:
		_, err = ssProtocol.ParseReqAddrType(decryptBuf[0:readLen])
	case ParseStageDomainLen:
		ssProtocol.ParseReqDomainLen(decryptBuf[0:readLen])
	case ParseStageAddrAnddPort:
		_, err = ssProtocol.ParseReqAddrAndPort(decryptBuf[0:readLen])
	case ParseStageHMAC:
		ssProtocol.ParseReqHMAC(decryptBuf[0:readLen])
	}

	return err
}

func ParseTcpReq(conn net.Conn, cryp *crypto.Crypto) (*SSProtocol, error) {

	ivLen := cryp.GetIVLen()
	ssProtocol := NewSSProcol(ivLen)

	parseTCPBuffer(conn, cryp, ssProtocol, ParseStageIV, ivLen)
	glog.V(5).Infof("read iv\r\n%s\r\n",
		util.DumpHex(ssProtocol.IV[:]))

	_, err := cryp.UpdataCipherStream(ssProtocol.IV, false)
	if err != nil {
		return nil, err
	}

	err = parseTCPBuffer(conn, cryp, ssProtocol, ParseStageAddrType, protocolAddrTypeLen)
	if err != nil {
		glog.Errorf("parse request addr type error %v\r\n", err)
		return nil, err
	}
	glog.V(5).Infof("read addr type %v\r\n",
		ssProtocol.AddrType)

	if ssProtocol.AddrType == addrTypeDomain {
		err = parseTCPBuffer(conn, cryp, ssProtocol, ParseStageDomainLen, protocolHostLen)
		if err != nil {
			glog.Errorf("parse request domain len  error %v\r\n", err)
			return nil, err
		}
		glog.V(5).Infof("read domain len\r\n%v\r\n",
			ssProtocol.HostLen)
	}

	err = parseTCPBuffer(conn, cryp, ssProtocol, ParseStageAddrAnddPort, ssProtocol.HostLen+protocolDstAddrPortLen)
	if err != nil {
		glog.Errorf("parse request addr and port error %v\r\n", err)
		return nil, err
	}
	glog.V(5).Infof("read domain and port %v:%v\r\n",
		ssProtocol.DstAddr.IP.String(), ssProtocol.DstAddr.Port)

	if ssProtocol.OneTimeAuth {
		err = parseTCPBuffer(conn, cryp, ssProtocol, ParseStageHMAC, protocolHMACLen)
		if err != nil {
			glog.Errorf("parse request hmac error %v\r\n", err)
			return nil, err
		}
		glog.V(5).Infof("read hmac %v\r\n", util.DumpHex(ssProtocol.HMAC[:]))
	}

	return ssProtocol, nil
}
