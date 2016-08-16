package shadowss

import (
	"net"
	"time"

	conn "shadowsocks-go/pkg/connection"
	encrypt "shadowsocks-go/pkg/connection"

	"github.com/golang/glog"
)

var udp bool

const dnsGoroutineNum = 64

type UDPListener struct {
	password string
	listener *net.UDPConn
}

//RunUDP implement a new udp listener
func RunUDP(password, method string, port int, timeout time.Duration) {
	cipher, err := encrypt.NewCipher(method, password)
	if err != nil {
		glog.Fatalf("Error generating cipher for udp port: %d %v\n", port, err)
		return
	}

	proxy := conn.NewProxy(port, cipher.Copy(), timeout)
	if proxy == nil {
		glog.Fatalf("listening upd port: %v error:%v\r\n", port, err)
		return
	}

	go proxy.RunProxy()
}

/*
func RunUDP(password, method string, port int, timeout time.Duration) {
	var cipher *encrypt.Cipher
	listenPort := strconv.Itoa(port)

	serverAddr, err := net.ResolveUDPAddr("udp", ":"+listenPort)
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		glog.Fatalf("listening upd port: %v error:%v\r\n", port, err)
		return
	}
	defer conn.Close()

	// passwdManager.addUDP(password, port, conn)
	// if err != nil {
	// 	glog.Fatalf("error listening udp port %v: %v\n", port, err)
	// 	return
	// }
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
*/
