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
	Del(user.User) error
	ClearAll() error
	IncrSize(u user.User, size int) error
	GetSize(u user.User) (int64, error)
	SetSize(u user.User, size int) error
	MarkUserOnline(u user.User) error
	IsUserOnline(u user.User) bool
	GetOnlineUsersCount(u []user.User) int
}
