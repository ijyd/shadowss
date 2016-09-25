package multiuser

import (
	"golib/pkg/util/network"
	"shadowsocks-go/pkg/multiuser/apiserverproxy"
	"shadowsocks-go/pkg/multiuser/users"
	"shadowsocks-go/pkg/proxyserver"

	"github.com/golang/glog"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/apiserverctl"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/etcdhelper"
	"cloud-keeper/pkg/watcher"
	"gofreezer/pkg/genericstoragecodec/options"
)

type MultiUser struct {
	etcdHandle  *etcdhelper.EtcdHelper
	proxyHandle *proxyserver.Servers
}

var schedule *MultiUser

func InitSchedule(options *options.StorageOptions, proxySrv *proxyserver.Servers) {
	schedule = NewMultiUser(options, proxySrv)

	obj, err := apiserverctl.GetAPIServers(schedule.etcdHandle)
	if err != nil {
		glog.Errorf("not found any api server in this cluster err %v\r\n", err)
		return
	}
	apiSrv := obj.(*api.APIServerList)
	var apiServerList []api.APIServerInfor
	for _, v := range apiSrv.Items {
		apiServerList = append(apiServerList, v.Spec.Server)
	}
	if len(apiServerList) == 0 {
		glog.Errorf("not found any api server %v in this cluster\r\n", apiSrv)
		return
	}
	glog.V(5).Infof("Got apiserver %+v\r\n", apiServerList)
	apiserverproxy.InitAPIServer(apiServerList)

	nodeName, err := network.ExternalMAC()
	if err != nil {
		glog.Errorf("got mac addr error %v\r\n", err)
		return
	}
	prifixKey := nodectl.PrefixNode + "/" + nodeName + nodectl.PrefixNodeUser
	KeepHealth(nodeName)

	userMgr := users.NewUsers(proxySrv, RefreshUser)
	go watcher.WatchNodeUsersLoop(prifixKey, schedule.etcdHandle, userMgr)

}

func KeepHealth(nodeName string) {
	location := "NewYork"
	accsrvid := 4516541
	accsrvname := "test1111"
	host, err := network.ExternalIP()
	ttl := uint64(0)
	if err != nil {
		return
	}

	obj, err := nodectl.GetNode(schedule.etcdHandle, nodeName)
	if err != nil {
		glog.Errorf("get node error:%v\r\n", err)
		return
	}

	node := obj.(*api.Node)
	if node.Name == nodeName {
		glog.Infof("our node %+v already exist \r\n", obj)
	} else {
		nodeHelper := &nodectl.NodeHelper{
			TTL:         ttl,
			Name:        nodeName,
			Host:        host,
			Location:    location,
			AccsrvID:    int64(accsrvid),
			AccsrvName:  accsrvname,
			Annotations: map[string]string{nodectl.NodeAnnotationUserCnt: "0"},
			//Labels:      map[string]string{nodect.NodeLablesDefault: "false"},
		}
		obj, err := nodectl.AddNodeToEtcdHelper(schedule.etcdHandle, nodeHelper)
		if err != nil {
			glog.Errorf("add node error %v\r\n", err)
			return
		}
		glog.V(5).Infof("Add node %+v err %v\r\n", obj, err)
	}

	// for {
	// 	select {
	// 	case <-time.After(time.Hour * 2):
	// 		nodectl.AddNodeHelper(schedule.etcdHandle, nodeName, host, location, int64(accsrvid), accsrvname)
	// 	default:
	// 		time.Sleep(1 * time.Hour)
	// 	}
	// }
}

func NewMultiUser(options *options.StorageOptions, proxySrv *proxyserver.Servers) *MultiUser {

	return &MultiUser{
		etcdHandle:  etcdhelper.NewEtcdHelper(options),
		proxyHandle: proxySrv,
	}
}
