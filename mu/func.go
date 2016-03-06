/**
@author orvice   https://github.com/orvice/shadowsocks-go
@author Lupino   https://github.com/Lupino/shadowsocks-auth
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cyfdecyf/leakybuf"
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var ssdebug ss.DebugLog

func getRequest(conn *ss.Conn) (host string, extra []byte, err error) {
	const (
		idType  = 0 // address type index
		idIP0   = 1 // ip addres start index
		idDmLen = 1 // domain address length index
		idDm0   = 2 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4   = 1 + net.IPv4len + 2 // 1addrType + ipv4 + 2port
		lenIPv6   = 1 + net.IPv6len + 2 // 1addrType + ipv6 + 2port
		lenDmBase = 1 + 1 + 2           // 1addrType + 1addrLen + 2port, plus addrLen
	)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port)
	buf := make([]byte, 260)
	var n int
	// read till we get possible domain length field
	ss.SetReadTimeout(conn)
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen]) + lenDmBase
	default:
		err = fmt.Errorf("addr type %d not supported", buf[idType])
		return
	}

	if n < reqLen { // rare case
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else if n > reqLen {
		// it's possible to read more than just the request head
		extra = buf[reqLen:n]
	}

	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch buf[idType] {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
	}
	// parse port
	port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

const logCntDelta = 100

var connCnt int
var nextLogConnCnt int = logCntDelta

func handleConnection(user user.User, conn *ss.Conn) {
	var host string
	var size = 0
	var raw_req_header, raw_res_header []byte
	var is_http = false
	var res_size = 0
	var req_chan = make(chan []byte)
	connCnt++ // this maybe not accurate, but should be enough
	if connCnt-nextLogConnCnt >= 0 {
		// XXX There's no xadd in the atomic package, so it's difficult to log
		// the message only once with low cost. Also note nextLogConnCnt maybe
		// added twice for current peak connection number level.
		Log.Debug("Number of client connections reaches %d\n", nextLogConnCnt)
		nextLogConnCnt += logCntDelta
	}

	// function arguments are always evaluated, so surround debug statement
	// with if statement
	Log.Debug(fmt.Sprintf("new client %s->%s\n", conn.RemoteAddr().String(), conn.LocalAddr()))
	closed := false
	defer func() {
		if ssdebug {
			Log.Debug(fmt.Sprintf("closed pipe %s<->%s\n", conn.RemoteAddr(), host))
		}
		connCnt--
		if !closed {
			conn.Close()
		}
	}()

	host, extra, err := getRequest(conn)
	if err != nil {
		Log.Error("error getting request", conn.RemoteAddr(), conn.LocalAddr(), err)
		return
	}
	Log.Info(fmt.Sprintf("[port-%d]connecting %s ", user.GetPort(), host))
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			Log.Error("dial error:", err)
		} else {
			Log.Error("error connecting to:", host, err)
		}
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()

	defer func() {
		if is_http {
			tmp_req_header := <-req_chan
			buffer := bytes.NewBuffer(raw_req_header)
			buffer.Write(tmp_req_header)
			raw_req_header = buffer.Bytes()
		}
		showConn(raw_req_header, raw_res_header, host, user, size, is_http)
		close(req_chan)
		if !closed {
			remote.Close()
		}
	}()

	// write extra bytes read from
	if extra != nil {
		// debug.Println("getRequest read extra data, writing to remote, len", len(extra))
		is_http, extra, _ = checkHttp(extra, conn)
		if strings.HasSuffix(host, ":80") {
			is_http = true
		}
		raw_req_header = extra
		res_size, err = remote.Write(extra)
		// size, err := remote.Write(extra)
		if err != nil {
			Log.Error("write request extra error:", err)
			return
		}
		// debug conn info
		Log.Debug(fmt.Sprintf("%d conn debug:  local addr: %s | remote addr: %s network: %s ", user.GetPort(),
			conn.LocalAddr().String(), conn.RemoteAddr().String(), conn.RemoteAddr().Network()))
		err = storage.IncrSize(user, res_size)
		if err != nil {
			Log.Error(err)
			return
		}
		err = storage.MarkUserOnline(user)
		if err != nil {
			Log.Error(err)
			return
		}
		Log.Debug(fmt.Sprintf("[port-%d] store size: %d", user.GetPort(), res_size))
	}
	Log.Debug(fmt.Sprintf("piping %s<->%s", conn.RemoteAddr(), host))
	/**
	go ss.PipeThenClose(conn, remote)
	ss.PipeThenClose(remote, conn)
	closed = true
	return
	**/
	go func() {
		_, raw_header := PipeThenClose(conn, remote, is_http, false, host, user)
		if is_http {
			req_chan <- raw_header
		}
	}()

	res_size, raw_res_header = PipeThenClose(remote, conn, is_http, true, host, user)
	size += res_size
	closed = true
	return
}

type PortListener struct {
	password string
	listener net.Listener
}

type PasswdManager struct {
	sync.Mutex
	portListener map[string]*PortListener
}

func (pm *PasswdManager) add(port, password string, listener net.Listener) {
	pm.Lock()
	pm.portListener[port] = &PortListener{password, listener}
	pm.Unlock()
}

func (pm *PasswdManager) get(port string) (pl *PortListener, ok bool) {
	pm.Lock()
	pl, ok = pm.portListener[port]
	pm.Unlock()
	return
}

