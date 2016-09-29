package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend"
	"cloud-keeper/pkg/controller/nodectl"
	"cloud-keeper/pkg/controller/userctl"
	"cloud-keeper/pkg/etcdhelper"
	"fmt"
	"gofreezer/pkg/api/prototype"
	"gofreezer/pkg/storage"
	"gofreezer/pkg/watch"
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

	nodeObj, err := nodectl.GetNode(ns.helper, user.Spec.NodeName)
	if err != nil {
		glog.Errorf("get node  error %v\r\n", err)
		return
	}

	node := nodeObj.(*api.Node)

	obj, err := userctl.GetUserService(ns.helper, user.Name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return
	}

	userSrv := obj.(*api.UserService)

	if len(userSrv.Name) == 0 {
		glog.Errorf("user %v not found\r\n", *user)
		return
	}

	if userSrv.Spec.NodeUserReference == nil {
		glog.Errorf("user %v invalid\r\n", *user)
		return
	}

	userSrv.Spec.NodeUserReference[user.Spec.NodeName] = api.NodeReferences{
		User: user.Spec.User,
		Host: node.Spec.Server.Host,
	}

	glog.V(5).Infof("add new user %+v for node %v\r\n", userSrv.Spec.NodeUserReference, user.Spec.NodeName)

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update user %+v failure %v\r\n", *userSrv, err)
	}

}

func (ns *NodeSchedule) UpdateNode(node *api.Node) {

	glog.V(5).Infof("udpate node %+v\r\n", *node)

	err := ns.UpdateNodeTraffic(node)
	if err != nil {
		glog.Errorf("update node traffic error %v\r\n", err)
	}

}

func (ns *NodeSchedule) UpdateNodeUser(user *api.NodeUser) {
	obj, err := userctl.GetUserService(ns.helper, user.Name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return
	}

	userSrv := obj.(*api.UserService)

	nodeName := user.Spec.NodeName
	nodeRefer, ok := userSrv.Spec.NodeUserReference[nodeName]
	if ok {
		err = ns.UpdateUserTraffic(nodeRefer.User)
		if err != nil {
			glog.Warningf("collected user(%v) traffic failure", user.Name)
		}

		glog.V(5).Infof("update user  %+v to %+v\r\n", *userSrv, *user)
		nodeRefer.User = user.Spec.User
		userSrv.Spec.NodeUserReference[nodeName] = nodeRefer
	} else {
		glog.Errorf("not found any user by name %v\r\n", user.Name)
		return
	}

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

	nodeName := user.Spec.NodeName
	nodeRefer, ok := userSrv.Spec.NodeUserReference[nodeName]
	if ok {
		err = ns.UpdateUserTraffic(nodeRefer.User)
		if err != nil {
			glog.Warningf("collected user(%v) traffic failure", user.Name)
		}
	} else {
		glog.Errorf("not found any user by name %v\r\n", user.Name)
		return
	}

	delete(userSrv.Spec.NodeUserReference, nodeName)
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

func (ns *NodeSchedule) DelUserFromNode(nodeName, userName string) error {
	return nodectl.DelNodeUsers(nodeName, ns.helper, userName)
}

//CleanUser to delete user for nodeuser
func (ns *NodeSchedule) CleanNodeUser(name string) error {
	obj, err := userctl.GetUserService(ns.helper, name)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return err
	}

	userSrv := obj.(*api.UserService)

	if userSrv.Spec.NodeUserReference == nil {
		glog.Errorf("not found any user by name %v\r\n", name)
		return fmt.Errorf("not have any node for this user")
	}

	for nodeName, nodeRefer := range userSrv.Spec.NodeUserReference {

		err := ns.DelUserFromNode(nodeName, nodeRefer.User.Name)
		if err != nil {
			glog.Errorf("del user %+v from node %+v error %v", nodeRefer.User.Name, nodeName, err)
		}
	}

	return nil
}

func (ns *NodeSchedule) BindUserToNode(nodeReference map[string]api.UserReferences) error {
	for nodeName, userRefer := range nodeReference {
		err := nodectl.AddNodeUserHelper(ns.helper, nodeName, userRefer)
		if err != nil {
			return fmt.Errorf("add user %v to node %v err %v", userRefer, nodeName, err)
		}
	}

	return nil
}

//AllocDefaultNode search default node for user
func (ns *NodeSchedule) AllocDefaultNode(name string) error {

	dbUserInfo, err := ns.be.GetUserByName(name)
	if err != nil || name != dbUserInfo.Name {
		return fmt.Errorf("invalid user %v error %v", dbUserInfo, err)
	}
	userInfo := *dbUserInfo

	nodeList := ns.findAPINode(name)
	if len(nodeList) == 0 {
		return fmt.Errorf("not have enough node")
	}

	//var successNode []string
	nodeReference := make(map[string]api.UserReferences)
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

		nodeReference[v] = user
	}

	err = ns.BindUserToNode(nodeReference)
	if err != nil {
		glog.Errorf("alloc default user err %v", err)
		return fmt.Errorf("alloc default user err %v", err)
	}

	return nil
}

