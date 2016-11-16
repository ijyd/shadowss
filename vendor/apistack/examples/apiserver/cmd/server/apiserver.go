package main

import (
	"fmt"
	"os"

	"gofreezer/pkg/util/flag"
	"gofreezer/pkg/util/logs"

	"apistack/examples/apiserver/cmd/server/app"
	"apistack/examples/apiserver/cmd/server/app/options"

	"github.com/spf13/pflag"
)

func main() {

	serverRunOptions := options.NewServerOption()
	// Parse command line flags.
	serverRunOptions.AddFlags(pflag.CommandLine)
	flag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := app.Run(serverRunOptions); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