func (pm *PasswdManager) del(port string) {
	pl, ok := pm.get(port)
	if !ok {
		return
	}
	pl.listener.Close()
	pm.Lock()
	delete(pm.portListener, port)
	pm.Unlock()
}

var passwdManager = PasswdManager{portListener: map[string]*PortListener{}}

func waitSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
		} else {
			// is this going to happen?
			Log.Printf("caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}

func runWithCustomMethod(user user.User) {
	// port, password string, Cipher *ss.Cipher
	port := strconv.Itoa(user.GetPort())
	password := user.GetPasswd()
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		Log.Error(fmt.Sprintf("error listening port %v: %v\n", port, err))
		os.Exit(1)
	}
	passwdManager.add(port, password, ln)
	cipher, err := user.GetCipher()
	if err != nil {
		return
	}
	Log.Info(fmt.Sprintf("server listening port %v ...\n", port))
	for {
		conn, err := ln.Accept()
		if err != nil {
			// listener maybe closed to update password
			Log.Debug(fmt.Sprintf("accept error: %v\n", err))
			return
		}
		// Creating cipher upon first connection.
		if cipher == nil {
			Log.Debug("creating cipher for port:", port)
			cipher, err = ss.NewCipher(user.GetMethod(), password)
			if err != nil {
				Log.Error(fmt.Sprintf("Error generating cipher for port: %s %v\n", port, err))
				conn.Close()
				continue
			}
		}
		go handleConnection(user, ss.NewConn(conn, cipher.Copy()))
	}
}

const bufSize = 4096
const nBuf = 2048

func PipeThenClose(src, dst net.Conn, is_http bool, is_res bool, host string, user user.User) (total int, raw_header []byte) {
	var pipeBuf = leakybuf.NewLeakyBuf(nBuf, bufSize)
	defer dst.Close()
	buf := pipeBuf.Get()
	// defer pipeBuf.Put(buf)
	var buffer = bytes.NewBuffer(nil)
	var is_end = false
	var size int

	for {
		SetReadTimeout(src)
		n, err := src.Read(buf)
		// read may return EOF with n > 0
		// should always process n > 0 bytes before handling error
		if n > 0 {
			if is_http && !is_end {
				buffer.Write(buf)
				raw_header = buffer.Bytes()
				lines := bytes.SplitN(raw_header, []byte("\r\n\r\n"), 2)
				if len(lines) == 2 {
					is_end = true
				}
			}

			size, err = dst.Write(buf[0:n])
			if is_res {
				err = storage.IncrSize(user, size)
				if err != nil {
					Log.Error(err)
				}
				Log.Debug(fmt.Sprintf("[port-%d] store size: %d", user.GetPort(), size))
			}
			total += size
			if err != nil {
				Log.Debug("write:", err)
				break
			}
		}
		if err != nil || n == 0 {
			// Always "use of closed network connection", but no easy way to
			// identify this specific error. So just leave the error along for now.
			// More info here: https://code.google.com/p/go/issues/detail?id=4373
			break
		}
	}
	return
}

var readTimeout time.Duration

func SetReadTimeout(c net.Conn) {
	if readTimeout != 0 {
		c.SetReadDeadline(time.Now().Add(readTimeout))
	}
}

func showConn(raw_req_header, raw_res_header []byte, host string, user user.User, size int, is_http bool) {
	if size == 0 {
		Log.Error(fmt.Sprintf("[port-%d]  Error: request %s cancel", user.GetPort(), host))
		return
	}
	if is_http {
		req, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw_req_header)))
		if req == nil {
			lines := bytes.SplitN(raw_req_header, []byte(" "), 2)
			Log.Debug(fmt.Sprintf("%s http://%s/ \"Unknow\" HTTP/1.1 unknow user-port: %d size: %d\n", lines[0], host, user.GetPort(), size))
			return
		}
		res, _ := http.ReadResponse(bufio.NewReader(bytes.NewReader(raw_res_header)), req)
		statusCode := 200
		if res != nil {
			statusCode = res.StatusCode
		}
		Log.Debug(fmt.Sprintf("%s http://%s%s \"%s\" %s %d  user-port: %d  size: %d\n", req.Method, req.Host, req.URL.String(), req.Header.Get("user-agent"), req.Proto, statusCode, user.GetPort(), size))
	} else {
		Log.Debug(fmt.Sprintf("CONNECT %s \"NONE\" NONE NONE user-port: %d  size: %d\n", host, user.GetPort(), size))
	}

}

func checkHttp(extra []byte, conn *ss.Conn) (is_http bool, data []byte, err error) {
	var buf []byte
	var methods = []string{"GET", "HEAD", "POST", "PUT", "TRACE", "OPTIONS", "DELETE"}
	is_http = false
	if extra == nil || len(extra) < 10 {
		buf = make([]byte, 10)
		if _, err = io.ReadFull(conn, buf); err != nil {
			return
		}
	}

	if buf == nil {
		data = extra
	} else if extra == nil {
		data = buf
	} else {
		buffer := bytes.NewBuffer(extra)
		buffer.Write(buf)
		data = buffer.Bytes()
	}

	for _, method := range methods {
		if bytes.HasPrefix(data, []byte(method)) {
			is_http = true
			break
		}
	}
	return
}
