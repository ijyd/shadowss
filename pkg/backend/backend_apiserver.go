package backend

import "shadowsocks-go/pkg/backend/db"

func (u *Backend) GetAPIServer() ([]db.APIServers, error) {
	return db.GetApiServers(u.StorageHandler)
}

func (u *Backend) CreateAPIServer(host string, port int64, isEnable bool) error {
	return db.CreateAPIServer(u.StorageHandler, host, port, isEnable)
}
