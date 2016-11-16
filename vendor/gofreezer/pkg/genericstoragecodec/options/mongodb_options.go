package options

import "github.com/spf13/pflag"

// AddMysqlStorageFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet
func (s *StorageOptions) AddMongoDBStorageFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&s.StorageConfig.MongoCfg.ServerList, "mongo-servers", s.StorageConfig.MongoCfg.ServerList, ""+
		"specify server to connented backend.eg:mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb.")

	fs.StringSliceVar(&s.StorageConfig.MongoCfg.AdminCred, "mongo-admin", s.StorageConfig.MongoCfg.AdminCred, ""+
		"specify admin cred eg:admindb,admin,123456.")

	fs.StringSliceVar(&s.StorageConfig.MongoCfg.GeneralCred, "mongo-user", s.StorageConfig.MongoCfg.GeneralCred, ""+
		"specify general cred for project eg:test,testUser,123456.")
}
