package config

import (
	"time"
)

var Conf = new(Config)

type Config struct {
	WebApi *WebApi
	Mysql  *MySql
	Redis  *Redis
	Base   *Base
}

func GetConf() *Config {
	return Conf
}

func (c *Config) SetMysql(m *MySql) {
	c.Mysql = m
}

func (c *Config) SetRedis(r *Redis) {
	c.Redis = r
}

func (c *Config) SetBase(b *Base) {
	c.Base = b
}

func (c *Config) SetWebApi(w *WebApi) {
	c.WebApi = w
}

type Base struct {
	N         float32       `goconf:"base:N"`
	IP        string        `goconf:"base:ip"`
	Client    string        `goconf:"base:client"`
	CheckTime time.Duration `goconf:"base:checktime"`
	SyncTime  time.Duration `goconf:"base:synctime"`
}

type WebApi struct {
	Url    string `goconf:"webapi:url"`
	Key    string `goconf:"webapi:key"`
	NodeId int    `goconf:"webapi:node_id"`
}

type MySql struct {
	Host  string `goconf:"mysql:host"`
	User  string `goconf:"mysql:user"`
	Pass  string `goconf:"mysql:pass"`
	Db    string `goconf:"mysql:db"`
	Table string `goconf:"mysql:table"`
}

type Redis struct {
	Host string `goconf:"redis:host"`
	Pass string `goconf:"redis:pass"`
	Db   int64  `goconf:"redis:db"`
}
