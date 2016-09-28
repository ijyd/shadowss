package multiuser

import (
	"encoding/json"
	"fmt"
	"golib/pkg/util/network"
	"io/ioutil"
	"shadowsocks-go/pkg/multiuser/apiserverproxy"
	"shadowsocks-go/pkg/multiuser/users"
	"shadowsocks-go/pkg/proxyserver"
	"time"

	"github.com/golang/glog"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/apiserverctl"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/etcdhelper"
	"cloud-keeper/pkg/watcher"
	"gofreezer/pkg/genericstoragecodec/options"
)

const (
	NodeDefaultTTL = 300
)

type MultiUser struct {
	etcdHandle  *etcdhelper.EtcdHelper
	proxyHandle *proxyserver.Servers
	nodeName    string
	nodeAttr    map[string]string
	userHandle  *users.Users
	ttl         uint64
}

var schedule *MultiUser

func InitSchedule(options *options.StorageOptions, proxySrv *proxyserver.Servers) {
	schedule = NewMultiUser(options, proxySrv)
	if schedule == nil {
		glog.Fatalf("create multi user failure\r\n")
		return
	}

	err := schedule.StartUp()
	if err != nil {
		glog.Fatalf("startup node failure %v\r\n", err)
	}

}

func NewMultiUser(options *options.StorageOptions, proxySrv *proxyserver.Servers) *MultiUser {

	nodeName, err := network.ExternalMAC()
	if err != nil {
		glog.Errorf("got mac addr error %v\r\n", err)
		return nil
	}

	fileName := string("./attr.json")
	config, err := ioutil.ReadFile(fileName)
	if err != nil {
		glog.Errorf("read node config file err %v \r\n", err)
		return nil
	}

	attr := make(map[string]string)

	err = json.Unmarshal(config, &attr)
	if err != nil {
		glog.Errorf("invalid node config field %v\r\n", err)
		return nil
	}

	_, ok := attr[api.NodeLablesChinaISP]
	if !ok {
		glog.Errorf("invalid node config field cnISP\r\n")
		return nil
	}

	_, ok = attr[api.NodeLablesUserSpace]
	if !ok {
		glog.Errorf("invalid node config field user space\r\n")
		return nil
	}
	_, ok = attr[api.NodeLablesVPSLocation]
	if !ok {
		glog.Errorf("invalid node config field vps location\r\n")
		return nil
	}

	_, ok = attr[api.NodeLablesVPSOP]
	if !ok {
		glog.Errorf("invalid node config field vps operator\r\n")
		return nil
	}

	_, ok = attr[api.NodeLablesVPSName]
	if !ok {
		glog.Errorf("invalid node config field vps name\r\n")
		return nil
	}

	_, ok = attr[api.NodeLablesVPSIP]
	if !ok {
		glog.Errorf("invalid node config field vps ip\r\n")
		return nil
	}

	return &MultiUser{
		etcdHandle:  etcdhelper.NewEtcdHelper(options),
		proxyHandle: proxySrv,
		nodeAttr:    attr,
		nodeName:    nodeName,
		ttl:         NodeDefaultTTL,
	}
}

func (mu *MultiUser) StartUp() error {
	obj, err := apiserverctl.GetAPIServers(schedule.etcdHandle)
	if err != nil {
		glog.Errorf("not found any api server in this cluster err %v\r\n", err)
		return err
	}

	apiSrv := obj.(*api.APIServerList)
	var apiServerList []api.APIServerInfor
	for _, v := range apiSrv.Items {
		apiServerList = append(apiServerList, v.Spec.Server)
	}
	if len(apiServerList) == 0 {
		glog.Errorf("not found any api server %v in this cluster\r\n", apiSrv)
		return fmt.Errorf("must have at least one node")
	}

	glog.V(5).Infof("Got apiserver %+v\r\n", apiServerList)
	apiserverproxy.InitAPIServer(apiServerList)

	prifixKey := nodectl.PrefixNode + "/" + mu.nodeName + nodectl.PrefixNodeUser
	go mu.KeepHealth()

	userMgr := users.NewUsers(mu.proxyHandle, RefreshUser)
	mu.userHandle = userMgr

	go watcher.WatchNodeUsersLoop(prifixKey, mu.etcdHandle, userMgr)

	return nil
}

