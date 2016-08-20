package users

import (
	"shadowsocks-go/pkg/storage"
	"shadowsocks-go/pkg/storage/storagebackend"
	"shadowsocks-go/pkg/storage/storagebackend/factory"
)

//User is a mysql users map
type User struct {
	ID             int    `column:"id"`
	Port           int    `column:"port"`
	Passwd         string `column:"passwd"`
	Method         string `column:"method"`
	Enable         int    `column:"enable"`
	TransferEnable int    `column:"transfer_enable"`
	U              int    `column:"u"`
	D              int    `column:"d"`
}

const (
	userTableName = "user"
)

var userTableField = []string{"id", "passwd", "port", "method", "enable", "transfer_enable", "u", "d"}

func newStorage(c storagebackend.Config) (storage.Interface, error) {
	return factory.Create(c)
}

func get(handle storage.Interface) ([]User, error) {
	var users []User
	err := handle.GetToList(userTableName, userTableField, &users)
	return users, err
}
