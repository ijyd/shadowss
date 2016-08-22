package connection

import (
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

//Proxy Information maintained for each client/server connection
type Proxy struct {
	port        int           //this proxy listen port
	cipher      *Cipher       //decode data read from client
	timeout     time.Duration //write or read timeoutt
	proxyConn   *net.UDPConn  //listener connection
	oneTimeAuth bool          //account configure one time auth

	//ClientDict Mapping from client addresses (as host:port) to connection
	ClientDict map[string]*Connection

	//dmutex Mutex used to serialize access to the dictionary
	dmutex *sync.Mutex
	//quit loop
	quit chan struct{}

	wg sync.WaitGroup

	UploadTraffic int64 //request upload traffic
}

type receive struct {
	readLen int
	cliAddr *net.UDPAddr
	err     error
	buffer  []byte
}

//NewProxy implement new a proxy. it listen on given port
func NewProxy(port int, cipher *Cipher, auth bool, timeout time.Duration) *Proxy {
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
	proxy.oneTimeAuth = auth

	proxy.ClientDict = make(map[string]*Connection)
	proxy.dmutex = new(sync.Mutex)
	proxy.quit = make(chan struct{})
	proxy.port = port

	return proxy
}

func (pxy *Proxy) dlock() {
	pxy.dmutex.Lock()
}

func (pxy *Proxy) dunlock() {
	pxy.dmutex.Unlock()
}

//Stop quit loop
func (pxy *Proxy) Stop() {
	close(pxy.quit)
	pxy.wg.Wait()
}

//Traffic quit loop
func (pxy *Proxy) Traffic() (int64, int64) {
	var upload, download int64
	for _, val := range pxy.ClientDict {
		download += val.DownloadTraffic
	}
	upload = pxy.UploadTraffic

	return upload, download
}

// func (pxy *Proxy) process() {
// 	var buffer [4096]byte
// 	n, cliaddr, err := pxy.proxyConn.ReadFromUDP(buffer[:])
// 	if err != nil {
// 		//glog.Errorf("Read from %s occur error:%v\r\n", cliaddr.String(), err)
// 		return
// 	}
//
// 	saddr := cliaddr.String()
//
// 	ssProtocol := Parse(buffer[:n], n, pxy.cipher)
// 	if ssProtocol == nil {
// 		glog.Warningln("parse empty request ignore \r\n")
// 		return
// 	}
//
// 	serverAddr := &net.UDPAddr{
// 		IP:   ssProtocol.DstAddr.IP,
// 		Port: ssProtocol.DstAddr.Port,
// 	}
//
// 	glog.V(5).Infof("from client %s to severr %s , read Read len(%d) buffer \r\n%s",
// 		cliaddr.String(), serverAddr.String(),
// 		n, util.DumpHex(ssProtocol.Data[:]))
//
// 	//check our local map for connection
// 	pxy.dlock()
// 	conn, found := pxy.ClientDict[saddr]
// 	if !found {
// 		conn = NewConnection(serverAddr, cliaddr)
// 		pxy.ClientDict[saddr] = conn
// 		pxy.dunlock()
//
// 		glog.V(5).Infof("Created new connection for client %s\n", saddr)
// 		// Fire up routine to manage new connection
// 		pxy.wg.Add(1)
// 		go RunConnection(conn, pxy.proxyConn, pxy.cipher, ssProtocol.RespHeader, &pxy.quit, pxy.wg)
// 	} else {
// 		glog.V(5).Infof("Found connection for client %s\n", saddr)
// 		pxy.dunlock()
// 	}
//
// 	// Relay to server
// 	_, err = conn.ServerConn.Write(ssProtocol.Data[:])
// 	if err != nil {
// 		glog.Warningln("write buffer into remote server failure:", err)
// 		return
// 	}
// }

func (pxy *Proxy) cleanup() {
	glog.Infof("Stop Proxy with %v\r\n", pxy.port)
	err := pxy.proxyConn.Close()
	if err != nil {
		glog.Errorf("connection close error %v\r\n", err)
	}

	for _, val := range pxy.ClientDict {
		close(val.Quit)
	}

	pxy.wg.Done()
}

func (pxy *Proxy) handleRequest(recv receive) {

	saddr := recv.cliAddr.String()
	cliaddr := recv.cliAddr
	n := recv.readLen

	if nil != recv.err && recv.err.(net.Error).Timeout() {
		return
	}

	ssProtocol, err := Parse(recv.buffer, n, pxy.cipher)
	if err != nil {
		glog.Warningf("parse request failure(%v) ignore \r\n", err)
		return
	}

	//check ota
	if pxy.oneTimeAuth {
		if ssProtocol.OneTimeAuth {
			authKey := append(ssProtocol.IV, pxy.cipher.key...)
			authData := append(ssProtocol.RespHeader, ssProtocol.Data...)

			hmac := util.HmacSha1(authKey, authData)
			if !bytes.Equal(ssProtocol.HMAC[:], hmac) {
				glog.Errorf("Unauthorized request\r\n")
				return
			}
		} else {
			glog.Warningf("invalid request with auth \r\n")
			return
		}
	} else {
		glog.V(5).Infof("this client(%d) not enable auth\r\n", pxy.port)
	}

	serverAddr := &net.UDPAddr{
		IP:   ssProtocol.DstAddr.IP,
		Port: ssProtocol.DstAddr.Port,
	}

	glog.V(5).Infof("from client %s to severr %s , read Read len(%d) buffer \r\n%s",
		saddr, serverAddr.String(),
		n, util.DumpHex(ssProtocol.Data[:]))

	//check our local map for connection
	pxy.dlock()
	conn, found := pxy.ClientDict[saddr]
	if !found {
		conn = NewConnection(serverAddr, cliaddr)
		if conn == nil {
			glog.Errorf("Created new connection remote network  is unreachable\n")
			pxy.dunlock()
			return
		}
		pxy.ClientDict[saddr] = conn
		pxy.dunlock()

		glog.V(5).Infof("Created new connection for client %s\n", saddr)
		// Fire up routine to manage new connection
		go conn.RunConnection(pxy.proxyConn, pxy.cipher, ssProtocol.RespHeader, pxy.wg)
	} else {
		glog.V(5).Infof("Found connection for client %s\n", saddr)
		pxy.dunlock()
	}

	// Relay to server
	_, err = conn.ServerConn.Write(ssProtocol.Data[:])
	if err != nil {
		glog.Warningln("write buffer into remote server failure:", err)
		return
	}
	pxy.UploadTraffic += int64(len(ssProtocol.Data[:]))
}

//RunProxy Routine to handle inputs to Proxy port
func (pxy *Proxy) RunProxy() {

	pxy.wg.Add(1)

	for {
		recvChan := make(chan receive, 1)
		go func() {
			var buffer [4096]byte
			n, cliaddr, err := pxy.proxyConn.ReadFromUDP(buffer[:])
			recvChan <- receive{n, cliaddr, err, buffer[0:n]}
		}()

		select {
		case <-pxy.quit:
			pxy.cleanup()
			return
		case recv := <-recvChan:
			pxy.handleRequest(recv)
		}
	}
}
