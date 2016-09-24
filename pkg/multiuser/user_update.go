package multiuser

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/nodectl"

	"github.com/golang/glog"
)

// func RefreshTraffic(config *config.ConnectionInfo, user *api.NodeUser, proxy *proxyserver.Servers) error {
// 	//update users traffic
// 	upload, download, err := proxy.GetTraffic(config)
// 	if err != nil {
// 		return err
// 	}
//
// 	totalUpload := int64(user.Spec.User.UploadTraffic) + upload
// 	totalDownlaod := int64(user.Spec.User.DownloadTraffic) + download
//
// 	Schedule.etcdHandle.UpdateNodeUsers()
// 	// return db.UpdateUserTraffic(u.StorageHandler, config.ID, totalUpload, totalDownlaod)
// 	return nil
// }

func RefreshUser(user *api.NodeUser, del bool) {
	if !del {
		//need update noe user port
		glog.V(5).Infof("update node user %+v", *user)
		err := nodectl.UpdateNodeUsersRefer(schedule.etcdHandle, user.Spec)
		if err != nil {
			glog.Errorf("update node user err %v \r\n", err)
		}
	}

	_, err := nodectl.UpdateNodeAnotationsUserCnt(schedule.etcdHandle, user.Spec.NodeName, del)
	if err != nil {
		glog.Errorf("update node anotation err %v \r\n", err)
	}

}
