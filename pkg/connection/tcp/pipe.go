package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"golang.org/x/net/context"

	"shadowsocks-go/pkg/util"

	"github.com/golang/glog"
)

type receive struct {
	buffer []byte
	err    error
}

type result struct {
	unloadTraffic   int64
	downloadTraffic int64
	err             error
}

func WriteToClient(client net.Conn, remote *remoteConnHelper, data []byte) (int, error) {
	originData := data
	// if remote.oneTimeAuth {
	// 	chunkID := remote.increaseChunkID()
	// 	originData = util.OtaReqChunkAuth(remote.iv, chunkID, data)
	// }

	ivLen := remote.cryp.GetIVLen()
	encryptBuf := make([]byte, len(originData)+ivLen)

	err := remote.cryp.Encrypt(encryptBuf[ivLen:], originData)
	if err != nil {
		return 0, err
	}

	glog.V(5).Infof("write to client iv : \r\n%s\r\n data: \r\n%s\r\n", util.DumpHex(remote.iv), util.DumpHex(originData[:]))

	var writeLen int
	if err == nil {
		copy(encryptBuf[0:ivLen], remote.iv)
		writeLen, err = client.Write(encryptBuf)
	}

	return writeLen, err
}

func recv(conn net.Conn, remote *remoteConnHelper, readLen int) ([]byte, error) {
	readBuf := make([]byte, readLen)
	decryptBuf := make([]byte, readLen)

	if n, err := io.ReadFull(conn, readBuf[:readLen]); err != nil {
		glog.Errorf("conn=%v read header error n=%v: %v", conn.RemoteAddr().String(), n, err)
		return nil, err
	}
	glog.Infof("Get Request encrypt Data \r\n%s\r\n", util.DumpHex(readBuf[:]))
	remote.cryp.Decrypt(decryptBuf[0:readLen], readBuf[0:readLen])
	glog.Infof("Get Request plainText Data \r\n%s\r\n", util.DumpHex(decryptBuf[:]))

	return decryptBuf[:readLen], nil
}

func readReqData(conn net.Conn, remote *remoteConnHelper) ([]byte, error) {
	const (
		dataLenLen   = 2
		hmacSha1Len  = 10
		dataStartIdx = dataLenLen + hmacSha1Len
	)

	readLen := dataLenLen + hmacSha1Len
	buffer, err := recv(conn, remote, readLen)
	if err != nil {
		return nil, err
	}

	dataLen := binary.BigEndian.Uint16(buffer[:dataLenLen])
	expectedHmacSha1 := buffer[dataLenLen:dataStartIdx]

	readLen = int(dataLen)
	dataBuf, err := recv(conn, remote, readLen)
	if err != nil {
		return nil, err
	}

	chunkIdBytes := make([]byte, 4)
	chunkID := remote.increaseChunkID()
	binary.BigEndian.PutUint32(chunkIdBytes, chunkID)
	actualHmacSha1 := util.HmacSha1(append(remote.iv, chunkIdBytes...), dataBuf)
	if !bytes.Equal(expectedHmacSha1, actualHmacSha1) {
		glog.Errorf("conn=%v read data hmac-sha1 mismatch, iv=%v chunkId=%v src=%v dst=%v len=%v expeced=%v actual=%v",
			conn, remote.iv, chunkID, conn.LocalAddr(), conn.RemoteAddr(), dataLen, expectedHmacSha1, actualHmacSha1)
		err := fmt.Errorf("not auth request")
		return nil, err
	}

	return dataBuf, nil
}

