package multiuser

import (
	"fmt"
	"strconv"
	"strings"

	"shadowss/pkg/api"
	"shadowss/pkg/util"

	"github.com/golang/glog"
)

type Node struct {
	APIServerHost string
	APIServerPort int64
}

func (mu *MultiUser) BuildNodeHelper() *api.Node {
	vpsIP, _ := mu.nodeAttr[api.NodeLablesVPSIP]
	vpsName, _ := mu.nodeAttr[api.NodeLablesVPSName]
	vpsLocation, _ := mu.nodeAttr[api.NodeLablesVPSLocation]
	vpsID, _ := mu.nodeAttr[api.NodeLablesVPSID]

	id, err := strconv.ParseInt(vpsID, 10, 64)
	if err != nil {
		id = 0
	}

	node := &api.Node{
		Spec: api.NodeSpec{
			Server: api.NodeServer{
				AccServerID:   id,
				AccServerName: vpsName,
				Location:      vpsLocation,
				Host:          vpsIP,
				Name:          mu.nodeName,
				Method:        "aes-256-cfb",
				EnableOTA:     true,
			},
		},
	}
	node.Name = mu.nodeName
	node.Annotations = map[string]string{
		api.NodeAnnotationUserCnt:    "0",
		api.NodeAnnotationVersion:    util.PrintVersion(),
		api.NodeAnnotationRefreshCnt: "0",
	}
	node.Labels = mu.nodeAttr

	return node

}

func (mu *MultiUser) refreshNode(loopcnt int64) {
	node := mu.BuildNodeHelper()
	if node == nil {
		glog.Errorf("invalid node configure\r\n")
		return
	}

	upload, download, usercnt, err := mu.CollectorAndUpdateUserTraffic()
	if err == nil {
		node.Spec.Server.Upload = upload
		node.Spec.Server.Download = download
	} else {
		glog.Warningf("collector user traffic error %v\r\n", err)
	}

	//it a bug, must update some field to keep node ttl?
	node.Annotations[api.NodeAnnotationRefreshCnt] = strconv.FormatInt(loopcnt, 10)
	node.Annotations[api.NodeAnnotationUserCnt] = strconv.FormatInt(usercnt, 10)
	node.Spec.Server.Status = 1
	err = UpdateNode(node, mu.ttl)
	if err != nil {
		glog.Errorf("refresh node(%+v) error :%v", *node, err)
	} else {
		glog.Infof("refresh node %+v\r\n", *node)
	}
}

func (mu *MultiUser) CollectorAndUpdateUserTraffic() (int64, int64, int64, error) {

	//userList := mu.userHandle.GetUsers()
	userList := mu.userHandle.GetUsersInfo()

	var upload, download, usercnt int64

	usercnt = int64(len(userList))

	for _, userInfo := range userList {
		userConfig := userInfo.ConnectInfo
		if userConfig.Name == string("") || strings.Contains(userConfig.Name, "apiproxy") {
			continue
		}

		nodeUser := &api.NodeUser{
			Spec: api.NodeUserSpec{
				User: api.UserReferences{
					ID:              userConfig.ID,
					Name:            userConfig.Name,
					Port:            int64(userConfig.Port),
					Method:          userConfig.EncryptMethod,
					Password:        userConfig.Password,
					EnableOTA:       userConfig.EnableOTA,
					UploadTraffic:   0,
					DownloadTraffic: 0,
				},
				NodeName: mu.nodeName,
				Phase:    api.NodeUserPhase(api.NodeUserPhaseUpdate),
			},
		}
		nodeUser.Kind = "NodeUser"
		nodeUser.APIVersion = "v1"
		nodeUser.Name = userConfig.Name

		mu.userHandle.UpdateTraffic(userConfig, nodeUser)
		nodeUser.Annotations = make(map[string]string)
		nodeUser.Annotations[api.UserFakeAnnotationLastActiveTime] = userInfo.LastActiveTime.String()

		err := UpdateNodeUserFromNode(nodeUser.Spec)
		if err != nil {
			glog.Errorf("update node user %+v err %v \r\n", nodeUser, err)
		} else {
			upload += nodeUser.Spec.User.UploadTraffic
			download += nodeUser.Spec.User.DownloadTraffic
		}
	}

	return upload, download, usercnt, nil
}

const (
	nodeUserLease = 10
)

//UpdateNodeUserFromNode this is support only for node call
func UpdateNodeUserFromNode(spec api.NodeUserSpec) error {
	nodeName := spec.NodeName
	userName := spec.User.Name
	//delete this node user if exist
	//nodectl.DelNodeUsers(nodeName, schedule.etcdHandle, userName)
	user := &api.NodeUser{
		Spec: spec,
	}
	user.Name = userName
	user.Spec.NodeName = nodeName

	user.Spec.Phase = api.NodeUserPhase(api.NodeUserPhaseUpdate)
	err := UpdateNodeUser(user)
	if err != nil {
		return fmt.Errorf("update nodeuser %v err %v", spec, err)
	}

	return nil
}

func RefreshUser(user *api.NodeUser, del bool) {
	if !del {
		//need update node user port
		err := UpdateNodeUserFromNode(user.Spec)
		if err != nil {
			glog.Errorf("update node user err %v \r\n", err)
		}
	}
}
