package config

import (
	"time"
)

var Conf = new(Config)

type Config struct {
	Mysql *MySql
	Redis *Redis
	Base  *Base
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

type Base struct {
	N         float32       `goconf:"base:N"`
	CheckTime time.Duration `goconf:"base:checktime"`
	SyncTime  time.Duration `goconf:"base:synctime"`
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
