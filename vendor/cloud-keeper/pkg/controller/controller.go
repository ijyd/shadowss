package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/cache"
	"cloud-keeper/pkg/controller/apiserverctl"
	"cloud-keeper/pkg/etcdhelper"
	"fmt"

	"golib/pkg/util/network"
	"strings"
	"time"

	"github.com/golang/glog"
)

const (
	maxNodeUserLockCacheSize int = 128
)

var AutoSchedule *NodeSchedule
var nodeUserLockCache *cache.LRUExpireCache

func getNodeUserLock(username string, val string) (string, bool) {
	value := val
	result, ok := nodeUserLockCache.Get(username)
	if !ok {
		nodeUserLockCache.Add(username, val, 10*time.Minute)
	} else {
		//covert result to value
		value = result.(string)
	}

	return value, ok
}

func releaseNodeUserLock(username string) {
	nodeUserLockCache.Remove(username)
}

func ControllerStart(helper *etcdhelper.EtcdHelper, be *backend.Backend, port int) error {
	AutoSchedule = NewNodeSchedule(helper, be)

	//add apiserver node
	has := apiserverctl.CheckLocalAPIServer(helper)
	glog.V(5).Infof("has local server %v", has)
	if !has {
		var hostList []string
		localExternalHost, err := network.ExternalIP()
		if err != nil {
			return err
		}
		hostList = append(hostList, localExternalHost)

		internetIP, err := network.ExternalInternetIP()
		if err != nil {
			return err
		}
		internetIP = strings.Replace(internetIP, "\n", "", -1)
		hostList = append(hostList, internetIP)

		_, err = apiserverctl.AddLocalAPIServer(be, helper, localExternalHost, hostList, port, uint64(0), true, true)
		if err != nil {
			return err
		}
	}

	nodeUserLockCache = cache.NewLRUExpireCache(maxNodeUserLockCacheSize)

	go manageNode(helper)
	go resumeUserEachMonth(AutoSchedule)

	return nil
}

//delete this user service.
func DeleteUserService(name string) error {
	_, ok := getNodeUserLock(name, name)
	glog.V(5).Infof("get user %v lock %v\r\n", name, ok)
	if ok {
		return fmt.Errorf("alloc for %v in progressing", name)
	}

	err := AutoSchedule.DelAllNodeUserByUser(name)
	if err != nil {
		glog.Errorf("del node user error %v \r\n", err)
	}
	err = AutoSchedule.DelUserService(name)

	releaseNodeUserLock(name)
	return err
}

func delUserServiceNode(nodeName, userName string) error {
	userSrv, err := AutoSchedule.checkUserServiceNode(userName, nodeName)
	if err != nil {
		return err
	}

	err = AutoSchedule.DelUserFromNode(nodeName, userSrv.Spec.NodeUserReference[nodeName].User)
	if err != nil {
		glog.Errorf("del user %+v from node %+v error %v", userName, nodeName, err)
	}

	err = AutoSchedule.DeleteUserServiceNode(userName, nodeName)
	if err != nil {
		glog.Errorf("del node %+v from user %+v error %v", userName, nodeName, err)
		return err
	}
	return nil
}

func DeleteUserServiceNode(nodeName, userName string) error {
	_, ok := getNodeUserLock(userName, userName)
	glog.V(5).Infof("get user %v lock %v\r\n", userName, ok)
	if ok {
		return fmt.Errorf("alloc for %v in progressing", userName)
	}

	err := delUserServiceNode(nodeName, userName)

	releaseNodeUserLock(userName)

	return err
}

func BindUserToNode(userName string, nodeReference map[string]api.UserReferences) error {
	_, ok := getNodeUserLock(userName, userName)
	glog.V(5).Infof("get user %v lock %v\r\n", userName, ok)
	if ok {
		return fmt.Errorf("alloc for %v in progressing", userName)
	}
	err := AutoSchedule.BindUserServiceWithNode(userName, nodeReference)
	glog.V(5).Infof("bind user %v to node %v done %v\r\n", userName, nodeReference, err)
	releaseNodeUserLock(userName)

	return err
}

func AllocDefaultNodeForUser(name string) error {

	err := AutoSchedule.AllocDefaultNode(name)
	if err != nil {
		glog.Errorf("alloc default api node for user error %v\r\n", err)
	}

	return err
}

func ReallocUserNodeByProperties(name string, properties map[string]string) error {
	_, ok := getNodeUserLock(name, name)
	if ok {
		return fmt.Errorf("alloc for %v in progressing", name)
	}

	AutoSchedule.DelAllNodeUserByUser(name)

	//must wait delete user done
	time.Sleep((nodeUserLease * 2) * time.Second)

	err := AutoSchedule.AllocNodeByUserProperties(name, properties)
	if err != nil {
		glog.Errorf("alloc user by properties error %v\r\n", err)
	}
	releaseNodeUserLock(name)

	return err
}

//GetAvailableNodeAPINode for api
func GetAvailableNodeAPINode(limit int) []dynamicNodeInfo {
	return AutoSchedule.getAvailableNodeAPINode(limit)
}
