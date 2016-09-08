package backend

import "shadowsocks-go/pkg/backend/db"

func (u *Backend) GetAPIServer() ([]db.APIServers, error) {
	return db.GetApiServers(u.StorageHandler)
}

func (u *Backend) CreateAPIServer(name string, host string, port int64, isEnable bool) error {
	return db.CreateAPIServer(u.StorageHandler, name, host, port, isEnable)
}

func (u *Backend) DeleteAPIServerByID(id int64) error {
	return db.DeleteAPIServerByID(u.StorageHandler, id)
}
