package main

import (
	muconfig "github.com/orvice/shadowsocks-go/mu/config"
	"gopkg.in/redis.v3"
)

var Redis = new(RedisClient)

type RedisClient struct{
	Client *redis.Client
}

func(r * RedisClient) SetClient(client *redis.Client){
	r.Client = client
}

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
	Redis.SetClient(client)
	return nil
}
