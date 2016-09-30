package main

import (
	"golib/pkg/util/flag"

	"github.com/golang/glog"

	"github.com/spf13/pflag"
)

func main() {

	AddFlags(pflag.CommandLine)
	flag.InitFlags()

	glog.Infof("test log output \r\n")
}

func AddFlags(fs *pflag.FlagSet) {

	var test string
	fs.StringVar(&test, "host", test, ""+
		"specific what host for api server")

}
