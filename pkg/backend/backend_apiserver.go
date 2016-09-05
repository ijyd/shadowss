package backend

import "shadowsocks-go/pkg/backend/db"

func (u *Backend) GetAPIServer() ([]db.APIServers, error) {
	return db.GetApiServers(u.StorageHandler)
}
