package main

import (
	"shadowsocks-go/pkg/config"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

type serverOption struct {
	configFile     string
	cpuCoreNum     int
	enableUDPRelay bool
}

func newServerOption() *serverOption {
	return &serverOption{
		configFile:     string(""),
		cpuCoreNum:     1,
		enableUDPRelay: false,
	}
}

func (s *serverOption) loadConfigFile() error {
	err := config.ServerCfg.Parse(s.configFile)
	if err != nil {
		glog.Fatalf("error reading %s: %v\n", s.configFile, err)
		return err
	}

	return nil
}

func (s *serverOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.configFile, "config-file", s.configFile, ""+
		"specify a configure file for server run. ")

	fs.IntVar(&s.cpuCoreNum, "cpu-core-num", s.cpuCoreNum, ""+
		"specify how many cpu core will be alloc for program")

	fs.BoolVar(&s.enableUDPRelay, "enable-udp-relay", s.enableUDPRelay, ""+
		"enable udp relay")
}
