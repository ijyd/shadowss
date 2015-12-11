package main

import (
	"github.com/orvice/shadowsocks-go/mu/user"
	log "github.com/Sirupsen/logrus"
)

func InitClient() {
	dbType := "mysql"
	dbuser := "sspanel"
	password := "sspanel"
	host := "127.0.0.1:3306"
	dbname := "sspanel"
	table := "user"
	client := user.NewMysqlClient()
	err := client.Boot(dbType, dbuser, password, host, dbname)
	if err != nil {
		log.Panic(err)
	}
	client.SetTable(table)
	user.SetClient(client)
}
