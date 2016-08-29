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

func readReqData(conn net.Conn, chunkId int, iv []byte) ([]byte, error) {
	const (
		dataLenLen   = 2
		hmacSha1Len  = 10
		dataStartIdx = dataLenLen + hmacSha1Len
	)

	buffer := make([]byte, dataLenLen+hmacSha1Len)
	if n, err := io.ReadFull(conn, buffer[:]); err != nil {
		glog.Errorf("conn=%v read header error n=%v: %v", conn, n, err)
		return nil, err
	}

	dataLen := binary.BigEndian.Uint16(buffer[:dataLenLen])
	expectedHmacSha1 := buffer[dataLenLen:dataStartIdx]

	dataBuf := make([]byte, dataLen)
	if n, err := io.ReadFull(conn, dataBuf); err != nil {
		glog.V(5).Infof("conn=%p  read data error n=%v: %v", conn, n, err)
		return nil, err
	}

	chunkIdBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(chunkIdBytes, uint32(chunkId))
	actualHmacSha1 := util.HmacSha1(append(iv, chunkIdBytes...), dataBuf)
	if !bytes.Equal(expectedHmacSha1, actualHmacSha1) {
		glog.Errorf("conn=%v read data hmac-sha1 mismatch, iv=%v chunkId=%v src=%v dst=%v len=%v expeced=%v actual=%v",
			conn, iv, chunkId, conn.LocalAddr(), conn.RemoteAddr(), dataLen, expectedHmacSha1, actualHmacSha1)
		err := fmt.Errorf("not auth request")
		return nil, err
	}

	return dataBuf, nil
}

func handleData(ctx context.Context, src, dst net.Conn, iv []byte) <-chan *result {
	i := 0
	var upload int64
	var download int64

	rst := make(chan *result)
	for {
		timeout := ctx.Value(timeoutKey)

		//handle client request
		recvClient := make(chan *receive, 1)
		go func() {
			var err error
			setReadTimeout(src, timeout.(time.Duration))

			buf, err := readReqData(src, i, iv)
			if err != nil {
				glog.V(5).Infof("conn=%p #%v read data error n=%v: %v", src, i, err)
			} else {
				i++
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
				glog.V(5).Infof("conn=%v  read %v data error  %v \r\n", src.RemoteAddr().String(), n, err)
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
			return rst
		case recvData := <-recvClient:
			if recvData.err != nil {
				rst <- &result{
					unloadTraffic:   upload,
					downloadTraffic: download,
					err:             recvData.err,
				}
				return rst
			} else {
				var err error
				var writeLen int

				if writeLen, err = dst.Write(recvData.buffer); err != nil {
					glog.V(5).Infof("conn=%p #%v write data error n=%v: %v", dst, i, writeLen, err)
				} else {
					upload += int64(writeLen)
				}
			}
		case recvData := <-recvRemote:
			if recvData.err != nil {
				rst <- &result{
					unloadTraffic:   upload,
					downloadTraffic: download,
					err:             recvData.err,
				}
				return rst
			} else {
				_, err := dst.Write(recvData.buffer[:])
				if err != nil {
					glog.Errorf("write err:%v", err)
				} else {
					download += int64(len(recvData.buffer))
				}

			}
		}
	}
}

func process(ctx context.Context, iv []byte, client, remote net.Conn) (<-chan int64, <-chan int64) {

	upload := make(chan int64)
	download := make(chan int64)

	for {

		//add context if timeout we assume this connection not read data anymore
		//need to close
		timeout := time.Minute * 30
		ctx = context.WithValue(ctx, timeoutKey, timeout)

		reqRst := make(chan *result, 1)
		go func() {
			rst := handleData(ctx, client, remote, iv)
			rstInfo := <-rst
			reqRst <- &result{
				unloadTraffic:   rstInfo.unloadTraffic,
				downloadTraffic: rstInfo.downloadTraffic,
				err:             rstInfo.err,
			}

		}()

		select {
		case <-ctx.Done():
			return upload, download
		case reqResult := <-reqRst:
			upload <- reqResult.unloadTraffic
			download <- reqResult.downloadTraffic
			return upload, download
		}
	}
}
