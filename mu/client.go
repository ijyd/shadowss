package main

import (
	"github.com/orvice/shadowsocks-go/mu/user"
)

func InitClient() {
	dbType := ""
	dbuser := ""
	password := ""
	host := "127.0.0.1:3306"
	dbname := ""
	table := "user"
	client := user.NewMysqlClient(dbType, dbuser, password, host, dbname,table)
	user.SetClient(client)
}
