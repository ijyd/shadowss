package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/controller/userctl"
	"fmt"

	"github.com/golang/glog"
)

/*out logic: 1. add user to node in nodeUserLease. 2. use phase to check what operators for obj*/
/*this will be  prevent many key in our etcd*/
const (
	nodeUserLease = 10
)

func (ns *NodeSchedule) NewNodeUserEvent(user *api.NodeUser) {
	switch user.Spec.Phase {
	case api.NodeUserPhaseAdd:
		glog.V(5).Infof("add new node user %v not need implement\r\n", user)
	case api.NodeUserPhaseDelete:
		glog.V(5).Infof("delete node user %v not need implement\r\n", user)
	case api.NodeUserPhaseUpdate:
		err := ns.UpdateUserDynamicDataByNodeUser(user)
		if err != nil {
			glog.Errorf("get user service error %v\r\n", err)
			return
		}
	default:
		glog.Warningf("invalid phase %v for user %v \r\n", user.Spec.Phase, *user)
	}
}

func (ns *NodeSchedule) UpdateNodeUserEvent(user *api.NodeUser) {
	glog.V(5).Infof("add new node user\r\n")
}

func (ns *NodeSchedule) DelNodeUserEvent(user *api.NodeUser) {
	glog.V(5).Infof("delete node user not need implement")
}

//DirectDeleteNodeUser direct delete node user from etcd
func (ns *NodeSchedule) DirectDeleteNodeUser(nodeName, userName string) error {
	return nodectl.DelNodeUsers(nodeName, ns.helper, userName)
}

//DelUserFromNode add delete nodeuser to node,then node will be delete it server
func (ns *NodeSchedule) DelUserFromNode(nodeName string, userRefer api.UserReferences) error {
	phase := api.NodeUserPhase(api.NodeUserPhaseDelete)
	err := nodectl.AddNodeUserHelper(ns.helper, nodeName, userRefer, phase, nodeUserLease)
	if err != nil {
		return fmt.Errorf("delete user %v from node %v err %v", userRefer, nodeName, err)
	}
	return err
}

//BindUserToNode add  nodeuser to node,then node will be server for this user
func (ns *NodeSchedule) BindUserToNode(nodeReference map[string]api.UserReferences) error {

	phase := api.NodeUserPhase(api.NodeUserPhaseAdd)
	for nodeName, userRefer := range nodeReference {
		err := nodectl.AddNodeUserHelper(ns.helper, nodeName, userRefer, phase, nodeUserLease)
		if err != nil {
			return fmt.Errorf("add user %v to node %v err %v", userRefer, nodeName, err)
		}
	}

	return nil
}

//SyncUserServiceToNodeUser when node startup sync UserService into node
func (ns *NodeSchedule) SyncUserServiceToNodeUser(node api.Node) {
	nodeName := node.Name
	userlist, err := userctl.GetUserServicesByNodeName(ns.helper, nodeName)
	if err != nil {
		glog.Errorf("sync user into node %v failure %v\r\n", nodeName, err)
		return
	}

	for _, user := range userlist.Items {
		nodeRefer, ok := user.Spec.NodeUserReference[nodeName]
		if ok {
			node2User := make(map[string]api.UserReferences)
			node2User[nodeName] = nodeRefer.User
			err = ns.BindUserToNode(node2User)
			if err != nil {
				glog.Warningf("sync user to node %v failure %v", nodeName, err)
			}
		}
	}
}
