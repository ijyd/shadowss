package backend

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/backend/db"
	"cloud-keeper/pkg/pagination"
)

func (u *Backend) GetAccounts(page pagination.Pager) ([]api.AccountDetail, error) {
	return db.GetAccounts(u.StorageHandler, page)
}

func (u *Backend) CreateAccount(acc api.AccountDetail) error {
	return db.CreateAccount(u.StorageHandler, acc)
}

func (u *Backend) DeleteAccount(name string) error {
	return db.DeleteAccount(u.StorageHandler, name)
}

func (u *Backend) GetAccountByname(name string) (*api.AccountDetail, error) {
	return db.GetAccountByname(u.StorageHandler, name)
}
