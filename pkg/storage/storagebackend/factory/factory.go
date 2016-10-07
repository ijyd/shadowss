package factory

// Create creates a storage backend based on given config.
import (
	"fmt"
	"shadowss/pkg/storage"
	"shadowss/pkg/storage/storagebackend"
)

//Create a storage interface
func Create(c storagebackend.Config) (storage.Interface, error) {
	switch c.Type {
	case storagebackend.StorageTypeUnset, storagebackend.StorageTypeMysql:
		return newMysqlStorage(c)
	default:
		return nil, fmt.Errorf("unknown storage type: %s", c.Type)
	}
}
