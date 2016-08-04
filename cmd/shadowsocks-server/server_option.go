package main

import (
	"errors"
	"os"
	"strconv"

	ss "shadowsocks-go/shadowsocks"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

var configFile string
var config *ss.Config

type serverOption struct {
	isPrintVersion  bool
	configFile      string
	userPassword    string
	userPort        string
	userConnTimeOut int
	userConnEcrypt  string
	cpuCoreNum      int
	enableDebug     bool
}

func newServerOption() *serverOption {
	return &serverOption{
		isPrintVersion:  false,
		configFile:      string(""),
		userPassword:    string("12345678"),
		userPort:        string(1080),
		userConnTimeOut: 60,
		userConnEcrypt:  string("aes-256-cfb"),
		cpuCoreNum:      1,
		enableDebug:     false,
	}
}

func enoughOptions(config *ss.Config) bool {
	return config.ServerPort != 0 && config.Password != ""
}

func unifyPortPassword(config *ss.Config) (err error) {
	if len(config.PortPassword) == 0 { // this handles both nil PortPassword and empty one
		if !enoughOptions(config) {
			glog.Fatalln("must specify both port and password")
			return errors.New("not enough options")
		}
		port := strconv.Itoa(config.ServerPort)
		config.PortPassword = map[string]string{port: config.Password}
	} else {
		if config.Password != "" || config.ServerPort != 0 {
			glog.Fatalln("given port_password, ignore server_port and password option")
		}
	}

	return nil
}

func (s *serverOption) loadUserConfig() (*ss.Config, error) {
	cmdConfig := ss.Config{
		Method: s.userConnEcrypt,
		PortPassword: map[string]string{
			s.userPort: s.userPassword,
		},
		Timeout: s.userConnTimeOut,
	}

	config, err := ss.ParseConfig(s.configFile)
	if err != nil {
		if !os.IsNotExist(err) {
			glog.Fatalf("error reading %s: %v\n", s.configFile, err)
			return nil, err
		}
		config = &cmdConfig
	} else {
		//force use command line configure
		ss.UpdateConfig(config, &cmdConfig)
	}

	if config.Method == "" {
		config.Method = "aes-256-cfb"
	}

	if err = ss.CheckCipherMethod(config.Method); err != nil {
		glog.Fatalf("check cipher configure error %v \r\n", err)
		return nil, err
	}

	if err = unifyPortPassword(config); err != nil {
		glog.Fatalf("check port and password error %v\r\n", err)
		return nil, err
	}

	return config, nil
}

func (s *serverOption) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.isPrintVersion, "print-version", s.isPrintVersion, ""+
		"only print version then shutdown")

	fs.StringVar(&s.configFile, "config-file", s.configFile, ""+
		"specify a configure file for server run. ")

	fs.StringVar(&s.userPassword, "user-conn-password", s.userPassword, ""+
		"specify a user password wait for client connection .")

	fs.StringVar(&s.userPort, "user-conn-port", s.userPort, ""+
		"specify a user port ,this options will used with user password")

	fs.IntVar(&s.userConnTimeOut, "user-conn-timeout", s.userConnTimeOut, ""+
		"specify user connection timeout options")

	fs.StringVar(&s.userConnEcrypt, "user-conn-encrypt", s.userConnEcrypt, ""+
		" .")

	fs.IntVar(&s.cpuCoreNum, "cpu-core-num", s.cpuCoreNum, ""+
		"specify how many cpu core will be alloc for program")

	fs.BoolVar(&s.enableDebug, "debug-mode", s.enableDebug, ""+
		"enable debug information")
}
