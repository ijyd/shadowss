package main

import (
	muconfig "github.com/orvice/shadowsocks-go/mu/config"
	"gopkg.in/redis.v3"
)

func InitRedis() error {
	client := redis.NewClient(&redis.Options{
		Addr:     muconfig.Conf.Redis.Host,
		Password: muconfig.Conf.Redis.Pass, // no password set
		DB:       muconfig.Conf.Redis.Db,   // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return err
	}
	Log.Info(pong)
	return nil
}
