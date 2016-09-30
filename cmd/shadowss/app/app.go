package app

import (
	"runtime"

	"shadowsocks-go/cmd/shadowss/app/options"
	"shadowsocks-go/pkg/multiuser"
	"shadowsocks-go/pkg/proxyserver"

	"github.com/golang/glog"
)

//Run start a api server
func Run(options *options.ServerOption) error {

	if options.CpuCoreNum > 0 {
		runtime.GOMAXPROCS(options.CpuCoreNum)
	}

	pxy := proxyserver.NewServers(options.EnableUDPRelay)
	err := options.LoadConfigFile()
	if err != nil {
		glog.Warning("load user configure error:\r\n", err)
	} else {
		pxy.Start()
	}

	//multiuser config
	multiuser.InitSchedule(options.EtcdStorageConfig, pxy)

	return nil
}
