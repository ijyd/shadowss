package connection

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
)

//Proxy Information maintained for each client/server connection
type Proxy struct {
	port      int           //this proxy listen port
	cipher    *Cipher       //decode data read from client
	timeout   time.Duration //write or read timeoutt
	proxyConn *net.UDPConn  //listener connection

	//ClientDict Mapping from client addresses (as host:port) to connection
	ClientDict map[string]*Connection

	//dmutex Mutex used to serialize access to the dictionary
	dmutex *sync.Mutex
}

//NewProxy implement new a proxy. it listen on given port
func NewProxy(port int, cipher *Cipher, timeout time.Duration) *Proxy {
	proxy := new(Proxy)
	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil
	}
	pudp, err := net.ListenUDP("udp", saddr)
	if err != nil {
		glog.Errorf("listen on %d failure : %v\r\n", port, err)
		return nil
	}

	glog.V(5).Infof("Proxy serving on port %d\n", port)
	proxy.proxyConn = pudp
	proxy.cipher = cipher
	proxy.timeout = timeout

	proxy.ClientDict = make(map[string]*Connection)
	proxy.dmutex = new(sync.Mutex)

	return proxy
}

func (pxy *Proxy) dlock() {
	pxy.dmutex.Lock()
}

func (pxy *Proxy) dunlock() {
	pxy.dmutex.Unlock()
}

//RunProxy Routine to handle inputs to Proxy port
func (pxy *Proxy) RunProxy() {

	for {
		var buffer [4096]byte
		n, cliaddr, err := pxy.proxyConn.ReadFromUDP(buffer[:])
		if err != nil {
			glog.Errorf("Read from %s occur error:%v\r\n", cliaddr.String(), err)
			continue
		}

		saddr := cliaddr.String()

		ssProtocol := Parse(buffer[:n], n, pxy.cipher)
		if ssProtocol == nil {
			glog.Warningln("parse empty request ignore \r\n")
			continue
		}

		serverAddr := &net.UDPAddr{
			IP:   ssProtocol.DstAddr.IP,
			Port: ssProtocol.DstAddr.Port,
		}

		glog.V(5).Infof("Read len(%d) buffer '%s' from client %s to severr %s\n",
			n, string(ssProtocol.Data[:]), cliaddr.String(), serverAddr.String())

		//check our local map for connection
		pxy.dlock()
		conn, found := pxy.ClientDict[saddr]
		if !found {
			conn = NewConnection(serverAddr, cliaddr)
			pxy.ClientDict[saddr] = conn
			pxy.dunlock()

			glog.V(5).Infof("Created new connection for client %s\n", saddr)
			// Fire up routine to manage new connection
			go RunConnection(conn, pxy.proxyConn, pxy.cipher, ssProtocol.RespHeader)
		} else {
			glog.V(5).Infof("Found connection for client %s\n", saddr)
			pxy.dunlock()
		}

		// Relay to server
		_, err = conn.ServerConn.Write(ssProtocol.Data[:])
		if err != nil {
			glog.Warningln("write buffer into remote server failure:", err)
			continue
		}
	}
}
