package main

import (
	"gofreezer/examples/dynamodb/app"
	"gofreezer/examples/dynamodb/app/options"
	"gofreezer/pkg/util/flag"

	"github.com/golang/glog"

	"github.com/spf13/pflag"
)

func main() {
	serverRunOptions := options.NewServerRunOptions()

	serverRunOptions.AddServerRunFlags(pflag.CommandLine)
	flag.InitFlags()

	if err := app.Run(serverRunOptions); err != nil {
		glog.Fatalf("Error in bringing up the server: %v", err)
	}
}
