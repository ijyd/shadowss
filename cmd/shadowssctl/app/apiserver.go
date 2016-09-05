package app

import (
	"fmt"
	"shadowsocks-go/cmd/shadowssctl/app/options"
	"shadowsocks-go/pkg/apiserver"
	"shadowsocks-go/pkg/backend"
)

//Run start a api server
func Run(options *options.ServerOption, be *backend.Backend) error {

	config := apiserver.Config{
		Host:          options.Host,
		Port:          int(options.Port),
		StorageClient: be,
	}

	serverHandler := apiserver.NewApiServer(config)
	if serverHandler == nil {
		return fmt.Errorf("api server init failure")
	}

	if err := serverHandler.Run(); err != nil {
		return err
	}

	return nil
}
