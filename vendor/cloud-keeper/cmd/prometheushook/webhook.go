package main

import (
	"fmt"
	"golib/pkg/util/flag"
	"os"

	"cloud-keeper/cmd/prometheushook/app"
	"cloud-keeper/cmd/prometheushook/app/options"

	"github.com/spf13/pflag"
)

func main() {

	serverRunOptions := options.NewServerOption()
	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)
	flag.InitFlags()

	// //glog.Infof("tst here\r\n")
	if err := app.Run(serverRunOptions); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
