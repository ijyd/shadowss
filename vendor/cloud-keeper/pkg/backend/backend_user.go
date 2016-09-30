package backend

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend/db"
	"cloud-keeper/pkg/pagination"
)

func (u *Backend) GetUserByID(id int64) (*api.UserInfo, error) {
	return db.GetUserByID(u.StorageHandler, id)
}

func (u *Backend) GetUserByName(name string) (*api.UserInfo, error) {
	return db.GetUserByName(u.StorageHandler, name)
}

func (u *Backend) GetUsersByLimit(limitRecord int64, skip int64) ([]api.UserInfo, error) {
	return db.GetUsersByLimit(u.StorageHandler, limitRecord, skip)
}

func (u *Backend) GetUserList(page pagination.Pager) ([]api.UserInfo, error) {
	return db.GetUserList(u.StorageHandler, page)
}

func (u *Backend) CreateUser(info api.UserInfo) error {
	return db.CreateUser(u.StorageHandler, info)
}

func (u *Backend) DeleteUserByName(name string) error {
	return db.DeleteUserByName(u.StorageHandler, name)
}

func (u *Backend) UpdateUserNode(userID int64, nodes string) error {
	return db.UpdateUserNode(u.StorageHandler, userID, nodes)
}

func (u *Backend) UpdateUserPort(userID int64, port int64) error {
	return db.UpdateUserPort(u.StorageHandler, userID, port)
}

func (u *Backend) UpdateUserTraffic(userID int64, upload, download int64) error {
	return db.UpdateUserTraffic(u.StorageHandler, userID, upload, download)
}