func (ns *NodeSchedule) AllocNodeByUserProperties(name string, properties map[string]string) error {

	dbUserInfo, err := ns.be.GetUserByName(name)
	if err != nil || name != dbUserInfo.Name {
		return fmt.Errorf("invalid user %v error %v", dbUserInfo, err)
	}
	userInfo := *dbUserInfo

	nodeList := ns.findNodeByUserProperties(name, properties)
	if len(nodeList) == 0 {
		return fmt.Errorf("not have enough node")
	}

	//var successNode []string
	nodeReference := make(map[string]api.UserReferences)
	userRefer := api.UserReferences{
		ID:        userInfo.ID,
		Name:      userInfo.Name,
		Port:      0,
		Method:    string("aes-256-cfb"),
		Password:  userInfo.Passwd,
		EnableOTA: true,
	}
	for _, v := range nodeList {
		//create a user for node
		nodeReference[v] = userRefer
	}

	err = ns.BindUserToNode(nodeReference)
	if err != nil {
		glog.Errorf("alloc default user err %v", err)
		return fmt.Errorf("alloc default user err %v", err)
	}

	return nil
}

func (ns *NodeSchedule) findAPINode(userName string) []string {

	var nodeName []string

	objList, err := nodectl.GetAllNodes(ns.helper)
	if err != nil {
		glog.Errorf("get all node failure %v\r\n", err)
		return nil
	}

	nodeList := objList.(*api.NodeList)
	glog.V(5).Infof("check all node %v \r\n", nodeList)

	for _, v := range nodeList.Items {
		userSpace, ok := v.Labels[api.NodeLablesUserSpace]
		if ok && userSpace == api.NodeUserSpaceAPI {
			nodeName = append(nodeName, v.Name)
		} else {
			glog.Warningf("node  has not the API value in userSpace %v \r\n", v)
		}
	}

	return nodeName
}

func (ns *NodeSchedule) findNodeByUserProperties(userName string, properties map[string]string) []string {

	var nodeName []string

	objList, err := nodectl.GetAllNodes(ns.helper)
	if err != nil {
		glog.Errorf("get all node failure %v\r\n", err)
		return nil
	}

	nodeList := objList.(*api.NodeList)
	userISP, _ := properties[api.NodeLablesChinaISP]

	glog.V(5).Infof("check user isp(%v) with node %v \r\n", userISP)

	for _, v := range nodeList.Items {
		//check this node is default node
		userSpace, ok := v.Labels[api.NodeLablesUserSpace]
		if ok && userSpace == api.NodeUserSpaceDefault {
			//check user number on this node
			cnt, ok := v.Annotations[nodectl.NodeAnnotationUserCnt]
			if ok {
				if cnt, err := strconv.ParseUint(cnt, 10, 32); err == nil && cnt < 80 {
					cnISP, ok := v.Labels[api.NodeLablesChinaISP]
					if ok && cnISP == userISP {
						nodeName = append(nodeName, v.Name)
					}
				}
			}
			nodeName = append(nodeName, v.Name)
		} else {
			glog.Warningf("not got Annotations with node %v \r\n", v)
		}

	}

	return nodeName
}

func (ns *NodeSchedule) UpdateNodeTraffic(node *api.Node) error {
	err := ns.be.UpdateNodeTraffic(node.Spec.Server.ID, node.Spec.Server.Upload, node.Spec.Server.Download)
	if err != nil {
		glog.Errorf("delete %+v error %v\r\n", node, err)
		return err
	}

	return nil
}

func manageNode(helper *etcdhelper.EtcdHelper) {
	watchKey := nodectl.PrefixNode
	ctx := prototype.NewContext()
	resourceVer := string("0")

	glog.V(5).Infof("watch at %v with resource %v", watchKey, resourceVer)
	watcher, err := helper.StorageCodec.Storage.WatchList(ctx, watchKey, resourceVer, storage.Everything)

	if err != nil {
		glog.Fatalf("Unexpected error: %v", err)
	}
	defer watcher.Stop()

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				glog.Errorf("Unexpected channel close")
				return
			}

			switch event.Type {
			case watch.Added:
				glog.V(5).Infof("Got Add  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.NewNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.NewNodeUser(gotObject)
					}
				}
			case watch.Modified:
				glog.V(5).Infof("Got modify  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.UpdateNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.UpdateNodeUser(gotObject)
					}
				}
			case watch.Deleted:
				glog.V(5).Infof("Got Deleted  got: %#v", event.Object)
				gotObject, ok := event.Object.(*api.Node)
				if ok {
					AutoSchedule.DelNode(gotObject)
				} else {
					gotObject, ok := event.Object.(*api.NodeUser)
					if ok {
						AutoSchedule.DelNodeUser(gotObject)
					}
				}
			case watch.Error:
				glog.V(5).Infof("Got Error  got: %#v", event.Object)
				return
			default:
				glog.Errorf("UnExpected: %#v, got: %#v", event.Type, event.Object)
			}

		}
	}
}
