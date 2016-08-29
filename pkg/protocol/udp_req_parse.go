package protocol

func ParseUDPReq(decBuffer []byte, byteLen int, ivLen int) (*SSProtocol, error) {
	ssProtocol := NewSSProcol(ivLen)

	var err error
	var tmpLen int
	var endIndex int

	parseLen := ssProtocol.ParseReqIV(decBuffer[0:ivLen])

	tmpLen, err = ssProtocol.ParseReqAddrType(decBuffer[parseLen : parseLen+protocolAddrTypeLen])
	if err != nil {
		return nil, err
	}

	parseLen += tmpLen

	if ssProtocol.AddrType == addrTypeDomain {
		tmpLen = ssProtocol.ParseReqDomainLen(decBuffer[parseLen : parseLen+protocolHostLen])
		parseLen += tmpLen
	}

	endIndex = parseLen + ssProtocol.HostLen + protocolDstAddrPortLen
	tmpLen, err = ssProtocol.ParseReqAddrAndPort(decBuffer[parseLen:endIndex])
	if err != nil {
		return nil, err
	}
	parseLen += tmpLen

	if ssProtocol.OneTimeAuth {
		endIndex = byteLen - protocolHMACLen
		tmpLen = ssProtocol.ParseUDPReqData(decBuffer[parseLen:endIndex])
		parseLen += tmpLen

		ssProtocol.ParseReqHMAC(decBuffer[parseLen:])
	} else {
		endIndex = byteLen
		tmpLen = ssProtocol.ParseUDPReqData(decBuffer[parseLen:endIndex])
	}

	return ssProtocol, nil
}
