package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"

	shadowss "shadowsocks-go/cmd/shadowss/app"
	"shadowsocks-go/pkg/backend"
	"shadowsocks-go/pkg/proxyserver"
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

	be := backend.NewBackend()
	be.AddFlags(pflag.CommandLine)

	flag.InitFlags()

	if serverRunOptions.CpuCoreNum > 0 {
		runtime.GOMAXPROCS(serverRunOptions.CpuCoreNum)
	}

	err := serverRunOptions.LoadConfigFile()
	if err != nil {
		glog.Fatalln("load user configure error:\r\n", err)
	}

	pxy := proxyserver.NewServers(serverRunOptions.EnableUDPRelay)
	pxy.Start()

	be.CreateUsersSync(pxy)
	waitSignal()

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
}
