package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"shadowsocks-go/cmd/shadowss"
	"shadowsocks-go/pkg/config"
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

func start(udp bool) {
	for _, v := range config.ServerCfg.Clients {
		glog.V(5).Infof("listen pair(%v:%v) with auth:%v", v.Port, v.Password, v.EnableOTA)
		go shadowss.Run(v.Password, v.EncryptMethod, v.Port, v.EnableOTA, time.Duration(v.Timeout)*time.Second)
		if udp {
			go shadowss.RunUDP(v.Password, v.EncryptMethod, v.Port, time.Duration(v.Timeout)*time.Second)
		}
	}
}

func main() {
	serverRunOptions := shadowss.NewServerOption()

	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)
	flag.InitFlags()

	if serverRunOptions.CpuCoreNum > 0 {
		runtime.GOMAXPROCS(serverRunOptions.CpuCoreNum)
	}

	glog.V(5).Infoln("tset")
	err := serverRunOptions.LoadConfigFile()
	if err != nil {
		glog.Fatalln("load user configure error:\r\n", err)
	}

	go start(serverRunOptions.EnableUDPRelay)

	waitSignal()
}
