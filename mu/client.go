package main

import (
	muconfig "github.com/orvice/shadowsocks-go/mu/config"
	"github.com/orvice/shadowsocks-go/mu/mysql"
	"github.com/orvice/shadowsocks-go/mu/user"
	"os"
)

var Client *mysql.Client

func InitMySqlClient() error {
	conf := muconfig.GetConf().Mysql
	client := new(mysql.Client)
	dbType := "mysql"
	dbuser := conf.User
	password := conf.Pass
	host := conf.Host
	dbname := conf.Db
	table := conf.Table

	err := client.Boot(dbType, dbuser, password, host, dbname)
	if err != nil {
		return err
	}
	client.SetTable(table)
	Client = client
	mysql.SetClient(client)
	user.SetClient(client)
	if err != nil {
		Log.Error(err)
		os.Exit(0)
	}
}
