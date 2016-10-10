package controller

import (
	"fmt"
	"time"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/userctl"
	"cloud-keeper/pkg/pagination"

	"github.com/golang/glog"
)

const (
	BillDay  = 1
	BillHour = 0
)

func (ns *NodeSchedule) cleanUserTraffic(users []api.UserInfo) {
	glog.V(5).Infof("clean user traffic %v\r\n", users)
	for _, user := range users {
		obj, err := userctl.GetUserService(ns.helper, user.Name)
		if err != nil {
			glog.Errorf("get user service error %v\r\n", err)
			continue
		}
		userSrv := obj.(*api.UserService)
		if userSrv.Spec.NodeUserReference == nil {
			continue
		}

		upload := int64(0)
		download := int64(0)
		totalUpload := user.TotalUploadTraffic
		totalDownload := user.TotalDownloadTraffic

		err = ns.be.UpdateUserTraffic(user.ID, totalUpload, totalDownload, upload, download)
		if err != nil {
			glog.Errorf("clean user traffic error %v\r\n", err)
		}

		if user.Status != 1 {
			err = userctl.UpdateUserServiceStatus(ns.helper, user.Name, true)
			if err != nil {
				glog.Errorf("update user service status to available status error %v \r\b", err)
			}

			err = ns.be.UpdateUserStatus(user.ID, true)
			if err != nil {
				glog.Errorf("update user status to available status error %v", err)
			}

			nodeRefer := make(map[string]api.UserReferences)
			for k, v := range userSrv.Spec.NodeUserReference {
				nodeRefer[k] = v.User
			}
			err = ns.BindUserToNode(nodeRefer)
			if err != nil {
				glog.Errorf("resume user bind  node error %v\r\n", err)
			}
		}

	}
}

func (ns *NodeSchedule) resumeUser() {
	reqPage := uint64(1)
	perPage := uint64(10)
	lastPage := uint64(1)

	pageParam := fmt.Sprintf("page=%v,perPage=%v", reqPage, perPage)

	glog.V(5).Infof("resume user start %v...\r\n", pageParam)
	pager, err := pagination.ParsePaginaton(pageParam)
	if err != nil {
		glog.Errorf("pagination error %v \r\n", err)
		return
	}

	user, err := ns.be.GetUserList(pager)
	if err != nil {
		glog.Errorf("got user list failure %v \r\n", err)
		return
	}

	if !pager.Empty() {
		has, last, _ := pager.LastPage()
		if has {
			lastPage = last
		}
	}

	//clean first page
	ns.cleanUserTraffic(user)
	if lastPage == 1 {
		return
	}

	for index := uint64(2); index <= lastPage; index++ {
		pageParam = fmt.Sprintf("page=%v,perPage=%v", index, perPage)
		pager, err = pagination.ParsePaginaton(pageParam)
		if err != nil {
			glog.Errorf("pagination error %v \r\n", err)
			continue
		}

		user, err := ns.be.GetUserList(pager)
		if err != nil {
			glog.Errorf("got user list failure %v \r\n", err)
			continue
		}
		ns.cleanUserTraffic(user)
	}

}

//resumeUserEachMonth clean user traffic and reconfig user node
func resumeUserEachMonth(ns *NodeSchedule) {
	year := time.Now().Year()
	month := time.Now().Year()

	var billYear, billMonth int
	if month == 12 {
		billYear = year + 1
		billMonth = 1
	} else {
		billYear = year
		billMonth = 1
	}

	bill := time.Date(billYear, time.Month(billMonth), BillDay, BillHour, 0, 0, 0, time.UTC)
	//test data
	// day := time.Now().Day()
	// hour := time.Now().Hour()
	// bill = time.Date(year, time.Month(month), day, hour, 10, 0, 0, time.UTC)

	duration := time.Since(bill)

	glog.Infof("install auto clean user traffic after %v \r\n", duration)

	time.AfterFunc(duration, func() {
		glog.Infof("clean trafifc and resume all users.....")
		ns.resumeUser()
	})

}

func (ns *NodeSchedule) exceedTrafficLimit(user api.UserInfo) error {
	err := userctl.UpdateUserServiceStatus(ns.helper, user.Name, false)
	if err != nil {
		glog.Errorf("user(%v) exceed traffic limit disable it failure %v\r\n", user.Name, err)
	}

	err = ns.be.UpdateUserStatus(user.ID, false)
	if err != nil {
		glog.Errorf("user(%v) exceed traffic limit disable it failure %v\r\n", user.Name, err)
	}

	err = ns.DelAllNodeUserByUser(user.Name)
	return err
}

func (ns *NodeSchedule) UpdateUserTraffic(userRefer api.UserReferences) error {
	userInfo, err := ns.be.GetUserByName(userRefer.Name)
	if err != nil {
		return err
	}

	upload := userInfo.UploadTraffic + userRefer.UploadTraffic
	download := userInfo.DownloadTraffic + userRefer.DownloadTraffic

	traffic := upload + download

	if traffic > userInfo.TrafficLimit {
		err = ns.exceedTrafficLimit(*userInfo)
		if err != nil {
			glog.Errorf("user %v exceed traffic(%v) limit disable it failure %v\r\n", userInfo, traffic, err)
		}
		return nil
	}

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
	if userSrv.Spec.NodeUserReference == nil {
		return nil, fmt.Errorf("not found user %v", userName)
	}

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
	if userSrv.Spec.NodeUserReference == nil {
		return fmt.Errorf("not found user %v", userName)
	}

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
	if userSrv.Spec.NodeUserReference == nil {
		return fmt.Errorf("not found user %v", userName)
	}

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
