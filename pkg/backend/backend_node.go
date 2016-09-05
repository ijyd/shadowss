package backend

import "shadowsocks-go/pkg/backend/db"

func (u *Backend) GetNodesByUID(uid int64) ([]db.Node, error) {
	return db.GetNodesByUserID(u.StorageHandler, uid)
}
