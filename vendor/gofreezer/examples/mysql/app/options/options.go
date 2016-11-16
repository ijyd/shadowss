package options

import (
	storageoptions "gofreezer/pkg/genericstoragecodec/options"

	"github.com/spf13/pflag"
)

// ServerRunOptions contains the options while running a generic api server.
type ServerRunOptions struct {
	Host        string
	Port        int32
	SwaggerPath string

	StorageConfig *storageoptions.StorageOptions
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		StorageConfig: storageoptions.NewStorageOptions(),
	}
}

// AddEtcdFlags adds flags related to etcd storage for a specific APIServer to the specified FlagSet
func (s *ServerRunOptions) AddServerRunFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Host, "host", s.Host, ""+
		"specific what host for api server")

	fs.Int32Var(&s.Port, "port", s.Port, ""+
		"specific what port for api server ")

	fs.StringVar(&s.SwaggerPath, "swagger-path", s.SwaggerPath, ""+
		"specific a path where found swagger index.html, if not will be disable swagger ui")

	s.StorageConfig.AddUniversalFlags(fs)
	s.StorageConfig.AddMysqlStorageFlags(fs)
}
