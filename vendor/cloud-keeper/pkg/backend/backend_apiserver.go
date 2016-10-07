package backend

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend/db"
	"cloud-keeper/pkg/pagination"
)

func (u *Backend) GetAPIServer(page pagination.Pager) ([]api.APIServerInfor, error) {
	return db.GetApiServers(u.StorageHandler, page)
}

func (u *Backend) CreateAPIServer(info api.APIServerInfor) error {
	return db.CreateAPIServer(u.StorageHandler, info)
}

func (u *Backend) DeleteAPIServerByName(name string) error {
	return db.DeleteAPIServerByName(u.StorageHandler, name)
}
