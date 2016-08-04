package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

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

func start() {
	for _, v := range config.ServerCfg.Clients {
		glog.V(5).Infof("listen pair(%v:%v) with auth:%v")
		go run(v.Password, v.EncryptMethod, v.Port, v.EnableOTA)
		if udp {
			go runUDP(v.Password, v.EncryptMethod, v.Port)
		}
	}
}

func main() {
	serverRunOptions := newServerOption()

	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)
	flag.InitFlags()

	if serverRunOptions.cpuCoreNum > 0 {
		runtime.GOMAXPROCS(serverRunOptions.cpuCoreNum)
	}

	err := serverRunOptions.loadConfigFile()
	if err != nil {
		glog.Fatalln("load user configure error:\r\n", err)
	}

	go start()

	waitSignal()
}
