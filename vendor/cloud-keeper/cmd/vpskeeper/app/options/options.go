package options

import (
	storageoptions "gofreezer/pkg/genericstoragecodec/options"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	SecurePort        int
	InsecurePort      int
	SwaggerPath       string
	LicenseFile       string
	TLSCertFile       string
	TLSPrivateKeyFile string

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
		InsecurePort:      8080,
		SecurePort:        0,
		EtcdStorageConfig: storageoptions.NewStorageOptions().WithEtcdOptions(),
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.IntVar(&s.InsecurePort, "insecure-port", s.InsecurePort, ""+
		"The port on which to serve unsecured, unauthenticated access. Default 8080.")

	fs.IntVar(&s.SecurePort, "secure-port", s.SecurePort, ""+
		"The port on which to serve HTTPS with authentication and authorization. If 0, "+
		"don't serve HTTPS at all.")

	fs.StringVar(&s.SwaggerPath, "swagger-path", s.SwaggerPath, ""+
		"specific a path where found swagger index.html, if not will be disable swagger ui")

	fs.StringVar(&s.LicenseFile, "license-file", s.LicenseFile, ""+
		"specific a file that contains license")

	fs.StringVar(&s.Storage.Type, "storage-type", s.Storage.Type, ""+
		"specify a storage backend for users ")

	fs.StringSliceVar(&s.Storage.ServerList, "server-list", s.Storage.ServerList, ""+
		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated")

	fs.StringVar(&s.TLSCertFile, "tls-cert-file", s.TLSCertFile, ""+
		"File containing x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert). If HTTPS serving is enabled, must configure this")

	fs.StringVar(&s.TLSPrivateKeyFile, "tls-private-key-file", s.TLSPrivateKeyFile,
		"File containing x509 private key matching --tls-cert-file.")

	s.EtcdStorageConfig.AddUniversalFlags(fs)
	s.EtcdStorageConfig.AddEtcdStorageFlags(fs)
}
