package app

import (
	"runtime"

	"shadowss/cmd/shadowss/app/options"
	"shadowss/pkg/multiuser"
	"shadowss/pkg/proxyserver"

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
	multiuser.InitSchedule(pxy, options.URL)

	return nil
}
