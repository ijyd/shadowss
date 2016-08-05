package shadowss

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	connection "shadowsocks-go/pkg/connection"
	encrypt "shadowsocks-go/pkg/connection"
	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

const (
	idType  = 0 // address type index
	idIP0   = 1 // ip addres start index
	idDmLen = 1 // domain address length index
	idDm0   = 2 // domain address start index

	typeIPv4 = 1 // type is ipv4 address
	typeDm   = 3 // type is domain address
	typeIPv6 = 4 // type is ipv6 address

	lenIPv4     = net.IPv4len + 2 // ipv4 + 2port
	lenIPv6     = net.IPv6len + 2 // ipv6 + 2port
	lenDmBase   = 2               // 1addrLen + 2port, plus addrLen
	lenHmacSha1 = 10
)

// PortListener is a listener
type PortListener struct {
	password string
	listener net.Listener
}

func getRequest(conn *connection.Conn, auth bool, timeout time.Duration) (host string, ota bool, err error) {
	glog.V(5).Infoln("getRequest from remote34342")

	connection.SetReadTimeout(conn, timeout)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 256(max length address) + 2(port) + 10(hmac-sha1)
	buf := make([]byte, 270)
	// read till we get possible domain length field
	if _, err = io.ReadFull(conn, buf[:idType+1]); err != nil {
		glog.Errorln("read buffer from remote connection error:", err)
		return
	}

	// if _, err = io.ReadFull(conn, buf[:64]); err != nil {
	// 	glog.Errorln("read buffer from remote connection error:", err)
	// 	return
	// }
	// glog.V(5).Infof("Got a Request string is byte:%v string:%v)\r\n", buf, string(buf))
	// return

	var reqStart, reqEnd int
	addrType := buf[idType]
	switch addrType & connection.AddrMask {
	case typeIPv4:
		reqStart, reqEnd = idIP0, idIP0+lenIPv4
	case typeIPv6:
		reqStart, reqEnd = idIP0, idIP0+lenIPv6
	case typeDm:
		glog.V(5).Infoln("Got a Domain Addr Type, read start(%v) end(%v)\r\n", idType+1, idDmLen+1)
		if _, err = io.ReadFull(conn, buf[idType+1:idDmLen+1]); err != nil {
			glog.Errorf("Read from remote err:%v\r\n", err)
			return
		}
		reqStart, reqEnd = idDm0, int(idDm0+buf[idDmLen]+lenDmBase)
	default:
		err = fmt.Errorf("addr type %d not supported", addrType&connection.AddrMask)
		return
	}

	if _, err = io.ReadFull(conn, buf[reqStart:reqEnd]); err != nil {
		return
	}

	glog.V(5).Infof("Got string from remote %v \r\n", buf[reqStart:reqEnd])
	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch addrType & connection.AddrMask {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
	}
	glog.V(5).Infof("Got host from remote: %v\r\n", host)
	// parse port
	port := binary.BigEndian.Uint16(buf[reqEnd-2 : reqEnd])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	// if specified one time auth enabled, we should verify this
	if auth || addrType&connection.OneTimeAuthMask > 0 {
		ota = true
		if _, err = io.ReadFull(conn, buf[reqEnd:reqEnd+lenHmacSha1]); err != nil {
			return
		}
		iv := conn.GetIv()
		key := conn.GetKey()
		actualHmacSha1Buf := util.HmacSha1(append(iv, key...), buf[:reqEnd])
		if !bytes.Equal(buf[reqEnd:reqEnd+lenHmacSha1], actualHmacSha1Buf) {
			err = fmt.Errorf("verify one time auth failed, iv=%v key=%v data=%v", iv, key, buf[:reqEnd])
			return
		}
	}
	return
}

const logCntDelta = 100

var connCnt int
var nextLogConnCnt int = logCntDelta

type isClosed struct {
	isClosed bool
}

func handleConnection(conn *connection.Conn, auth bool, timeout time.Duration) {
	var host string

	connCnt++ // this maybe not accurate, but should be enough
	if connCnt-nextLogConnCnt >= 0 {
		// XXX There's no xadd in the atomic package, so it's difficult to log
		// the message only once with low cost. Also note nextLogConnCnt maybe
		// added twice for current peak connection number level.
		log.Printf("Number of client connections reaches %d\n", nextLogConnCnt)
		nextLogConnCnt += logCntDelta
	}

	// function arguments are always evaluated, so surround debug statement
	// with if statement
	glog.V(5).Infof("new client %s->%s\n", conn.RemoteAddr().String(), conn.LocalAddr())

	closed := false
	defer func() {
		glog.V(5).Infof("closed pipe %s<->%s\n", conn.RemoteAddr(), host)
		connCnt--
		if !closed {
			conn.Close()
		}
	}()

	host, ota, err := getRequest(conn, auth, timeout)
	if err != nil {
		glog.Errorf("error getting request %v<->%v err:%v", conn.RemoteAddr(), conn.LocalAddr(), err)
		return
	}
	glog.V(5).Infof("connection host:%v ota:%v \r\n", host, ota)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			glog.Errorf("dial error:%v\r\n", err)
		} else {
			glog.Errorf(" connecting to:%v occur err:%v", host, err)
		}
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()

	glog.V(5).Infof("piping %s<->%s ota=%v connOta=%v", conn.RemoteAddr(), host, ota, conn.IsOta())

	if ota {
		go connection.PipeThenCloseOta(conn, remote, timeout)
	} else {
		go connection.PipeThenClose(conn, remote, timeout)
	}
	connection.PipeThenClose(remote, conn, timeout)
	closed = true
	return
}

// func updatePasswd() {
// 	log.Println("updating password")
// 	newconfig, err := ss.ParseConfig(configFile)
// 	if err != nil {
// 		log.Printf("error parsing config file %s to update password: %v\n", configFile, err)
// 		return
// 	}
// 	oldconfig := config
// 	config = newconfig
//
// 	if err = unifyPortPassword(config); err != nil {
// 		return
// 	}
// 	for port, passwd := range config.PortPassword {
// 		passwdManager.updatePortPasswd(port, passwd, config.Auth)
// 		if oldconfig.PortPassword != nil {
// 			delete(oldconfig.PortPassword, port)
// 		}
// 	}
// 	// port password still left in the old config should be closed
// 	for port, _ := range oldconfig.PortPassword {
// 		log.Printf("closing port %s as it's deleted\n", port)
// 		passwdManager.del(port)
// 	}
// 	log.Println("password updated")
// }

func Run(password, method string, port int, auth bool, timeout time.Duration) {
	portStr := strconv.Itoa(port)
	ln, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		log.Printf("error listening port %v: %v\n", port, err)
		os.Exit(1)
	}
	passwdManager.add(password, port, ln)
	var cipher *encrypt.Cipher
	log.Printf("server listening port %v ...\n", port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// listener maybe closed to update password
			glog.V(5).Infof("accept error: %v\n", err)
			return
		}
		// Creating cipher upon first connection.
		if cipher == nil {
			glog.Infof("creating cipher for port :%v and method: %v\r\n", port, method)
			cipher, err = encrypt.NewCipher(method, password)
			if err != nil {
				glog.Errorf("Error generating cipher for port: %s %v\n", port, err)
				conn.Close()
				continue
			}
		}
		go handleConnection(connection.NewConn(conn, cipher.Copy()), auth, timeout)
	}
}
