package users

import (
	"cloud-keeper/pkg/api"
	"gofreezer/pkg/runtime"
	"shadowsocks-go/pkg/config"
	"shadowsocks-go/pkg/proxyserver"
	"time"

	"github.com/golang/glog"
)

type RefreshUser func(user *api.NodeUser, del bool)

type Users struct {
	proxyHandle *proxyserver.Servers
	refresh     RefreshUser
}

func NewUsers(proxyserver *proxyserver.Servers, refresh RefreshUser) *Users {
	return &Users{
		proxyHandle: proxyserver,
		refresh:     refresh,
	}
}

func (u *Users) CoverUserToConfig(user *api.NodeUser) *config.ConnectionInfo {
	return &config.ConnectionInfo{
		ID:            int64(user.Spec.User.ID),
		Host:          string("0.0.0.0"),
		Port:          int(user.Spec.User.Port),
		EncryptMethod: user.Spec.User.Method,
		Password:      user.Spec.User.Password,
		EnableOTA:     user.Spec.User.EnableOTA,
		Timeout:       60,
	}
}

func (u *Users) UpdateTraffic(config *config.ConnectionInfo, user *api.NodeUser) error {
	//update users traffic
	upload, download, err := u.proxyHandle.GetTraffic(config)
	if err != nil {
		return err
	}

	totalUpload := int64(user.Spec.User.UploadTraffic) + upload
	totalDownlaod := int64(user.Spec.User.DownloadTraffic) + download

	user.Spec.User.DownloadTraffic = totalDownlaod
	user.Spec.User.UploadTraffic = totalUpload

	return nil
}

func (u *Users) AddObj(obj runtime.Object) {
	nodeUser := obj.(*api.NodeUser)

	config := u.CoverUserToConfig(nodeUser)
	glog.V(5).Infof("add user %v \r\n", config)
	u.proxyHandle.StartWithConfig(config)

	time.Sleep(1 * time.Second)
	port, err := u.proxyHandle.GetListenPort(config)
	if err != nil {
		glog.Errorf("Get listen port failure %v\r\n", err)
	} else {
		nodeUser.Spec.User.Port = int64(port)
		u.refresh(nodeUser, false)
	}
}

func (u *Users) ModifyObj(obj runtime.Object) {
	// nodeUser := obj.(*api.NodeUser)
	// config := coverUserToConfig(nodeUser)
	//
	// u.proxyHandle.StopServer(config)
	//
	// u.updateTraffic(config, nodeUser)
	//
	// //release this server re add
	// u.proxyHandle.CleanUpServer(config)
	//
	// time.Sleep(time.Duration(500) * time.Microsecond)
	// u.proxyHandle.StartWithConfig(config)
	//
	// port, err := u.proxyHandle.GetListenPort(config)
	// if err != nil {
	// 	glog.Errorf("Get listen port failure %v\r\n", err)
	// } else {
	// 	nodeUser.Spec.User.Port = int64(port)
	// 	u.refresh(nodeUser, false)
	// }
	//not support this for node user.delete it then add
}

func (u *Users) DelObj(obj runtime.Object) {
	nodeUser := obj.(*api.NodeUser)
	config := u.CoverUserToConfig(nodeUser)
	u.proxyHandle.StopServer(config)

	u.UpdateTraffic(config, nodeUser)

	u.proxyHandle.CleanUpServer(config)

	u.refresh(nodeUser, true)
}
