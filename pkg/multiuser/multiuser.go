package multiuser

import (
	"encoding/json"
	"fmt"
	"golib/pkg/util/network"
	"io/ioutil"
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
	nodeName    string
	nodeAttr    map[string]string
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

	// nodeName, err := network.External()
	// if err != nil {
	// 	glog.Errorf("got mac addr error %v\r\n", err)
	// 	return nil
	// }

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
	mu.KeepHealth()

	userMgr := users.NewUsers(mu.proxyHandle, RefreshUser)
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
	ttl := uint64(0)

	obj, err := nodectl.GetNode(schedule.etcdHandle, mu.nodeName)
	if err != nil {
		glog.Errorf("get node error:%v\r\n", err)
		return
	}

	node := obj.(*api.Node)
	if node.Name == mu.nodeName {
		glog.Infof("our node %+v already exist \r\n", obj)
	} else {
		nodeHelper := mu.BuildNodeHelper(ttl)
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

	// for {
	// 	select {
	// 	case <-time.After(time.Hour * 2):
	// 		nodectl.AddNodeHelper(schedule.etcdHandle, nodeName, host, location, int64(accsrvid), accsrvname)
	// 	default:
	// 		time.Sleep(1 * time.Hour)
	// 	}
	// }
}
