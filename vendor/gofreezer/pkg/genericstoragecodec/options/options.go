package options

import (
	"gofreezer/pkg/storage/storagebackend"

	"github.com/spf13/pflag"
)

const (
	DefaultDeserializationCacheSize = 50000
)

// StorageOptions contains the options while running a generic storage.
type StorageOptions struct {
	DefaultStorageMediaType string
	EtcdServersOverrides    []string
	StorageConfig           storagebackend.Config
}

func NewStorageOptions() *StorageOptions {
	return &StorageOptions{
		DefaultStorageMediaType: "application/json",
	}
}

func (s *StorageOptions) AddUniversalFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.StorageConfig.Type, "storage-backend", s.StorageConfig.Type,
		"The storage backend for persistence. Options: 'etcd3' (default), 'etcd2', 'mysql'.")

	fs.StringVar(&s.DefaultStorageMediaType, "storage-media-type", s.DefaultStorageMediaType, ""+
		"The media type to use to store objects in storage. Defaults to application/json. "+
		"Some resources may only support a specific media type and will ignore this setting.")
}

func (o *StorageOptions) WithEtcdOptions() *StorageOptions {
	o.StorageConfig = storagebackend.Config{
		Prefix: DefaultEtcdPathPrefix,
		DeserializationCacheSize: DefaultDeserializationCacheSize,
	}
	return o
}
