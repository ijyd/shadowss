package main

import (
	"github.com/Terry-Mao/goconf"
	muconfig "github.com/orvice/shadowsocks-go/mu/config"
)

func InitConfig() error {
	conf := goconf.New()
	if err := conf.Parse("config.conf"); err != nil {
		return err
	}
	mysql := new(muconfig.MySql)
	redis := new(muconfig.Redis)
	if err := conf.Unmarshal(mysql); err != nil {
		return err
	}
	if err := conf.Unmarshal(redis); err != nil {
		return err
	}
	Log.Info(mysql, redis)
	muconfig.Conf.SetMysql(mysql)
	muconfig.Conf.SetRedis(redis)
	return nil
}
