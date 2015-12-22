package main

import (
	"github.com/orvice/shadowsocks-go/mu/user"
)

var storage Storage

func SetStorage(s Storage) {
	storage = s
}

type Storage interface {
	GetUserInfo(user.User) (user.UserInfo, error)
	StoreUser(user.UserInfo) error
	Exists(user.User) (bool, error)
}
