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
}

type remoteConnHelper struct {
	cryp        *crypto.Crypto //every to remote request have diff iv
	server      net.Conn       //proxy to remote connnection
	chunkID     uint32         //record request count for this remote
	iv          []byte         //store first request iv for continue request use
	oneTimeAuth bool           //for this remote oneTimeAuth flag
}

type connector struct {
	clientConn net.Conn
	remoteConn map[string]*remoteConnHelper
}

func (r *remoteConnHelper) increaseChunkID() uint32 {
	chunkID := r.chunkID
	r.chunkID += 1
	return chunkID
}

//NewTCPServer create a TCPServer
func NewTCPServer(cfg *config.ConnectionInfo) *TCPServer {
	return &TCPServer{
		Config:     cfg,
		quit:       make(chan struct{}),
		clientDict: make(map[string]*connector),
		mapMutex:   new(sync.Mutex),
		dataMutex:  new(sync.Mutex),
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

func (tcpSrv *TCPServer) parseRequest(client net.Conn, cryp *crypto.Crypto) (*protocol.SSProtocol, error) {
	ssProtocol, err := protocol.ParseTcpReq(client, cryp)
	if err != nil {
		return nil, err
	}

	if tcpSrv.Config.EnableOTA {
		reqHeader := make([]byte, len(ssProtocol.RespHeader))
		copy(reqHeader, ssProtocol.RespHeader)
		reqHeader[0] = ssProtocol.AddrType | (protocol.AddrOneTimeAuthFlag)

		result := ssProtocol.CheckHMAC(cryp.Key[:], reqHeader)
		if !result {
			glog.Errorln("invalid not auth request")
			return nil, fmt.Errorf("not auth")
		}
	}

	return ssProtocol, nil
}

func (tcpSrv *TCPServer) handleRequest(ctx context.Context, client net.Conn, reqAddr string) {

	lock(tcpSrv.mapMutex)
	conn, found := tcpSrv.clientDict[reqAddr]
	if !found {

		conn = &connector{
			clientConn: client,
			remoteConn: make(map[string]*remoteConnHelper),
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
		glog.V(5).Infof("close pipe %v\r\n", reqAddr)
		conn.clientConn.Close()

		lock(tcpSrv.mapMutex)
		delete(tcpSrv.clientDict, reqAddr)
		unlock(tcpSrv.mapMutex)
	}()

	//new cyrpto for this remote
	crypto, err := crypto.NewCrypto(tcpSrv.Config.EncryptMethod, tcpSrv.Config.Password)
	if err != nil {
		glog.Errorf("create crypto error :%v", err)
		return
	}
	ssProtocol, err := tcpSrv.parseRequest(conn.clientConn, crypto)
	if err != nil {
		glog.Errorf("get invalid request  from %s error %v\r\n", reqAddr, err)
		return
	}

	var host string
	var remote *remoteConnHelper

	remoteAddr := &net.TCPAddr{
		IP:   ssProtocol.DstAddr.IP,
		Port: ssProtocol.DstAddr.Port,
	}

	type request struct {
		client net.Conn
		remote *remoteConnHelper
		iv     []byte
		err    error
	}
	reqChan := make(chan request, 1)
	go func() {
		lock(tcpSrv.mapMutex)
		host = remoteAddr.String()
		remote, found = conn.remoteConn[host]
		if !found {
			remoteSrv, err := net.Dial("tcp", remoteAddr.String())
			if err != nil {
				glog.Errorf(" connecting to:%v occur err:%v", host, err)
				reqChan <- request{
					err: fmt.Errorf("remote connecting failure"),
				}
				return
			}

			//reset encrypt stream
			_, err = crypto.UpdataCipherStream(ssProtocol.IV, true)
			if err != nil {
				reqChan <- request{
					err: err,
				}
				return
			}

			remote = &remoteConnHelper{
				cryp:        crypto,
				server:      remoteSrv,
				iv:          ssProtocol.IV,
				chunkID:     0,
				oneTimeAuth: tcpSrv.Config.EnableOTA,
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

	var wg sync.WaitGroup
	processDone := make(chan struct{})
	for {
		select {
		case <-ctx.Done():
			glog.V(5).Infof("handle %s read requet will be done\n", reqAddr)
			wg.Wait()
			return
		case <-processDone:
			glog.V(5).Infof("handle %s read requet process done \n", reqAddr)
			wg.Wait()
			return
		case req := <-reqChan:
			if req.err != nil {
				glog.Errorf("handle %s read requet error: %v\n", reqAddr, req.err)
				err := req.err
				if err.Error() == "not auth" {
					continue
				} else if err == io.EOF {
					glog.Errorf("handle %s read requet error: %v will be return\n", reqAddr, req.err)
					return
				} else if err.Error() == "remote connecting failure" {
					continue
				} else {
					glog.Errorf("not implement error %v\r\n", err)
					return
				}
			}

			wg.Add(1)
			go func() {
				glog.Infof("handle %s read process %v->%v\n", reqAddr, req.client.RemoteAddr().String(), req.remote.server.RemoteAddr().String())
				upload, download := process(ctx, req.client, req.remote)

				glog.Infof("handle %s read process %v->%v done\n", reqAddr, req.client.RemoteAddr().String(), req.remote.server.RemoteAddr().String())
				lock(tcpSrv.dataMutex)
				tcpSrv.uploadTraffic += <-upload
				tcpSrv.downloadTraffic += <-download
				unlock(tcpSrv.dataMutex)

				remote.server.Close()
				lock(tcpSrv.mapMutex)
				delete(conn.remoteConn, host)
				unlock(tcpSrv.mapMutex)

				wg.Done()
				close(processDone)
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
	type accepted struct {
		conn net.Conn
		err  error
	}
	for {
		c := make(chan accepted, 1)
		go func() {
			glog.V(5).Infof("wait for accept on %v \r\n", port)
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

			glog.V(5).Infof("accept remote client: %v\n", reqAddr)
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
