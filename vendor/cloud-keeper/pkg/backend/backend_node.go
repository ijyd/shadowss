package backend

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend/db"
	"cloud-keeper/pkg/pagination"
)

func (u *Backend) GetNodesByUID(uid int64) ([]api.NodeServer, error) {
	return db.GetNodesByUserID(u.StorageHandler, uid)
}

func (u *Backend) GetNodes(page pagination.Pager) ([]api.NodeServer, error) {
	return db.GetNodes(u.StorageHandler, page)
}

func (u *Backend) GetNodeByName(name string) (*api.NodeServer, error) {
	return db.GetNodeByName(u.StorageHandler, name)
}

func (u *Backend) CreateNode(detail api.NodeServer) error {
	return db.CreateNode(u.StorageHandler, detail)
}

func (u *Backend) DeleteNode(name string) error {
	return db.DeleteNode(u.StorageHandler, name)
}

func (u *Backend) UpdateNode(detail api.NodeServer) error {
	return db.UpdateNode(u.StorageHandler, detail)
}

func (u *Backend) UpdateNodeTraffic(nodeID int64, totalUpload, totalDownload, upload, download int64) error {
	return db.UpdateNodeTraffic(u.StorageHandler, nodeID, totalUpload, totalDownload, upload, download)
}

func (u *Backend) UpdateNodeStatus(nodeID int64, status bool) error {
	return db.UpdateNodeStatus(u.StorageHandler, nodeID, status)
}
