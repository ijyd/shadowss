package connection

import (
	"net"

	"github.com/golang/glog"
)

//Connection implement for every conection
type Connection struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}

//NewConnection Generate a new connection by opening a UDP connection to the server
func NewConnection(srvAddr, cliAddr *net.UDPAddr) *Connection {
	conn := new(Connection)
	conn.ClientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if err != nil {
		glog.Errorf("Dial to remote server %s err:%v\r\n", srvAddr.String(), err)
		return nil
	}

	conn.ServerConn = srvudp
	return conn
}

//RunConnection Go routine which manages connection from server to single client
func RunConnection(conn *Connection, write *net.UDPConn, cipher *Cipher, respHeader []byte) {
	for {
		// Read from server
		buffer := make([]byte, 2048)
		n, err := conn.ServerConn.Read(buffer[0:])
		if err != nil {
			glog.Errorf("read %s->%s failure %v \r\n", conn.ClientAddr.String(), conn.ServerConn.RemoteAddr().String(), err)
			continue
		}

		// Relay it to client
		respHeaderLen := len(respHeader)
		resp := make([]byte, n+respHeaderLen)
		copy(resp[:], respHeader)
		copy(resp[respHeaderLen:], buffer[0:n])
		encBuff, err := encodeUDPResp(resp[:], n+respHeaderLen, cipher)
		_, err = write.WriteToUDP(encBuff[:], conn.ClientAddr)
		if err != nil {
			glog.Errorf("write local->%s failure %v \r\n", conn.ClientAddr.String(), err)
			continue
		}
	}
}
