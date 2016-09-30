package factory

// Create creates a storage backend based on given config.
import (
	"fmt"

	"gofreezer/pkg/storage"
	"gofreezer/pkg/storage/storagebackend"
)

// DestroyFunc is to destroy any resources used by the storage returned in Create() together.
type DestroyFunc func()

//Create a storage interface
func Create(c storagebackend.Config) (storage.Interface, DestroyFunc, error) {
	switch c.Type {
	case storagebackend.StorageTypeUnset, storagebackend.StorageTypeETCD3:
		// TODO: We have the following features to implement:
		// - Support secure connection by using key, cert, and CA files.
		// - Honor "https" scheme to support secure connection in gRPC.
		// - Support non-quorum read.
		return newETCD3Storage(c)
	default:
		return nil, nil, fmt.Errorf("unknown storage type: %s", c.Type)
	}
}
