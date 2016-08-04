package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"shadowsocks-go/shadowsocks/util"
	"shadowsocks-go/shadowsocks/util/flag"

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
	serverRunOptions := newServerOption()

	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)
	flag.InitFlags()

	if serverRunOptions.isPrintVersion {
		util.PrintVersion()
		os.Exit(0)
	}

	if serverRunOptions.cpuCoreNum > 0 {
		runtime.GOMAXPROCS(serverRunOptions.cpuCoreNum)
	}

	config, err := serverRunOptions.loadUserConfig()
	if err != nil {
		glog.Fatalln("load user configure error:\r\n", err)
	}

	for port, password := range config.PortPassword {
		glog.Infof("listen pair(%v:%v) with auth:%v", port, password, config.Auth)
		go run(port, password, config.Method, config.Auth)
		if udp {
			go runUDP(port, password, config.Method)
		}
	}

	waitSignal()
}
