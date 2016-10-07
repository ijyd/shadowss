package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/userctl"
	"fmt"

	"github.com/golang/glog"
)

func (ns *NodeSchedule) UpdateUserTraffic(userRefer api.UserReferences) error {
	userInfo, err := ns.be.GetUserByName(userRefer.Name)
	if err != nil {

	}

	upload := userInfo.UploadTraffic + userRefer.UploadTraffic
	download := userInfo.DownloadTraffic + userRefer.DownloadTraffic
	totalUpload := userInfo.TotalUploadTraffic + userRefer.UploadTraffic
	totalDownload := userInfo.TotalDownloadTraffic + userRefer.DownloadTraffic

	return ns.be.UpdateUserTraffic(userRefer.ID, totalUpload, totalDownload, upload, download)
}

func (ns *NodeSchedule) DelUserService(name string) error {
	return userctl.DelUserService(ns.helper, name)
}

func (ns *NodeSchedule) checkUserServiceNode(userName, nodeName string) (*api.UserService, error) {
	obj, err := userctl.GetUserService(ns.helper, userName)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return nil, err
	}
	userSrv := obj.(*api.UserService)

	_, ok := userSrv.Spec.NodeUserReference[nodeName]
	if ok {
		return userSrv, nil
	} else {
		glog.Errorf("not found any user by name %v\r\n", userName)
		return nil, fmt.Errorf("not found user %v by node %v \r\n", userName, nodeName)
	}
}

func (ns *NodeSchedule) DeleteUserServiceNode(userName, nodeName string) error {
	userSrv, err := ns.checkUserServiceNode(userName, nodeName)
	if err != nil {
		return err
	}

	delete(userSrv.Spec.NodeUserReference, nodeName)

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update %+v error %v\r\n", userSrv, err)
		return err
	}
	return nil
}

//UpdateNewUserServiceSpec Replace olde spec with new spec
func (ns *NodeSchedule) UpdateNewUserServiceSpec(nodes map[string]api.NodeReferences, userName string) error {
	obj, err := userctl.GetUserService(ns.helper, userName)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return err
	}
	userSrv := obj.(*api.UserService)

	userSrv.Spec.NodeUserReference = nodes

	glog.V(5).Infof("update user %+v for node %+v\r\n", userName, nodes)

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update user %+v failure %v\r\n", *userSrv, err)
		return err
	}

	return nil
}

func (ns *NodeSchedule) BindUserServiceWithNode(userName string, nodes map[string]api.UserReferences) error {
	obj, err := userctl.GetUserService(ns.helper, userName)
	if err != nil {
		glog.Errorf("get user service error %v\r\n", err)
		return err
	}
	userSrv := obj.(*api.UserService)

	for nodeName, userRefer := range nodes {
		node, err := ns.GetActiveNodeByName(nodeName)
		if err != nil {
			glog.Errorf("check node(%v) error %v\r\n", nodeName, err)
			return err
		}

		nodeRefer := api.NodeReferences{
			Host: node.Spec.Server.Host,
			User: userRefer,
		}
		//update user spec
		userSrv.Spec.NodeUserReference[nodeName] = nodeRefer
	}

	glog.V(5).Infof("update user %+v for node %+v\r\n", userName, nodes)

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update user %+v failure %v\r\n", *userSrv, err)
		return err
	}

	return ns.BindUserToNode(nodes)
}

//UpdateUserDynamicData dynamic data contains: user traffic, user actual port
//need sync user service.
func (ns *NodeSchedule) UpdateUserDynamicDataByNodeUser(nodeUser *api.NodeUser) error {

	userName := nodeUser.Name
	nodeName := nodeUser.Spec.NodeName
	newUserRefer := nodeUser.Spec.User

	return ns.UpdateUserDynamicData(userName, nodeName, newUserRefer)
}

func (ns *NodeSchedule) UpdateUserDynamicData(userName, nodeName string, userRefer api.UserReferences) error {
	newUserRefer := userRefer

	userSrv, err := ns.checkUserServiceNode(userName, nodeName)
	if err != nil {
		return err
	}

	nodeUserRefer, ok := userSrv.Spec.NodeUserReference[nodeName]
	if ok {
		err = ns.UpdateUserTraffic(newUserRefer)
		if err != nil {
			glog.Warningf("collected user(%v) traffic failure", userName)
		}
	} else {
		glog.Errorf("not found any user by name %v\r\n", userName)
		return fmt.Errorf("not found user %v by node %v \r\n", userName, nodeName)
	}

	glog.V(5).Infof("update user  %+v to %+v\r\n", *userSrv, userRefer)
	nodeUserRefer.User = newUserRefer
	userSrv.Spec.NodeUserReference[nodeName] = nodeUserRefer

	err = userctl.UpdateUserService(ns.helper, userSrv)
	if err != nil {
		glog.Errorf("update %+v error %v\r\n", userSrv, err)
		return err
	}

	return nil
}

func (ns *NodeSchedule) DelAllNodeUserByUser(name string) error {
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
		err := ns.DelUserFromNode(nodeName, nodeRefer.User)
		if err != nil {
			glog.Errorf("del user %+v from node %+v error %v", nodeRefer.User.Name, nodeName, err)
		}
	}

	return nil
}
