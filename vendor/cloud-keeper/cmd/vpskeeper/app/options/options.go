package options

import (
	genericoptions "apistack/pkg/genericapiserver/options"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	LicenseFile string

	GenericServerRunOptions *genericoptions.ServerRunOptions
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		GenericServerRunOptions: genericoptions.NewServerRunOptions().WithEtcdOptions(),
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.LicenseFile, "license-file", s.LicenseFile, ""+
		"specific a file that contains license")

	s.GenericServerRunOptions.AddUniversalFlags(fs)
	s.GenericServerRunOptions.AddEtcdStorageFlags(fs)
	s.GenericServerRunOptions.AddMysqlStorageFlags(fs)
}
