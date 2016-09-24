package options

import (
	storageoptions "gofreezer/pkg/genericstoragecodec/options"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	Host        string
	Port        int32
	SwaggerPath string
	LicenseFile string

	Storage           storageOption
	EtcdStorageConfig *storageoptions.StorageOptions
}

type storageOption struct {
	Type string

	// ServerList is the list of storage servers to connect with.eg for mysql user@host:port/dbname?param1=value
	ServerList []string
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		Host:              "",
		Port:              80,
		EtcdStorageConfig: storageoptions.NewStorageOptions().WithEtcdOptions(),
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.Host, "host", s.Host, ""+
		"specific what host for api server")

	fs.Int32Var(&s.Port, "port", s.Port, ""+
		"specific what port for api server ")

	fs.StringVar(&s.SwaggerPath, "swagger-path", s.SwaggerPath, ""+
		"specific a path where found swagger index.html, if not will be disable swagger ui")

	fs.StringVar(&s.LicenseFile, "license-file", s.LicenseFile, ""+
		"specific a file that contains license")

	fs.StringVar(&s.Storage.Type, "storage-type", s.Storage.Type, ""+
		"specify a storage backend for users ")

	fs.StringSliceVar(&s.Storage.ServerList, "server-list", s.Storage.ServerList, ""+
		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated")

	s.EtcdStorageConfig.AddUniversalFlags(fs)
	s.EtcdStorageConfig.AddEtcdStorageFlags(fs)
}
