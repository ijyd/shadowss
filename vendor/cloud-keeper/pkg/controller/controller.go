package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller/apiserverctl"
	"cloud-keeper/pkg/etcdhelper"
	"golib/pkg/util/network"

	"github.com/golang/glog"
)

const (
	perNodeUserLimit = 30
)

var AutoSchedule *NodeSchedule

func ControllerStart(helper *etcdhelper.EtcdHelper, be *backend.Backend, port int) error {
	AutoSchedule = NewNodeSchedule(helper, be)

	//add apiserver node
	has := apiserverctl.CheckLocalAPIServer(helper)
	glog.V(5).Infof("has local server %v", has)
	if !has {
		host, err := network.ExternalIP()
		if err != nil {
			return err
		}
		_, err = apiserverctl.AddLocalAPIServer(be, helper, host, port, uint64(0), true, true)
		if err != nil {
			return err
		}
	}

	go manageNode(helper)

	return nil
}

func DeleteUserAllNode(name string) error {
	err := AutoSchedule.CleanNodeUser(name)
	if err != nil {
		return err
	}

	return AutoSchedule.DelUserService(name)
}

func DeleteUserNode(nodeName, userName string) error {
	err := AutoSchedule.DelUserFromNode(nodeName, userName)
	if err != nil {
		glog.Errorf("del user %+v from node %+v error %v", userName, nodeName, err)
	}

	return err
}

func BindUserToNode(nodeReference map[string]api.UserReferences) error {
	return AutoSchedule.BindUserToNode(nodeReference)
}

func AllocDefaultNodeForUser(name string) error {
	err := AutoSchedule.AllocDefaultNode(name)
	if err != nil {
		glog.Errorf("alloc default api node for user error %v\r\n", err)
	}
	return err
}

func ReallocUserNodeByProperties(name string, properties map[string]string) error {
	AutoSchedule.CleanNodeUser(name)
	err := AutoSchedule.AllocNodeByUserProperties(name, properties)
	if err != nil {
		glog.Errorf("alloc user by properties error %v\r\n", err)
	}

	return err
}
