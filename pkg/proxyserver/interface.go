package proxyserver

import (
	"shadowss/pkg/config"
	"time"
)

//ProxyServer implement a proxy server interface
type ProxyServer interface {
	Run()
	Stop()
	Compare(*config.ConnectionInfo) bool
	Traffic() (int64, int64)
	GetListenPort() int
	GetConfig() config.ConnectionInfo
	GetLastActiveTime() time.Time
}
