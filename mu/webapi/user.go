package web

import (
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

type User struct {
	id             int
	port           int
	passwd         string
	method         string
	enable         int
	transferEnable int `json:"transfer_enable"`
	u              int
	d              int
}

func (u User) GetPort() int {
	return u.id
}

func (u User) GetPasswd() string {
	return u.passwd
}

func (u User) GetMethod() string {
	return u.method
}

func (u User) IsEnable() bool {
	return true
}

func (u User) GetCipher() (*ss.Cipher, error) {
	return ss.NewCipher(u.method, u.passwd)
}

func (u User) UpdateTraffic(storageSize int) error {
	return nil
}

func (u User) GetUserInfo() user.UserInfo {
	var user user.UserInfo
	return user
}
