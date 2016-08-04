package main

import (
	"os"

	muconfig "shadowsocks/shadowsocks-go/mu/config"
	"shadowsocks/shadowsocks-go/mu/mysql"
	"shadowsocks/shadowsocks-go/mu/user"
	webapi "shadowsocks/shadowsocks-go/mu/webapi"
)

func InitMySqlClient() error {
	initMysqlClientConfig()
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
	mysql.SetClient(client)
	user.SetClient(client)
	if err != nil {
		Log.Error(err)
		os.Exit(0)
	}
	return nil
}

func InitWebApi() error {
	err := initWebApiConfig()
	if err != nil {
		Log.Error(err)
		os.Exit(0)
	}
	conf := muconfig.GetConf().WebApi
	webapi.SetClient(webapi.NewClient())
	webapi.SetBaseUrl(conf.Url)
	webapi.SetKey(conf.Key)
	webapi.SetNodeId(conf.NodeId)
	user.SetClient(webapi.GetClient())
	return nil
}