func handleData(ctx context.Context, src net.Conn, remote *remoteConnHelper) <-chan *result {
	var upload int64
	var download int64
	timeout := ctx.Value(timeoutKey)

	glog.V(5).Infof("timeout value %v\r\n", timeout)

	rst := make(chan *result, 1)
	for {
		dst := remote.server

		//handle client request
		recvClient := make(chan *receive, 1)
		go func() {
			var err error
			//setReadTimeout(src, timeout.(time.Duration))

			buf, err := readReqData(src, remote)
			if err != nil {
				glog.V(5).Infof("conn=%v read data error n=%v: %v", src.RemoteAddr().String(), remote.chunkID, err)
			}

			recvClient <- &receive{buffer: buf, err: err}
		}()

		//handle remote server
		recvRemote := make(chan *receive, 1)
		go func() {
			readBuf := make([]byte, 4108)

			setReadTimeout(dst, timeout.(time.Duration))

			n, err := dst.Read(readBuf)
			if err != nil {
				glog.V(5).Infof("dst conn=%v  read %v data error  %v \r\n", dst.RemoteAddr().String(), n, err)
			}
			recvRemote <- &receive{buffer: readBuf[0:n], err: err}

		}()

		select {
		case <-ctx.Done():
			rst <- &result{
				unloadTraffic:   upload,
				downloadTraffic: download,
				err:             nil,
			}
			glog.V(5).Infof("handle %s<->%s data will be done\n", src.RemoteAddr().String(), dst.RemoteAddr().String())
			return rst
		case recvData := <-recvClient:
			if recvData.err != nil {
				glog.V(5).Infof("handle %s<->%s data occur error %v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), recvData.err)
				rst <- &result{
					unloadTraffic:   upload,
					downloadTraffic: download,
					err:             recvData.err,
				}
				glog.V(5).Infof("handle %s<->%s data occur error done%v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), recvData.err)
				return rst
			} else {
				var err error
				var writeLen int

				glog.V(5).Infof("handle %s<->%s data write buffer to remote \r\n%s\r\n",
					src.RemoteAddr().String(),
					dst.RemoteAddr().String(), util.DumpHex(recvData.buffer[:]))
				if writeLen, err = dst.Write(recvData.buffer); err != nil {
					glog.V(5).Infof("conn=%s  write data error n=%v: %v", dst.RemoteAddr().String(), writeLen, err)
				} else {
					upload += int64(writeLen)
				}
			}
		case recvData := <-recvRemote:
			if recvData.err != nil {
				glog.V(5).Infof("handle %s<->%s data occur error %v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), recvData.err)
				rst <- &result{
					unloadTraffic:   upload,
					downloadTraffic: download,
					err:             recvData.err,
				}
				return rst
			} else {
				glog.V(5).Infof("handle %s<->%s data write buffer to client \r\n%s\r\n",
					src.RemoteAddr().String(), dst.RemoteAddr().String(), util.DumpHex(recvData.buffer[:]))

				n, err := WriteToClient(src, remote, recvData.buffer[:])
				if err != nil {
					glog.Errorf("write err:%v", err)
				} else {
					download += int64(n)
				}

			}
		}
	}
}

func process(ctx context.Context, client net.Conn, remote *remoteConnHelper) (<-chan int64, <-chan int64) {

	upload := make(chan int64, 1)
	download := make(chan int64, 1)

	//add context if timeout we assume this connection not read data anymore
	//need to close
	timeout := time.Duration(1000) * time.Millisecond
	ctx = context.WithValue(ctx, timeoutKey, timeout)

	for {

		reqRst := make(chan *result, 1)
		go func() {
			glog.V(5).Infof("handle %s<->%s process done in normal status 1\n", client.RemoteAddr().String(), remote.server.RemoteAddr().String())

			rstInfo := <-handleData(ctx, client, remote)
			reqRst <- &result{
				unloadTraffic:   rstInfo.unloadTraffic,
				downloadTraffic: rstInfo.downloadTraffic,
				err:             rstInfo.err,
			}
			glog.V(5).Infof("handle %s<->%s process done in normal status\n", client.RemoteAddr().String(), remote.server.RemoteAddr().String())
		}()

		select {
		case <-ctx.Done():
			glog.V(5).Infof("handle %s<->%s process will be done\n", client.RemoteAddr().String(), remote.server.RemoteAddr().String())
			return upload, download
		case reqResult := <-reqRst:
			upload <- reqResult.unloadTraffic
			download <- reqResult.downloadTraffic
			glog.V(5).Infof("handle %s<->%s process done in normal status\n", client.RemoteAddr().String(), remote.server.RemoteAddr().String())
			return upload, download
		}
	}
}
