package options

import (
	"shadowss/pkg/config"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

type ServerOption struct {
	ConfigFile         string
	CpuCoreNum         int
	EnableUDPRelay     bool
	URL                string
	MaxTCPConnPerPort  int
	APIProxyListenPort int
	APIProxyPassword   string
	APIServerToken     string
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		ConfigFile:     string(""),
		CpuCoreNum:     1,
		EnableUDPRelay: false,
	}
}

func (s *ServerOption) LoadConfigFile() error {
	glog.V(5).Infof("Parse %s file\r\n", s.ConfigFile)
	err := config.ServerCfg.Parse(s.ConfigFile)
	if err != nil {
		glog.Errorf("error reading %s: %v\n", s.ConfigFile, err)
		return err
	}

	return nil
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.ConfigFile, "config-file", s.ConfigFile, ""+
		"specify a configure file for server run. ")

	fs.IntVar(&s.CpuCoreNum, "cpu-core-num", s.CpuCoreNum, ""+
		"specify how many cpu core will be alloc for program")

	fs.BoolVar(&s.EnableUDPRelay, "enable-udp-relay", s.EnableUDPRelay, ""+
		"enable udp relay")

	fs.StringVar(&s.URL, "apiserver-url", s.URL, ""+
		"specify a api server url. ")

	fs.IntVar(&s.MaxTCPConnPerPort, "max-tcp-conn-per-port", 300, ""+
		"specify how many tcp connection per port")

	fs.IntVar(&s.APIProxyListenPort, "api-proxy-port", 12345, ""+
		"specify a port for api proxy")

	fs.StringVar(&s.APIProxyPassword, "api-proxy-password", s.APIProxyPassword, ""+
		"specify password for api proxy. ")

	fs.StringVar(&s.APIServerToken, "apiserver-token", s.APIServerToken, ""+
		"specify a api server token. ")

}
