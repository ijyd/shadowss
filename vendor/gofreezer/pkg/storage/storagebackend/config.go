package storagebackend

import "gofreezer/pkg/runtime"

const (
	//StorageTypeUnset not used any storage backend
	StorageTypeUnset = ""

	//StorageTypeETCD3 used etcd3 as a storabe backend
	StorageTypeETCD3 = "etcd3"
	StorageTypeETCD2 = "etcd2"
)

// //StorageType where found users
// type StorageType string

// Config is configuration for creating a storage backend.
// type Config struct {
// 	// Type defines the type of storage backend, e.g. "etcd2", etcd3". Default ("") is "etcd2".
// 	Type string
//
// 	// ServerList is the list of storage servers to connect with.eg for mysql user@host:port/dbname?param1=value
// 	ServerList []string
// }

// EtcdConfig is configuration for creating a storage backend.
type Config struct {
	// Type defines the type of storage backend, e.g. "etcd2", etcd3". Default ("") is "etcd2".
	Type string
	// Prefix is the prefix to all keys passed to storage.Interface methods.
	Prefix string
	// ServerList is the list of storage servers to connect with.
	ServerList []string
	// TLS credentials
	KeyFile  string
	CertFile string
	CAFile   string
	// Quorum indicates that whether read operations should be quorum-level consistent.
	Quorum bool
	// DeserializationCacheSize is the size of cache of deserialized objects.
	// Currently this is only supported in etcd2.
	// We will drop the cache once using protobuf.
	DeserializationCacheSize int

	Codec runtime.Codec
}
