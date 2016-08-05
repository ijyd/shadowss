package shadowss

import (
	"log"
	"net"
	"time"

	"shadowsocks-go/pkg/connection"
	encrypt "shadowsocks-go/pkg/connection"

	"github.com/golang/glog"
)

var udp bool

const dnsGoroutineNum = 64

type UDPListener struct {
	password string
	listener *net.UDPConn
}

func RunUDP(password, method string, port int, timeout time.Duration) {
	var cipher *encrypt.Cipher
	listenPort := port

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv6zero,
		Port: listenPort,
	})
	defer conn.Close()
	if err != nil {
		glog.Fatalf("listening upd port: %v error:%v\r\n", port, err)
	}

	passwdManager.addUDP(password, port, conn)
	if err != nil {
		glog.Fatalf("error listening udp port %v: %v\n", port, err)
		return
	}
	cipher, err = encrypt.NewCipher(method, password)
	if err != nil {
		log.Printf("Error generating cipher for udp port: %s %v\n", port, err)
		conn.Close()
	}

	UDPConn := connection.NewUDPConn(conn, cipher, timeout)
	for {
		UDPConn.ReadAndHandleUDPReq()
	}
}
