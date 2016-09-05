package backend

import "shadowsocks-go/pkg/backend/db"

func (u *Backend) AddUserToken(token string, uid int64, macAddr string) error {
	return db.CreateToken(u.StorageHandler, token, uid, macAddr)
}
