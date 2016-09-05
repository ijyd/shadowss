package main

import (
	"fmt"
	"os"

	"shadowsocks-go/cmd/shadowssctl/app"
	"shadowsocks-go/cmd/shadowssctl/app/options"
	"shadowsocks-go/pkg/backend"
	"shadowsocks-go/pkg/util/flag"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func main() {

	serverRunOptions := options.NewServerOption()
	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)

	be := backend.NewBackend()
	be.AddFlags(pflag.CommandLine)

	flag.InitFlags()

	glog.V(5).Infoln("tset")
	if err := app.Run(serverRunOptions, be); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
