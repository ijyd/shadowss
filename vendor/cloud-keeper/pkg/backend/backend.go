package backend

import (
	"fmt"

	"cloud-keeper/pkg/backend/db"

	"golib/pkg/storage"
	"golib/pkg/storage/storagebackend"

	"github.com/golang/glog"
)

//Backend implement backend Information storage and users connection manager
type Backend struct {
	StorageConfig  storagebackend.Config
	StorageHandler storage.Interface
}

func NewBackend(typ string, servers []string) *Backend {
	u := &Backend{
		StorageConfig: storagebackend.Config{
			Type:       typ,
			ServerList: servers,
		},
	}
	return u
}

// func (u *Backend) AddFlags(fs *pflag.FlagSet) {
//
// 	fs.StringVar(&u.StorageConfig.Type, "storage-type", u.StorageConfig.Type, ""+
// 		"specify a storage backend for users ")
//
// 	fs.StringSliceVar(&u.StorageConfig.ServerList, "server-list", u.StorageConfig.ServerList, ""+
// 		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated")
// }

func (u *Backend) CreateStorage() error {

	var err error
	glog.V(3).Infof("New storage  %s with server %v\r\n", u.StorageConfig.Type, u.StorageConfig.ServerList)

	if len(u.StorageConfig.Type) != 0 && len(u.StorageConfig.ServerList) != 0 {
		u.StorageHandler, err = db.NewStorage(u.StorageConfig)
	} else {
		err = fmt.Errorf("not configure any backend")
	}

	return err
}
