package app

import (
	"fmt"

	"gofreezer/examples/common/apiext"
	"gofreezer/examples/common/apiext/v1"
	"gofreezer/examples/common/apiserver"
	"gofreezer/examples/dynamodb/app/options"
	"gofreezer/examples/dynamodb/app/resthandle"
	"gofreezer/pkg/genericstoragecodec"
	"gofreezer/pkg/storage/awsdynamodb"
)

func Run(options *options.ServerRunOptions) error {

	storageVersion := v1.SchemeGroupVersion

	storageCodec, err := genericstoragecodec.NewGenericStorageCodec(options.StorageConfig, apiext.Codecs, storageVersion)
	if err != nil {
		return fmt.Errorf("new storage codec error %v", err)
	}

	configsrv := apiserver.Config{
		Host:        options.Host,
		Port:        int(options.Port),
		SwaggerPath: options.SwaggerPath,
	}

	serverHandler := apiserver.NewApiServer(configsrv)
	if serverHandler == nil {
		return fmt.Errorf("api server init failure")
	}

	resthandle.GenericStorage = storageCodec
	resthandle.StorageCodec = resthandle.GenericStorage.Storage.(awsdynamodb.Interface)

	if err := serverHandler.Run(resthandle.InstallWS); err != nil {
		return err
	}

	return nil
}
