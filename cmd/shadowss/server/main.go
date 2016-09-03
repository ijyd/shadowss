package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"shadowsocks-go/cmd/shadowss"
	"shadowsocks-go/pkg/proxyserver"
	"shadowsocks-go/pkg/users"
	"shadowsocks-go/pkg/util/flag"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func waitSignal() {
	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
			//updatePasswd()
		} else {
			// is this going to happen?
			log.Printf("caught signal %v, exit", sig)
			os.Exit(0)
		}
	}
}

func main() {
	serverRunOptions := shadowss.NewServerOption()
	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)

	users := users.NewUsers()
	users.AddFlags(pflag.CommandLine)

	flag.InitFlags()

	if serverRunOptions.CpuCoreNum > 0 {
		runtime.GOMAXPROCS(serverRunOptions.CpuCoreNum)
	}

	glog.V(5).Infoln("tset")
	err := serverRunOptions.LoadConfigFile()
	if err != nil {
		glog.Fatalln("load user configure error:\r\n", err)
	}

	pxy := proxyserver.NewServers(serverRunOptions.EnableUDPRelay)
	pxy.Start()

	users.CreateUsersSync(pxy)
	waitSignal()
}
