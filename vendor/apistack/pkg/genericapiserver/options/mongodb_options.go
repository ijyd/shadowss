package options

import "github.com/spf13/pflag"

// AddMysqlStorageFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet
func (s *ServerRunOptions) AddMongoDBStorageFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&s.StorageConfig.Mongodb.ServerList, "mongo-servers", s.StorageConfig.Mongodb.ServerList, ""+
		"specify server to connented backend.eg:mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb.")

	fs.StringSliceVar(&s.StorageConfig.Mongodb.AdminCred, "mongo-admin", s.StorageConfig.Mongodb.AdminCred, ""+
		"specify admin cred eg:admindb,admin,123456.")

	fs.StringSliceVar(&s.StorageConfig.Mongodb.GeneralCred, "mongo-user", s.StorageConfig.Mongodb.GeneralCred, ""+
		"specify general cred for project eg:test,testUser,123456.")
}
