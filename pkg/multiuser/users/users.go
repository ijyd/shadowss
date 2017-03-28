package users

import (
	"fmt"
	"math/rand"
	"time"

	"shadowss/pkg/api"
	"shadowss/pkg/config"
	muconfig "shadowss/pkg/multiuser/config"
	"shadowss/pkg/proxyserver"

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

func (u *Users) GetUsers() []config.ConnectionInfo {
	return u.proxyHandle.GetUsersConfig()
}

func (u *Users) GetUsersInfo() []proxyserver.UserInfo {
	return u.proxyHandle.GetUsersInfor()
}

func (u *Users) StartAPIProxy(nodename string) error {
	generateRandID := int64(100000000)
	id := rand.Int63n(generateRandID)
	port := muconfig.GetDefaultAPIProxyListenPort()
	name := fmt.Sprintf("%s_apiproxy", nodename)
	config := &config.ConnectionInfo{
		ID:            (id + generateRandID) * 100,
		Host:          string("0.0.0.0"),
		Port:          port,
		EncryptMethod: string("aes-256-cfb"),
		Password:      string("48c8591290877f737202ad20c06780e9"),
		EnableOTA:     true,
		Timeout:       60,
		Name:          name,
		MaxConnection: 100,
	}

	u.proxyHandle.StartWithConfig(config)

	return nil
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
		Name:          user.Name,
		MaxConnection: muconfig.GetMaxTCPConnectionPerPort(),
	}
}

func (u *Users) StartUserSrv(config *config.ConnectionInfo, nodeUser *api.NodeUser) {
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

func (u *Users) UpdateTraffic(config *config.ConnectionInfo, user *api.NodeUser) error {
	//update users traffic
	upload, download, err := u.proxyHandle.GetTraffic(config)
	if err != nil {
		return err
	}

	totalUpload := upload
	totalDownlaod := download

	user.Spec.User.DownloadTraffic = totalDownlaod
	user.Spec.User.UploadTraffic = totalUpload

	return nil
}

func (u *Users) AddUsers(nodeUser *api.NodeUser) {

	config := u.CoverUserToConfig(nodeUser)

	exist, equal := u.proxyHandle.CheckServer(config)
	if exist {
		if !equal {
			glog.Infof("add user exist(%+v). restart  config(%+v) \r\n", *nodeUser, config)
			u.proxyHandle.StopServer(config)
			u.UpdateTraffic(config, nodeUser)
			u.proxyHandle.CleanUpServer(config)
			u.StartUserSrv(config, nodeUser)
		}

	} else {
		glog.Infof("add user(%+v) config(%+v) \r\n", *nodeUser, config)
		u.StartUserSrv(config, nodeUser)
	}
}

func (u *Users) DelUsers(nodeUser *api.NodeUser) {
	config := u.CoverUserToConfig(nodeUser)
	glog.Infof("del user(%+v) config(%+v)\r\n", *nodeUser, config)
	u.proxyHandle.StopServer(config)

	u.UpdateTraffic(config, nodeUser)

	u.proxyHandle.CleanUpServer(config)

	u.refresh(nodeUser, true)
}

func (u *Users) ListUserLoop(nodename string) {

	var loopcnt int
loop:
	loopcnt++
	err := u.WatchUserLoop(nodename)

	rand.Seed(time.Now().UTC().UnixNano())
	waitTime := rand.Intn(60) + 10
	glog.Errorf("watch user loop quit %v, restart after %v(s) \r\n", err, waitTime)
	time.Sleep(time.Second * time.Duration((waitTime)))
	goto loop

}
