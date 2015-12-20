package main

import (
	// "github.com/orvice/shadowsocks-go/mu/user"
	"github.com/orvice/shadowsocks-go/mu/mysql"
)

var Client *mysql.Client

func InitMySqlClient() error {
	client := new(mysql.Client)
	dbType := "mysql"
	dbuser := "sspanel"
	password := "sspanel"
	host := "localhost:3306"
	dbname := "sspanel"
	table := "user"

	err := client.Boot(dbType, dbuser, password, host, dbname)
	if err != nil {
		return err
	}
	client.SetTable(table)
	Client = client
	return nil
}
