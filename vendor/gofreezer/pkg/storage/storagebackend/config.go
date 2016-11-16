package storagebackend

import "gofreezer/pkg/runtime"

const (
	//StorageTypeUnset not used any storage backend
	StorageTypeUnset = ""

	//StorageTypeETCD3 used etcd3 as a storabe backend
	StorageTypeETCD3       = "etcd3"
	StorageTypeETCD2       = "etcd2"
	StorageTypeMysql       = "mysql"
	StorageTypeMongoDB     = "mongo"
	StorageTypeAWSDynamodb = "awsdynamodb"
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
	// Type defines the type of storage backend, e.g. "etcd2", etcd3", mysql. Default ("") is "etcd3".
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

	//if backend is sql type.need give default storage version
	StorageVersion string

	//append backend config.
	//mongodb extend config
	Mongodb MongoExtendConfig
	//aws dynamodb config
	AWSDynamoDB AWSDynamoDBConfig
	//mysql config
	Mysql MysqlConfig
}

type MongoExtendConfig struct {
	//holds options for establishing a session with a MongoDB cluster
	ServerList []string
	//admin credentials:db,user,pwd
	AdminCred []string
	//normal user credentials:db,user,pwd
	GeneralCred []string
}

type AWSDynamoDBConfig struct {
	Region string
	Table  string
}

type MysqlConfig struct {
	// ServerList is the list of storage servers to connect with.
	ServerList []string
}
