package storagebackend

const (
	//StorageTypeUnset not used any storage backend
	StorageTypeUnset = ""
	//StorageTypeMysql used mysql as a storabe backend
	StorageTypeMysql = "mysql"
)

// //StorageType where found users
// type StorageType string

// Config is configuration for creating a storage backend.
type Config struct {
	// Type defines the type of storage backend, e.g. "etcd2", etcd3". Default ("") is "etcd2".
	Type string

	// ServerList is the list of storage servers to connect with.eg for mysql user@host:port/dbname?param1=value
	ServerList []string
}
