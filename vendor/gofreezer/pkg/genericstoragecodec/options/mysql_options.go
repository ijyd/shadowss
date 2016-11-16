package options

import (
	"github.com/spf13/pflag"
)

// AddMysqlStorageFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet
func (s *StorageOptions) AddMysqlStorageFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&s.StorageConfig.ServerList, "mysql-servers", s.StorageConfig.ServerList, ""+
		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated.")
}
