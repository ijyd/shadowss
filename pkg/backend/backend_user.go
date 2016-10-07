package backend

import (
	"fmt"
	"time"

	"shadowss/pkg/backend/db"
	"shadowss/pkg/config"
	"shadowss/pkg/proxyserver"

	"github.com/golang/glog"
)

//CreateUsersSync create a routine for sync users from backend storage
func (u *Backend) CreateUsersSync(proxySrv *proxyserver.Servers) {

	var err error
	glog.V(3).Infof("New storage  %s with server %v\r\n", u.StorageConfig.Type, u.StorageConfig.ServerList)

	if len(u.StorageConfig.Type) != 0 && len(u.StorageConfig.ServerList) != 0 {
		u.StorageHandler, err = db.NewStorage(u.StorageConfig)
		if err != nil {
			glog.Errorf("Create backend error:%v\r\n", err)
			return
		}
		u.ProxyServer = proxySrv
		go runSync(u)
	} else {
		glog.Warningf("have not any backend \r\n")
	}
}

func (u *Backend) GetUserByName(name string) (*db.User, error) {
	return db.GetUser(u.StorageHandler, name)
}

func (u *Backend) GetUserByID(id int64) (*db.User, error) {
	return db.GetUserByID(u.StorageHandler, id)
}

func (u *Backend) getUsers() ([]db.User, error) {
	node, err := db.Getnodes(u.StorageHandler)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, fmt.Errorf("this node not found\r\n")
	}
	glog.V(5).Infof("Got node %+v \r\n", node)

	return db.GetServUsers(u.StorageHandler, node.StartUserID, node.EndUserID)
}

func (u *Backend) refreshTraffic(config *config.ConnectionInfo, user *db.User) error {
	//update users traffic
	upload, download, err := u.ProxyServer.GetTraffic(config)
	if err != nil {
		return err
	}

	totalUpload := int64(user.UploadTraffic) + upload
	totalDownlaod := int64(user.DownloadTraffic) + download

	return db.UpdateUserTraffic(u.StorageHandler, config.ID, totalUpload, totalDownlaod)
}

func coverUserToConfig(user *db.User) *config.ConnectionInfo {
	return &config.ConnectionInfo{
		ID:            int64(user.ID),
		Host:          string("0.0.0.0"),
		Port:          int(user.Port),
		EncryptMethod: user.Method,
		Password:      user.Passwd,
		EnableOTA:     user.Enable != 0,
		Timeout:       60,
	}
}

func runSync(be *Backend) {
	for {
		go func() {
			// sync users
			userList, err := be.getUsers()
			if err != nil {
				glog.Errorf("Get users failure:%v\r\n", err)
			} else {
				glog.V(5).Infof("Got users %v\r\n", userList)
				for _, v := range userList {
					config := coverUserToConfig(&v)
					exist, equal := be.ProxyServer.CheckServer(config)
					if !exist {
						//force add for new item
						be.ProxyServer.StartWithConfig(config)
					} else {
						if !equal {
							be.ProxyServer.StopServer(config)
							be.refreshTraffic(config, &v)
							//release this server re add
							be.ProxyServer.CleanUpServer(config)

							time.Sleep(time.Duration(500) * time.Microsecond)
							be.ProxyServer.StartWithConfig(config)
						} else {
							be.refreshTraffic(config, &v)
						}
					}
				}
			}
		}()

		time.Sleep(time.Duration(be.SyncInterval) * time.Second)
	}
}
