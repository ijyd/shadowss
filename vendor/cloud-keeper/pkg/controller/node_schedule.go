package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/controller/userctl"
	"cloud-keeper/pkg/etcdhelper"
	"fmt"
	"strconv"

	"github.com/golang/glog"
)

type NodeSchedule struct {
	helper *etcdhelper.EtcdHelper
	be     *backend.Backend
}

func NewNodeSchedule(helper *etcdhelper.EtcdHelper, be *backend.Backend) *NodeSchedule {
	return &NodeSchedule{
		helper: helper,
		be:     be,
	}
}

func (ns *NodeSchedule) NewNode(node *api.Node) {
	if node.Name == "" {
		glog.Errorf("invalid node name\r\n")
		return
	}

	_, err := nodectl.GetNode(ns.helper, node.Name)
	if err != nil {
		glog.Errorf("get node  error %v\r\n", err)
		return
	}

	var mysql bool
	_, err = nodectl.GetNodeFromDB(ns.be, node.Name)
	if err != nil && err.Error() == "not found" {
		mysql = true
	} else if err != nil {
		glog.Errorf("get node from db error %v\r\n", err)
		return
	} else {
		mysql = false
	}

	if mysql {
		_, err := nodectl.AddNode(ns.be, nil, node, 0, false, true)
		if err != nil {
			glog.Errorf("add node to db error %v\r\n", err)
			return
		}
	}

}

func (ns *NodeSchedule) NewNodeUser(user *api.NodeUser) {

	if user.Spec.NodeName == "" {
		glog.Errorf("ignore this user not have node name\r\n")
		return
	}

	obj, err := userctl.GetUserService(ns.helper, user.Name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return
	}

	userSrv := obj.(*api.UserService)

	var notfound bool
	if userSrv.Spec.NodeUserReference == nil {
		userSrv.Spec.NodeUserReference = make(map[string]api.UserReferences)
		notfound = true
	}

	userSrv.Spec.NodeUserReference[user.Spec.NodeName] = user.Spec.User

	if notfound {
		err = userctl.AddUserServiceHelper(ns.helper, user.Name, userSrv.Spec.NodeUserReference)
	} else {
		err = userctl.UpdateUserService(ns.helper, userSrv)
	}

}

func (ns *NodeSchedule) UpdateNode(node *api.Node) {

	glog.V(5).Infof("udpate node %v\r\n", node)

	_, err := nodectl.UpdateNode(ns.be, nil, node, false, true)
	if err != nil {
		glog.Errorf("update node %+v error %v\r\n", node, err)
	}
}

func (ns *NodeSchedule) UpdateNodeUser(user *api.NodeUser) {
	obj, err := userctl.GetUserService(ns.helper, user.Name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return
	}

	userSrv := obj.(*api.UserService)

	if userSrv.Spec.NodeUserReference == nil {
		glog.Errorf("not found any user by name %v\r\n", user.Name)
		return
	}

	userSrv.Spec.NodeUserReference[user.Spec.NodeName] = user.Spec.User

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update %+v error %v\r\n", userSrv, err)
	}

}

func (ns *NodeSchedule) DelNodeUser(user *api.NodeUser) {
	obj, err := userctl.GetUserService(ns.helper, user.Name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return
	}

	userSrv := obj.(*api.UserService)

	if userSrv.Spec.NodeUserReference == nil {
		glog.Errorf("not found any user by name %v\r\n", user.Name)
		return
	}

	delete(userSrv.Spec.NodeUserReference, user.Name)

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update %+v error %v\r\n", userSrv, err)
	}

}

func (ns *NodeSchedule) DelNode(node *api.Node) {
	err := nodectl.DelNode(ns.be, nil, node.Name, false, true)
	if err != nil {
		glog.Errorf("delete %+v error %v\r\n", node, err)
	}
}

//AllocNode you need to delete node first
//will update user NodeReferences
func (ns *NodeSchedule) AllocNode(user *api.User) error {
	//check if user have node
	// if v, _ := ns.user2node[user.Spec.DetailInfo.Name]; v.nodeCnt != 0 {
	// 	return fmt.Errorf("already have node")
	// }

	userInfo := user.Spec.DetailInfo
	dbUserInfo, err := ns.be.GetUserByName(user.Name)
	if err != nil || userInfo.Name != dbUserInfo.Name {
		return fmt.Errorf("invalid user %v error %v", dbUserInfo, err)
	}
	userInfo = *dbUserInfo

	nodeList := ns.findIdleNode(user.Name, false)
	if len(nodeList) == 0 {
		return fmt.Errorf("not have enough node")
	}

	var successNode []string
	for _, v := range nodeList {
		//create a user for node

		user := api.UserReferences{
			ID:        userInfo.ID,
			Name:      userInfo.Name,
			Port:      0,
			Method:    string("aes-256-cfb"),
			Password:  userInfo.Passwd,
			EnableOTA: true,
		}

		glog.V(5).Infof("Add user %+v to node %v\r\n", user, v)
		err = nodectl.AddNodeUserHelper(ns.helper, v, user)
		if err != nil {
			glog.Errorf("add user %+v to node %+v error %v", user, v, err)
			return fmt.Errorf("add user %v to node %v err %v", user, v, err)
		} else {
			successNode = append(successNode, v)
		}
	}

	return nil
}

func (ns *NodeSchedule) findIdleNode(userName string, defaultNode bool) []string {

	var nodeName []string

	objList, err := nodectl.GetAllNodes(ns.helper)
	if err != nil {
		glog.Errorf("get all node failure %v\r\n", err)
		return nil
	}

	nodeList := objList.(*api.NodeList)
	glog.V(5).Infof("check all node %v \r\n", nodeList)

	for _, v := range nodeList.Items {
		cnt, ok := v.Annotations[nodectl.NodeAnnotationUserCnt]
		if ok {
			if cnt, err := strconv.ParseUint(cnt, 10, 32); err == nil && cnt < 50 {
				nodeName = append(nodeName, v.Name)
			}
		} else {
			glog.Warningf("not got Annotations with node %v \r\n", v)
		}
	}

	return nodeName
}
