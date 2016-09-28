package controller

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/controller/userctl"
)

func (ns *NodeSchedule) UpdateUserTraffic(userRefer api.UserReferences) error {
	return ns.be.UpdateUserTraffic(userRefer.ID, userRefer.UploadTraffic, userRefer.DownloadTraffic)
}

func (ns *NodeSchedule) DelUserService(name string) error {
	return userctl.DelUserService(ns.helper, name)
}