func (mu *MultiUser) BuildNodeHelper(ttl uint64) *nodectl.NodeHelper {
	vpsIP, _ := mu.nodeAttr[api.NodeLablesVPSIP]
	vpsName, _ := mu.nodeAttr[api.NodeLablesVPSName]
	vpsLocation, _ := mu.nodeAttr[api.NodeLablesVPSLocation]

	nodeHelper := &nodectl.NodeHelper{
		TTL:         ttl,
		Name:        mu.nodeName,
		Host:        vpsIP,
		Location:    vpsLocation,
		AccsrvID:    int64(0),
		AccsrvName:  vpsName,
		Annotations: map[string]string{nodectl.NodeAnnotationUserCnt: "0"},
		Labels:      mu.nodeAttr,
	}

	return nodeHelper

}

func (mu *MultiUser) KeepHealth() {

	obj, err := nodectl.GetNode(schedule.etcdHandle, mu.nodeName)
	if err != nil {
		glog.Errorf("get node error:%v\r\n", err)
		return
	}

	node := obj.(*api.Node)
	if node.Name == mu.nodeName {
		glog.Infof("our node %+v already exist \r\n", obj)
		nodectl.UpdateNode(nil, mu.etcdHandle, node, true, false)
	} else {
		nodeHelper := mu.BuildNodeHelper(mu.ttl)
		if nodeHelper == nil {
			glog.Errorf("invalid node configure\r\n")
			return
		}

		obj, err := nodectl.AddNodeToEtcdHelper(mu.etcdHandle, nodeHelper)
		if err != nil {
			glog.Errorf("add node error %v\r\n", err)
			return
		}
		glog.V(5).Infof("Add node %+v err %v\r\n", obj, err)
	}

	expireTime := time.Duration(mu.ttl - 100)

	for {
		select {
		case <-time.After(time.Second * expireTime):
			upload, download, err := mu.CollectorAndUpdateUserTraffic()
			if err == nil {
				node.Spec.Server.Upload += upload
				node.Spec.Server.Download += download
				nodectl.UpdateNode(nil, mu.etcdHandle, node, true, false)
				glog.V(5).Infof("refresh node %v\r\n", *node)
			} else {
				glog.Warningf("collector user traffic error %v\r\n", err)
			}
		}
	}
}

func (mu *MultiUser) CollectorAndUpdateUserTraffic() (int64, int64, error) {

	obj, err := nodectl.GetNodeAllUsers(schedule.etcdHandle, mu.nodeName)
	if err != nil {
		return 0, 0, err
	}
	nodeUserList := obj.(*api.NodeUserList)

	var upload, download int64
	for _, nodeUser := range nodeUserList.Items {
		config := mu.userHandle.CoverUserToConfig(&nodeUser)
		mu.userHandle.UpdateTraffic(config, &nodeUser)
		err := nodectl.UpdateNodeUsersRefer(mu.etcdHandle, nodeUser.Spec)
		if err != nil {
			glog.Errorf("update node user err %v \r\n", err)
		} else {
			upload += nodeUser.Spec.User.UploadTraffic
			download += nodeUser.Spec.User.DownloadTraffic
		}
	}

	return upload, download, nil
}

func RefreshUser(user *api.NodeUser, del bool) {
	if !del {
		//need update noe user port
		glog.V(5).Infof("update node user %+v", *user)
		err := nodectl.UpdateNodeUsersRefer(schedule.etcdHandle, user.Spec)
		if err != nil {
			glog.Errorf("update node user err %v \r\n", err)
		}
	}

	_, err := nodectl.UpdateNodeAnotationsUserCnt(schedule.etcdHandle, user.Spec.NodeName, del)
	if err != nil {
		glog.Errorf("update node anotation err %v \r\n", err)
	}

}
