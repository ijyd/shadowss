package app

import "cloud-keeper/cmd/prometheushook/app/options"

//Run start a api server
func Run(options *options.ServerOption) error {

	apisrv := APIServer{
		InsecurePort:      options.InsecurePort,
		SecurePort:        options.SecurePort,
		TLSCertFile:       options.TLSCertFile,
		TLSPrivateKeyFile: options.TLSPrivateKeyFile,
	}

	if err := apisrv.Run(); err != nil {
		return err
	}

	return nil
}
