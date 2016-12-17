package options

import (
	genericoptions "apistack/pkg/genericapiserver/options"

	"github.com/spf13/pflag"
)

type ServerOption struct {
	LicenseFile string

	GenericServerRunOptions *genericoptions.ServerRunOptions
	Storage                 *genericoptions.StorageOptions
	SecureServing           *genericoptions.SecureServingOptions
	InsecureServing         *genericoptions.ServingOptions
	Authentication          *genericoptions.BuiltInAuthenticationOptions
	Authorization           *genericoptions.BuiltInAuthorizationOptions
}

func NewServerOption() *ServerOption {
	return &ServerOption{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		Storage:                 genericoptions.NewStorageOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		Authentication:          genericoptions.NewBuiltInAuthenticationOptions().WithAll(),
		Authorization:           genericoptions.NewBuiltInAuthorizationOptions(),
	}
}

func (s *ServerOption) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&s.LicenseFile, "license-file", s.LicenseFile, ""+
		"specific a file that contains license")

	s.GenericServerRunOptions.AddUniversalFlags(fs)
	s.Storage.NewEtcdOptions().AddFlags(fs)
	s.Storage.NewMysqlOptions().AddFlags(fs)
	s.Storage.NewDynamoDBOptions().AddFlags(fs)
	s.SecureServing.AddFlags(fs)
	s.SecureServing.AddDeprecatedFlags(fs)
	s.InsecureServing.AddFlags(fs)
	s.InsecureServing.AddDeprecatedFlags(fs)
	s.Authentication.AddFlags(fs)
	s.Authorization.AddFlags(fs)

}
