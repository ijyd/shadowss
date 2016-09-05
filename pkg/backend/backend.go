package backend

import (
	"fmt"

	"shadowsocks-go/pkg/backend/db"
	"shadowsocks-go/pkg/proxyserver"
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

//Backend implement backend Information storage and users connection manager
type Backend struct {
	StorageConfig  storagebackend.Config
	StorageHandler storage.Interface
	SyncInterval   int
	//ProxyServer hold on a proxy server to sys user
	ProxyServer *proxyserver.Servers
}

func NewBackend() *Backend {
	u := &Backend{}
	return u
}

func (u *Backend) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&u.StorageConfig.Type, "storage-type", u.StorageConfig.Type, ""+
		"specify a storage backend for users ")

	fs.StringSliceVar(&u.StorageConfig.ServerList, "server-list", u.StorageConfig.ServerList, ""+
		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated")

	fs.IntVar(&u.SyncInterval, "sync-user-interval", u.SyncInterval, ""+
		"specify a interval for sync user from backend storage")
}

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
