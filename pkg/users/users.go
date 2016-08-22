package users

import (
	"time"

	"shadowsocks-go/pkg/config"
	"shadowsocks-go/pkg/proxyserver"
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

//Users implement user Information storage and users connection manager
type Users struct {
	StorageConfig  storagebackend.Config
	StorageHandler storage.Interface
	SyncInterval   int
	//ProxyServer hold on a proxy server to sys user
	ProxyServer *proxyserver.Servers
}

func NewUsers() *Users {
	u := &Users{}
	return u
}

func (u *Users) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&u.StorageConfig.Type, "storage-type", u.StorageConfig.Type, ""+
		"specify a storage backend for users ")

	fs.StringSliceVar(&u.StorageConfig.ServerList, "server-list", u.StorageConfig.ServerList, ""+
		"specify server to connented backend.eg:user:password@tcp(host:port)/dbname, comma separated")

	fs.IntVar(&u.SyncInterval, "sync-user-interval", u.SyncInterval, ""+
		"specify a interval for sync user from backend storage")
}

//CreateUsersSync create a routine for sync users from backend storage
func (u *Users) CreateUsersSync(proxySrv *proxyserver.Servers) {

	var err error
	glog.V(3).Infof("New storage  %s with server %v\r\n", u.StorageConfig.Type, u.StorageConfig.ServerList)

	if len(u.StorageConfig.Type) != 0 && len(u.StorageConfig.ServerList) != 0 {
		u.StorageHandler, err = newStorage(u.StorageConfig)
		if err != nil {
			glog.Errorf("Create backend error:%v\r\n", err)
		}
	}

	u.ProxyServer = proxySrv

	go runSync(u)
}

func (u *Users) getUsers() ([]User, error) {
	return get(u.StorageHandler)
}

func (u *Users) refreshTraffic(config *config.ConnectionInfo, user *User) error {
	//update users traffic
	upload, download, err := u.ProxyServer.GetTraffic(config)
	if err != nil {
		return err
	}

	totalUpload := int64(user.UploadTraffic) + upload
	totalDownlaod := int64(user.DownloadTraffic) + download

	return updateTraffic(u.StorageHandler, int(config.ID), totalUpload, totalDownlaod)
}

func coverUserToConfig(user *User) *config.ConnectionInfo {
	return &config.ConnectionInfo{
		ID:            int64(user.ID),
		Host:          string("0.0.0.0"),
		Port:          user.Port,
		EncryptMethod: user.Method,
		Password:      user.Passwd,
		EnableOTA:     user.Enable != 0,
		Timeout:       60,
	}
}

func runSync(users *Users) {
	for {
		go func() {
			// sync users
			userList, err := users.getUsers()
			if err != nil {
				glog.Errorf("Get users failure:%v\r\n", err)
			} else {
				glog.V(5).Infof("Got users %v\r\n", userList)
				for _, v := range userList {
					config := coverUserToConfig(&v)
					exist, equal := users.ProxyServer.CheckServer(config)
					if !exist {
						//force add for new item
						users.ProxyServer.StartWithConfig(config)
					} else {
						if !equal {
							//re add modify server
							users.ProxyServer.StopServer(config)
							time.Sleep(time.Duration(500) * time.Microsecond)
							users.ProxyServer.StartWithConfig(config)
						}
					}
					users.refreshTraffic(config, &v)
				}
			}
		}()

		time.Sleep(time.Duration(users.SyncInterval) * time.Second)
	}
}
