package backend

import "cloud-keeper/pkg/backend/db"

func (u *Backend) AddUserToken(token string, uid int64, name string) error {
	return db.CreateToken(u.StorageHandler, token, uid, name)
}

func (u *Backend) GetTokenIDByUID(uid int64) (int64, error) {
	return db.GetTokenIDByUID(u.StorageHandler, uid)
}

func (u *Backend) GetUserIDByToken(token string) (int64, error) {
	return db.CheckToken(u.StorageHandler, token)
}

func (u *Backend) UpdateToken(token string, id int64) error {
	return db.UpdateToken(u.StorageHandler, token, id)
}
