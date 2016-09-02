package tcp

import (
	"io"
	"net"
	"time"

	"shadowsocks-go/pkg/connection/tcp/ssclient"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func SetReadTimeout(c net.Conn, timeout time.Duration) {
	if timeout != 0 {
		c.SetReadDeadline(time.Now().Add(timeout))
	}
}

// PipeData copies data from src to dst.
func PipeData(ctx context.Context, src *ssclient.Client, dst net.Conn, timeout time.Duration) (int64, int64) {
	var upload int64
	var download int64

	result := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			result <- err
		}()

		for {
			SetReadTimeout(src, timeout)

			dataBuf, parseErr := src.ParseReqData()
			err = parseErr
			if err != nil && err == io.EOF {
				break
			} else {
				var writeLen int
				if writeLen, err = dst.Write(dataBuf); err != nil {
					glog.V(5).Infof("conn=%p  write data error n=%v: %v", dst, writeLen, err)
					break
				} else {
					upload += int64(writeLen)
				}
			}
		}
	}()

	go func() {
		var err error
		defer func() {
			result <- err
		}()

		buf := make([]byte, 5000)
		for {
			SetReadTimeout(dst, timeout)
			var readLen int
			readLen, err = dst.Read(buf)
			if readLen > 0 {
				if _, err = src.Write(buf[0:readLen]); err != nil {
					glog.Errorf("write err:%v", err)
					break
				}
				download += int64(readLen)
			}
			if err != nil && err == io.EOF {
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			glog.V(5).Infof("handle %s<->%s data will be done\n", src.RemoteAddr().String(), dst.RemoteAddr().String())
			return upload, download
		case rst := <-result:
			glog.Infof("handle %s<->%s data will be done by break %v\n", src.RemoteAddr().String(), dst.RemoteAddr().String(), rst)
			return upload, download
		default:
			time.Sleep(1 * time.Second)
		}

	}

}
