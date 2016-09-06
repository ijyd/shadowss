package options

import "github.com/spf13/pflag"

type ServerOption struct {
	Host        string
	Port        int32
	SwaggerPath string
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		Host: "",
		Port: 80,
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.Host, "host", s.Host, ""+
		"specific what host for api server")

	fs.Int32Var(&s.Port, "port", s.Port, ""+
		"specific what port for api server ")

	fs.StringVar(&s.SwaggerPath, "swagger-path", s.SwaggerPath, ""+
		"specific a path where found swagger index.html, if not will be disable swagger ui")
}
