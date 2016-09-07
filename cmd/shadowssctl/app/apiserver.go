package app

import (
	"fmt"
	"io/ioutil"
	"shadowsocks-go/cmd/shadowssctl/app/options"
	"shadowsocks-go/pkg/apiserver"
	"shadowsocks-go/pkg/backend"
	"shadowsocks-go/pkg/permissions"

	"github.com/golang/glog"
)

func checkLicense(file string) bool {
	if len(file) == 0 {
		file = "./.license"
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		glog.Errorf("Not find license file\r\n")
		return false
	}

	result := permissions.PermissionsCheck(string(data))

	return result

}

//Run start a api server
func Run(options *options.ServerOption, be *backend.Backend) error {

	config := apiserver.Config{
		Host:          options.Host,
		Port:          int(options.Port),
		StorageClient: be,
		SwaggerPath:   options.SwaggerPath,
	}

	if checkLicense(options.LicenseFile) == false {
		return fmt.Errorf("not allow on this server, please contact administrator")
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
