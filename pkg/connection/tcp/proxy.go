package tcp

import (
	"fmt"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"

	"shadowsocks-go/pkg/config"
	"shadowsocks-go/pkg/crypto"
	"shadowsocks-go/pkg/protocol"

	"github.com/golang/glog"
)

const (
	timeoutKey = string("timeout")
)

//TCPServer maintain a listener
type TCPServer struct {
	Config          *config.ConnectionInfo
	quit            chan struct{}
	uploadTraffic   int64 //request upload traffic
	downloadTraffic int64 //request download traffic
	//ClientDict Mapping from client addresses (as host:port) to connection
	clientDict map[string]*connector

	//mutex Mutex used to serialize access to the dictionary
	mapMutex  *sync.Mutex
	dataMutex *sync.Mutex
	cryp      *crypto.Crypto //decode data read from client
}

type accepted struct {
	conn net.Conn
	err  error
}

type connector struct {
	clientConn net.Conn
	remoteConn map[string]net.Conn //proxy to remote connnection
}

//NewTCPServer create a TCPServer
func NewTCPServer(cfg *config.ConnectionInfo) *TCPServer {
	crypto, err := crypto.NewCrypto(cfg.EncryptMethod, cfg.Password)
	if err != nil {
		glog.Errorf("create crypto error :%v", err)
		return nil
	}

	return &TCPServer{
		Config:     cfg,
		quit:       make(chan struct{}),
		clientDict: make(map[string]*connector),
		mapMutex:   new(sync.Mutex),
		dataMutex:  new(sync.Mutex),
		cryp:       crypto,
	}
}

//Stop implement quit go routine
func (tcpSrv *TCPServer) Stop() {
	glog.V(5).Infof("tcp server close %v\r\n", tcpSrv.Config)
	close(tcpSrv.quit)
}

//Traffic ollection traffic for client,return upload traffic and download traffic
func (tcpSrv *TCPServer) Traffic() (int64, int64) {
	lock(tcpSrv.dataMutex)
	upload := tcpSrv.uploadTraffic
	download := tcpSrv.downloadTraffic

	tcpSrv.uploadTraffic = 0
	tcpSrv.downloadTraffic = 0
	unlock(tcpSrv.dataMutex)

	return upload, download
}

func (tcpSrv *TCPServer) handleRequest(ctx context.Context, client net.Conn, reqAddr string) {

	lock(tcpSrv.mapMutex)
	conn, found := tcpSrv.clientDict[reqAddr]
	if !found {
		conn = &connector{
			clientConn: client,
			remoteConn: make(map[string]net.Conn),
		}
		glog.V(5).Infof("Created new connection for client %s\n", reqAddr)

		tcpSrv.clientDict[reqAddr] = conn
		unlock(tcpSrv.mapMutex)
	} else {
		glog.V(5).Infof("Found connection for client %s this means is our loop already server for this client \r\n", reqAddr)
		unlock(tcpSrv.mapMutex)
		return
	}

	defer func() {
		conn.clientConn.Close()

		lock(tcpSrv.mapMutex)
		delete(tcpSrv.clientDict, reqAddr)
		unlock(tcpSrv.mapMutex)
	}()

	var wg sync.WaitGroup

	var host string
	var remote net.Conn
	for {
		type request struct {
			client net.Conn
			remote net.Conn
			iv     []byte
			err    error
		}

		reqChan := make(chan request, 1)
		go func() {
			ssProtocol, err := protocol.ParseTcpReq(conn.clientConn, tcpSrv.cryp)
			if err != nil {
				glog.Errorf("read a eof maybe remote close socket %v\r\n", err)
				reqChan <- request{
					err: err,
				}
			}

			if tcpSrv.Config.EnableOTA {
				result := ssProtocol.CheckHMAC(tcpSrv.cryp.Key[:])
				if !result {
					glog.Errorln("invalid not auth request")
					reqChan <- request{
						err: fmt.Errorf("not auth"),
					}
				}
			}

			remoteAddr := &net.TCPAddr{
				IP:   ssProtocol.DstAddr.IP,
				Port: ssProtocol.DstAddr.Port,
			}

			lock(tcpSrv.mapMutex)
			host = remoteAddr.String()
			remote, found = conn.remoteConn[host]
			if !found {
				remote, err = net.Dial("tcp", remoteAddr.String())
				if err != nil {
					glog.Errorf(" connecting to:%v occur err:%v", host, err)
					reqChan <- request{
						err: fmt.Errorf("remote connecting failure"),
					}

				}
				conn.remoteConn[host] = remote

				glog.V(5).Infof("Created new remote connection for client %s\n", host)
			}
			unlock(tcpSrv.mapMutex)

			reqChan <- request{
				client: conn.clientConn,
				remote: remote,
				iv:     ssProtocol.IV,
			}
		}()

		select {
		case <-ctx.Done():
			wg.Wait()
			return
		case req := <-reqChan:
			if req.err != nil {
				err := req.err
				if err.Error() == "not auth" {
					continue
				} else if err == io.EOF {
					return
				} else if err.Error() == "remote connecting failure" {
					continue
				} else {
					glog.Errorf("un implement error %v\r\n", err)
					continue
				}
			}

			wg.Add(1)
			go func() {
				upload, download := process(ctx, req.iv, req.client, req.remote)

				lock(tcpSrv.dataMutex)
				tcpSrv.uploadTraffic += <-upload
				tcpSrv.downloadTraffic += <-download
				unlock(tcpSrv.dataMutex)

				remote.Close()
				lock(tcpSrv.mapMutex)
				delete(conn.remoteConn, host)
				unlock(tcpSrv.mapMutex)

				wg.Done()
			}()
		}
	}

}

//Run start a tcp listen for user
func (tcpSrv *TCPServer) Run() {
	port := tcpSrv.Config.Port

	portStr := strconv.Itoa(port)
	ln, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		glog.Errorf("tcp server(%v) error: %v\n", port, err)
	}

	var ctx context.Context
	ctx, cancel := context.WithCancel(context.TODO())

	defer func() {
		cancel()
		ln.Close()
	}()

	var wg sync.WaitGroup
	for {
		c := make(chan accepted, 1)
		go func() {
			glog.V(5).Infoln("wait for accept")
			var conn net.Conn
			conn, err = ln.Accept()
			c <- accepted{conn: conn, err: err}
		}()

		select {
		case <-tcpSrv.quit:
			glog.Infof("Receive Quit singal for %s\r\n", port)
			wg.Wait()
			return
		case accept := <-c:
			if accept.err != nil {
				glog.V(5).Infof("accept error: %v\n", accept.err)
				continue
			}
			reqAddr := accept.conn.RemoteAddr().String()

			go func() {
				wg.Add(1)
				tcpSrv.handleRequest(ctx, accept.conn, reqAddr)
				wg.Done()
			}()

		}
	}
}

func (tcpSrv *TCPServer) Compare(client *config.ConnectionInfo) bool {
	return reflect.DeepEqual(*tcpSrv.Config, *client)
}

func lock(mutex *sync.Mutex) {
	mutex.Lock()
}

func unlock(mutex *sync.Mutex) {
	mutex.Unlock()
}

func setReadTimeout(c net.Conn, timeout time.Duration) {
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
}
