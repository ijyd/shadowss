package main

import (
	"github.com/Terry-Mao/goconf"
	muconfig "shadowsocks/shadowsocks-go/mu/config"
)

var (
	conf = goconf.New()
)

func initMysqlClientConfig() error {
	mysql := new(muconfig.MySql)
	if err := conf.Unmarshal(mysql); err != nil {
		return err
	}
	muconfig.Conf.SetMysql(mysql)
	return nil
}

func initWebApiConfig() error {
	webapi := new(muconfig.WebApi)
	if err := conf.Unmarshal(webapi); err != nil {
		return err
	}
	muconfig.Conf.SetWebApi(webapi)
	return nil
}

func InitConfig() error {

	if err := conf.Parse("config.conf"); err != nil {
		return err
	}

	redis := new(muconfig.Redis)
	base := new(muconfig.Base)

	if err := conf.Unmarshal(redis); err != nil {
		return err
	}

	if err := conf.Unmarshal(base); err != nil {
		return err
	}
	muconfig.Conf.SetRedis(redis)
	muconfig.Conf.SetBase(base)
	return nil
}
