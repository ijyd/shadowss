package app

import (
	"fmt"
	"io/ioutil"

	"cloud-keeper/cmd/vpskeeper/app/options"
	"cloud-keeper/pkg/apiserver"
	"cloud-keeper/pkg/backend"

	"golib/pkg/permissions"
)

func checkLicense(file string) bool {
	if len(file) == 0 {
		file = "./.license"
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	result := permissions.PermissionsHandler.PermissionsCheck(string(data))

	return result

}

//Run start a api server
func Run(options *options.ServerOption) error {

	be := backend.NewBackend(options.Storage.Type, options.Storage.ServerList)

	config := apiserver.Config{
		Host:               options.Host,
		Port:               int(options.Port),
		StorageClient:      be,
		SwaggerPath:        options.SwaggerPath,
		EtcdStorageOptions: options.EtcdStorageConfig,
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
