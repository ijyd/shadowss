package app

import (
	"fmt"
	"gofreezer/examples/etcd/app/api"
	"gofreezer/examples/etcd/app/api/v1"
	"gofreezer/examples/etcd/app/apiserver"
	"gofreezer/examples/etcd/app/options"
	"gofreezer/pkg/genericstoragecodec"
)

func Run(options *options.ServerRunOptions) error {

	// config := storagebackend.Config{
	// 	Type:       "etcd3",
	// 	Prefix:     "/registry",
	// 	ServerList: []string{"http://192.168.60.100:2379"},
	// 	Quorum:     false,
	// 	DeserializationCacheSize: 0,
	// }

	storageVersion := v1.SchemeGroupVersion
	//memVersion := api.SchemeGroupVersion

	// codec, err := NewStorageCodec()
	//
	// config.Codec = codec
	//
	// StorageHandle, _, err := factory.Create(config)
	// if err != nil {
	// 	return err
	// }

	storageCodec, err := genericstoragecodec.NewGenericStorageCodec(options.StorageConfig, api.Codecs, storageVersion)
	if err != nil {
		return fmt.Errorf("new storage codec error %v", err)
	}

	configsrv := apiserver.Config{
		Host:         options.Host,
		Port:         int(options.Port),
		StorageCodec: storageCodec,
		SwaggerPath:  options.SwaggerPath,
	}

	serverHandler := apiserver.NewApiServer(configsrv)
	if serverHandler == nil {
		return fmt.Errorf("api server init failure")
	}

	if err := serverHandler.Run(); err != nil {
		return err
	}

	return nil
}
