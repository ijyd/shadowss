package options

import (
	genericoptions "apistack/pkg/genericapiserver/options"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		GenericServerRunOptions: genericoptions.NewServerRunOptions().WithEtcdOptions(),
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	s.GenericServerRunOptions.AddUniversalFlags(fs)
	s.GenericServerRunOptions.AddEtcdStorageFlags(fs)
	s.GenericServerRunOptions.AddMysqlStorageFlags(fs)
	s.GenericServerRunOptions.AddDynamoDBStorageFlags(fs)
}
