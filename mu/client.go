package main

import (
	"github.com/orvice/shadowsocks-go/mu/user"
)

func InitClient() {
	client := user.NewMysqlClient()
	user.SetClient(client)
}
