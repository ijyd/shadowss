package config

import ()

var Conf = new(Config)

type Config struct {
	Mysql *MySql
	Redis *Redis
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
